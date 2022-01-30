package cpu

// TODO Should mmu be moved into its own package perhaps?

const (
	pageFlagValid    uint32 = 0x01 // the virtual address is valid
	pageFlagRead            = 0x02 // indicates that the processor is allowed to read data from this address
	pageFlagWrite           = 0x04 // indicates that the processor is allowed to write data to this address
	pageFlagExec            = 0x08 // indicates that the processor is allowed to fetch instructions from this address
	pageFlagUser            = 0x10 // indicates that the processor can access this page in user mode
	pageFlagGlobal          = 0x20 // whether this page is globally mapped into all address spaces (probably unused here)
	pageFlagAccessed        = 0x40 // whether this page has been accessed since the access bit was last cleared
	pageFlagDirty           = 0x80 // whether this page has been written to since the dirty bit was last cleared
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

type mmu struct {
}

// Translates the address and returns the flags that apply if the address is valid
// This method should return (false, false, 0, 0) when the address is invalid
// If a translation is valid, but the page is missing (when paging is implemented), the function should return (true, false, vAddr, flags)
// If a translation is valid and the page is present, the function should return (true, true, vAddr, flags)
func (c *Core) Translate(vAddr uint32) (hit bool, pAddr uint32, flags uint32) {
	return true, vAddr, pageFlagValid | pageFlagRead | pageFlagWrite | pageFlagExec
}

func newMMU() mmu {
	return mmu{}
}
