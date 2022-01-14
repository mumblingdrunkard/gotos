package system

import "gotos/cpu"

func TrapHandler(c *cpu.Core) {
	// get trap reason
	reason := c.MCAUSE()
	switch reason {
	case cpu.TrapEcallUMode:
		handleUModeEcall(c)
	case cpu.TrapBreakpoint:
		handleBreakpoint(c)
	}
}

// System calls
func handleUModeEcall(c *cpu.Core) {
	const (
		sys_exit = 1
		sys_id   = 2
	)

	// get the call-number
	callNumber := c.GetIRegister(cpu.RegA7)
	switch callNumber {
	case sys_exit:
		sysExit(c)
	case sys_id:
		sysId(c)
	}
}

func handleBreakpoint(c *cpu.Core) {

}

func sysExit(c *cpu.Core) {
	// Check job queue

	// If job queue is empty, halt the core
	c.HaltIfRunning()
}

func sysId(c *cpu.Core) {
	c.SetIRegister(cpu.RegA0, c.MHARTID())
	// TODO increment pc
	// This is done by setting the pc to the address *after* the one that caused the trap
	// c.MEPC() will yield the address of the instruction that caused the trap
	c.SetPC(c.MEPC() + 4)
}
