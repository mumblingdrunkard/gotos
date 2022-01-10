package cpu

type ExceptionCode uint32
type InterruptExceptionCode uint32

// Exception codes when Interrupt is 1/true (see risc-v privileged spec table 8.6)
const (
	TRAP_SUPERVISOR_SOFTWARE_INTERRUPT         InterruptExceptionCode = 1  // not used
	TRAP_VIRTUAL_SUPERVISOR_SOFTWARE_INTERRUPT                        = 2  // not used
	TRAP_MACHINE_SOFTWARE_INTERRUPT                                   = 3  // not used
	TRAP_SUPERVISOR_TIMER_INTERRUPT                                   = 5  // not used
	TRAP_VIRTUAL_SUPERVISOR_TIMER_INTERRUPT                           = 6  // not used
	TRAP_MACHINE_TIMER_INTERRUPT                                      = 7  // not used
	TRAP_SUPERVISOR_EXTERNAL_INTERRUPT                                = 9  // not used
	TRAP_VIRTUAL_SUPERVISOR_EXTERNAL_INTERRUPT                        = 10 // not used
	TRAP_MACHINE_EXTERNAL_INTERRUPT                                   = 11
)

// Exception codes when Interrupt is 0/false (see risc-v privileged spec table 8.6)
const (
	TRAP_INSTRUCTION_ADDRESS_MISALIGNED  ExceptionCode = 0
	TRAP_INSTRUCTION_ACCESS_FAULT                      = 1
	TRAP_ILLEGAL_INSTRUCTION                           = 2
	TRAP_BREAKPOINT                                    = 3
	TRAP_LOAD_ADDRESS_MISALIGNED                       = 4
	TRAP_LOAD_ACCESS_FAULT                             = 5
	TRAP_STORE_OR_AMO_ADDRESS_MISALIGNED               = 6
	TRAP_STORE_OR_AMO_ACCESS_FAULT                     = 7
	TRAP_ECALL_UMODE                                   = 8
	TRAP_ECALL_HSMODE                                  = 9  // not used
	TRAP_ECALL_VSMODE                                  = 10 // not used
	TRAP_ECALL_MMODE                                   = 11 // not used
	TRAP_INSTRUCTION_PAGE_FAULT                        = 12
	TRAP_LOAD_PAGE_FAULT                               = 13
	TRAP_STORE_OR_AMO_PAGE_FAULT                       = 15
)

func (c *Core) trap(interrupt bool, icode InterruptExceptionCode, ecode ExceptionCode) {
	// TODO
}
