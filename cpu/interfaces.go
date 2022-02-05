package cpu

type System interface {
	HandleTrap(c *Core)
	HandleBoot(c *Core)
	Memory() *Memory
	ReservationSets() *ReservationSets
}
