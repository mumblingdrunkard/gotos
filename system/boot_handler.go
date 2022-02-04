package system

import "gotos/cpu"

// This file should contain system startup function
// This function should set up all registers and whatnot to prepare the core to start running programs.

func (s *System) HandleBoot(c *cpu.Core) {
	next := s.Scheduler.Pop()
	if next != nil {
		s.contextSwitch(c, nil, next)
	} else {
		c.HaltIfRunning()
	}
}
