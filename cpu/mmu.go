package cpu

// TODO Should mmu be moved into its own package perhaps?

const (
	memFlagRead    uint8 = 0x01 // indicates that the processor is allowed to read data from this address
	memFlagWrite         = 0x02 // indicates that the processor is allowed to write data to this address
	memFlagExec          = 0x04 // indicates that the processor is allowed to fetch instructions from this address
	memFlagNoCache       = 0x08 // if I ever get around to doing MMIO
)

//
// Very simple mmu
// MMU likely requires at least 4 base-bound pairs for general purpose use, but likely should have 8 or more.
// The segments I propose are:
//
// Name       | Flag | Purpose
// -----------|------|--------
// SegProgram | X    | Instruction memory
// SegData    |  R   | Read only data
// SegHeap    |  RW  | Program heap (can explicitly be resized by the program)
// SegStack   |  RW  | Program stack (should imlicitly grow on demand)
//
// This layout does not allow for uncached I/O and it would have to be done all through syscalls.
//
// Paging is the better option, but I think students should be gradually introduced to the concepts
// Might even do well with just a single segment at first, then two, then four, etc. etc..
//

type mmu struct {
	base uint32
	size uint32
}

// Translates the address and returns the flags that apply if the address is valid
// This method should return (false, false, 0, 0) when the address is invalid
// If a translation is valid, but the page is missing (when paging is implemented), the function should return (true, false, vAddr, flags)
// If a translation is valid and the page is present, the function should return (true, true, vAddr, flags)
func (m *mmu) translateAndCheck(vAddr uint32) (valid bool, present bool, pAddr uint32, flags uint8) {
	if pAddr >= m.size {
		return false, false, 0, 0
	}
	return true, true, vAddr + m.base, memFlagRead | memFlagWrite | memFlagExec
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
