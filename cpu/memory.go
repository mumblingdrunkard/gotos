package cpu

import (
	"fmt"
	"sync"
)

const (
	// 4 MiB of memory ought to be enough
	MemorySize = 1024 * 1024 * 4
)

// Memory is a structure that contains a mutex and a large array of bytes.
type Memory struct {
	sync.Mutex
	data [MemorySize]uint8
}

// WriteRaw will write len(data) number of bytes into m.data from offset
// and out.
//   No address translation happens.
//   Trying to write out of range will return an error and no data will
// be written.
//   Maybe this should panic?
func (m *Memory) WriteRaw(address uint32, data []uint8) (error, int) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.data)-len(data)) {
		return fmt.Errorf("Address out of range!"), 0
	}

	copy(m.data[address:], data)

	return nil, len(data)
}

// ReadRaw n number of bytes from address and out
//   No address translation happens.
//   Trying to read out of range will return an error and no data will
// be read.
//   Maybe this should panic?
func (m *Memory) ReadRaw(address, n uint32) (error, []uint8) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.data))-n {
		return fmt.Errorf("Address out of range!"), nil
	}

	bytes := make([]uint8, n)
	copy(bytes, m.data[address:address+n])

	return nil, bytes
}

// NewMemory creates a new instance of `Memory`
func NewMemory() Memory {
	return Memory{}
}
