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
	// get the call-number

}
