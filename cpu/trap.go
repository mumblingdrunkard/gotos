package cpu

// TODO
// Before calling trap(reason), some types of traps may require additional setup.
//
// TODO Page faults for example, should give sufficient information about what action caused the fault so that it can be properly handled.
//      This could be the virtual address that was attempted to be loaded.
//
// TODO External interrupts should likely also come with additional information about what raised the interrupt.
//
// TODO Machine timer interrupts should include information about what the timer interrupts for.
//
// TODO AddressMisaligned traps should likely also include information about what the address being accessed was.
//

// Exception codes when Interrupt is 0/false (see risc-v privileged spec table 8.6)
const (
	// --- Exception reasons ---
	TrapInstructionAddressMisaligned uint32 = 0x00000000 // kill
	TrapInstructionAccessFault              = 0x00000001 // kill
	TrapIllegalInstruction                  = 0x00000002 // kill
	TrapBreakpoint                          = 0x00000003 // continue (software)
	TrapLoadAddressMisaligned               = 0x00000004 // kill
	TrapLoadAccessFault                     = 0x00000005 // kill
	TrapStoreAddressMisaligned              = 0x00000006 // kill
	TrapStoreAccessFault                    = 0x00000007 // handle (might be the stack is too small)
	TrapEcallUMode                          = 0x00000008 // handle (software)
	// TrapEcallHSMode = 0x00000009 // not used
	// TrapEcallVSMode = 0x0000000A // not used
	// TrapEcallMMode  = 0x0000000B // not used
	TrapInstructionPageFault = 0x0000000C // handle (software)
	TrapLoadPageFault        = 0x0000000D // handle (software)
	TrapStorePageFault       = 0x0000000F // handle (software)

	// --- Interrupt reasons ---
	// TrapSupervisorSoftwareInterrupt        = 0x80000001 // not used
	// TrapVirtualSupervisorSoftwareInterrupt = 0x80000002 // not used
	// TrapMachineSoftwareInterrupt           = 0x80000003 // not used
	// TrapSupervisorTimerInterrupt           = 0x80000005 // not used
	// TrapVirtualSupervisorTimerInterrupt    = 0x80000006 // not used
	TrapMachineTimerInterrupt = 0x80000007 // handle (software)
	// TrapSupervisorExternalInterrupt        = 0x80000009 // not used
	// TrapVirtualSupervisorExternalInterrupt = 0x8000000A // not used
	TrapMachineExternalInterrupt = 0x8000000B // handle (software)
)

func (c *Core) trap(reason uint32) {
	c.csr[Csr_MCAUSE] = reason
	c.csr[Csr_MEPC] = c.pc
	c.jumped = true
	c.system.HandleTrap(c)
	c.pc = c.csr[Csr_MEPC]
}
