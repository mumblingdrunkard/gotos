package system

import (
	"fmt"
	"gotos/cpu"
)

func (s *System) syscall(c *cpu.Core, number uint32) {
	const (
		sys_exit = 1
	)

	switch number {
	case sys_exit:
		s.syscall_exit(c)
	}
}

//
// syscalls
//

func (s *System) syscall_exit(c *cpu.Core) {
	exitValue := c.GetIRegister(cpu.Reg_A0)
	coreId := c.GetCSR(cpu.Csr_MHARTID)
	fmt.Printf("[core %d] : Process %d exited with value 0x%08X = %d\n", coreId, s.running[coreId], exitValue, exitValue)
	c.HaltIfRunning()
}
