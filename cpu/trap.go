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

// Exception codes when Interrupt is 0/false (see risc-v privileged spec table 8.6)
const (
	// Exception reasons
	TrapInstructionAddressMisaligned uint32 = 0x00000000
	TrapInstructionAccessFault              = 0x00000001
	TrapIllegalInstruction                  = 0x00000002
	TrapBreakpoint                          = 0x00000003
	TrapLoadAddressMisaligned               = 0x00000004
	TrapLoadAccessFault                     = 0x00000005
	TrapStoreAddressMisaligned              = 0x00000006
	TrapStoreAccessFault                    = 0x00000007
	TrapEcallUMode                          = 0x00000008
	trapEcallHSMode                         = 0x00000009 // not used
	trapEcallVSMode                         = 0x0000000A // not used
	trapEcallMMode                          = 0x0000000B // not used
	TrapInstructionPageFault                = 0x0000000C
	TrapLoadPageFault                       = 0x0000000D
	TrapStorePageFault                      = 0x0000000F

	// Interrupt reasons
	trapSupervisorSoftwareInterrupt        = 0x80000001 // not used
	trapVirtualSupervisorSoftwareInterrupt = 0x80000002 // not used
	trapMachineSoftwareInterrupt           = 0x80000003 // not used
	trapSupervisorTimerInterrupt           = 0x80000005 // not used
	trapVirtualSupervisorTimerInterrupt    = 0x80000006 // not used
	TrapMachineTimerInterrupt              = 0x80000007
	trapSupervisorExternalInterrupt        = 0x80000009 // not used
	trapVirtualSupervisorExternalInterrupt = 0x8000000A // not used
	TrapMachineExternalInterrupt           = 0x8000000B
)

func (c *Core) trap(reason uint32) {
	c.mcause = reason
	c.mepc = c.pc
	c.jumped = true
	c.trapFn(c)
}
