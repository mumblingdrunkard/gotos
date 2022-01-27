package cpu

// TODO Should mmu be moved into its own package perhaps?

const (
	mmuFlagValid       uint8 = 0x01 // the virtual address is valid
	mmuFlagRead              = 0x02 // indicates that the processor is allowed to read data from this address
	mmuFlagWrite             = 0x04 // indicates that the processor is allowed to write data to this address
	mmuFlagExec              = 0x08 // indicates that the processor is allowed to fetch instructions from this address
	mmuFlagUModeAccess       = 0x10
	mmuFlagGlobal            = 0x20
	mmuFlagAccessed          = 0x40
	mmuFlagDirty             = 0x80
)

// Two schemes to manage the A and D bits are permitted:
// - When a virtual page is accessed and the A bit is clear, or is written and
//   the D bit is clear, a page-fault exception is raised.
// - ...
//
// ...
//
//     ---
//         The A and D bits are never cleared by the implementation.
//     If the supervisor software does not rely on accessed and/or dirty bits,
//     e.g. if it does not swap memory pages to secondary storage or if the
//     pages are being used to map I/O space, it should always set them to 1 in
//     the PTE to improve performance.
//     ---

type page struct {
	frameNumber uint32
	flags       uint8
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
	return true, vAddr + c.mc.mmu.base, mmuFlagValid | mmuFlagRead | mmuFlagWrite | mmuFlagExec
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
