// This file contains implementation of Cache and methods to access it.
// Caching helps reduce lock contention when multiple processors are running.
// This introduces the problem of cache coherence, which is creatively solved in other ways.

package cpu

import (
	"encoding/binary"
	"math/rand"
)

// TODO check alignments

// Cache constants
const (
	CACHE_LINE_LENGTH      = 64   // must be equal to 2^CACHE_LINE_OFFSET_BITS (^ being power, not xor)
	CACHE_LINE_OFFSET_BITS = 6    // how many bits of the address are used for the offset within a cache-line
	CACHE_LINE_OFFSET_MASK = 0x3F // used to extract the offset in a cache line from an address
	CACHE_LINE_COUNT       = 64   // how many cache lines the cache contains by default
)

// Cache flags
const (
	CACHE_F_NONE  uint8 = 0x00 // no flags
	CACHE_F_ALL         = 0xFF // all flags
	CACHE_F_DIRTY       = 0x01 // dirty flag - triggers a cache flush when replaced
	CACHE_F_STALE       = 0x02 // stale flag - triggers/or should trigger a refresh when a stale line is accessed
)

// A CacheLine contains
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

// Attempts to load a byte (uint8) from cache.
// If the cache line containing the byte is present and not marked stale, returns
// (true, byte), oherwise (false, 0).
func (c *Cache) LoadByte(address uint32) (present bool, byte uint8) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false, 0
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		return true, line.data[offset]
	}
	return false, 0
}

// Attempts to load a halfword (uint16) from cache.
// If the cache line containing the halfword is present and not marked stale, returns
// (true, halfword), oherwise (false, 0).
func (c *Cache) LoadHalfWord(address uint32) (bool, uint16) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false, 0
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		if c.endian == ENDIAN_BIG {
			return true, binary.BigEndian.Uint16(line.data[offset : offset+2])
		} else {
			return true, binary.LittleEndian.Uint16(line.data[offset : offset+2])
		}
	}
	return false, 0
}

// Attempts to load a word (uint32) from cache.
// If the cache line containing the word is present and not marked stale, returns
// (true, word), oherwise (false, 0).
func (c *Cache) LoadWord(address uint32) (bool, uint32) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false, 0
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		if c.endian == ENDIAN_BIG {
			return true, binary.BigEndian.Uint32(line.data[offset : offset+4])
		} else {
			return true, binary.LittleEndian.Uint32(line.data[offset : offset+4])
		}
	}
	return false, 0
}

// Attempts to load a doubleword (uint64) from cache.
// If the cache line containing the doubleword is present and not marked stale, returns
// (true, doubleword), oherwise (false, 0).
func (c *Cache) LoadDoubleWord(address uint32) (bool, uint64) {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false, 0
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		if c.endian == ENDIAN_BIG {
			return true, binary.BigEndian.Uint64(line.data[offset : offset+8])
		} else {
			return true, binary.LittleEndian.Uint64(line.data[offset : offset+8])
		}
	}
	return false, 0
}

// Attempts to store a byte (uint8) to cache.
// If the cache line containing the byte is present and not marked stale, stores
// the byte and returns true, false otherwise.
func (c *Cache) StoreByte(address uint32, b uint8) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		line.data[offset] = b
		line.flags |= CACHE_F_DIRTY
		return true
	}
	return false
}

// Attempts to store a halfword (uint16) to cache.
// If the cache line containing the halfword is present and not marked stale, stores
// the halfword and returns true, false otherwise.
func (c *Cache) StoreHalfWord(address uint32, hw uint16) bool {
	// TODO
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 2)
		if c.endian == ENDIAN_BIG {
			binary.BigEndian.PutUint16(bytes, hw)
		} else {
			binary.LittleEndian.PutUint16(bytes, hw)
		}
		copy(line.data[offset:], bytes)
		line.flags |= CACHE_F_DIRTY
		return true
	}
	return false
}

// Attempts to store a word (uint32) to cache.
// If the cache line containing the word is present and not marked stale, stores
// the word and returns true, false otherwise.
func (c *Cache) StoreWord(address uint32, w uint32) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 4)
		if c.endian == ENDIAN_BIG {
			binary.BigEndian.PutUint32(bytes, w)
		} else {
			binary.LittleEndian.PutUint32(bytes, w)
		}
		copy(line.data[offset:], bytes)
		line.flags |= CACHE_F_DIRTY
		return true
	}
	return false
}

