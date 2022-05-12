// this file contains structs and methods related to reservation sets which are
// used to achieve atomics in RISC-V.

package cpu

import "sync"

// reservation is a small struct that contains a boolean and a value.
//   `valid` marks whether this reservation is valid or not.
//   `value` is the actual value of the reservation and subsumes the bytes of
// the reservation.
type reservation struct {
	valid bool
	value uint32
}

// ReservationSets holds active reservations for all cores and whether they are
// valid or not.
type ReservationSets struct {
	sync.Mutex
	reservations [CoresMax]reservation
}

// unsafeInvalidate invalidates a matching reservation on all cores
func (rs *ReservationSets) unsafeInvalidate(addr uint32) {
	for i := range rs.reservations {
		if rs.reservations[i].value == addr && rs.reservations[i].valid {
			rs.reservations[i].valid = false
		}
	}
}

// unsafeRegister will register a reservation in a given set.
//   The mutex should be held when this is called.
func (rs *ReservationSets) unsafeRegister(set int, addr uint32) {
	rs.reservations[set].value = addr
	rs.reservations[set].valid = true
}

// unsafeInvalidateSingle will check and invalidate a reservation for a set.
//   If the reservation matches `addr` this function returns true,
// otherwise it returns false.
//   In either case, the reservation is invalidated.
func (rs *ReservationSets) unsafeInvalidateSingle(set int, addr uint32) bool {
	if rs.reservations[set].value == addr && rs.reservations[set].valid {
		rs.reservations[set].valid = false
		return true
	}
	rs.reservations[set].valid = false
	return false
}

// NewReservationSets creates a new instance of `ReservationSets`
func NewReservationSets() ReservationSets {
	return ReservationSets{}
}
