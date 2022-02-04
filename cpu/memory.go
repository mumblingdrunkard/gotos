package cpu

import (
	"fmt"
	"sync"
)

const (
	// 4 MiB of memory ought to be enough
	MEMORY_SIZE = 1024 * 1024 * 1
)

type Memory struct {
	sync.Mutex
	data [MEMORY_SIZE]uint8
}

// Write len(data) number of bytes into m.data from offset and out
func (m *Memory) Write(address uint32, data []uint8) (error, int) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.data)-len(data)) {
		return fmt.Errorf("Address out of range!"), 0
	}

	copy(m.data[address:], data)

	return nil, len(data)
}

// Read n number of bytes from address and out
func (m *Memory) Read(address, n uint32) (error, []uint8) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.data))-n {
		return fmt.Errorf("Address out of range!"), nil
	}

	bytes := make([]uint8, n)
	copy(bytes, m.data[address:address+n])

	return nil, bytes
}

func (m *Memory) Size() uint32 {
	return MEMORY_SIZE
}

func NewMemory() Memory {
	return Memory{}
}
