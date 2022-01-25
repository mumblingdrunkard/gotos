package cpu

// TODO raise exceptions when addresses are misaligned

// TODO remove helper functions to not pollute namespace, only used once each anyway
func max(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func maxu(a, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func minu(a, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func (c *Core) lr_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	addr := c.reg[rs1]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	// update rset
	c.mc.rsets.Lock()
	c.mc.mem.Lock()
	success, w := c.unsafeLoadThroughWord(addr)
	if success {
		c.reg[rd] = w
	}
	c.mc.mem.Unlock()
	if success {
		m := c.mc.rsets.lookup[c.csr[Csr_MHARTID]]
		(*m)[pAddr] = true
	}
	c.mc.rsets.Unlock()
}

func (c *Core) sc_w(inst uint32) {
	// decode instruction
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)
	// TODO: Do usual checks

	c.mc.rsets.Lock()
	// check rset
	if _, ok := (*c.mc.rsets.lookup[c.csr[Csr_MHARTID]])[pAddr]; ok {
		c.mc.mem.Lock()
		success := c.unsafeStoreThroughWord(addr, c.reg[rs2])
		c.mc.mem.Unlock()

		if success {
			// invalidate entries on all harts
			c.reg[rd] = 0
			for i := range c.mc.rsets.sets {
				delete(*c.mc.rsets.lookup[i], pAddr)
			}
		}
	} else {
		// failed
		c.reg[rd] = 1
	}

	// Regardless of success or failure, executing an SC.W instruction invalidates any reservation held by this hart.
	delete(*c.mc.rsets.lookup[c.csr[Csr_MHARTID]], pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoswap_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		if c.unsafeStoreThroughWord(addr, src) {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amoadd_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		if c.unsafeStoreThroughWord(addr, src+w) {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LR in all cores
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amoand_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.unsafeLoadThroughWord(addr)
	c.unsafeStoreThroughWord(addr, src&c.reg[rd])
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, src&w)
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amoor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, src|w)
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amoxor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, src^w)
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amomax_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, uint32(max(int32(src), int32(w))))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amomaxu_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, maxu(src, w))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amomin_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, uint32(min(int32(src), int32(w))))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}

func (c *Core) amominu_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, _, pAddr, _ := c.mc.mmu.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, w := c.unsafeLoadThroughWord(addr)
	if lsuccess {
		ssuccess := c.unsafeStoreThroughWord(addr, minu(src, w))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}
