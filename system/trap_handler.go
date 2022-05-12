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

// HandleTrap is required to satisfy the `cpu.System` interface.
//   It takes a pointer to a `cpu.Core` as an argument and should use
//   internal information from the core to appropriately handle the trap.
func (s *System) HandleTrap(c *cpu.Core) {
	// get trap reason
	reason := c.GetCSR(cpu.Csr_MCAUSE)
	switch reason {
	case cpu.TrapInstructionAddressMisaligned:
		s.handleInstructionAddressMisaligned(c)
	case cpu.TrapInstructionAccessFault:
		s.handleInstructionAccessFault(c)
	case cpu.TrapIllegalInstruction:
		s.handleIllegalInstruction(c)
	case cpu.TrapBreakpoint:
		s.handleBreakpoint(c)
	case cpu.TrapLoadAddressMisaligned:
		s.handleLoadAddressMisaligned(c)
	case cpu.TrapLoadAccessFault:
		s.handleLoadAccessFault(c)
	case cpu.TrapStoreAddressMisaligned:
		s.handleStoreAddressMisaligned(c)
	case cpu.TrapStoreAccessFault:
		s.handleStoreAccessFault(c)
	case cpu.TrapEcallUMode:
		s.handleEcallUMode(c)
	case cpu.TrapInstructionPageFault:
		s.handleInstructionPageFault(c)
	case cpu.TrapLoadPageFault:
		s.handleLoadPageFault(c)
	case cpu.TrapStorePageFault:
		s.handleStorePageFault(c)
	case cpu.TrapMachineTimerInterrupt:
		s.handleMachineTimerInterrupt(c)
	case cpu.TrapMachineExternalInterrupt:
		s.handleMachineExternalInterrupt(c)
	}
}

//
// Trap handlers
//

func (s *System) handleInstructionAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleInstructionAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleIllegalInstruction(c *cpu.Core) {
	fmt.Printf("[core %d]: Illegal Instruction. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleLoadAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Load Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleLoadAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Load Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleStoreAddressMisaligned(c *cpu.Core) {
	fmt.Printf("[core %d]: Store Address Misaligned. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleStoreAccessFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Store Access Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleBreakpoint(c *cpu.Core) {
	fmt.Printf("[core %d]: Breakpoint.\n", c.GetCSR(cpu.Csr_MHARTID))
	c.SetCSR(cpu.Csr_MEPC, c.GetCSR(cpu.Csr_MEPC)+4) // return to the instruction after the breakpoint
}

func (s *System) handleMachineTimerInterrupt(c *cpu.Core) {
	// switch to the next process if available
	old := &PCB{}
	next := s.Scheduler.Pop()

	if next != nil {
		s.swtch(c, old, next)
		s.Scheduler.Push(old)
	}

	c.SetCounter(timeSlice)
}

func (s *System) handleMachineExternalInterrupt(c *cpu.Core) {
	by, code := c.InterruptInfo()
	fmt.Printf("[core %d]: Machine External Interrupt - code %d, by %d\n", c.GetCSR(cpu.Csr_MHARTID), code, by)
	if code == 1 {
		c.Stop()
		return
	}
	c.Halt()
}

func (s *System) handleEcallUMode(c *cpu.Core) {
	// syscall number is placed in register a0
	number := c.GetIRegister(cpu.Reg_A0)
	s.syscall(c, number)
}

func (s *System) handleInstructionPageFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Instruction Page Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleLoadPageFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Load Page Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}

func (s *System) handleStorePageFault(c *cpu.Core) {
	fmt.Printf("[core %d]: Store Page Fault. **mtval** : %08X\n", c.GetCSR(cpu.Csr_MHARTID), c.GetCSR(cpu.Csr_MTVAL))
	c.Halt()
}
