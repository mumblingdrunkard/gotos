package system

import (
	"gotos/cpu"
	"sync"
	"sync/atomic"
)

// System implements the `System` interface from the `cpu` package.
type System struct {
	// --- system cores ---
	cores []cpu.Core

	// --- System interface ---
	memory     cpu.Memory
	rsets      cpu.ReservationSets
	interrupts cpu.InterruptMatrix
	wgAwake    sync.WaitGroup
	wgRunning  sync.WaitGroup

	// --- other fields ---
	running   []uint32  // keeps track of which PID is running on which core
	Scheduler Scheduler // acts as the system scheduler
}

// Memory returns a pointer to the system Memory and is part of the cpu.System
// interface
func (s *System) Memory() *cpu.Memory {
	return &s.memory
}

// ReservationSets returns a pointer to the system ReservationSets and is part
// of the cpu.System interface
func (s *System) ReservationSets() *cpu.ReservationSets {
	return &s.rsets
}

// InterruptMatrix returns a pointer to the system InterruptMatrix and is part
// of the cpu.System interface
func (s *System) InterruptMatrix() *cpu.InterruptMatrix {
	return &s.interrupts
}

func (s *System) WgAwake() *sync.WaitGroup {
	return &s.wgAwake
}

func (s *System) WgRunning() *sync.WaitGroup {
	return &s.wgRunning
}

func (s *System) RaiseInterrupt(coreID, code uint32) {
	for !atomic.CompareAndSwapUint32(&s.interrupts[coreID][cpu.CoresMax], 0, code) {
	}
}

// creates a new system with `n` cores
func NewSystem(n int) *System {
	sys := &System{
		// necessary
		cores:  make([]cpu.Core, n),
		memory: cpu.NewMemory(),
		rsets:  cpu.NewReservationSets(),

		// other
		running: make([]uint32, n),
	}

	for i := range sys.cores {
		sys.cores[i] = cpu.NewCore(uint32(i), sys)
	}

	return sys
}

// Run will start all cores and run them until they halt, then send a signal to
// all cores that they should stop, then wait for all cores to stop before
// finally returning with no value.
func (s *System) Run() {
	s.Start()
	s.WaitHalt()
	s.Stop()
}

// Boot will cause all cores on the system to run the boot routine
func (s *System) Boot() {
	for i := range s.cores {
		s.cores[i].Boot()
	}
}

// Dump will print the registers of all cores in the system
func (s *System) Dump() {
	for i := range s.cores {
		s.cores[i].DumpRegisters()
	}
}

// StepAndDump will call the `UnsafeStep` function on all cores in the system
// before dumping their registers.
//   This function is likely best used with a single core.
func (s *System) StepAndDump() {
	for i := range s.cores {
		s.cores[i].Step()
	}
	s.Dump()
}

// Start will start all cores in the system
func (s *System) Start() {
	for i := range s.cores {
		s.cores[i].Start()
	}
}

// Stop will raise an interrupt on each core with code 1 which should cause the
// core to eventually stop.
//   Stop then waits for all cores to finish stopping before finally returning.
func (s *System) Stop() {
	for i := range s.cores {
		s.RaiseInterrupt(uint32(i), 1)
	}

	s.WaitStop()
}

// WaitHalt will wait for all cores on the system to enter the halted state.
func (s *System) WaitHalt() {
	s.wgRunning.Wait()
}

// WaitStop will wait for all cores on the system to enter the stopped state.
func (s *System) WaitStop() {
	s.wgAwake.Wait()
}
