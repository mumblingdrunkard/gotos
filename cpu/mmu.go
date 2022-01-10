package cpu

const (
	MEM_F_READ    uint8 = 0x01 // indicates that the processor is allowed to read data from this address
	MEM_F_WRITE         = 0x02 // indicates that the processor is allowed to write data to this address
	MEM_F_EXEC          = 0x04 // indicates that the processor is allowed to fetch instructions from this address
	MEM_F_NOCACHE       = 0x08 // if I ever get around to doing MMIO
)

// Very simple MMU
type MMU struct {
	base uint32
	size uint32
}

// Translates the address and verifies that it is legal for the hart to access
// This method should return (false, false, 0, 0) when the address is invalid
// If a translation is valid, but the page is missing (when paging is implemented), the function should return (true, false, vAddr, flags)
// If a translation is valid and the page is present, the function should return (true, true, vAddr, flags)
func (m *MMU) TranslateAndCheck(vAddr uint32) (valid bool, present bool, pAddr uint32, flags uint8) {
	if pAddr >= m.size {
		return false, false, 0, 0
	}
	return true, true, vAddr + m.base, MEM_F_READ | MEM_F_WRITE | MEM_F_EXEC
}

func NewMMU() MMU {
	return MMU{}
}
