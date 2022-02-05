package main

import (
	"gotos/system"
)

func main() {
	// create a system with 1 core
	sys := system.NewSystem(1)

	// create a simple batch scheduler queue (FIFO)
	fifo := &system.FIFO{}
	sys.Scheduler = fifo

	// load the program
	sys.Load("c-programs/the-answer/main.text", 0x04000, 0x06000, 0)

	// run the system
	sys.Run()
}
