package system

import "gotos/cpu"

func TrapHandler(c *cpu.Core, reason cpu.TrapReason) {
	switch reason {
	case cpu.TrapEcallUMode:
		handleUModeEcall(c)
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

func sysExit(c *cpu.Core) {
	c.UnsafeSetState(cpu.CoreStateHalting)
}

func sysId(c *cpu.Core) {
	c.SetIRegister(cpu.RegA0, c.Id())
}
