package system

import "gotos/cpu"

type Scheduler interface {
	// should put a PCB into the scheduler queue
	Push(*PCB)

	// should remove the next PCB from the queue and return it
	Pop() *PCB
}

const (
	timeSlice uint64 = 1000000
)

func (s *System) contextSwitch(c *cpu.Core, oldPCB *PCB, newPCB *PCB) {
	// Have to invalidate cache on context switches so work can be resumed on a different core
	c.CacheWritebackAndInvalidate()
	c.InstructionCacheInvalidate()

	coreId := c.GetCSR(cpu.Csr_MHARTID)

	if oldPCB != nil {
		ireg := c.GetIRegisters()
		copy(oldPCB.IReg[:], ireg[:])

		freg := c.GetFRegisters()
		copy(oldPCB.FReg[:], freg[:])

		pc := c.GetCSR(cpu.Csr_MEPC)
		oldPCB.PC = pc

		oldPCB.PID = s.running[coreId]
	}

	c.SetIRegisters(newPCB.IReg)
	c.SetFRegisters(newPCB.FReg)
	c.SetPC(newPCB.PC)
	s.running[coreId] = newPCB.PID
}
