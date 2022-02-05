package system

import "gotos/cpu"

// This file should contain system startup function
// This function should set up all registers and whatnot to prepare the core to start running programs.

func (s *System) HandleBoot(c *cpu.Core) {
	c.SetPC(0)
	c.SetIRegister(cpu.Reg_SP, cpu.MemorySize)
}