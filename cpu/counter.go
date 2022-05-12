package cpu

type counter struct {
	enable bool
	value  uint64
}

// SetCounter will enable the counter and set it to interrupt the core
// in v cycles.
//   It should not be set too low or LR/SC pairs may never succeed.
//   We recommend an absolute lower limit of 100, though it can and
// should be set much higher.
//   A value of 300'000 equates to approximately 5ms in optimal
// conditions.
func (c *Core) SetCounter(v uint64) {
	c.counter.enable = true
	c.counter.value = v
}
