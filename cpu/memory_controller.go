package cpu

import "encoding/binary"

type MemoryController struct {
	iCache Cache
	dCache Cache
	mem    *Memory
	mmu    MMU
}

func (mc *MemoryController) LoadInstruction(address uint32) (err error, inst uint32) {
	err, address, flags := mc.mmu.Translate(address)

	// TODO check errors and shit

	// ugly solution for now
	if err != nil {
		return
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
		// TODO check flags
	}

	// check if address is misaligned
	// this shuold _never_ happen
	if address&0x3 != 0 {
		panic("Instruction address misaligned!")
	}

	if hit, instruction := mc.iCache.LoadWord(address); !hit {
		// fmt.Println("Cache miss!")
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		// fmt.Printf("Loading line #%04X\n", lineNumber)
		mc.mem.Lock()
		mc.iCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		// fmt.Printf("Loaded line #%04X\n", lineNumber)
		_, instruction := mc.iCache.LoadWord(address)
		inst = instruction
	} else {
		inst = instruction
	}

	// fmt.Printf("inst: %08X\n", inst)
	return
}

// Return the byte stored at
func (mc *MemoryController) LoadByte(address uint32) (error, uint8) {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err, 0
	}

	if flags&F_READ == 0 {
	}

	if flags&F_NOCACHE != 0 {
	}

	if hit, b := mc.dCache.LoadByte(address); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, b := mc.dCache.LoadByte(address)
		return nil, b
	} else {
		return nil, b
	}
}

func (mc *MemoryController) LoadHalfWord(address uint32) (error, uint16) {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err, 0
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit, hw := mc.dCache.LoadHalfWord(address); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, hw := mc.dCache.LoadHalfWord(address)
		return nil, hw
	} else {
		return nil, hw
	}
}

func (mc *MemoryController) LoadWord(address uint32) (error, uint32) {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err, 0
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit, w := mc.dCache.LoadWord(address); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		_, w := mc.dCache.LoadWord(address)
		return nil, w
	} else {
		return nil, w
	}
}

// Return the byte stored at
func (mc *MemoryController) StoreByte(address uint32, b uint8) error {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit := mc.dCache.StoreByte(address, b); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreByte(address, b)
	}

	return nil
}

func (mc *MemoryController) StoreHalfWord(address uint32, hw uint16) error {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit := mc.dCache.StoreHalfWord(address, hw); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreHalfWord(address, hw)
	}

	return nil
}

func (mc *MemoryController) StoreWord(address uint32, w uint32) error {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit := mc.dCache.StoreWord(address, w); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.data[:])
		mc.mem.Unlock()
		mc.dCache.StoreWord(address, w)
	}

	return nil
}

func (mc *MemoryController) UnsafeLoadThroughWord(address uint32) (error, uint32) {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return err, 0
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	var value uint32
	if mc.mem.endian == BIG {
		value = binary.BigEndian.Uint32(mc.mem.data[address : address+4])
	} else {
		value = binary.LittleEndian.Uint32(mc.mem.data[address : address+4])
	}
	// TODO: Verify the integrity of this
	// store loaded value into cache if it's cached
	// should the entire cache line just be invalidated instead perhaps?
	mc.dCache.StoreWordNoDirty(address, value)

	return nil, value
}

func (mc *MemoryController) UnsafeStoreThroughWord(address uint32, w uint32) error {
	err, address, flags := mc.mmu.Translate(address)

	if err != nil {
		return err
	}

	if flags&F_WRITE == 0 {
	}

	var bytes [4]uint8

	if mc.mem.endian == BIG {
		binary.BigEndian.PutUint32(bytes[:], w)
	} else {
		binary.LittleEndian.PutUint32(bytes[:], w)
	}

	copy(mc.mem.data[address:], bytes[:])

	// also update cache
	mc.dCache.StoreWordNoDirty(address, w) // May be uncached, ignore

	return nil
}

func (mc *MemoryController) FlushCache() {
	// Flush data cache
	mc.mem.Lock()
	mc.dCache.FlushAll(mc.mem.data[:])
	mc.mem.Unlock()
}

func (mc *MemoryController) InvalidateCache() {
	// TODO
	mc.dCache.InvalidateAll()
}

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
