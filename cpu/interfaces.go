package cpu

import "sync"

// System is the interface that other packages should implement to
// communicate with Cores.
type System interface {
	// HandleTrap will be called when a core wants to trap
	//   See cpu.trap() for usage
	HandleTrap(*Core)
	// HandleBoot will be called when a core initially starts
	//   See cpu.run() for usage
	HandleBoot(*Core)

	// Shared resources
	Memory() *Memory
	ReservationSets() *ReservationSets
	InterruptMatrix() *InterruptMatrix

	// Tracking resources
	WgAwake() *sync.WaitGroup
	WgRunning() *sync.WaitGroup
}
