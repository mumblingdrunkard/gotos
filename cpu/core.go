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

const (
	EndianLittle Endian = 0
	EndianBig           = 1
)

type counter struct {
	enable bool
	value  uint64
}

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	// The big core mutex (bcm) ensures that only one goroutine is inside the fetch-decode-execute loop at any one time
	bcm    sync.Mutex
	state  atomic.Value // can be HALTED, HALTING, or RUNNING
	jumped bool
	mc     memoryController
	// normal registers (save on context switch)
	reg  [32]uint32 // normal registers
	freg [32]uint64 // fp registers
	pc   uint32     // program counter
	// CSRs
	csr [4096]uint32 // control and status registers
	// function pointers
	trapFn func(*Core)
	bootFn func(*Core)
	// counters, used for predictable scheduling
	counter counter
	// timers?
	// miscellaneous
	retired uint64 // number of instructions executed/retired
}

func (c *Core) fetch() (bool, uint32) {
	return c.loadInstruction(c.pc)
}

func (c *Core) run(wg *sync.WaitGroup) {
	c.bcm.Lock()

	// Start running core in loop
	if !c.state.CompareAndSwap(coreStateHalted, coreStateRunning) {
		panic("Attempted to call `run(...)` on a core that was not in the HALTED state")
	}

	c.bootFn(c)

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

func (c *Core) UnsafeReset() {
	// Reset registers
	for i := range c.reg {
		c.reg[i] = 0
	}

	// Initialize reg[2] with memory size
	c.reg[2] = c.mc.mmu.size

	c.pc = 0
	// c.state.Store(HALTED)
}

// Halts the core and waits for it to halt before returning
// WARNING: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndWait() {
	if c.state.Load() == coreStateHalted {
		return
	}

	if !c.state.CompareAndSwap(coreStateRunning, coreStateHalting) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	c.Wait() // Wait for core to halt
}

// Halts the core, but leaves it to the caller to sync
// WARNING: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndSync(wg *sync.WaitGroup) {
	if !c.state.CompareAndSwap(coreStateRunning, coreStateHalting) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

func (c *Core) HaltIfRunning() bool {
	return c.state.CompareAndSwap(coreStateRunning, coreStateHalting)
}

func (c *Core) UnsafeStep() {
	// Interrupts
	if c.counter.enable {
		if c.counter.value == 0 {
			c.counter.enable = false
			c.trap(TrapMachineTimerInterrupt)
			return
		}
		c.counter.value -= 1
	}

	success, inst := c.fetch()
	if success {
		c.execute(inst)
		c.retired += 1
	} else {
		return // retry the fetch next time
	}

	// Only increment pc if the processor did not trap or perform a jump
	// This means that branches and jumps don't need to jump to the address *before* the intended target.
	// This also means that for most traps/exceptions, the instruction will be retried.
	// This is helpful for stuff like page-faults that may occur.
	// It also means that when something like ECALL or EBREAK is performed, there may be a need to manually increment the program counter.
	if !c.jumped {
		c.pc += 4
	}
	c.jumped = false
}

func NewCoreWithMemoryAndReservationSets(m *Memory, rs *ReservationSets, id uint32) (c Core) {
	c = Core{
		retired: 0,
		mc:      newMemoryController(m, rs),
	}

	c.csr[Csr_MHARTID] = id

	c.state.Store(coreStateHalted)

	c.UnsafeReset()

	return
}

func (c *Core) InstructionsRetired() uint64 {
	return c.retired
}

func (c *Core) State() interface{} {
	return c.state.Load()
}

func (c *Core) DumpRegisters() {
	fmt.Println("Register dump")
	fmt.Printf("pc: %X\n", c.pc)
	for i, r := range c.reg {
		fmt.Printf("[%02d]: %08X\n", i, r)
	}
}

// --- Getters and setters ---

func (c *Core) GetCSR(number int) uint32 {
	return c.csr[number]
}

func (c *Core) SetCSR(number int, val uint32) {
	c.csr[number] = val
}

// TODO/WARNING: should perhaps not be provided as a "real" trap would mangle the pc.
// Use MEPC instead.
func (c *Core) GetPC() uint32 {
	return c.pc
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

func (c *Core) SetBootHandler(handler func(*Core)) {
	c.bootFn = handler
}

func (c *Core) SetTrapHandler(handler func(*Core)) {
	c.trapFn = handler
}
