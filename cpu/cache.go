// This file contains implementation of Cache and methods to access it.
// Caching helps reduce lock contention when multiple processors are running.
// This introduces the problem of cache coherence, which is creatively solved in other ways.

// WARNING: man physically addressed caches really suck when they may be written back
// May have to store virtual line numbers as well and translate to physical address before writing back
// This may cause several cache misses.

package cpu

import (
	"encoding/binary"
	"math/rand"
)

// Cache constants
const (
	cacheLineLength     = 64   // must be equal to 2^CACHE_LINE_OFFSET_BITS (^ being power, not xor)
	cacheLineOffsetBits = 6    // how many bits of the address are used for the offset within a cache-line
	cacheLineOffsetMask = 0x3F // used to extract the offset in a cache line from an address
	cacheLineCount      = 64   // how many cache lines the cache contains by default
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

type cache struct {
	lines  [cacheLineCount]cacheLine
	lookup map[uint32]*cacheLine
	size   int
	endian Endian
}

func (c *cache) load(address, width uint32) (bool, uint64) {
	lineNumber := address >> cacheLineOffsetBits

	// Misaligned load from cache will always fail
	if address&(width-1) != 0 {
		return false, 0
	}

	if line, present := c.lookup[lineNumber]; present {
		if line.flags&cacheFlagStale != 0 {
			return false, 0
		}

		offset := address & cacheLineOffsetMask

		if width == 1 {
			return true, uint64(line.data[offset])
		}

		if c.endian == EndianBig {
			switch width {
			case 2:
				return true, uint64(binary.BigEndian.Uint16(line.data[offset : offset+2]))
			case 4:
				return true, uint64(binary.BigEndian.Uint32(line.data[offset : offset+4]))
			case 8:
				return true, binary.BigEndian.Uint64(line.data[offset : offset+8])
			default:
				panic("Invalid load width")
			}
		} else {
			switch width {
			case 2:
				return true, uint64(binary.LittleEndian.Uint16(line.data[offset : offset+2]))
			case 4:
				return true, uint64(binary.LittleEndian.Uint32(line.data[offset : offset+4]))
			case 8:
				return true, binary.LittleEndian.Uint64(line.data[offset : offset+8])
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
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&cacheFlagStale != 0 {
			return false
		}

		offset := address & cacheLineOffsetMask

		if width == 1 {
			line.data[offset] = uint8(v)
		}

		var bytes [8]uint8
		if c.endian == EndianBig {
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
		copy(line.data[offset:], bytes[0:width])
		line.flags |= cacheFlagDirty
		return true
	}
	return false
}

// Attempts to store a word (uint32) to cache, but does not set the dirty bit if
// the store is completed.
// If the cache line containing the word is present and not marked stale, stores
// the word and returns true, false otherwise.
func (c *cache) storeWordNoDirty(address uint32, w uint32) bool {
	lineNumber := address >> cacheLineOffsetBits
	if line, present := c.lookup[lineNumber]; present {
		if line.flags&cacheFlagStale != 0 {
			return false
		}

		offset := address & cacheLineOffsetMask
		bytes := make([]uint8, 4)
		if c.endian == EndianBig {
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
// If this line is dirty, the line is first written back before
// it is replaced with the new data.
//
// Returns true if a replacement/refresh was performed, false otherwise.
func (c *cache) replaceRandom(lineNumber uint32, flags uint8, src []uint8) bool {
	// if line is already present
	if _, present := c.lookup[lineNumber]; present {
		line := c.lookup[lineNumber]

		if line.flags&cacheFlagStale == 0 { // if line isn't stale
			return false // line is not updated
		}

		address := lineNumber << cacheLineOffsetBits
		copy(line.data[:], src[address:address+cacheLineLength])
		line.flags &= (cacheFlagDirty ^ cacheFlagStale ^ cacheFlagAll) // remove stale and dirty bits

		return true
	}

	// if cache isn't full, just take the next open space
	// pick a random line to eject from cache
	target := rand.Intn(cacheLineCount)
	if c.size < cacheLineCount {
		target = c.size
		defer func() { c.size += 1 }()
	}

	eject := &c.lines[target]
	address := eject.number << cacheLineOffsetBits
	delete(c.lookup, eject.number) // delete the old entry, no-op if there is no old entry

	// check if it's dirty, if so, writeback
	if eject.flags&cacheFlagDirty != 0 {
		copy(src[address:], eject.data[:])
	}

	// writes the line to cache
	address = lineNumber << cacheLineOffsetBits
	copy(eject.data[:], src[address:address+cacheLineLength])
	eject.flags = flags
	eject.number = lineNumber

	// update the lookup
	c.lookup[lineNumber] = eject // add new entry

	return true
}

// writes a cache-line if present
func (c *cache) writebackLine(lineNumber uint32, dst []uint8) {
	if _, present := c.lookup[lineNumber]; present {
		line := c.lookup[lineNumber]

		if line.flags&cacheFlagDirty != 0 {
			address := line.number << cacheLineOffsetBits
			copy(dst[address:], line.data[:])
			// clear the dirty bit
			line.flags &= (cacheFlagDirty ^ cacheFlagAll)
		}
	}
}

// Writes back all cache-lines back to `dst`.
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

// Invalidates a cache-line if present
func (c *cache) invalidateLine(lineNumber uint32) {
	if _, present := c.lookup[lineNumber]; present {
		line := c.lookup[lineNumber]

		line.flags |= cacheFlagStale // set stale bit, forces refresh on next access
	}
}

// Marks all cache-lines as stale.
// Helpful if you want to guarantee that you can see updates made from other cores.
func (c *cache) invalidateAll() {
	for i := range c.lines {
		line := &c.lines[i]
		line.flags |= cacheFlagStale
	}
}

func newCache(endian Endian) cache {
	return cache{
		lookup: make(map[uint32]*cacheLine, cacheLineCount),
		endian: endian,
	}
}
