package cpu

import "fmt"

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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)
	// TODO: Do usual checks

	// update rset
	c.mc.rsets.Lock()
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.mc.mem.Unlock()
	m := c.mc.rsets.lookup[c.id]
	(*m)[pAddr] = true
	c.mc.rsets.Unlock()
}

func (c *Core) sc_w(inst uint32) {
	// decode instruction
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)
	// TODO: Do usual checks

	c.mc.rsets.Lock()
	// check rset
	if _, ok := (*c.mc.rsets.lookup[c.id])[pAddr]; ok {
		fmt.Println("SC.W success!")
		c.mc.mem.Lock()
		c.UnsafeStoreThroughWord(addr, c.reg[rs2])
		c.mc.mem.Unlock()

		// invalidate entries on all harts
		c.reg[rd] = 0
		for i := range c.mc.rsets.sets {
			delete(*c.mc.rsets.lookup[i], pAddr)
		}
	} else {
		// failed
		c.reg[rd] = 1
	}

	// Regardless of success or failure, executing an SC.W instruction invalidates any reservation held by this hart.
	delete(*c.mc.rsets.lookup[c.id], pAddr)
	c.mc.rsets.Unlock()
}

func (c *Core) amoswap_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	src := c.reg[rs2]
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, src)
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, src+c.reg[rd])
	c.mc.mem.Unlock()

	// Invalidate LRs
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, src&c.reg[rd])
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, src|c.reg[rd])
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, src^c.reg[rd])
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, uint32(max(int32(src), int32(c.reg[rd]))))
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, maxu(src, c.reg[rd]))
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, uint32(min(int32(src), int32(c.reg[rd]))))
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
	_, _, pAddr, _ := c.mc.mmu.TranslateAndCheck(addr)

	c.mc.rsets.Lock() // always lock rsets before mc.mem to avoid deadlock
	c.mc.mem.Lock()
	_, c.reg[rd] = c.UnsafeLoadThroughWord(addr)
	c.UnsafeStoreThroughWord(addr, minu(src, c.reg[rd]))
	c.mc.mem.Unlock()

	// Invalidate LRs
	for i := range c.mc.rsets.sets {
		delete(*c.mc.rsets.lookup[i], pAddr)
	}
	c.mc.rsets.Unlock()
}
