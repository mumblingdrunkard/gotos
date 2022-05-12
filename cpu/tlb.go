// This file contains the `tlb` struct that is used to cache
// virtual memory translations.

package cpu

const (
	// valid vpis will never have the leftmost bit set
	// 9-bit asid, 22-bit vpn, 32-bit pte = 63 bits
	tlbInvalidEntry = 0xFFFFFFFFFFFFFFFF
)

type tlb struct {
	entries [tlbSize]uint64
}

// load will attempt to locate a given virtual page index
func (t *tlb) load(vpi uint32) (bool, uint32) {
	for i := uint32(0); i < tlbProbeDepth; i++ {
		v := t.entries[(vpi+i*i)&0xFF]
		if uint32(v>>32) == vpi {
			return true, uint32(v)
		}
	}
	return false, 0
}

func (t *tlb) store(vpi, pte uint32) bool {
	for i := uint32(0); i < tlbProbeDepth; i++ {
		v := t.entries[(vpi+i*i)&0xFF]
		if v == tlbInvalidEntry {
			t.entries[(vpi+i*i)&0xFF] = (uint64(vpi) << 32) | uint64(pte)
			return true
		}
	}

	// if all are filled, invalidate all entries
	for i := uint32(0); i < tlbProbeDepth; i++ {
		t.entries[(vpi+i*i)&0xFF] = tlbInvalidEntry
	}

	t.entries[vpi&0xFF] = (uint64(vpi) << 32) | uint64(pte)
	return true
}

func (t *tlb) invalidateAll() {
	for i := range t.entries {
		t.entries[i] = tlbInvalidEntry
	}
}

func newTLB() tlb {
	t := tlb{}
	for i := range t.entries {
		t.entries[i] = tlbInvalidEntry
	}
	return t
}
