package cpu

import "sync"

type reservation struct {
	valid bool
	value uint32
}

type ReservationSets struct {
	sync.Mutex
	reservations []reservation
}

// Invalidate the reservation accross all cores
func (rs *ReservationSets) unsafeCompareAndInvalidateAllReservations(addr uint32) {
	for i := range rs.reservations {
		if rs.reservations[i].value == addr && rs.reservations[i].valid {
			rs.reservations[i].valid = false
		}
	}
}

// Register a reservation in a given set
func (rs *ReservationSets) unsafeRegisterReservation(set int, addr uint32) {
	rs.reservations[set].value = addr
	rs.reservations[set].valid = true
}

// Check and invalidate a reservation
func (rs *ReservationSets) unsafeCompareAndInvalidateReservation(set int, addr uint32) bool {
	if rs.reservations[set].value == addr && rs.reservations[set].valid {
		rs.reservations[set].valid = false
		return true
	}
	return false
}

func NewReservationSets(n int) (rs ReservationSets) {
	rs = ReservationSets{
		reservations: make([]reservation, n),
	}

	return
}
