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
	iCache cache // instruction cache
	dCache cache // data cache
	tCache cache // translation cache
	tlb0   tlb   // level 0 tlb - normal pages

	// metrics
	cacheMisses uint64
	accesses    uint64
}

func newMemoryController() memoryController {
	return memoryController{
		dCache: newCache(),
		iCache: newCache(),
		tCache: newCache(),
		tlb0:   newTLB(),
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

	if vAddr&0x3 != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapInstructionAddressMisaligned)
		return false, 0
	}

	success, pAddr := c.translate(vAddr, accessTypeInstructionFetch)

	if !success {
		return false, 0
	}

	if !cacheEnable {
		c.system.Memory().Lock()
		defer c.system.Memory().Unlock()
		return true, binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])
	}

	if hit, instruction := c.mc.iCache.load(pAddr, 4); !hit {
		c.mc.cacheMisses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.system.Memory().Lock()
		c.mc.iCache.replace(lineNumber, cacheFlagNone, c.system.Memory().data[:])
		c.system.Memory().Unlock()
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

	if vAddr&(width-1) != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapLoadAddressMisaligned)
		return false, 0
	}

	success, pAddr := c.translate(vAddr, accessTypeLoad)

	if !success {
		return false, 0
	}

	if !cacheEnable {
		c.system.Memory().Lock()
		defer c.system.Memory().Unlock()

		if width == 1 {
			return true, uint64(c.system.Memory().data[pAddr])
		}

		switch width {
		case 2:
			return true, uint64(binary.LittleEndian.Uint16(c.system.Memory().data[pAddr : pAddr+4]))
		case 4:
			return true, uint64(binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4]))
		case 8:
			return true, binary.LittleEndian.Uint64(c.system.Memory().data[pAddr : pAddr+4])
		default:
			panic("Invalid load width")
		}
	}

	if hit, v := c.mc.dCache.load(pAddr, width); !hit {
		c.mc.cacheMisses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.system.Memory().Lock()
		c.mc.dCache.replace(lineNumber, cacheFlagNone, c.system.Memory().data[:])
		c.system.Memory().Unlock()
		_, v := c.mc.dCache.load(pAddr, width)
		return true, v
	} else {
		return true, v
	}
}

func (c *Core) store(vAddr, width uint32, v uint64) bool {
	c.mc.accesses++

	if vAddr&(width-1) != 0 { // address alignment
		c.csr[Csr_MTVAL] = vAddr
		c.trap(TrapStoreAddressMisaligned)
		return false
	}

	success, pAddr := c.translate(vAddr, accessTypeStore)

	if !success {
		return false
	}

	if !cacheEnable {
		c.system.Memory().Lock()
		defer c.system.Memory().Unlock()

		if width == 1 {
			c.system.Memory().data[pAddr] = uint8(v)
			return true
		}

		var bytes [8]uint8

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

		copy(c.system.Memory().data[pAddr:], bytes[:width])
		return true
	}

	if hit := c.mc.dCache.store(pAddr, width, v); !hit {
		c.mc.cacheMisses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.system.Memory().Lock()
		c.mc.dCache.replace(lineNumber, cacheFlagNone, c.system.Memory().data[:])
		c.system.Memory().Unlock()
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

// load a word without translating the address first. Used for the page table walker
func (c *Core) loadWordPhysical(pAddr uint32) (bool, uint32) {
	c.mc.accesses++
	if !cacheEnable {
		c.system.Memory().Lock()
		defer c.system.Memory().Unlock()
		return true, binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])
	}

	if hit, v := c.mc.tCache.load(pAddr, 4); !hit {
		c.mc.cacheMisses++
		lineNumber := pAddr >> cacheLineOffsetBits
		c.system.Memory().Lock()
		c.mc.tCache.replace(lineNumber, cacheFlagNone, c.system.Memory().data[:])
		c.system.Memory().Unlock()
		_, v := c.mc.tCache.load(pAddr, 4)
		return true, uint32(v)
	} else {
		return true, uint32(v)
	}
}

// Writes the data cache to memory
func (c *Core) DataCacheWriteback() {
	c.system.Memory().Lock()
	c.mc.dCache.writebackAll(c.system.Memory().data[:])
	c.system.Memory().Unlock()
}

// Invalidates the data cache
func (c *Core) DataCacheInvalidate() {
	c.mc.dCache.invalidateAll()
}

// Invalidates the instruction cache
func (c *Core) InstructionCacheInvalidate() {
	c.mc.iCache.invalidateAll()
}

func (c *Core) TLBInvalidate() {
	c.mc.tlb0.invalidateAll()
}

func (c *Core) TranslationCacheInvalidate() {
	c.mc.tCache.invalidateAll()
}

func (c *Core) SignalVirtualMemoryUpdates() {
	c.vmaUpd.Store(true)
}

// reads n bytes from the (possibly virtual) address addr and out
func (c *Core) Read(addr, n uint32) (error, []uint8) {

	return c.system.Memory().Read(addr, n)
}

func (c *Core) Write(addr uint32, data []uint8) (error, int) {
	return c.system.Memory().Write(addr, data)
}

// Writeback and invalidate the data cache
func (c *Core) CacheWritebackAndInvalidate() {
	c.DataCacheWriteback()
	c.DataCacheInvalidate()
}

func (c *Core) Misses() uint64 {
	return c.mc.cacheMisses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
}
