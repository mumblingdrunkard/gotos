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

const (
	cacheEnable = true
)

// Attempts to load a 4 byte instruction stored at virtual address `vAddr`.
// If successful, returns `true, <instruction>`, `false, 0` otherwise.
func (c *Core) loadInstruction(vAddr uint32) (bool, uint32) {
	c.mc.accesses++
	var inst uint32
	_, pAddr, flags := c.Translate(vAddr)

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

	if !cacheEnable {
		c.mc.mem.Lock()
		defer c.mc.mem.Unlock()
		if c.mc.mem.endian == EndianBig {
			return true, binary.BigEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4])
		} else {
			return true, binary.LittleEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4])
		}
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
	_, pAddr, flags := c.Translate(vAddr)

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

	if !cacheEnable {
		c.mc.mem.Lock()
		defer c.mc.mem.Unlock()

		if width == 1 {
			return true, uint64(c.mc.mem.data[pAddr])
		}

		if c.mc.mem.endian == EndianBig {
			switch width {
			case 2:
				return true, uint64(binary.BigEndian.Uint16(c.mc.mem.data[pAddr : pAddr+2]))
			case 4:
				return true, uint64(binary.BigEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4]))
			case 8:
				return true, binary.BigEndian.Uint64(c.mc.mem.data[pAddr : pAddr+8])
			default:
				panic("Invalid load width")
			}
		} else {
			switch width {
			case 2:
				return true, uint64(binary.LittleEndian.Uint16(c.mc.mem.data[pAddr : pAddr+4]))
			case 4:
				return true, uint64(binary.LittleEndian.Uint32(c.mc.mem.data[pAddr : pAddr+4]))
			case 8:
				return true, binary.LittleEndian.Uint64(c.mc.mem.data[pAddr : pAddr+4])
			default:
				panic("Invalid load width")
			}
		}
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
	_, pAddr, flags := c.Translate(vAddr)

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

	if !cacheEnable {
		c.mc.mem.Lock()
		defer c.mc.mem.Unlock()

		if width == 1 {
			c.mc.mem.data[pAddr] = uint8(v)
			return true
		}

		var bytes [8]uint8

		if c.mc.mem.endian == EndianBig {
			switch width {
			case 2:
				binary.BigEndian.PutUint16(bytes[:], uint16(v))
			case 4:
				binary.BigEndian.PutUint32(bytes[:], uint32(v))
			case 8:
				binary.BigEndian.PutUint64(bytes[:], v)
			default:
				panic("Invalid store width")
			}
		} else {
			switch width {
			case 2:
				binary.LittleEndian.PutUint16(bytes[:], uint16(v))
			case 4:
				binary.LittleEndian.PutUint32(bytes[:], uint32(v))
			case 8:
				binary.LittleEndian.PutUint64(bytes[:], v)
			default:
				panic("Invalid store width")
			}
		}

		copy(c.mc.mem.data[pAddr:], bytes[:width])
		return true
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

// Requires memory to be locked before calling
func (c *Core) unsafeLoadAtomic(vAddr, width uint32) (bool, uint64) {
	c.mc.accesses++
	_, pAddr, flags := c.Translate(vAddr)

	if flags&pageFlagValid == 0 { // address was invalid
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if flags&pageFlagWrite == 0 { // permissions
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAccessFault)
		return false, 0
	}

	if pAddr&(width-1) != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	lineNumber := pAddr >> cacheLineOffsetBits

	// writebackLine if it is present and dirty
	c.mc.dCache.writebackLine(lineNumber, c.mc.mem.data[:])

	// invalidateLine if it is present
	c.mc.dCache.invalidateLine(lineNumber)

	// cache replace will refresh the line if present, or eject a random line
	c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])

	// load value from cache
	return c.mc.dCache.load(pAddr, width)
}

// Requires memory to be locked before calling
func (c *Core) unsafeStoreAtomic(vAddr, width uint32, v uint64) bool {
	c.mc.accesses++
	_, pAddr, flags := c.Translate(vAddr)

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

	lineNumber := pAddr >> cacheLineOffsetBits

	// ensure line is in cache
	c.mc.dCache.replaceRandom(lineNumber, cacheFlagNone, c.mc.mem.data[:])

	// write data to cache
	c.mc.dCache.store(pAddr, width, v)

	// writeback line
	c.mc.dCache.writebackLine(lineNumber, c.mc.mem.data[:])

	return true
}

// Writes the data cache to memory
func (c *Core) CacheWriteback() {
	c.mc.mem.Lock()
	c.mc.dCache.writebackAll(c.mc.mem.data[:])
	c.mc.mem.Unlock()
}

// Invalidates the data cache
func (c *Core) CacheInvalidate() {
	c.mc.dCache.invalidateAll()
}

// Invalidates the instruction cache
func (c *Core) InstructionCacheInvalidate() {
	c.mc.iCache.invalidateAll()
}

func (c *Core) ReadMemory(addr, bytes uint32) (error, []uint8) {
	return c.mc.mem.Read(addr, bytes)
}

func (c *Core) WriteMemory(addr uint32, data []uint8) (error, int) {
	return c.mc.mem.Write(addr, data)
}

// Writeback and invalidate the data cache
func (c *Core) CacheWritebackAndInvalidate() {
	c.CacheWriteback()
	c.CacheInvalidate()
}

func (c *Core) Misses() uint64 {
	return c.mc.misses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
}
