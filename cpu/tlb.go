package cpu

const (
	tlbSize = 16
)

type tlb struct {
	lookup map[uint32]uint32
	size   int
}

func (t *tlb) load(vpn uint32) (bool, uint32) {
	if pte, present := t.lookup[vpn]; present {
		return true, pte
	}

	return false, 0
}

func (t *tlb) store(vpn, pte uint32) bool {
	if _, present := t.lookup[vpn]; present {
		return false
	}
	t.lookup[vpn] = pte
	return false
}

func (t *tlb) replaceRandom(vpn, pte uint32) bool {
	if _, present := t.lookup[vpn]; present {
		return false
	}

	if t.size < tlbSize {
		t.lookup[vpn] = pte
		t.size += 1
		return false
	}

	for k := range t.lookup {
		delete(t.lookup, k)
		t.lookup[vpn] = pte
		return true
	}

	return false
}

func (t *tlb) flushAll() {
	t.lookup = make(map[uint32]uint32)
	t.size = 0
}

func newTLB() tlb {
	return tlb{
		lookup: make(map[uint32]uint32, tlbSize),
	}
}
