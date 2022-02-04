package system

import (
	"gotos/cpu"
	"sync"
)

type System struct {
	cores     []cpu.Core
	running   []uint32 // keep track of which PID is running on which core
	Scheduler Scheduler
	memory    cpu.Memory
	rsets     cpu.ReservationSets
}

func (s *System) Memory() *cpu.Memory {
	return &s.memory
}

func (s *System) ReservationSets() *cpu.ReservationSets {
	return &s.rsets
}

// creates a new system with `n` cores
func NewSystem(n int) *System {
	sys := &System{
		cores:   make([]cpu.Core, n),
		running: make([]uint32, n),
		memory:  cpu.NewMemory(),
		rsets:   cpu.NewReservationSets(n),
	}

	for i := range sys.cores {
		sys.cores[i] = cpu.NewCore(uint32(i))
		sys.cores[i].SetSystem(sys)
	}

	return sys
}

// runs until all cores halt by themselves
func (s *System) Run() {
	s.Start()
	s.Wait()
}

func (s *System) Boot() {
	s.cores[0].UnsafeBoot()
}

func (s *System) StepAndDump() {
	s.cores[0].UnsafeStep()
	s.cores[0].DumpRegisters()
}

func (s *System) Start() {
	var wg sync.WaitGroup
	for i := range s.cores {
		s.cores[i].StartAndSync(&wg)
	}
	wg.Wait()
}

func (s *System) Stop() {
	for i := range s.cores {
		s.cores[i].HaltIfRunning()
	}

	s.Wait()
}

func (s *System) Wait() {
	for i := range s.cores {
		s.cores[i].Wait()
	}
}
