package cpu

import "encoding/binary"

func (c *Core) lr_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapLoadAddressMisaligned)
		return
	}

	success, pAddr := c.translate(addr, accessTypeLoad)

	if !success {
		return
	}

	pLine := pAddr >> cacheLineOffsetBits

	// update rset
	c.system.ReservationSets().Lock()
	c.system.Memory().Lock()
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])
	c.mc.dCache.store(pAddr, 4, uint64(w)) // attempt to update value in cache, don't care about success
	c.system.Memory().Unlock()
	if success {
		c.reg[rd] = w
	}
	if success {
		hid := int(c.csr[Csr_MHARTID])
		c.system.ReservationSets().unsafeRegisterReservation(hid, pLine)
	}
	c.system.ReservationSets().Unlock()
}

func (c *Core) sc_w(inst uint32) {
	// decode instruction
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]

	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}

	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock()
	// check rset
	hid := int(c.csr[Csr_MHARTID])
	if c.system.ReservationSets().unsafeCompareAndInvalidateReservation(hid, pLine) {
		var bytes [4]uint8
		binary.LittleEndian.PutUint32(bytes[:], c.reg[rs2])

		c.system.Memory().Lock()
		copy(c.system.Memory().data[pAddr:], bytes[:])
		c.mc.dCache.store(pAddr, 4, uint64(c.reg[rs2])) // attempt to update value in cache, don't care about success
		c.system.Memory().Unlock()

		if success {
			c.reg[rd] = 0
			// invalidate entries on all harts
			c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
		}
	} else {
		// failed
		c.reg[rd] = 1
	}

	// Regardless of success or failure, executing an SC.W instruction invalidates any reservation held by this hart.
	c.system.ReservationSets().unsafeCompareAndInvalidateReservation(hid, pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amoswap_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}

	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := src

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amoadd_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := w + src

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LR in all cores
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amoand_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := w & src

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amoor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := w | src

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amoxor_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := w ^ src

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amomax_w(inst uint32) {
	max := func(a, b int32) int32 {
		if a > b {
			return a
		} else {
			return b
		}
	}

	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := uint32(max(int32(w), int32(src)))

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amomaxu_w(inst uint32) {
	maxu := func(a, b uint32) uint32 {
		if a > b {
			return a
		} else {
			return b
		}
	}

	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := maxu(w, src)

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amomin_w(inst uint32) {
	min := func(a, b int32) int32 {
		if a < b {
			return a
		} else {
			return b
		}
	}

	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := uint32(min(int32(w), int32(src)))

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}

func (c *Core) amominu_w(inst uint32) {
	minu := func(a, b uint32) uint32 {
		if a < b {
			return a
		} else {
			return b
		}
	}

	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	addr := c.reg[rs1]
	// check alignment
	if addr&3 != 0 {
		c.csr[Csr_MTVAL] = addr
		c.trap(TrapStoreAddressMisaligned)
		return
	}

	src := c.reg[rs2]
	success, pAddr := c.translate(addr, accessTypeStore)

	if !success {
		return
	}
	pLine := pAddr >> cacheLineOffsetBits

	c.system.ReservationSets().Lock() // always lock rsets before system.Memory() to avoid deadlock
	c.system.Memory().Lock()
	// read bytes directly from memory
	w := binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])

	// calculate new value
	res := minu(w, src)

	// write value back to memory
	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], res)
	copy(c.system.Memory().data[pAddr:], bytes[:])

	// update cache
	c.mc.dCache.store(pAddr, 4, uint64(res))

	// store old value in rd
	c.reg[rd] = w
	c.system.Memory().Unlock()

	// Invalidate LRs
	c.system.ReservationSets().unsafeCompareAndInvalidateAllReservations(pLine)
	c.system.ReservationSets().Unlock()
}
