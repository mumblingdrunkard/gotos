package cpu

import "encoding/binary"

type MemoryController struct {
	iCache   Cache
	dCache   Cache
	mem      *Memory
	mmu      MMU
	misses   uint64
	accesses uint64
}

func (mc *MemoryController) LoadInstruction(address uint32) (bool, uint32) {
	mc.accesses++
	var inst uint32
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_INSTRUCTION_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_INSTRUCTION_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_EXEC == 0 { // permissions
		// TODO TRAP_INSTRUCTION_ACCESS_FAULT
		return false, 0
	}

	if address&0x3 != 0 { // address alignment
		// TODO TRAP_INSTRUCTION_ADDRESS_MISALIGNED
		return false, 0
	}

	// TODO Check if address is uncached
	if hit, instruction := mc.iCache.LoadWord(address); !hit {
		// fmt.Println("Cache miss!")
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		// fmt.Printf("Loading line #%04X\n", lineNumber)
		mc.mem.Lock()
		mc.iCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		// fmt.Printf("Loaded line #%04X\n", lineNumber)
		_, instruction := mc.iCache.LoadWord(address)
		inst = instruction
	} else {
		inst = instruction
	}

	return true, inst
}

// Return the byte stored at
func (mc *MemoryController) LoadByte(address uint32) (bool, uint8) {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_LOAD_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_READ == 0 { // permissions
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if hit, b := mc.dCache.LoadByte(address); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, b := mc.dCache.LoadByte(address)
		return true, b
	} else {
		return true, b
	}
}

func (mc *MemoryController) LoadHalfWord(address uint32) (bool, uint16) {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_LOAD_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_READ == 0 { // permissions
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if address&0x1 != 0 { // address alignment
		// TODO TRAP_LOAD_ADDRESS_MISALIGNED
		return false, 0
	}

	if hit, hw := mc.dCache.LoadHalfWord(address); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, hw := mc.dCache.LoadHalfWord(address)
		return true, hw
	} else {
		return true, hw
	}
}

func (mc *MemoryController) LoadWord(address uint32) (bool, uint32) {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_LOAD_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_READ == 0 { // permissions
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if address&0x3 != 0 { // address alignment
		// TODO TRAP_LOAD_ADDRESS_MISALIGNED
		return false, 0
	}

	if hit, w := mc.dCache.LoadWord(address); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, w := mc.dCache.LoadWord(address)
		return true, w
	} else {
		return true, w
	}
}

func (mc *MemoryController) LoadDoubleWord(address uint32) (bool, uint64) {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_LOAD_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_READ == 0 { // permissions
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if address&0x7 != 0 { // address alignment
		// TODO TRAP_LOAD_ADDRESS_MISALIGNED
		return false, 0
	}

	if hit, dw := mc.dCache.LoadDoubleWord(address); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, dw := mc.dCache.LoadDoubleWord(address)
		return true, dw
	} else {
		return true, dw
	}
}

// Return the byte stored at
func (mc *MemoryController) StoreByte(address uint32, b uint8) bool {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if !present { // possible page fault
		// TODO TRAP_STORE_OR_AMO_PAGE_FAULT
		return false
	}

	if flags&MEM_F_WRITE == 0 { // permissions
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if hit := mc.dCache.StoreByte(address, b); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreByte(address, b)
	}

	return true
}

func (mc *MemoryController) StoreHalfWord(address uint32, hw uint16) bool {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if !present { // possible page fault
		// TODO TRAP_STORE_OR_AMO_PAGE_FAULT
		return false
	}

	if flags&MEM_F_WRITE == 0 { // permissions
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if address&0x1 != 0 { // address alignment
		// TODO TRAP_STORE_OR_AMO_ADDRESS_MISALIGNED
		return false
	}

	if hit := mc.dCache.StoreHalfWord(address, hw); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreHalfWord(address, hw)
	}

	return false
}

