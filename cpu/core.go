package cpu

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	XLEN  = 32
	MXLEN = 32
)

const (
	CoreStateRunning = 2 // Core is running and executing instructions
	CoreStateHalting = 1 // Core is running and executing instructions, but will turn off in the next cycle
	CoreStateHalted  = 0 // Core is halted and will stay halted until Start() is called
)

// Register mnemonics
const (
	RegZero = 0  // Hard-wired zero
	RegRA   = 1  // Return address
	RegSP   = 2  // Stack pointer
	RegGP   = 3  // Global pointer
	RegTP   = 4  // Thread pointer
	RegT0   = 5  // Temporary/alternate link register
	RegT1   = 6  // Temporaries
	RegT2   = 7  //
	RegS0   = 8  // Saved register/frame pointer
	RegFP   = 8  //
	RegS1   = 9  // Saved register
	RegA0   = 10 // Function arguments/return values
	RegA1   = 11 //
	RegA2   = 12 //
	RegA3   = 13 //
	RegA4   = 14 //
	RegA5   = 15 //
	RegA6   = 16 //
	RegA7   = 17 //
	RegS2   = 18 // Saved registers
	RegS3   = 19 //
	RegS4   = 20 //
	RegS5   = 21 //
	RegS6   = 22 //
	RegS7   = 23 //
	RegS8   = 24 //
	RegS9   = 25 //
	RegS10  = 26 //
	RegS11  = 27 //
	RegT3   = 28 // Temporaries
	RegT4   = 29 //
	RegT5   = 30 //
	RegT6   = 31 //
)

const (
	EndianLittle Endian = 0
	EndianBig           = 1
)

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	// The big core mutex (bcm) ensures that only one goroutine is inside the fetch-decode-execute loop at any one time
	bcm     sync.Mutex
	state   atomic.Value // can be HALTED, HALTING, or RUNNING
	jumped  bool
	retired uint64       // number of instructions executed/retired
	reg     [32]uint32   // normal registers
	csr     [4096]uint32 // control and status registers
	freg    [32]uint64   // fp registers
	pc      uint32       // program counter
	mc      memoryController
	trapFn  func(*Core)
	mhartid uint32
	mtval   uint32 // should be set before a trap when certain faults occur
	mcause  uint32
	mepc    uint32
}

func (c *Core) fetch() (bool, uint32) {
	success, inst := c.loadInstruction(c.pc)

	if !success {
		return false, 0
	}

	return true, inst
}

func (c *Core) UnsafeSetMemBase(base uint32) {
	c.mc.mmu.base = base
}

func (c *Core) UnsafeSetMemSize(size uint32) {
	c.mc.mmu.size = size
	c.reg[2] = c.mc.mmu.size
}

func (c *Core) run(wg *sync.WaitGroup) {
	c.bcm.Lock()

	// Start running core in loop
	if !c.state.CompareAndSwap(CoreStateHalted, CoreStateRunning) {
		panic("Attempted to call `run(...)` on a core that was not in the HALTED state")
	}

	wg.Done() // core has switched to running state

	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == CoreStateHalting {
			c.state.Store(CoreStateHalted)
			break
		}

		// CAS alternative (-20% performance hit)
		// if c.state.CompareAndSwap(CoreStateHalting, CoreStateHalted) {
		// 	break
		// }

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
	if c.state.Load() == CoreStateHalted {
		return
	}

	if !c.state.CompareAndSwap(CoreStateRunning, CoreStateHalting) {
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
	if !c.state.CompareAndSwap(CoreStateRunning, CoreStateHalting) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

func (c *Core) HaltIfRunning() bool {
	return c.state.CompareAndSwap(CoreStateRunning, CoreStateHalting)
}

func (c *Core) UnsafeStep() {
	// TODO:  interrupts
	success, inst := c.fetch()

	if success {
		c.execute(inst)
		c.retired += 1
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
		mhartid: id,
		retired: 0,
		mc:      newMemoryController(m, rs),
	}

	c.state.Store(CoreStateHalted)

	c.UnsafeReset()

	return
}

func (c *Core) InstructionsRetired() uint64 {
	return c.retired
}

func (c *Core) Misses() uint64 {
	return c.mc.misses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
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

func (c *Core) SetBootHandler(handler func(*Core)) {

}

func (c *Core) SetTrapHandler(handler func(*Core)) {
	c.trapFn = handler
}

func (c *Core) SetPC(pc uint32) {
	c.pc = pc
}

func (c *Core) MHARTID() uint32 {
	return c.mhartid
}

func (c *Core) MCAUSE() uint32 {
	return c.mcause
}

func (c *Core) MTVAL() uint32 {
	return c.mtval
}

func (c *Core) MEPC() uint32 {
	return c.mepc
}
