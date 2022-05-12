// this file contains functions to work with interrupts across processors
//   Interrupts are limited in Gotos, but are adequate for the most
// important functionality such as TLB shootdowns.

package cpu

import (
	"sync"
	"sync/atomic"
)

// InterruptMatrix is an (N+1)×(N+1) matrix.
//   The last column is used by the system to interrupt a core
type InterruptMatrix [CoresMax + 1][CoresMax + 1]uint32

// checkInterrupts will scan through interrupts from other cores
func (c *Core) checkInterrupts() bool {
	thisCoreID := c.csr[Csr_MHARTID]
	codes := &c.system.InterruptMatrix()[thisCoreID]
	interrupted := false

	// if ipi is disabled, only check system interrupt
	if !ipiEnable {
		if code := atomic.LoadUint32(&codes[CoresMax]); code != 0 {
			c.interruptedBy = uint32(CoresMax)
			c.interruptCode = code
			c.trap(TrapMachineExternalInterrupt)
			atomic.StoreUint32(&codes[CoresMax], 0)
			interrupted = true
		}
		return interrupted
	}

	for i := 0; i < CoresMax+1; i++ {
		if code := atomic.LoadUint32(&codes[i]); uint32(i) != thisCoreID && code != 0 {
			c.interruptedBy = uint32(i)
			c.interruptCode = code
			c.trap(TrapMachineExternalInterrupt)
			atomic.StoreUint32(&codes[i], 0)
			interrupted = true
		}
	}

	return interrupted
}

// RaiseInterrupt will raise an interrupt on another core
// raise an interrupt on a different core.
//   This will spin and check own interrupts if other core is already
// busy with an interrupt from this core.
//   It is invalid to have more than one interrupt active from
// any core at a time.
//   An interrupt is considered active from when CAS returns true,
// until a response is received, if a response is expected.
func (c *Core) RaiseInterrupt(coreID, code uint32) {
	if !ipiEnable {
		panic("Tried to perform IPI with ipiEnable set to false")
	}

	if coreID >= CoresMax {
		panic("Invalid interrupt target")
	}

	thisCoreID := c.csr[Csr_MHARTID]
	if coreID == thisCoreID {
		panic("A core can't raise an interrupt on itself!")
	}
	ptr := &c.system.InterruptMatrix()[coreID][thisCoreID]
	// use CAS here so that core cannot accidentally raise two interrupts on
	// the same core this ensures that interrupts don't get lost
	for !atomic.CompareAndSwapUint32(ptr, 0, code) {
		c.checkInterrupts()
	}
}

// InterruptInfo returns a tuple containing information about the
// latest received interrupt.
//   `by` is the ID of the core that sent the interrupt
//   `code` is an interrupt code
func (c *Core) InterruptInfo() (by, code uint32) {
	return c.interruptedBy, c.interruptCode
}

// AwaitInterruptResponse waits for a response for the latest
// interrupt raised by this core.
//   It is an error to not await a response when one is expected.
//   It is an error to await a response when one is not expected.
func (c *Core) AwaitInterruptResponse() uint32 {
	if !ipiEnable {
		panic("Tried to wait for IPI response with ipiEnable set to false")
	}

	thisID := c.csr[Csr_MHARTID]
	ptr := &c.system.InterruptMatrix()[thisID][thisID]
	for atomic.LoadUint32(ptr) == 0 {
		c.checkInterrupts()
	}
	code := atomic.LoadUint32(ptr)
	atomic.StoreUint32(ptr, 0)
	return code
}

// RespondInterrupt responds to the latest interrupt.
//   It is an error to not respond when a response is expected.
//   It is an error to respond when a response is not expected.
func (c *Core) RespondInterrupt(code uint32) {
	by, _ := c.InterruptInfo()
	ptr := &c.system.InterruptMatrix()[by][by]
	atomic.StoreUint32(ptr, code)
}

// SafelyAcquire is a helper method to acquire mutexes in a way that
// doesn't break interrupts.
//   It will attempt to acquire the Mutex and if it can't, it will
// spin and check interrupts until it can.
func (c *Core) SafelyAcquire(mu *sync.Mutex) {
	for !mu.TryLock() {
		c.checkInterrupts()
	}
}
