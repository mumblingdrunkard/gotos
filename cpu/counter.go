package cpu

type counter struct {
	enable bool
	value  uint64
}

func (c *Core) EnableCounter() {
	c.counter.enable = true
}

func (c *Core) SetCounter(v uint64) {
	c.counter.value = v
}
