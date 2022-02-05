package system

import (
	"gotos/cpu"
	"sync"
)

type System struct {
	Cores  []cpu.Core
	memory cpu.Memory
	rsets  cpu.ReservationSets
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
		Cores:  make([]cpu.Core, n),
		memory: cpu.NewMemory(),
		rsets:  cpu.NewReservationSets(n),
	}

	for i := range sys.Cores {
		sys.Cores[i] = cpu.NewCore(uint32(i))
		sys.Cores[i].SetSystem(sys)
	}

	return sys
}

// runs until all cores halt by themselves
func (s *System) Run() {
	s.Start()
	s.Wait()
}

func (s *System) Boot() {
	s.Cores[0].UnsafeBoot()
}

func (s *System) StepAndDump() {
	s.Cores[0].UnsafeStep()
	s.Cores[0].DumpRegisters()
}

func (s *System) Start() {
	var wg sync.WaitGroup
	for i := range s.Cores {
		s.Cores[i].StartAndSync(&wg)
	}
	wg.Wait()
}

func (s *System) Stop() {
	for i := range s.Cores {
		s.Cores[i].HaltIfRunning()
	}

	s.Wait()
}

func (s *System) Wait() {
	for i := range s.Cores {
		s.Cores[i].Wait()
	}
}
