package system

type PCB struct {
	ID        int
	PC        uint32
	Registers [32]uint32
}
