// This file contains implementation of Cache and methods to access it.
// Caching helps reduce lock contention when multiple processors are running.
// This introduces the problem of cache coherence, which is creatively solved in other ways.

package cpu

import (
	"encoding/binary"
)

// Cache constants
const (
	cacheLineLength     = 64   // must be equal to 2^CACHE_LINE_OFFSET_BITS (^ being power, not xor)
	cacheLineOffsetBits = 6    // how many bits of the address are used for the offset within a cache-line
	cacheLineOffsetMask = 0x3F // used to extract the offset in a cache line from an address
	cacheLineCount      = 256  // how many cache lines the cache contains by default
	cacheTryDepth       = 2
	cacheInvalidEntry   = 0xFFFFFFFF
)

// Cache flags
const (
	cacheFlagNone  uint8 = 0x00 // no flags
	cacheFlagAll         = 0xFF // all flags
	cacheFlagDirty       = 0x01 // dirty flag - triggers a cache flush when replaced
	cacheFlagStale       = 0x02 // stale flag - triggers/or should trigger a refresh when a stale line is accessed
)

// A cacheLine contains
type cacheLine struct {
	number uint32
	flags  uint8
	data   [cacheLineLength]uint8
}

// cache is implemented as a hash-map with quadratic probing.
// It will try probing a set number of times before giving up.
type cache struct {
	lines [cacheLineCount]cacheLine
}

func (c *cache) load(address, width uint32) (bool, uint64) {
	// Misaligned load from cache will always fail
	if address&(width-1) != 0 {
		return false, 0
	}

	lineNumber := address >> cacheLineOffsetBits
	offset := address & cacheLineOffsetMask

	for i := uint32(0); i < cacheTryDepth; i++ {
		try := (lineNumber + i*i) & 0xFF
		if c.lines[try].number == lineNumber {
			if c.lines[try].flags&cacheFlagStale != 0 {
				return false, 0
			}

			switch width {
			case 2:
				return true, uint64(binary.LittleEndian.Uint16(c.lines[try].data[offset : offset+2]))
			case 4:
				return true, uint64(binary.LittleEndian.Uint32(c.lines[try].data[offset : offset+4]))
			case 8:
				return true, binary.LittleEndian.Uint64(c.lines[try].data[offset : offset+8])
			default:
				panic("Invalid load width")
			}
		}
	}

	return false, 0
}

func (c *cache) store(address, width uint32, v uint64) bool {
	// misaligned store will always miss
	if address&(width-1) != 0 {
		return false
	}

	lineNumber := address >> cacheLineOffsetBits
	offset := address & cacheLineOffsetMask

	for i := uint32(0); i < cacheTryDepth; i++ {
		try := (lineNumber + i*i) & 0xFF
		if c.lines[try].number == lineNumber {
			if c.lines[try].flags&cacheFlagStale != 0 {
				return false
			}

			var bytes [8]uint8

			switch width {
			case 1:
				bytes[0] = uint8(v)
			case 2:
				binary.LittleEndian.PutUint16(bytes[:], uint16(v))
			case 4:
				binary.LittleEndian.PutUint32(bytes[:], uint32(v))
			case 8:
				binary.LittleEndian.PutUint64(bytes[:], v)
			default:
				panic("Invalid store width")
			}

			copy(c.lines[try].data[offset:], bytes[0:width])
			c.lines[try].flags |= cacheFlagDirty
			return true
		}
	}

	return false
}

// attempts to load a line into cache
// if it cannot find a vacant slot after a given number of probings, it ejects
// all visited slots and loads the line into the first of these slots.
func (c *cache) replace(lineNumber uint32, flags uint8, src []uint8) bool {
	for i := uint32(0); i < cacheTryDepth; i++ {
		try := (lineNumber + i*i) & 0xFF
		if c.lines[try].number == lineNumber {
			if c.lines[try].flags&cacheFlagStale == 0 {
				return false
			}

			address := lineNumber << cacheLineOffsetBits
			copy(c.lines[try].data[:], src[address:address+cacheLineLength])
			c.lines[try].flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits
			return true
		}
	}

	// line was not present, try to find a place for it
	for i := uint32(0); i < cacheTryDepth; i++ {
		try := (lineNumber + i*i) & 0xFF
		if c.lines[try].number == cacheInvalidEntry {
			address := lineNumber << cacheLineOffsetBits
			c.lines[try].number = lineNumber
			copy(c.lines[try].data[:], src[address:address+cacheLineLength])
			c.lines[try].flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits
			return true
		}
	}

	// no vacant slots, invalidate all slots
	for i := uint32(0); i < cacheTryDepth; i++ {
		try := (lineNumber + i*i) & 0xFF
		if c.lines[try].flags&cacheFlagDirty == 1 {
			address := c.lines[try].number << cacheLineOffsetBits
			copy(src[address:], c.lines[try].data[:])
		}
		c.lines[try].number = cacheInvalidEntry
		c.lines[try].flags = 0
	}

	// finally, a vacant space is guaranteed
	try := lineNumber & 0xFF
	address := lineNumber << cacheLineOffsetBits
	c.lines[try].number = lineNumber
	copy(c.lines[try].data[:], src[address:address+cacheLineLength])
	c.lines[try].flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits
	return true
}

// Writes all cache-lines back to `dst`.
// Helpful if you want all other cores to see changes made by this core.
func (c *cache) writebackAll(dst []uint8) int {
	written := 0
	for i := range c.lines {
		line := &c.lines[i]

		// check if it's dirty, if so, flush
		if line.flags&cacheFlagDirty != 0 {
			address := line.number << cacheLineOffsetBits
			copy(dst[address:], line.data[:])
			written += 1
			// clear the dirty bit
			line.flags &= (cacheFlagDirty ^ cacheFlagAll)
		}
	}
	return written
}

// Marks all cache-lines as stale.
// Helpful if you want to guarantee that you can see updates made from other cores.
func (c *cache) invalidateAll() {
	for i := range c.lines {
		line := &c.lines[i]
		line.flags |= cacheFlagStale
	}
}

func newCache() cache {
	c := cache{}
	for i := range c.lines {
		c.lines[i].number = cacheInvalidEntry
		c.lines[i].flags = 0
	}
	return c
}
