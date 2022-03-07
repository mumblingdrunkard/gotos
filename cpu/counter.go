package cpu

type counter struct {
	enable bool
	value  uint64
}

// lowest valid value for the counter should be 32
// 0 would make the counter underflow in the next cycle
// 1 would make the counter trigger in the next cycle
// 2 lets the processor make process (albeit slowly because of excessive interrupts) in most cases
// 32 guarantees that LR/SC loops can make progress as well
func (c *Core) SetCounter(v uint64) {
	c.counter.enable = true
	c.counter.value = v
}
