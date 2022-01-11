package cpu

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	STATE_RUNNING = 2 // Core is running and executing instructions
	STATE_HALTING = 1 // Core is running and executing instructions, but will turn off in the next cycle
	STATE_HALTED  = 0 // Core is halted and will stay halted until Start() is called
)

// Register mnemonics
const (
	REG_ZERO = 0  // Hard-wired zero
	REG_RA   = 1  // Return address
	REG_SP   = 2  // Stack pointer
	REG_GP   = 3  // Global pointer
	REG_TP   = 4  // Thread pointer
	REG_T0   = 5  // Temporary/alternate link register
	REG_T1   = 6  // Temporaries
	REG_T2   = 7  //
	REG_S0   = 8  // Saved register/frame pointer
	REG_FP   = 8  //
	REG_S1   = 9  // Saved register
	REG_A0   = 10 // Function arguments/return values
	REG_A1   = 11 //
	REG_A2   = 12 //
	REG_A3   = 13 //
	REG_A4   = 14 //
	REG_A5   = 15 //
	REG_A6   = 16 //
	REG_A7   = 17 //
	REG_S2   = 18 // Saved registers
	REG_S3   = 19 //
	REG_S4   = 20 //
	REG_S5   = 21 //
	REG_S6   = 22 //
	REG_S7   = 23 //
	REG_S8   = 24 //
	REG_S9   = 25 //
	REG_S10  = 26 //
	REG_S11  = 27 //
	REG_T3   = 28 // Temporaries
	REG_T4   = 29 //
	REG_T5   = 30 //
	REG_T6   = 31 //
)

const (
	ENDIAN_LITTLE Endian = 0
	ENDIAN_BIG           = 1
)

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	id      int
	state   atomic.Value // can be HALTED, HALTING, or RUNNING
	retired uint64       // number of instructions executed/retired
	reg     [32]uint32   // normal registers
	csr     [4096]uint32 // control and status registers
	fdirty  bool         // fp register file dirty bit
	freg    [32]uint64   // fp registers
	pc      uint32       // program counter
	mc      memoryController
}

func (c *Core) fetch() uint32 {
	success, inst := c.loadInstruction(c.pc)

	if !success {
		panic("Failed to load instruction")
	}

	return inst
}

func (c *Core) UnsafeSetMemBase(base uint32) {
	c.mc.mmu.base = base
}

func (c *Core) UnsafeSetMemSize(size uint32) {
	c.mc.mmu.size = size
	c.reg[2] = c.mc.mmu.size
}

func (c *Core) run(wg *sync.WaitGroup) {
	// Start running core in loop
	if !c.state.CompareAndSwap(STATE_HALTED, STATE_RUNNING) {
		panic("Attempted to call `run()` on a core that was not in the HALTED state")
	}

	wg.Done() // core has switched to running state

	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == STATE_HALTING {
			c.state.Store(STATE_HALTED)
			break
		}

		c.UnsafeStep()
	}

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

	// TODO: Initialize reg[2] with memory size
	c.reg[2] = c.mc.mmu.size

	c.pc = 0
	// c.state.Store(HALTED)
}

// Halts the core and waits for it to halt before returning
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndWait() {
	if c.state.Load() == STATE_HALTED {
		return
	}

	if !c.state.CompareAndSwap(STATE_RUNNING, STATE_HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	c.Wait() // Wait for core to halt
}

// Halts the core, but leaves it to the caller to sync
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndSync(wg *sync.WaitGroup) {
	if !c.state.CompareAndSwap(STATE_RUNNING, STATE_HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

// Not really needed anymore, use core.Wait() instead
func (c *Core) SyncOnHalt(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

func NewCore() (c Core) {
	c = Core{}

	c.UnsafeReset()

	return
}

func (c *Core) UnsafeStep() {
	// TODO: Check for interrupts

	c.retired += 1
	inst := c.fetch()

	// fmt.Printf("executing: %08X\n", inst)

	c.execute(inst)
	// if (opcode != BRANCH) && (opcode != JAL) && (opcode != JALR) {
	c.pc += 4
	// }
}

func NewCoreWithMemoryAndReservationSets(m *Memory, rs *ReservationSets, id int) (c Core) {
	c = Core{
		id:      id,
		retired: 0,
		mc:      newMemoryController(m, rs),
	}

	c.state.Store(STATE_HALTED)

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
