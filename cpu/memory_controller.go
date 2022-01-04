package cpu

type MemoryController struct {
	iCache Cache
	dCache Cache
	mem    *Memory
	mmu    MMU
}

func (mc *MemoryController) LoadInstruction(address uint32) (err error, inst uint32) {
	err, address, flags := mc.mmu.Translate(address)

	// ugly solution for now
	if err != nil {
		return
	}

	// Don't worry about flags for now
	if flags&F_READ == 0 {
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
		mc.iCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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

	// Don't worry about flags for now
	if flags&F_READ == 0 {
	}

	if hit, b := mc.dCache.LoadByte(address); !hit {
		lineNumber := address >> CACHE_LINE_OFFSET_BITS
		mc.mem.Lock()
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
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
		mc.dCache.ReplaceRandom(lineNumber, F_NONE, mc.mem.Data[:])
		mc.mem.Unlock()
		mc.dCache.StoreWord(address, w)
	}

	return nil
}

func NewMemoryController(m *Memory) MemoryController {
	return MemoryController{
		dCache: NewCache(m.endian),
		iCache: NewCache(m.endian),
		mem:    m,
		mmu:    NewMMU(),
	}
}
