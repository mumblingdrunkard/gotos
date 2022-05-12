package system

import (
	"fmt"
	"gotos/cpu"
)

//
// syscalls
//

func (s *System) syscall(c *cpu.Core, number uint32) {
	const (
		sys_exit   = 1
		sys_id     = 2
		sys_getpid = 6
		sys_putint = 8
		sys_yield  = 10
	)

	switch number {
	case sys_exit:
		s.syscall_exit(c)
	case sys_id:
		s.sysId(c)
	case sys_getpid:
		s.sysGetPID(c)
	case sys_putint:
		s.sysPutInt(c)
	case sys_yield:
		s.sysYield(c)
	}
}

func getArgs(c *cpu.Core) [7]uint32 {
	var args [7]uint32
	args[0] = c.GetIRegister(cpu.Reg_A1)
	args[1] = c.GetIRegister(cpu.Reg_A2)
	args[2] = c.GetIRegister(cpu.Reg_A3)
	args[3] = c.GetIRegister(cpu.Reg_A4)
	args[4] = c.GetIRegister(cpu.Reg_A5)
	args[5] = c.GetIRegister(cpu.Reg_A6)
	args[6] = c.GetIRegister(cpu.Reg_A7)
	return args
}

func returnValue(c *cpu.Core, v uint32) {
	c.SetIRegister(cpu.Reg_A0, v)
}

func (s *System) syscall_exit(c *cpu.Core) {
	// TODO sys_exit just halts the core
	// Instead, extract return value and store somewhere safe.
	// Start the next process in the queue
	// Halt if the queue is empty

	value := c.GetIRegister(cpu.Reg_A1)
	coreId := c.GetCSR(cpu.Csr_MHARTID)
	fmt.Printf("[core %d]: Process %d exited with value 0x%08X = %d\n", coreId, s.running[coreId], value, value)

	// run the next process if available
	next := s.Scheduler.Pop()
	if next != nil {
		s.swtch(c, nil, next)
		c.SetCounter(timeSlice)
	} else {
		fmt.Printf("[core %d]: No more pcb's in queue!\n", c.GetCSR(cpu.Csr_MHARTID))
		c.Halt()
	}
}

func (s *System) sysYield(c *cpu.Core) {
	// run the next process if available
	// Set the address that the process should return to when it starts up again
	trapAddr := c.GetCSR(cpu.Csr_MEPC)
	c.SetCSR(cpu.Csr_MEPC, trapAddr+4)

	next := s.Scheduler.Pop()
	if next != nil {
		old := &PCB{}
		s.swtch(c, old, next)
		s.Scheduler.Push(old)
		c.SetCounter(timeSlice)
	} // else do nothing
}

func (s *System) sysId(c *cpu.Core) {
	c.SetIRegister(cpu.Reg_A0, c.GetCSR(cpu.Csr_MHARTID))
	// return to program
	// This is done by setting the pc to the address *after* the one that caused the trap
	// c.MEPC() will contain the address of the instruction that caused the trap
	c.SetCSR(cpu.Csr_MEPC, c.GetCSR(cpu.Csr_MEPC)+4)
}

func (s *System) sysGetPID(c *cpu.Core) {
	c.SetIRegister(cpu.Reg_A0, c.GetCSR(cpu.Csr_MHARTID))
	c.SetCSR(cpu.Csr_MEPC, c.GetCSR(cpu.Csr_MEPC)+4)
}

func (s *System) sysPutInt(c *cpu.Core) {
	fmt.Println(c.GetIRegister(cpu.Reg_A1))
	c.SetCSR(cpu.Csr_MEPC, c.GetCSR(cpu.Csr_MEPC)+4)
}
