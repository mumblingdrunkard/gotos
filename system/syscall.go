package system

import (
	"fmt"
	"gotos/cpu"
)

func (s *System) syscall(c *cpu.Core) {
	number := c.GetIRegister(cpu.Reg_A0)
	switch number {
	case 1: // `exit` system call
		fmt.Printf("[core %d] Process exited with value %08x = %d\n",
			c.GetCSR(cpu.Csr_MHARTID),
			c.GetIRegister(cpu.Reg_A1),
			c.GetIRegister(cpu.Reg_A1),
		)
		c.HaltIfRunning() // stop the processor
	}
}
