package system

type pcb struct {
	ireg [32]uint32
	freg [32]uint64
	pc   uint32
	// TODO page table

}
