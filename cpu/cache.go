// This file contains implementation of Cache and methods to access it.
// Caching helps reduce lock contention when multiple processors are running.
// This introduces the problem of cache coherence, which is creatively solved in other ways.

package cpu

import (
	"encoding/binary"
)

// Cache flags, cache lines can be dirty or stale
const (
	cacheInvalidEntry       = 0xFFFFFFFF
	cacheFlagNone     uint8 = 0x00 // no flags
	cacheFlagAll            = 0xFF // all flags
	cacheFlagDirty          = 0x01 // dirty flag - triggers a cache flush when replaced
	cacheFlagStale          = 0x02 // stale flag - triggers/or should trigger a refresh when a stale line is accessed
)

// cacheLine consists of a number, some flags, and the data
type cacheLine struct {
	number uint32
	flags  uint8
	data   [cacheLineLength]uint8
}

// cache is implemented as a hash-map with quadratic probing.
//   It will try probing a set number of times before giving up.
//   When it gives up, the candidate slots should be flushed/set stale
// and the new line should be brought in from main memory at the first
// position.
type cache struct {
	lines [cacheLineCount]cacheLine
}

// load will load a value from cache with a given width.
//   It will fail/miss if the correct line is not in cache.
//   Panics if address is not aligned to `width` bytes
func (c *cache) load(address, width uint32) (bool, uint64) {
	// Misaligned load from cache will always fail
	if address&(width-1) != 0 {
		panic("misaligned access to cache")
	}

	lineNumber := address >> cacheLineOffsetBits
	offset := address & cacheLineOffsetMask

	for i := uint32(0); i < cacheProbeDepth; i++ {
		try := (lineNumber + i*i) % cacheLineCount
		if c.lines[try].number == lineNumber {
			if c.lines[try].flags&cacheFlagStale != 0 {
				return false, 0
			}

			switch width {
			case 1:
				return true, uint64(c.lines[try].data[offset])
			case 2:
				return true, uint64(binary.LittleEndian.Uint16(c.lines[try].data[offset : offset+2]))
			case 4:
				return true, uint64(
					binary.LittleEndian.Uint32(c.lines[try].data[offset : offset+4]))
			case 8:
				return true, binary.LittleEndian.Uint64(c.lines[try].data[offset : offset+8])
			default:
				panic("Invalid load width")
			}
		}
	}

	return false, 0
}

// store will store a value into cache with a given width.
//   It will fail/miss if the correct line is not in cache.
//   Panics if address is not aligned to `width` bytes
func (c *cache) store(address, width uint32, v uint64) bool {
	// misaligned store will always miss
	if address&(width-1) != 0 {
		panic("misaligned access to cache")
	}

	lineNumber := address >> cacheLineOffsetBits
	offset := address & cacheLineOffsetMask

	for i := uint32(0); i < cacheProbeDepth; i++ {
		try := (lineNumber + i*i) % cacheLineCount
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

// replace brings a line into cache at one of the candidate slots.
//   Candidate slots are all slots that might hold the given line; e.g if the
// probe depth is 2, there are 2 slots that any cache line might fit into.
//   If all candidate slots are filled, they will be flushed and the line will
// be placed in the first slot.
//   If the line is already in cache, it will be refreshed by this.
func (c *cache) replace(lineNumber uint32, flags uint8, src []uint8) bool {
	for i := uint32(0); i < cacheProbeDepth; i++ {
		try := (lineNumber + i*i) % cacheLineCount
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
	for i := uint32(0); i < cacheProbeDepth; i++ {
		try := (lineNumber + i*i) % cacheLineCount
		if c.lines[try].number == cacheInvalidEntry {
			address := lineNumber << cacheLineOffsetBits
			c.lines[try].number = lineNumber
			copy(c.lines[try].data[:], src[address:address+cacheLineLength])
			c.lines[try].flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits
			return true
		}
	}

	// no vacant slots, invalidate all candidate slots
	for i := uint32(0); i < cacheProbeDepth; i++ {
		try := (lineNumber + i*i) % cacheLineCount
		if c.lines[try].flags&cacheFlagDirty == 1 {
			address := c.lines[try].number << cacheLineOffsetBits
			copy(src[address:], c.lines[try].data[:])
		}
		c.lines[try].number = cacheInvalidEntry
		c.lines[try].flags = 0
	}

	// finally, a vacant space is guaranteed
	try := lineNumber % cacheLineCount
	address := lineNumber << cacheLineOffsetBits
	c.lines[try].number = lineNumber
	copy(c.lines[try].data[:], src[address:address+cacheLineLength])
	c.lines[try].flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits
	return true
}

// writebackAll writes all cache-lines back to `dst`.
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

// invalidateAll marks all cache-lines as stale.
// Helpful if you want to guarantee that you can see updates made from other cores.
func (c *cache) invalidateAll() {
	for i := range c.lines {
		line := &c.lines[i]
		line.flags |= cacheFlagStale
	}
}

// newCache returns a new cache filled with invalid (open) entries.
func newCache() cache {
	c := cache{}
	for i := range c.lines {
		c.lines[i].number = cacheInvalidEntry
		c.lines[i].flags = 0
	}
	return c
}
