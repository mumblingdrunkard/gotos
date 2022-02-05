package main

import (
	"gotos/system"
)

func main() {
	// create a system with 1 core
	sys := system.NewSystem(1)

	sys.Memory().Write(0, []uint8{
		0x93, 0x02, 0x10, 0x01,
		0x13, 0x03, 0x90, 0x01,
		0x33, 0x85, 0x62, 0x00,
	})

	// run the system
	for i := 0; i < 3; i++ {
		sys.StepAndDump()
	}
}
