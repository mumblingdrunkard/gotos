package cpu

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	coreStateRunning = 2 // Core is running and executing instructions
	coreStateHalting = 1 // Core is running and executing instructions, but will turn off in the next cycle
	coreStateHalted  = 0 // Core is halted and will stay halted until Start() is called
)

// Register mnemonics
const (
	Reg_ZERO = 0  // Hard-wired zero
	Reg_RA   = 1  // Return address
	Reg_SP   = 2  // Stack pointer
	Reg_GP   = 3  // Global pointer
	Reg_TP   = 4  // Thread pointer
	Reg_T0   = 5  // Temporary/alternate link register
	Reg_T1   = 6  // Temporaries
	Reg_T2   = 7  //
	Reg_S0   = 8  // Saved register/frame pointer
	Reg_FP   = 8  //
	Reg_S1   = 9  // Saved register
	Reg_A0   = 10 // Function arguments/return values
	Reg_A1   = 11 //
	Reg_A2   = 12 //
	Reg_A3   = 13 //
	Reg_A4   = 14 //
	Reg_A5   = 15 //
	Reg_A6   = 16 //
	Reg_A7   = 17 //
	Reg_S2   = 18 // Saved registers
	Reg_S3   = 19 //
	Reg_S4   = 20 //
	Reg_S5   = 21 //
	Reg_S6   = 22 //
	Reg_S7   = 23 //
	Reg_S8   = 24 //
	Reg_S9   = 25 //
	Reg_S10  = 26 //
	Reg_S11  = 27 //
	Reg_T3   = 28 // Temporaries
	Reg_T4   = 29 //
	Reg_T5   = 30 //
	Reg_T6   = 31 //
)

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	bcm sync.Mutex
	// The big core mutex (bcm) ensures that only one goroutine is inside the fetch-decode-execute loop at any one time
	state            atomic.Value // can be HALTED, HALTING, or RUNNING
	vmaUpd           atomic.Value
	jumped           bool
	iAmLockingMemory bool
	counter          counter
	system           System
	mc               memoryController
	// normal registers (save on context switch)
	reg  [32]uint32 // normal registers
	freg [32]uint64 // fp registers
	pc   uint32     // program counter
	// CSRs
	csr [4096]uint32 // control and status registers
	// function pointers
	// counters, used for predictable scheduling
	// timers?
	// miscellaneous
}

func (c *Core) fetch() (bool, uint32) {
	return c.loadInstruction(c.pc)
}

func (c *Core) UnsafeBoot() {
	c.system.HandleBoot(c)
}

func (c *Core) run(wg *sync.WaitGroup) {
	c.bcm.Lock()

	// Start running core in loop
	if !c.state.CompareAndSwap(coreStateHalted, coreStateRunning) {
		panic("Attempted to call `run(...)` on a core that was not in the HALTED state")
	}

	c.system.HandleBoot(c)

	wg.Done() // core has switched to running state

	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == coreStateHalting {
			c.state.Store(coreStateHalted)
			break
		}

		c.UnsafeStep()
	}

	c.bcm.Unlock()
	c.Done() // core is done
}

// makes sure the core has at least entered the running state before returning
// It is an error to call Start on a core that is already started
func (c *Core) StartAndWait() {
	var wg sync.WaitGroup
	wg.Add(1)
	c.Add(1) // caller can't accidentally wait on core that hasn't entered loop
	go c.run(&wg)
	wg.Wait()
}

// Leaves it up to the caller when to sync
// It is an error to call Start on a core that is already started
func (c *Core) StartAndSync(wg *sync.WaitGroup) {
	wg.Add(1)
	c.Add(1) // caller can't accidentally wait on core
	go c.run(wg)
}

func (c *Core) HaltIfRunning() bool {
	return c.state.CompareAndSwap(coreStateRunning, coreStateHalting)
}

func (c *Core) UnsafeStep() {
	if c.counter.enable {
		if c.counter.value == 0 {
			c.counter.enable = false
			c.trap(TrapMachineTimerInterrupt)
			return
		}
		c.counter.value -= 1
	}
	c.jumped = false

	if c.vmaUpd.Load() == true {
		c.vmaUpd.Store(false)
		c.TranslationCacheInvalidate()
		c.TLBInvalidate()
	}

	success, inst := c.fetch()
	if success {
		c.execute(inst)
	} else {
		return // retry the fetch in the next cycle/step
	}

	if !c.jumped {
		c.pc += 4
	}
}

func NewCore(id uint32) (c Core) {
	c = Core{
		mc: newMemoryController(),
	}

	c.csr[Csr_MHARTID] = id

	c.state.Store(coreStateHalted)

	return
}

func (c *Core) InstructionsRetired() uint64 {
	retired := uint64(c.csr[Csr_INSTRET]) | (uint64(c.csr[Csr_INSTRETH]) << 32)
	return retired
}

func (c *Core) State() interface{} {
	return c.state.Load()
}

func (c *Core) DumpRegisters() {
	fmt.Printf("\n=== Register dump for core %d ===\n", c.csr[Csr_MHARTID])
	fmt.Printf("pc: %X\n", c.pc)

	fmt.Println("Integer registers")
	// prints all registers that have non-zero values.
	// this makes output cleaner
	for i, r := range c.reg {
		if r == 0 {
			// fmt.Printf("[%02d]: \n", i)
		} else {
			fmt.Printf("[%02d]: %08X\n", i, r)
		}
	}
}

// --- Getters and setters ---

func (c *Core) GetCSR(number int) uint32 {
	return c.csr[number]
}

func (c *Core) SetCSR(number int, val uint32) {
	c.csr[number] = val
}

func (c *Core) SetPC(pc uint32) {
	c.pc = pc
}

func (c *Core) GetIRegister(number int) uint32 {
	return c.reg[number]
}

func (c *Core) SetIRegister(number int, value uint32) {
	c.reg[number] = value
}

func (c *Core) GetFRegister(number int) uint64 {
	return c.freg[number]
}

func (c *Core) SetFRegister(number int, value uint64) {
	c.freg[number] = value
}

func (c *Core) GetIRegisters() [32]uint32 {
	var a [32]uint32
	copy(a[:], c.reg[:])
	return a
}

func (c *Core) SetIRegisters(a [32]uint32) {
	copy(c.reg[:], a[:])
}

func (c *Core) GetFRegisters() [32]uint64 {
	var a [32]uint64
	copy(a[:], c.freg[:])
	return a
}

func (c *Core) SetFRegisters(a [32]uint64) {
	copy(c.freg[:], a[:])
}

func (c *Core) SetSystem(system System) {
	c.system = system
}
