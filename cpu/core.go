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

func (core *Core) fetch() uint32 {
	core.pc += 4
	return 0
}

func (core *Core) execute(inst uint32) {
	// Register 0 is hardwired with all 0s
	core.reg[0] = 0
	core.reg[3]++ // counts instructions executed
}

func (core *Core) run() {
	// Start running core in loop
	if !core.state.CompareAndSwap(HALTED, RUNNING) {
		panic("Attempted to call run() on a core that was not in the HALTED state")
	}

	for {
		// Test that state is HALTING, swap to HALTED if so, then break
        // CompareAndSwap makes this very slow so we use Load instead
		if core.state.Load() == HALTING {
            core.state.Store(HALTED)
			break
		}

		inst := core.fetch()
		core.execute(inst)
	}
}

func (core *Core) StartAndWait() {
	go core.run()

	// Spin wait for core to go out of halted state
	// This ensures the core has actually started when the function returns
	for {
		if core.state.Load() == RUNNING {
			break
		}
	}
}

func (core *Core) StartAndSync(wg *sync.WaitGroup) {
	go core.run()

	// Spin wait for core to go out of halted state
	// This ensures the core has actually started when the function returns
	wg.Add(1)
	go func() {
		for {
			if core.state.Load() == RUNNING {
				break
			}
		}
		wg.Done()
	}()
}

// UnsafeX methods are not safe to call while the core is running
func (core *Core) UnsafeGetState() State {
	return State{
		reg: core.reg,
		pc:  core.pc,
	}
}

func (core *Core) UnsafeReset() {
	// Reset registers
	for i := range core.reg {
		core.reg[i] = 0
	}

	// TODO: Initialize reg[2] with memory size

	core.pc = 0
	core.state.Store(HALTED)
}

func (core *Core) HaltAndWait() {
	if !core.state.CompareAndSwap(RUNNING, HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	// Spin wait for core to go into halted state
	for {
		if core.state.Load() == HALTED {
			break
		}
	}

	fmt.Println("Successfully halted core!")
}

func (core *Core) HaltAndSync(wg *sync.WaitGroup) {
	if core.state.CompareAndSwap(RUNNING, HALTING) {
	}

	// Spin wait for core to go into halted state
	wg.Add(1)
	go func() {
		for {
			if core.state.Load() == HALTED {
				break
			}
		}
		wg.Done()
	}()
}

func NewCore() (core Core) {
	core = Core{}

	core.UnsafeReset()

	return
}
