package system

import (
	"fmt"
	"gotos/cpu"
)

// These traps should be handled
//
// TrapInstructionAddressMisaligned
// TrapInstructionAccessFault
// TrapIllegalInstruction
// TrapBreakpoint
// TrapLoadAddressMisaligned
// TrapLoadAccessFault
// TrapStoreAddressMisaligned
// TrapStoreAccessFault
// TrapEcallUMode
// TrapInstructionPageFault
// TrapLoadPageFault
// TrapStorePageFault
//
// TrapMachineTimerInterrupt
// TrapMachineExternalInterrupt

func TrapHandler(c *cpu.Core) {
	// get trap reason
	reason := c.GetCSR(cpu.Csr_MCAUSE)
	switch reason {
	case cpu.TrapEcallUMode:
		handleEcallUMode(c)
	case cpu.TrapInstructionAddressMisaligned:
		handleInstructionAddressMisaligned(c)
	case cpu.TrapBreakpoint:
		handleBreakpoint(c)
	}
}

//
// Trap handlers
//

func handleInstructionAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleInstructionAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleIllegalInstruction(c *cpu.Core) {
	fmt.Printf("[core %d]: Illegal Instruction. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleLoadAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Load Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleLoadAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Load Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleStoreAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Store Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleStoreAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Store Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.HaltIfRunning()
}

func handleBreakpoint(c *cpu.Core) {
	fmt.Printf("[core %d]: Breakpoint.\n", c.GetCSR(cpu.Csr_MHARTID))
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4) // return to the instruction after the breakpoint
}

func handleEcallUMode(c *cpu.Core) {
	// syscall number is placed in register a7
	number := c.GetIRegister(cpu.Reg_A7)
	syscall(c, number)
}

func handleInstructionPageFault(c *cpu.Core) {
	// TODO
}

func handleLoadPageFault(c *cpu.Core) {
	// TODO
}

func handleStorePageFault(c *cpu.Core) {
	// TODO
}
