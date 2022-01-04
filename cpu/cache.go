package cpu

import (
	"encoding/binary"
	"math/rand"
)

const (
	CACHE_LINE_LENGTH      = 128
	CACHE_LINE_OFFSET_BITS = 7
	CACHE_LINE_OFFSET_MASK = 0x7F
	CACHE_LINE_COUNT       = 64
)

const (
	F_NONE  uint8 = 0x00
	F_DIRTY       = 0x01
)

type CacheLine struct {
	number uint32
	flags  uint8
	data   [CACHE_LINE_LENGTH]uint8
}

type Endian uint8

type Cache struct {
	lines  [CACHE_LINE_COUNT]CacheLine
	lookup map[uint32]*CacheLine
	size   int
	endian Endian
}

func (c *Cache) LoadByte(address uint32) (present bool, byte uint8) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		return true, line.data[offset]
	}
	return false, 0
}

func (c *Cache) LoadHalfWord(address uint32) (bool, uint16) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		if c.endian == BIG {
			return true, binary.BigEndian.Uint16(line.data[offset : offset+2])
		} else {
			return true, binary.LittleEndian.Uint16(line.data[offset : offset+2])
		}
	}
	return false, 0
}

func (c *Cache) LoadWord(address uint32) (bool, uint32) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		if c.endian == BIG {
			return true, binary.BigEndian.Uint32(line.data[offset : offset+4])
		} else {
			return true, binary.LittleEndian.Uint32(line.data[offset : offset+4])
		}
	}
	return false, 0
}

func (c *Cache) StoreByte(address uint32, b uint8) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		line.data[offset] = b
		line.flags |= F_DIRTY
		return true
	}
	return false
}

func (c *Cache) StoreHalfWord(address uint32, hw uint16) bool {
	// TODO
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 2)
		if c.endian == BIG {
			binary.BigEndian.PutUint16(bytes, hw)
		} else {
			binary.LittleEndian.PutUint16(bytes, hw)
		}
		copy(line.data[offset:offset+2], bytes)
		line.flags |= F_DIRTY
		return true
	}
	return false
}

func (c *Cache) StoreWord(address uint32, w uint32) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 4)
		if c.endian == BIG {
			binary.BigEndian.PutUint32(bytes, w)
		} else {
			binary.LittleEndian.PutUint32(bytes, w)
		}
		copy(line.data[offset:offset+4], bytes)
		line.flags |= F_DIRTY
		return true
	}
	return false
}

func (c *Cache) ReplaceRandom(lineNumber uint32, flags uint8, src []uint8) bool {
	if _, present := c.lookup[lineNumber]; present {
		return false
	}

	// if cache isn't full, just take the next open space
	// pick a random line to eject from cache
	target := rand.Intn(CACHE_LINE_COUNT)
	if c.size < CACHE_LINE_COUNT {
		target = c.size
		defer func() { c.size += 1 }()
	}

	eject := &c.lines[target]
	address := eject.number << CACHE_LINE_OFFSET_BITS
	delete(c.lookup, eject.number) // delete the old entry, no-op if there is no old entry

	// check if it's dirty, if so, flush
	// if the cache isn't full, this shouldn'd be dirty
	if eject.flags&F_DIRTY != 0 {
		copy(src[address:], eject.data[:])
	}

	// writes the line to cache
	address = lineNumber << CACHE_LINE_OFFSET_BITS
	copy(eject.data[:], src[address:address+CACHE_LINE_LENGTH])
	// fmt.Println("FRESH CACHE LINE")
	// for i := 0; i < 64; i += 4 {
	// 	fmt.Printf("%08X\n", binary.LittleEndian.Uint32(eject.data[i:i+4]))
	// }
	eject.flags = flags
	eject.number = lineNumber

	// update the lookup
	c.lookup[lineNumber] = eject // add new entry

	return true
}

func (c *Cache) FlushAll(src []uint8) int {
	flushed := 0
	for i := range c.lines {
		line := &c.lines[i]

		// check if it's dirty, if so, flush
		// if the cache isn't full, this shouldn'd be dirty
		if line.flags&F_DIRTY != 0 {
			address := line.number << CACHE_LINE_OFFSET_BITS
			copy(src[address:], line.data[:])
			flushed += 1
		}
	}
	return flushed
}

func NewCache(endian Endian) Cache {
	return Cache{
		lookup: make(map[uint32]*CacheLine, CACHE_LINE_COUNT),
		endian: endian,
	}
}