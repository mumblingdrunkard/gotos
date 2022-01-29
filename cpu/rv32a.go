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
	_, pAddr, _ := c.translateAndCheck(addr)

	// update rset
	c.mc.rsets.Lock()
	c.mc.mem.Lock()
	success, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
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
	_, pAddr, _ := c.translateAndCheck(addr)
	// TODO: Do usual checks

	c.mc.rsets.Lock()
	// check rset
	if _, ok := (*c.mc.rsets.lookup[c.csr[Csr_MHARTID]])[pAddr]; ok {
		c.mc.mem.Lock()
		success := c.unsafeStoreAtomic(addr, 4, uint64(c.reg[rs2]))
		c.mc.mem.Unlock()

		if success {
			c.reg[rd] = 0
			// invalidate entries on all harts
			c.mc.rsets.unsafeInvalidateAll(pAddr)
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
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		if c.unsafeStoreAtomic(addr, 4, uint64(src)) {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoadd_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		if c.unsafeStoreAtomic(addr, 4, uint64(src+w)) {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LR in all cores
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoand_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(src&w))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(src|w))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoxor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(src^w))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amomax_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(max(int32(src), int32(w))))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amomaxu_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(maxu(src, w)))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amomin_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(min(int32(src), int32(w))))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amominu_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, pAddr, _ := c.translateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	lsuccess, v := c.unsafeLoadAtomic(addr, 4)
	w := uint32(v)
	if lsuccess {
		ssuccess := c.unsafeStoreAtomic(addr, 4, uint64(minu(src, w)))
		if ssuccess {
			c.reg[rd] = w
		}
	}
	c.mc.mem.Unlock()

	// Invalidate LRs
	c.mc.rsets.unsafeInvalidateAll(pAddr)
	c.mc.rsets.Unlock()
}
