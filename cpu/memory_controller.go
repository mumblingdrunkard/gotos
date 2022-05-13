// this file contains structs and methods that manage memory

package cpu

import (
	"encoding/binary"
)

type memoryController struct {
	iCache cache // instruction cache
	dCache cache // data cache
	tlb    tlb   // level 0 tlb - normal pages
}

// newMemoryController returns a new memory controller, complete with data
// cache, instruction cache, and a tlb.
func newMemoryController() memoryController {
	return memoryController{
		dCache: newCache(),
		iCache: newCache(),
		tlb:    newTLB(),
	}
}

// loadInstruction attempts to load a 4 byte instruction stored at virtual
// address `vAddr`.
//   If successful, returns `true, instruction`, `false, 0` otherwise.
func (c *Core) loadInstruction(vAddr uint32) (bool, uint32) {
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

// load will attempt to load `width` bytes from the virtual address `vAddr`.
//   `vAddr` has to be aligned on a `width` byte boundary.
//   On success, returns `true, v` where `v` is a uint64 and the requested data
// is right-aligned in `v`.
//   On failure, returns `false, 0`.
func (c *Core) load(vAddr, width uint32) (bool, uint64) {
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

// store will attempt to store `width` bytes to the virtual address `vAddr`.
//   `vAddr` has to be aligned on a `width` byte boundary.
//   The value passed to store should be right-aligned in `v`.
//   On success, returns `true`.
//   On failure, returns `false`.
func (c *Core) store(vAddr, width uint32, v uint64) bool {
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
		lineNumber := pAddr >> cacheLineOffsetBits
		c.system.Memory().Lock()
		c.mc.dCache.replace(lineNumber, cacheFlagNone, c.system.Memory().data[:])
		c.system.Memory().Unlock()
		c.mc.dCache.store(pAddr, width, v)
	}

	return true
}

// loadByte attempts to load a single byte from the virtual address `vAddr`.
//   This function is a wrapper for Core.load
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

// SC_W is intended to be used by a system to invalidate reservations that a
// core may be holding.
func (c *Core) SC_W() {
	c.system.ReservationSets().Lock()
	c.system.ReservationSets().unsafeInvalidateSingle(int(c.csr[Csr_MHARTID]), 0)
	c.system.ReservationSets().Unlock()
}

// FENCE invalidates the data cache.
func (c *Core) FENCE() {
	c.system.Memory().Lock()
	c.mc.dCache.writebackAll(c.system.Memory().data[:])
	c.system.Memory().Unlock()
	c.mc.dCache.invalidateAll()
}

// FENCE_I flushes written data and invalidates the instruction cache.
func (c *Core) FENCE_I() {
	c.system.Memory().Lock()
	c.mc.dCache.writebackAll(c.system.Memory().data[:])
	c.system.Memory().Unlock()
	c.mc.iCache.invalidateAll()
}

const (
	SFENCE_VMA_ALL       uint32 = 0
	SFENCE_VMA_ASID             = 1
	SFENCE_VMA_ADDR             = 2
	SFENCE_VMA_ASID_ADDR        = 3
)

// SFENCE_VMA invalidates translation caches (tlb).
func (c *Core) SFENCE_VMA(asid, vAddr, flag uint32) {
	// TODO discriminate on asid and vAddr depending on the flags
	c.mc.tlb.invalidateAll()
}

// AtomicStoreWordPhysicalUncached will atomically store a single word `w` to
// the physical address `pAddr`.
//   Misaligned access causes this function to fail with `false`.
//   Otherwise, the memory is locked, the word is written, and this function
// returns `true`.
//
//   This function, along with Core.AtomicLoadWordPhysicalUncached should only
// be used by the system when atomic access is required and access should be
// uncached (such as when modifying the page table).
//   In all other instances, systems should use the Core.Read and Core.Write
// functions to access the virtual address space currently in use.
func (c *Core) AtomicStoreWordPhysicalUncached(pAddr, w uint32) bool {
	if pAddr&0x3 != 0 { // misaligned access
		return false
	}

	var bytes [4]uint8
	binary.LittleEndian.PutUint32(bytes[:], w)
	c.system.Memory().Lock()
	copy(c.system.Memory().data[pAddr:], bytes[:])
	c.system.Memory().Unlock()
	return true
}

// AtomicLoadWordPhysicalUncached will atomically load a single word from the
// physical address `pAddr`.
//   Misaligned access causes this function to fail with `false, 0`.
//   Otherwise, the memory is locked, the word is written, and this function
// returns `true, w` where `w` is the word.
//
//   This function, along with Core.AtomicStoreWordPhysicalUncached should only
// be used by the system when atomic access is required and access should be
// uncached (such as when modifying the page table).
//   In all other instances, systems should use the Core.Read and Core.Write
// functions to access the virtual address space currently in use.
func (c *Core) AtomicLoadWordPhysicalUncached(pAddr uint32) (bool, uint32) {
	if pAddr&0x3 != 0 { // misaligned access
		return false, 0
	}

	c.system.Memory().Lock()
	defer c.system.Memory().Unlock()
	return true, binary.LittleEndian.Uint32(c.system.Memory().data[pAddr : pAddr+4])
}

// Read reads `n` bytes from the core's current address space (might be
// virtual).
func (c *Core) Read(addr, n uint32) (error, []uint8) {
	panic("Core.Read not implemented!")
}

// Write writes len(data) bytes from the core's current address space
// (might be virtual).
func (c *Core) Write(addr uint32, data []uint8) (error, int) {
	panic("Core.Write not implemented!")
}
