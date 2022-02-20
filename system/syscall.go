package system

import (
	"fmt"
	"gotos/cpu"
)

func (s *System) syscall(c *cpu.Core) {
	number := c.GetIRegister(cpu.Reg_A7)
	switch number {
	case 1:
		s.syscallExit(c)
	case 513:
		s.syscallSayhello(c)
	}
}

func (s *System) syscallExit(c *cpu.Core) {
	fmt.Println("Process exited")
	c.HaltIfRunning() // stop the processor
}

func (s *System) syscallSayhello(c *cpu.Core) {
	fmt.Println("Hello!")
	cameFrom := c.GetCSR(cpu.Csr_MEPC)
	c.SetPC(cameFrom + 4)
}
