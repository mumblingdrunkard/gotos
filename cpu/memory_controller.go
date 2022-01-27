// RISC-V privileged specification says:
//
// ---
//
// If **mtval** is written with a nonzero value when a breakpoint, address misaligned, access-fault, or page-fault exception occurs on an instruction fetch, load, or store, then **mtval** will contain the faulting virtual address.
//
// If **mtval** is written with a nonzero value when a misaligned load or store causes an access-fault or page-fault exception occurs, then **mtval** will contain the virtual address o the portion of the access that caused the fault.
//
// ---
//
// So **mtval** should contain the virtual address that caused the fault.

package cpu

import (
	"encoding/binary"
)

type memoryController struct {
	iCache   cache
	dCache   cache
	mem      *Memory          // RAM (possibly shared)
	rsets    *ReservationSets // Reservation sets (possibly shared)
	mmu      mmu
	misses   uint64
	accesses uint64
}

func newMemoryController(m *Memory, rs *ReservationSets) memoryController {
	return memoryController{
		dCache: newCache(m.endian),
		iCache: newCache(m.endian),
		rsets:  rs,
		mem:    m,
		mmu:    newMMU(),
	}
}

// Attempts to load a 4 byte instruction stored at virtual address `vAddr`.
// If successful, returns `true, <instruction>`, `false, 0` otherwise.
func (c *Core) loadInstruction(vAddr uint32) (bool, uint32) {
	c.mc.accesses++
	var inst uint32
	_, pAddr, flags := c.translateAndCheck(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionPageFault)
		return false, 0
	}

	if flags&pageFlagExec == 0 { // physical address is not marked executable
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAccessFault)
		return false, 0
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAddressMisaligned)
		return false, 0
	}

	if hit, instruction := c.mc.iCache.load(pAddr, 4); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.iCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, instruction := c.mc.iCache.load(pAddr, 4)
		inst = uint32(instruction)
	} else {
		inst = uint32(instruction)
	}

	return true, inst
}

// loads up to 8 bytes
func (c *Core) load(vAddr, width uint32) (bool, uint64) {
	c.mc.accesses++
	_, pAddr, flags := c.translateAndCheck(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
		return false, 0
	}

	if flags&pageFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if pAddr&(width-1) != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	if hit, v := c.mc.dCache.load(pAddr, width); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, v := c.mc.dCache.load(pAddr, width)
		return true, v
	} else {
		return true, v
	}
}

func (c *Core) store(vAddr, width uint32, v uint64) bool {
	c.mc.accesses++
	_, pAddr, flags := c.translateAndCheck(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&pageFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if pAddr&(width-1) != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	if hit := c.mc.dCache.store(pAddr, width, v); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		c.mc.dCache.store(pAddr, width, v)
	}

	return true
}

// Return the byte stored at
func (c *Core) loadByte(vAddr uint32) (bool, uint8) {
	success, v := c.load(vAddr, 1)
	return success, uint8(v)
}

func (c *Core) loadHalfWord(vAddr uint32) (bool, uint16) {
	success, v := c.load(vAddr, 2)
	return success, uint16(v)
}

func (c *Core) loadWord(vAddr uint32) (bool, uint32) {
	success, v := c.load(vAddr, 4)
	return success, uint32(v)
}

func (c *Core) loadDoubleWord(vAddr uint32) (bool, uint64) {
	success, v := c.load(vAddr, 8)
	return success, v
}

func (c *Core) storeByte(vAddr uint32, b uint8) bool {
	return c.store(vAddr, 1, uint64(b))
}

func (c *Core) storeHalfWord(vAddr uint32, hw uint16) bool {
	return c.store(vAddr, 2, uint64(hw))
}

func (c *Core) storeWord(vAddr uint32, w uint32) bool {
	return c.store(vAddr, 4, uint64(w))
}

func (c *Core) storeDoubleWord(vAddr uint32, dw uint64) bool {
	return c.store(vAddr, 8, dw)
}

// Loads a memory straight from memory, bypassing the cache.
func (c *Core) unsafeLoadThroughWord(vAddr uint32) (bool, uint32) {
	c.mc.accesses++
	_, pAddr, flags := c.translateAndCheck(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&pageFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	var value uint32
	if c.mc.mem.endian == EndianBig {
		value = binary.BigEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4])
	} else {
		value = binary.LittleEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4])
	}
	c.mc.dCache.storeWordNoDirty(pAddr, value)

	return true, value
}

// Stores a word straight to memory, bypassing cache.
func (c *Core) unsafeStoreThroughWord(vAddr uint32, w uint32) bool {
	c.mc.accesses++
	_, pAddr, flags := c.translateAndCheck(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&pageFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	var bytes [4]uint8

	if c.mc.mem.endian == EndianBig {
		binary.BigEndian.PutUint32(bytes[:], w)
	} else {
		binary.LittleEndian.PutUint32(bytes[:], w)
	}

	copy(c.mc.mem.data[pAddr:], bytes[:])

	// also update cache
	c.mc.dCache.storeWordNoDirty(pAddr, w) // May be uncached, ignore

	return true
}

// Writes the data cache to memory
func (c *Core) cacheWriteback() {
	c.mc.mem.Lock()
	c.mc.dCache.writebackAll(c.mc.mem.data[:])
	c.mc.mem.Unlock()
}

// Invalidates the data cache
func (c *Core) cacheInvalidate() {
	c.mc.dCache.invalidateAll()
}

// Invalidates the instruction cache
func (c *Core) instructionCacheInvalidate() {
	c.mc.iCache.invalidateAll()
}

// Writeback and invalidate the data cache
func (c *Core) cacheWritebackAndInvalidate() {
	c.cacheWriteback()
	c.cacheInvalidate()
}

func (c *Core) Misses() uint64 {
	return c.mc.misses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
}
