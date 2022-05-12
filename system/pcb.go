package system

type PCB struct {
	IReg   [32]uint32
	FReg   [32]uint64
	PC     uint32
	PID    uint32
	PTable uint32
}
