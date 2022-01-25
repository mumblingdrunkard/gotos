package cpu

func (c *Core) fence_i(inst uint32) {
	c.instructionCacheInvalidate()
}
