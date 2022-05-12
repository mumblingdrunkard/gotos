// This file contains the system boot handler that is required by the
// cpu.System interface.

package system

import "gotos/cpu"

// HandleBoot handles the boot-up process of a core.
func (s *System) HandleBoot(c *cpu.Core) {
	next := s.Scheduler.Pop()
	if next != nil {
		s.swtch(c, nil, next)
		c.SetCounter(timeSlice)
	} else {
		c.Halt()
	}
}
