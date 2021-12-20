package cpu

import (
	"fmt"
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
	state atomic.Value
	reg   [32]uint32
	pc    uint32
}

type State struct {
	reg [32]uint32
	pc  uint32
}

func (state *State) Reg() [32]uint32 {
	// TODO: stuff
	return state.reg
}

func (state *State) Pc() uint32 {
	return state.pc
}

func (c *Core) fetch() uint32 {
	c.pc += 4
	return 0
}

func (c *Core) execute(inst uint32) {
	// Register 0 is hardwired with all 0s
	c.reg[0] = 0
	c.reg[3]++ // counts instructions executed
}

func (c *Core) run() {
	// Start running core in loop
	if !c.state.CompareAndSwap(HALTED, RUNNING) {
		panic("Attempted to call run() on a core that was not in the HALTED state")
	}

	for {
		// Test that state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == HALTING {
			c.state.Store(HALTED)
			break
		}

		inst := c.fetch()
		c.execute(inst)
	}
}

func (c *Core) StartAndWait() {
	go c.run()

	// Spin wait for core to go out of halted state
	// This ensures the core has actually started when the function returns
	for {
		if c.state.Load() == RUNNING {
			break
		}
	}
}

func (c *Core) StartAndSync(wg *sync.WaitGroup) {
	go c.run()

	// Spin wait for core to go out of halted state
	// This ensures the core has actually started when the function returns
	wg.Add(1)
	go func() {
		for {
			if c.state.Load() == RUNNING {
				break
			}
		}
		wg.Done()
	}()
}

// Gets state of processor
func (c *Core) UnsafeGetState() State {
	return State{
		reg: c.reg,
		pc:  c.pc,
	}
}

func (c *Core) UnsafeReset() {
	// Reset registers
	for i := range c.reg {
		c.reg[i] = 0
	}

	// TODO: Initialize reg[2] with memory size

	c.pc = 0
	c.state.Store(HALTED)
}

func (c *Core) HaltAndWait() {
	if !c.state.CompareAndSwap(RUNNING, HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	// Spin wait for core to go into halted state
	for {
		if c.state.Load() == HALTED {
			break
		}
	}

	fmt.Println("Successfully halted core!")
}

func (c *Core) HaltAndSync(wg *sync.WaitGroup) {
	if c.state.CompareAndSwap(RUNNING, HALTING) {
	}

	// Spin wait for core to go into halted state
	wg.Add(1)
	go func() {
		for {
			if c.state.Load() == HALTED {
				break
			}
		}
		wg.Done()
	}()
}

func NewCore() (c Core) {
	c = Core{}

	c.UnsafeReset()

	return
}