// Attempts to store a doubleword (uint32) to cache.
// If the cache line containing the doubleword is present and not marked stale, stores
// the doubleword and returns true, false otherwise.
func (c *Cache) StoreDoubleWord(address uint32, dw uint64) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 4)
		if c.endian == ENDIAN_BIG {
			binary.BigEndian.PutUint64(bytes, dw)
		} else {
			binary.LittleEndian.PutUint64(bytes, dw)
		}
		copy(line.data[offset:], bytes)
		line.flags |= CACHE_F_DIRTY
		return true
	}
	return false
}

// Attempts to store a word (uint32) to cache, but does not set the dirty bit if
// the store is completed.
// If the cache line containing the word is present and not marked stale, stores
// the word and returns true, false otherwise.
func (c *Cache) StoreWordNoDirty(address uint32, w uint32) bool {
	lineNumber := address >> CACHE_LINE_OFFSET_BITS
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&CACHE_F_STALE != 0 {
			return false
		}

		offset := address & CACHE_LINE_OFFSET_MASK
		bytes := make([]uint8, 4)
		if c.endian == ENDIAN_BIG {
			binary.BigEndian.PutUint32(bytes, w)
		} else {
			binary.LittleEndian.PutUint32(bytes, w)
		}
		copy(line.data[offset:], bytes)
		return true
	}
	return false
}

// Loads a cache-line into cache.
// Behaviour depends on the state of the cache.
// If the cache is not full: selects the next open slot in the line-storage and
// loads the cache-line into it.
// This is done by copying the line data from `src` (usually mc.mem)
//
// If the cache-line is already present, but it is stale: refreshes that line.
//
// If the cache is full and the cache-line is not present, selects a random line
// to be evicted.
// If this line is dirty, the line is first flushed (written back to `src`) before
// it is replaced with the new data.
//
// Returns true if a replacement/refresh was performed, false otherwise.
func (c *Cache) ReplaceRandom(lineNumber uint32, flags uint8, src []uint8) bool {
	// if line is already present
	if _, present := c.lookup[lineNumber]; present {
		line := c.lookup[lineNumber]

		if line.flags&CACHE_F_STALE == 0 { // if line isn't stale
			return false // line is not updated
		}

		address := lineNumber << CACHE_LINE_OFFSET_BITS
		copy(line.data[:], src[address:address+CACHE_LINE_LENGTH])
		line.flags &= (CACHE_F_DIRTY ^ CACHE_F_STALE ^ CACHE_F_ALL) // remove stale and dirty bits

		return true
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
	if eject.flags&CACHE_F_DIRTY != 0 {
		copy(src[address:], eject.data[:])
	}

	// writes the line to cache
	address = lineNumber << CACHE_LINE_OFFSET_BITS
	copy(eject.data[:], src[address:address+CACHE_LINE_LENGTH])
	eject.flags = flags
	eject.number = lineNumber

	// update the lookup
	c.lookup[lineNumber] = eject // add new entry

	return true
}

// Flushes all cache-lines back to `src`.
// Helpful if you want all other cores to see changes made by this core.
func (c *Cache) FlushAll(src []uint8) int {
	flushed := 0
	for i := range c.lines {
		line := &c.lines[i]

		// check if it's dirty, if so, flush
		// if the cache isn't full, this shouldn'd be dirty
		if line.flags&CACHE_F_DIRTY != 0 {
			address := line.number << CACHE_LINE_OFFSET_BITS
			copy(src[address:], line.data[:])
			flushed += 1
			// clear the dirty bit
			line.flags &= (CACHE_F_DIRTY ^ CACHE_F_ALL)
		}
	}
	return flushed
}

// Marks all cache-lines as stale.
// Helpful if you want to guarantee that you can see updates made from other cores.
func (c *Cache) InvalidateAll() {
	for i := range c.lines {
		line := &c.lines[i]
		line.flags |= CACHE_F_STALE
	}
}

func NewCache(endian Endian) Cache {
	return Cache{
		lookup: make(map[uint32]*CacheLine, CACHE_LINE_COUNT),
		endian: endian,
	}
}
