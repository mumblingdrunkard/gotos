package system

import "gotos/cpu"

type Scheduler interface {
	// should put a PCB into the scheduler queue
	Push(*PCB)

	// should remove the next PCB from the queue and return it
	Pop() *PCB
}

const (
	timeSlice uint64 = 100000000000
)

// TODO, don't just restore PCB from current state
// instead, hold a pointer to the PCB of the currently running thread on each core
// restore state to that PCB instead

func (s *System) swtch(c *cpu.Core, oldPCB *PCB, newPCB *PCB) {
	// Have to invalidate cache on context switches so work can be resumed on a different core
	c.FENCE()
	c.FENCE_I()

	coreId := c.GetCSR(cpu.Csr_MHARTID)

	if oldPCB != nil {
		ireg := c.GetIRegisters()
		copy(oldPCB.IReg[:], ireg[:])

		freg := c.GetFRegisters()
		copy(oldPCB.FReg[:], freg[:])

		pc := c.GetCSR(cpu.Csr_MEPC)
		oldPCB.PC = pc

		oldPCB.PID = s.running[coreId]
		oldPCB.PTable = (c.GetCSR(cpu.Csr_SATP) & 0x003FFFFF) << 12
	}

	c.SetIRegisters(newPCB.IReg)
	c.SetFRegisters(newPCB.FReg)
	c.SetCSR(cpu.Csr_SATP, 0x80000000|(newPCB.PTable>>12)|newPCB.PID<<22)
	c.SetCSR(cpu.Csr_MEPC, newPCB.PC)

	c.SFENCE_VMA(0, 0, 0)

	s.running[coreId] = newPCB.PID
}
