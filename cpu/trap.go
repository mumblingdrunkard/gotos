package cpu

type TrapReason uint32

// Exception codes when Interrupt is 0/false (see risc-v privileged spec table 8.6)
const (
	// Exception reasons
	TrapInstructionAddressMisaligned TrapReason = 0x000
	TrapInstructionAccessFault                  = 0x001
	TrapIllegalInstruction                      = 0x002
	TrapBreakpoint                              = 0x003
	TrapLoadAddressMisaligned                   = 0x004
	TrapLoadAccessFault                         = 0x005
	TrapStoreAddressMisaligned                  = 0x006
	TrapStoreAccessFault                        = 0x007
	TrapEcallUMode                              = 0x008
	trapEcallHSMode                             = 0x009 // not used
	trapEcallVSMode                             = 0x00A // not used
	trapEcallMMode                              = 0x00B // not used
	TrapInstructionPageFault                    = 0x00C
	TrapLoadPageFault                           = 0x00D
	TrapStorePageFault                          = 0x00F

	// Interrupt reasons
	trapSupervisorSoftwareInterrupt        = 0x801 // not used
	trapVirtualSupervisorSoftwareInterrupt = 0x802 // not used
	trapMachineSoftwareInterrupt           = 0x803 // not used
	trapSupervisorTimerInterrupt           = 0x805 // not used
	trapVirtualSupervisorTimerInterrupt    = 0x806 // not used
	TrapMachineTimerInterrupt              = 0x807
	trapSupervisorExternalInterrupt        = 0x809 // not used
	trapVirtualSupervisorExternalInterrupt = 0x80A // not used
	TrapMachineExternalInterrupt           = 0x80B
)

func (c *Core) trap(reason TrapReason) {
	// TODO
	c.trapFn(c, reason)
}
