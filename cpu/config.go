package cpu

// enable the F extension
const xFEnable = false

// enable the D extension
const xDEnable = false

// enable the Zicsr extension
const xZicsrEnable = false

// enable interprocessor interrupts for normal operation; if set to false,
// interrupts are only received from the system and can not be raised on, or
// received from other cores
const ipiEnable = false

// CoresMax is the maximum number of cores we support
// This just simplifies the interrupt implementation
const CoresMax = 4

// Cache constants, cache size, line count, etc.
const (
	cacheEnable         = true
	cacheLineLength     = 64   // must be equal to exp(2, CACHE_LINE_OFFSET_BITS)
	cacheLineOffsetBits = 6    // how many bits of the address are used for the offset
	cacheLineOffsetMask = 0x3F // used to extract the offset in a cache line from an address
	cacheLineCount      = 256  // how many cache lines the cache contains by default
	cacheProbeDepth     = 2
)

const (
	tlbSize       = 256
	tlbProbeDepth = 3
)
