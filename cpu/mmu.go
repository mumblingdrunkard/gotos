package cpu

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

type accessType int

const (
	accessTypeLoad             accessType = 0
	accessTypeStore                       = 1
	accessTypeInstructionFetch            = 2
)

// Translates the address and returns the flags that apply if the address is valid
// This method should return (false, false, 0, 0) when the address is invalid
// If a translation is valid, but the page is missing (when paging is implemented), the function should return (true, false, vAddr, flags)
// If a translation is valid and the page is present, the function should return (true, true, vAddr, flags)
func (c *Core) translate(vAddr uint32, at accessType) (success bool, pAddr uint32) {
	// get the satp register
	satp := c.csr[Csr_SATP]
	if satp&0x80000000 == 0 { // bare mode, no translation or protection
		return true, vAddr
	}

	return false, 0
}
