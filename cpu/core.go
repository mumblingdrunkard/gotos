package cpu

import (
	"fmt"
	"gotos/memory"
	"sync"
	"sync/atomic"
)

const (
	RUNNING = 2 // Core is running and executing instructions
	HALTING = 1 // Core is running and executing instructions, but will turn off in the next cycle
	HALTED  = 0 // Core is halted and will stay halted until Start() is called
)

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	state  atomic.Value
	cycles uint64
	reg    [32]uint32
	pc     uint32
	// Core doesn't have exclusive ownership of memory so we hold a pointer/reference
	// to it instead
	mem *memory.Memory
}

type State struct {
	cycles uint64
	reg    [32]uint32
	pc     uint32
}

func (state *State) Reg() [32]uint32 {
	return state.reg
}

func (state *State) Pc() uint32 {
	return state.pc
}

func (state *State) Cycles() uint64 {
	return state.cycles
}

func (c *Core) fetch() uint32 {
	err, inst := c.mem.LoadWord(int(c.pc))

	if err != nil {
		panic(err)
	}

	return inst
}

func (c *Core) run(wg *sync.WaitGroup) {
	// Start running core in loop
	if !c.state.CompareAndSwap(HALTED, RUNNING) {
		panic("Attempted to call `run()` on a core that was not in the HALTED state")
	}

	wg.Done() // core has switched to running state

	c.Add(1) // core is running
	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == HALTING {
			c.state.Store(HALTED)
			break
		}

		c.cycles += 1
		inst := c.fetch()
		c.execute(inst)
		opcode := inst & 0x7f
		if (opcode != BRANCH) && (opcode != JAL) && (opcode != JALR) && (opcode != SYSTEM) {
			c.pc += 4
		}
	}
	c.Done() // core is done
}

// makes sure the core has at least entered the running state before returning
// It is an error to call Start on a core that is already started
func (c *Core) StartAndWait() {
	var wg sync.WaitGroup
	wg.Add(1)
	go c.run(&wg)
	wg.Wait()
}

// Leaves it up to the caller when to sync
// It is an error to call Start on a core that is already started
func (c *Core) StartAndSync(wg *sync.WaitGroup) {
	wg.Add(1)
	go c.run(wg)
}

// Gets state of processor
func (c *Core) UnsafeGetState() State {
	return State{
		cycles: c.cycles,
		reg:    c.reg,
		pc:     c.pc,
	}
}

func (c *Core) UnsafeReset() {
	// Reset registers
	for i := range c.reg {
		c.reg[i] = 0
	}

	// TODO: Initialize reg[2] with memory size
	c.reg[2] = c.mem.Size()

	c.pc = 0
	// c.state.Store(HALTED)
}

// Halts the core and waits for it to halt before returning
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndWait() {
	if !c.state.CompareAndSwap(RUNNING, HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	c.Wait() // Wait for core to halt

	fmt.Println("Successfully halted core!")
}

// Halts the core, but leaves it to the caller to sync
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndSync(wg *sync.WaitGroup) {
	if !c.state.CompareAndSwap(RUNNING, HALTING) {
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

func (c *Core) Step() {
	inst := c.fetch()
	c.execute(inst)
	opcode := inst & 0x7f
	if (opcode != BRANCH) && (opcode != JAL) && (opcode != JALR) {
		c.pc += 4
	}
}

func NewCoreWithMemory(m *memory.Memory) (c Core) {
	c = Core{
		cycles: 0,
		mem:    m,
	}

	c.state.Store(HALTED)

	c.UnsafeReset()

	return
}
