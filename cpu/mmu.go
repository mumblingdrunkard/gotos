package cpu

// TODO Should mmu be moved into its own package perhaps?

const (
	mmuFlagRead    uint8 = 0x01 // indicates that the processor is allowed to read data from this address
	mmuFlagWrite         = 0x02 // indicates that the processor is allowed to write data to this address
	mmuFlagExec          = 0x04 // indicates that the processor is allowed to fetch instructions from this address
	mmuFlagPresent       = 0x08 // the virtual address maps to a physical one, if this flag is 0, the returned address should be 0
	mmuFlagValid         = 0x10 // the virtual address is valid
	mmuFlagNoCache       = 0x20 // if I ever get around to doing MMIO
)

type page struct {
	frameNumber uint32
	flags       uint8
}

// `NNNNNFF0` where NNNNN is the physical frame number that it maps to and FF are the memory flags
type tlbEntry uint32

type tlb struct {
}

type mmu struct {
	base uint32
	size uint32
}

// Translates the address and returns the flags that apply if the address is valid
// This method should return (false, false, 0, 0) when the address is invalid
// If a translation is valid, but the page is missing (when paging is implemented), the function should return (true, false, vAddr, flags)
// If a translation is valid and the page is present, the function should return (true, true, vAddr, flags)
func (c *Core) translateAndCheck(vAddr uint32) (hit bool, pAddr uint32, flags uint8) {
	if pAddr >= c.mc.mmu.size {
		return true, 0, 0
	}
	return true, vAddr + c.mc.mmu.base, mmuFlagRead | mmuFlagWrite | mmuFlagExec | mmuFlagPresent
}

func newMMU() mmu {
	return mmu{}
}

func (c *Core) UnsafeSetMemBase(base uint32) {
	c.mc.mmu.base = base
}

func (c *Core) UnsafeSetMemSize(size uint32) {
	c.mc.mmu.size = size
	c.reg[2] = c.mc.mmu.size
}
