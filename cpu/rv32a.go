package cpu

func (c *Core) lr_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	addr := c.reg[rs1]
	_, pAddr, _ := c.mc.mmu.Translate(addr)
	// TODO: Do usual checks

	// update rset
	c.rsets.Lock()
	_, c.reg[rd] = c.mc.LoadThroughWord(addr)
	c.rsets.sets[c.id][pAddr] = true
	c.rsets.Unlock()
}

func (c *Core) sc_w(inst uint32) {
	// decode instruction
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	_, pAddr, _ := c.mc.mmu.Translate(addr)
	// TODO: Do usual checks

	c.rsets.Lock()

	// check rset
	if _, ok := c.rsets.sets[c.id][pAddr]; ok {
		c.mc.StoreThroughWord(addr, c.reg[rs2])

		// invalidate entries on all harts
		c.reg[rd] = 0
		for i := range c.rsets.sets {
			delete(c.rsets.sets[i], pAddr)
		}
	} else {
		// failed
		c.reg[rd] = 1
	}

	// Regardless of success or failure, executing an SC.W instruction invalidates any reservation held by this hart.
	delete(c.rsets.sets[c.id], pAddr)

	c.rsets.Unlock()
}

func (c *Core) amoswap_w(inst uint32) {
}

func (c *Core) amoadd_w(inst uint32) {
}

func (c *Core) amoand_w(inst uint32) {
}

func (c *Core) amoor_w(inst uint32) {
}

func (c *Core) amoxor_w(inst uint32) {
}

func (c *Core) amomax_w(inst uint32) {
}

func (c *Core) amomaxu_w(inst uint32) {
}

func (c *Core) amomin_w(inst uint32) {
}

func (c *Core) amominu_w(inst uint32) {
}
