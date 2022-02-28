package cpu

const (
	tlbSize = 256
	// valid vpis will never have the leftmost bit set
	tlbInvalidEntry = 0xFFFFFFFFFFFFFFFF
	tlbTryDepth     = 2
)

type tlb struct {
	entries [tlbSize]uint64
}

func (t *tlb) load(vpi uint32) (bool, uint32) {
	for i := uint32(0); i < tlbTryDepth; i++ {
		v := t.entries[(vpi+i*i)&0xFF]
		if uint32(v>>32) == vpi {
			return true, uint32(v)
		}
	}
	return false, 0
}

func (t *tlb) store(vpi, pte uint32) bool {
	for i := uint32(0); i < tlbTryDepth; i++ {
		v := t.entries[(vpi+i*i)&0xFF]
		if v == tlbInvalidEntry {
			t.entries[(vpi+i*i)&0xFF] = (uint64(vpi) << 32) | uint64(pte)
			return true
		}
	}

	// if all are filled, invalidate all entries
	for i := uint32(0); i < tlbTryDepth; i++ {
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