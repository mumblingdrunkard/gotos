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
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAccessFault)
		return false, 0
	}

	if flags&memFlagExec == 0 { // physical address is not marked executable
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAccessFault)
		return false, 0
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAddressMisaligned)
		return false, 0
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionPageFault)
		return false, 0
	}

	if hit, instruction := c.mc.iCache.loadWord(pAddr); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.iCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, instruction := c.mc.iCache.loadWord(pAddr)
		inst = instruction
	} else {
		inst = instruction
	}

	return true, inst
}

// Return the byte stored at
func (c *Core) loadByte(vAddr uint32) (bool, uint8) {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&memFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
		return false, 0
	}

	if hit, b := c.mc.dCache.loadByte(pAddr); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, b := c.mc.dCache.loadByte(pAddr)
		return true, b
	} else {
		return true, b
	}
}

func (c *Core) loadHalfWord(vAddr uint32) (bool, uint16) {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&memFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if !present { // possible page fault
		// TODO
		// When a page fault trap happens, nothing should happen to processor state.
		// Rather, the page should be fetched and loaded and the instruction should be retried.
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
		return false, 0
	}

	if pAddr&0x1 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	if hit, hw := c.mc.dCache.loadHalfWord(pAddr); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, hw := c.mc.dCache.loadHalfWord(pAddr)
		return true, hw
	} else {
		return true, hw
	}
}

func (c *Core) loadWord(vAddr uint32) (bool, uint32) {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&memFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if !present { // possible page fault
		// TODO
		// When a page fault trap happens, nothing should happen to processor state.
		// Rather, the page should be fetched and loaded and the instruction should be retried.
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
		return false, 0
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	if hit, w := c.mc.dCache.loadWord(pAddr); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, w := c.mc.dCache.loadWord(pAddr)
		return true, w
	} else {
		return true, w
	}
}

func (c *Core) loadDoubleWord(vAddr uint32) (bool, uint64) {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&memFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if !present { // possible page fault
		// TODO
		// When a page fault trap happens, nothing should happen to processor state.
		// Rather, the page should be fetched and loaded and the instruction should be retried.
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
		return false, 0
	}

	if pAddr&0x7 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	if hit, dw := c.mc.dCache.loadDoubleWord(pAddr); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		_, dw := c.mc.dCache.loadDoubleWord(pAddr)
		return true, dw
	} else {
		return true, dw
	}
}

// Return the byte stored at
func (c *Core) storeByte(vAddr uint32, b uint8) bool {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&memFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStorePageFault)
		return false
	}

	if hit := c.mc.dCache.storeByte(pAddr, b); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		c.mc.dCache.storeByte(pAddr, b)
	}

	return true
}

func (c *Core) storeHalfWord(vAddr uint32, hw uint16) bool {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&memFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStorePageFault)
		return false
	}

	if pAddr&0x1 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	if hit := c.mc.dCache.storeHalfWord(pAddr, hw); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		c.mc.dCache.storeHalfWord(pAddr, hw)
	}

	return false
}

func (c *Core) storeWord(vAddr uint32, w uint32) bool {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&memFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStorePageFault)
		return false
	}

	if pAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	if hit := c.mc.dCache.storeWord(pAddr, w); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		c.mc.dCache.storeWord(pAddr, w)
	}

	return true
}

func (c *Core) storeDoubleWord(vAddr uint32, dw uint64) bool {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&memFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStorePageFault)
		return false
	}

	if pAddr&0x7 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	if hit := c.mc.dCache.storeDoubleWord(pAddr, dw); !hit {
		c.mc.misses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.mc.mem.Lock()
		c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])
		c.mc.mem.Unlock()
		c.mc.dCache.storeDoubleWord(pAddr, dw)
	}

	return true
}

// Loads a memory straight from memory, bypassing the cache.
func (c *Core) unsafeLoadThroughWord(vAddr uint32) (bool, uint32) {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&memFlagRead == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if !present { // possible page fault
		// TODO
		// When a page fault trap happens, nothing should happen to processor state.
		// Rather, the page should be fetched and loaded and the instruction should be retried.
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadPageFault)
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
	// TODO: Verify the integrity of this
	// store loaded value into cache if it's cached
	// should the entire cache line just be invalidated instead perhaps?
	c.mc.dCache.storeWordNoDirty(pAddr, value)

	return true, value
}

// Stores a word straight to memory, bypassing cache.
func (c *Core) unsafeStoreThroughWord(vAddr uint32, w uint32) bool {
	c.mc.accesses++
	valid, present, pAddr, flags := c.mc.mmu.translateAndCheck(vAddr)

	if !valid { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if flags&memFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAccessFault)
		return false
	}

	if !present { // possible page fault
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStorePageFault)
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

// Flushes the data cache to memory
func (c *Core) flushCache() {
	c.mc.mem.Lock()
	c.mc.dCache.flushAll(c.mc.mem.data[:])
	c.mc.mem.Unlock()
}

// Invalidates the data cache
func (c *Core) invalidateCache() {
	c.mc.dCache.invalidateAll()
}

// Invalidates the instruction cache
func (c *Core) invalidateInstructionCache() {
	c.mc.iCache.invalidateAll()
}

// Flush and invalidate the data cache
func (c *Core) flushAndInvalidateCache() {
	c.flushCache()
	c.invalidateCache()
}

func (c *Core) Misses() uint64 {
	return c.mc.misses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
}
