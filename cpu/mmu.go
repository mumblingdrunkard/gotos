package cpu

const (
	F_READ    uint8 = 0x01
	F_WRITE         = 0x02
	F_NOCACHE       = 0x04 // if I ever get around to doing MMIO
	F_EXEC          = 0x08
)

type MMU struct {
	base uint32
	size uint32
}

func (m *MMU) Translate(vAddr uint32) (err error, pAddr uint32, flags uint8) {
	return nil, vAddr + m.base, F_READ | F_WRITE
}

func NewMMU() MMU {
	return MMU{}
}
