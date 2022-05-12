package cpu

import (
	"sync"
	"time"
)

type coreState int

const (
	coreStateNopLoop  coreState = 3
	coreStateRunning            = 2 // Core is running and executing instructions
	coreStateStopping           = 1 // Core is running and executing instructions, but will turn off in the next cycle
	coreStateStopped            = 0 // Core is halted and will stay halted until Start() is called
)

// Core contains internal resources
type Core struct {
	// The big core mutex (bcm) ensures that only one goroutine is inside the
	// fetch-decode-execute loop at any one time
	bcm    sync.Mutex
	system System

	// state tracks the state of the core at any given time
	state coreState

	// jumped is a flag to tell if the most recently executed instruction
	// caused a jump
	jumped bool

	// normal registers (save on context switch)
	reg  [32]uint32 // normal registers
	freg [32]uint64 // fp registers
	pc   uint32     // program counter

	// To prevent atomic loads every cycle, a local counter is kept so
	// interrupts are only checked every N cycles.
	// Atomic operations in the main loop have a big performance penalty.
	interruptCounter uint32

	// interrupt info
	interruptedBy uint32
	interruptCode uint32

	// counter can be set to interrupt the core in N cycles
	counter counter

	// mc controls access to memory and manages caches
	mc memoryController

	// Control Status Registers, see chapter 2 in the RISC-V privileged specification.
	csr [4096]uint32
}

// Boot will run the boot method from the System interface
func (c *Core) Boot() {
	c.system.HandleBoot(c)
	c.pc = c.csr[Csr_MEPC]
}

// Step will check counters and interrupts and, if it does not trap,
// it will attempt to fetch the next instruction and execute it.
func (c *Core) Step() {
	// When a core halts, we can't just stop it or it might miss
	// interrupts.
	//   Instead, we enter a no-op loop where interrupts are checked and
	// the core sleeps for a bit.
	//   This keeps happening until the core is interrupted and a
	// handler is run to take it out of the halt loop.
	if c.state == coreStateNopLoop {
		c.checkInterrupts()
		time.Sleep(time.Millisecond) // reduce system load by sleeping
		return
	}

	// check timer
	if c.counter.enable {
		if c.counter.value == 0 {
			c.counter.enable = false
			c.trap(TrapMachineTimerInterrupt)
			return
		}
		c.counter.value -= 1
	}

	c.interruptCounter++
	// check IPIs every 100 cycles
	if c.interruptCounter >= 100 {
		c.interruptCounter = 0 // reset the counter
		if c.checkInterrupts() {
			return
		}
	}

	// --- normal instruction flow ---

	c.jumped = false

	// load and execute instruction
	success, inst := c.loadInstruction(c.pc)
	if !success {
		return
	}

	c.execute(inst)

	// increment program counter if previous instruction didn't jump
	if !c.jumped {
		c.pc += 4
	}
}

// run will place the core into a run loop.
//   When it exits, it will signal `c.system.WgAwake().Done()`.
//   In this way, it is possible to wait on the WaitGroup returned
// by WgAwake() to wait for the processor to stop.
func (c *Core) run() {
	defer c.system.WgAwake().Done()
	c.bcm.Lock()

	c.state = coreStateRunning

	c.Boot()

	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state == coreStateStopping {
			c.state = coreStateStopped
			break
		}

		c.Step()
	}

	c.bcm.Unlock()
}

// Start will start a core running in its own goroutine.
//   This is the main method to start a core
func (c *Core) Start() {
	c.system.WgAwake().Add(1)
	c.system.WgRunning().Add(1)
	go c.run()
}

// Stop will eventually transition the core into the stopped state.
//   This causes the goroutine to exit and the core will no longer check
// interrupts.
//   Before the core is stopped, make sure it releases all held
// resources and will not be interrupted by other cores, or they might
// deadlock.
func (c *Core) Stop() {
	if c.state == coreStateRunning {
		c.system.WgRunning().Done()
	}
	c.state = coreStateStopping
}

// Halt will put the core into the halted state.
//   In this state, the core is still active and checking interrupts,
// but is not actively executing instructions.
//   Before the core is transitioned into this state, make sure it
// releases all held resources.
func (c *Core) Halt() bool {
	if c.state == coreStateRunning {
		c.system.WgRunning().Done() // core is done
		c.state = coreStateNopLoop
		return true
	}
	return false
}

// NewCore creates a new core with a given id and system.
//   `sys` must be a System with at least `id + 1` cores and `id` must
// be unique among all cores that reference `sys`.
//   `id` must be less than `CoresMax`.
func NewCore(id uint32, sys System) (c Core) {
	c = Core{
		mc:     newMemoryController(),
		system: sys,
	}

	c.csr[Csr_MHARTID] = id
	c.state = coreStateStopped

	return
}

// --- Getters and setters ---

// GetCSR will give the value of a named CSR.
//   `name` must be one of the constants defined in `zicsr.go`.
func (c *Core) GetCSR(name Csr) uint32 {
	return c.csr[name]
}

// SetCSR will set the value of a named CSR.
//   `name` must be one of the constants defined in `zicsr.go`.
//   `val` is a 32-bit unsigned integer
func (c *Core) SetCSR(name Csr, val uint32) {
	c.csr[name] = val
}

// GetIRegister gets the value of a named integer register.
//   `name` must be one of the constants defined in `register.go`.
func (c *Core) GetIRegister(name Reg) uint32 {
	return c.reg[name]
}

// SetIRegister sets the value of a named integer register.
//   `name` must be one of the constants defined in `register.go`.
func (c *Core) SetIRegister(name Reg, value uint32) {
	c.reg[name] = value
}

// GetFRegister gets the value of a named floating-point register.
//   `name` must be one of the constants defined in `rv32f.go`.
func (c *Core) GetFRegister(name FReg) uint64 {
	return c.freg[name]
}

// SetFRegister sets the value of a named floating-point register.
//   `name` must be one of the constants defined in `rv32f.go`.
func (c *Core) SetFRegister(name FReg, value uint64) {
	c.freg[name] = value
}

// GetIRegisters gets all integer registers as an array of unsigned
// 32-bit integers.
func (c *Core) GetIRegisters() [32]uint32 {
	var a [32]uint32
	copy(a[:], c.reg[:])
	return a
}

// SetIRegisters sets all integer registers given an array of unsigned
// 32-bit integers.
func (c *Core) SetIRegisters(a [32]uint32) {
	copy(c.reg[:], a[:])
}

// GetFRegisters gets all floating-point registers as an array of
// unsigned 64-bit integers.
func (c *Core) GetFRegisters() [32]uint64 {
	var a [32]uint64
	copy(a[:], c.freg[:])
	return a
}

// SetFRegisters sets all floating-point registers given an array of
// unsigned 64-bit integers.
func (c *Core) SetFRegisters(a [32]uint64) {
	copy(c.freg[:], a[:])
}
