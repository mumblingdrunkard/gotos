package cpu

// Exception codes (see risc-v privileged spec table 8.6)
const (
	// --- Exception reasons ---
	TrapInstructionAddressMisaligned uint32 = 0x00000000
	TrapInstructionAccessFault              = 0x00000001
	TrapIllegalInstruction                  = 0x00000002
	TrapBreakpoint                          = 0x00000003
	TrapLoadAddressMisaligned               = 0x00000004
	TrapLoadAccessFault                     = 0x00000005
	TrapStoreAddressMisaligned              = 0x00000006
	TrapStoreAccessFault                    = 0x00000007
	TrapEcallUMode                          = 0x00000008
	TrapInstructionPageFault                = 0x0000000C
	TrapLoadPageFault                       = 0x0000000D
	TrapStorePageFault                      = 0x0000000F

	// --- Interrupt reasons ---
	TrapMachineTimerInterrupt    = 0x80000007
	TrapMachineExternalInterrupt = 0x8000000B
)

// trap() sets up the parts of the trap that are common for every trap.
//   Other setup has to be done independently depending on the trap
// reason.
func (c *Core) trap(reason uint32) {
	c.csr[Csr_MCAUSE] = reason
	c.csr[Csr_MEPC] = c.pc
	c.jumped = true
	c.system.HandleTrap(c)
	c.pc = c.csr[Csr_MEPC]
}