func (mc *MemoryController) StoreWord(address uint32, w uint32) bool {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if !present { // possible page fault
		// TODO TRAP_STORE_OR_AMO_PAGE_FAULT
		return false
	}

	if flags&MEM_F_WRITE == 0 { // permissions
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if address&0x3 != 0 { // address alignment
		// TODO TRAP_STORE_OR_AMO_ADDRESS_MISALIGNED
		return false
	}

	if hit := mc.dCache.StoreWord(address, w); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreWord(address, w)
	}

	return true
}

func (mc *MemoryController) StoreDoubleWord(address uint32, dw uint64) bool {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if !present { // possible page fault
		// TODO TRAP_STORE_OR_AMO_PAGE_FAULT
		return false
	}

	if flags&MEM_F_WRITE == 0 { // permissions
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if address&0x7 != 0 { // address alignment
		// TODO TRAP_STORE_OR_AMO_ADDRESS_MISALIGNED
		return false
	}

	if hit := mc.dCache.StoreDoubleWord(address, dw); !hit {
		mc.misses++
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, CACHE_F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreDoubleWord(address, dw)
	}

	return true
}

// Loads a memory straight from memory, bypassing the cache.
func (mc *MemoryController) UnsafeLoadThroughWord(address uint32) (bool, uint32) {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if !present { // possible page fault
		// TODO TRAP_LOAD_PAGE_FAULT
		return false, 0
	}

	if flags&MEM_F_READ == 0 { // permissions
		// TODO TRAP_LOAD_ACCESS_FAULT
		return false, 0
	}

	if address&0x3 != 0 { // address alignment
		// TODO TRAP_LOAD_ADDRESS_MISALIGNED
		return false, 0
	}

	var value uint32
	if mc.mem.endian == ENDIAN_BIG {
		value = binary.BigEndian.Uint32(mc.mem.data[address : address+4])
	} else {
		value = binary.LittleEndian.Uint32(mc.mem.data[address : address+4])
	}
	// TODO: Verify the integrity of this
	// store loaded value into cache if it's cached
	// should the entire cache line just be invalidated instead perhaps?
	mc.dCache.StoreWordNoDirty(address, value)

	return true, value
}

// Stores a word straight to memory, bypassing cache.
func (mc *MemoryController) UnsafeStoreThroughWord(address uint32, w uint32) bool {
	mc.accesses++
	valid, present, address, flags := mc.mmu.TranslateAndCheck(address)

	if !valid { // address was invalid
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if !present { // possible page fault
		// TODO TRAP_STORE_OR_AMO_PAGE_FAULT
		return false
	}

	if flags&MEM_F_WRITE == 0 { // permissions
		// TODO TRAP_STORE_OR_AMO_ACCESS_FAULT
		return false
	}

	if address&0x3 != 0 { // address alignment
		// TODO TRAP_STORE_OR_AMO_ADDRESS_MISALIGNED
		return false
	}

	var bytes [4]uint8

	if mc.mem.endian == ENDIAN_BIG {
		binary.BigEndian.PutUint32(bytes[:], w)
	} else {
		binary.LittleEndian.PutUint32(bytes[:], w)
	}

	copy(mc.mem.data[address:], bytes[:])

	// also update cache
	mc.dCache.StoreWordNoDirty(address, w) // May be uncached, ignore

	return true
}

// Flushes the data cache to memory
func (mc *MemoryController) FlushCache() {
	mc.mem.Lock()
	mc.dCache.FlushAll(mc.mem.data[:])
	mc.mem.Unlock()
}

// Invalidates the data cache
func (mc *MemoryController) InvalidateCache() {
	mc.dCache.InvalidateAll()
}

// Invalidates the instruction cache
func (mc *MemoryController) InvalidateInstructionCache() {
	mc.iCache.InvalidateAll()
}

// Flush and invalidate the data cache
func (mc *MemoryController) FlushAndInvalidateCache() {
	mc.FlushCache()
	mc.InvalidateCache()
}

func NewMemoryController(m *Memory) MemoryController {
	return MemoryController{
		dCache: NewCache(m.endian),
		iCache: NewCache(m.endian),
		mem:    m,
		mmu:    NewMMU(),
	}
}
