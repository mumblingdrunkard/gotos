package cpu

type ExceptionCode uint32

const (
	ADDRESS_MISALIGNED ExceptionCode = 0
	ACCESS_FAULT                     = 1
)

func (c *Core) exception(code ExceptionCode) {
	// Usually kill the offender and schedule the next process
}
