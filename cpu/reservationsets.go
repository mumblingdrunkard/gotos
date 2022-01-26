package cpu

import "sync"

type ReservationSets struct {
	sync.Mutex
	sets   []map[uint32]bool
	lookup []*map[uint32]bool
}

func (rs *ReservationSets) unsafeInvalidateAll(addr uint32) {
	for i := range rs.sets {
		delete(*rs.lookup[i], addr)
	}
}

func (rs *ReservationSets) unsafeInvalidate(set int, addr uint32) {
	delete(*rs.lookup[set], addr)
}

func NewReservationSets(n int) (rs ReservationSets) {
	rs = ReservationSets{
		sets:   make([]map[uint32]bool, n),
		lookup: make([]*map[uint32]bool, n),
	}

	for i := 0; i < n; i++ {
		rs.sets[i] = make(map[uint32]bool)
		rs.lookup[i] = &rs.sets[i]
	}

	return
}
