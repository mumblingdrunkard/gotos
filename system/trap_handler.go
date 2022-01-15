package system

import (
	"fmt"
	"gotos/cpu"
)

func TrapHandler(c *cpu.Core) {
	// get trap reason
	reason := c.MCAUSE()
	switch reason {
	case cpu.TrapEcallUMode:
		handleUModeEcall(c)
	case cpu.TrapInstructionAddressMisaligned:
		handleInstructionAddressMisaligned(c)
	case cpu.TrapBreakpoint:
		handleBreakpoint(c)
	}
}

//
// Trap handlers
//

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

func handleInstructionAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Address Misaligned. **mtval** : %08X\n", c.MHARTID(), c.MTVAL())
	c.HaltIfRunning()
}

//
// syscalls
//

func sysExit(c *cpu.Core) {
	// TODO sys_exit just halts the core
	// Instead, extract return value and store somewhere safe.
	// Start the next process in the queue
	// Halt if the queue is empty
	c.HaltIfRunning()
}

func sysId(c *cpu.Core) {
	c.SetIRegister(cpu.RegA0, c.MHARTID())
	// TODO increment pc
	// This is done by setting the pc to the address *after* the one that caused the trap
	// c.MEPC() will yield the address of the instruction that caused the trap
	c.SetPC(c.MEPC() + 4)
}
