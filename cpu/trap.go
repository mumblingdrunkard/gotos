package cpu

type ExceptionCode uint32
type InterruptExceptionCode uint32

// Exception codes when Interrupt is 1 (see risc-v privileged spec table 8.6)
const (
	RESERVED_0                            InterruptExceptionCode = 0
	SUPERVISOR_SOFTWARE_INTERRUPT                                = 1
	VIRTUAL_SUPERVISOR_SOFTWARE_INTERRUPT                        = 2
	MACHINE_SOFTWARE_INTERRUPT                                   = 3
	SUPERVISOR_TIMER_INTERRUPT                                   = 5
	VIRTUAL_SUPERVISOR_TIMER_INTERRUPT                           = 6
	MACHINE_TIMER_INTERRUPT                                      = 7
	SUPERVISOR_EXTERNAL_INTERRUPT                                = 9
	VIRTUAL_SUPERVISOR_EXTERNAL_INTERRUPT                        = 10
	MACHINE_EXTERNAL_INTERRUPT                                   = 11
	SUPERVISOR_GUEST_EXTERNAL_INTERRUPT                          = 12
)

// Exception codes when Interrupt is 0 (see risc-v privileged spec table 8.6)
const (
	INSTRUCTION_ADDRESS_MISALIGNED  ExceptionCode = 0
	INSTRUCTION_ACCESS_FAULT                      = 1
	ILLEGAL_INSTRUCTION                           = 2
	BREAKPOINT                                    = 3
	LOAD_ADDRESS_MISALIGNED                       = 2
	LOAD_ACCESS_FAULT                             = 5
	STORE_OR_AMO_ADDRESS_MISALIGNED               = 6
	STORE_OR_AMO_ACCESS_FAULT                     = 7
	ECALL_UMODE                                   = 8
	ECALL_HSMODE                                  = 9
	ECALL_VSMODE                                  = 10
	ECALL_MMODE                                   = 11
	INSTRUCTION_PAGE_FAULT                        = 12
	LOAD_PAGE_FAULT                               = 13
	STORE_OR_AMO_PAGE_FAULT                       = 15
	INSTRUCTION_GUEST_PAGE_FAULT                  = 20
	LOAD_GUEST_PAGE_FAULT                         = 21
	VIRTUAL_INSTRUCTION                           = 22
	STORE_AMO_GUEST_PAGE_FAULT                    = 23
)

func (c *Core) trap() {

}
