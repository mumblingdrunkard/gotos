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
		sys_exit = 1
	)

	switch number {
	case sys_exit:
		s.syscall_exit(c)
	}
}

func getArgs(c *cpu.Core) [6]uint32 {
	var args [6]uint32
	args[0] = c.GetIRegister(cpu.Reg_A1)
	args[1] = c.GetIRegister(cpu.Reg_A2)
	args[2] = c.GetIRegister(cpu.Reg_A3)
	args[3] = c.GetIRegister(cpu.Reg_A4)
	args[4] = c.GetIRegister(cpu.Reg_A5)
	args[5] = c.GetIRegister(cpu.Reg_A6)
	return args
}

func returnValue(c *cpu.Core, v uint32) {
	c.SetIRegister(cpu.Reg_A0, v)
}

func (s *System) syscall_exit(c *cpu.Core) {
	exitValue := c.GetIRegister(cpu.Reg_A0)
	coreId := c.GetCSR(cpu.Csr_MHARTID)
	fmt.Printf("[core %d] : Process %d exited with value 0x%08X = %d\n", coreId, s.running[coreId], exitValue, exitValue)
	c.HaltIfRunning()
}
