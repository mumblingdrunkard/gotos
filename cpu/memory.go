package cpu

import (
	"fmt"
	"sync"
)

const (
	// 4 MiB of memory ought to be enough
	MEMORY_SIZE = 1024 * 1024 * 4
)

type Memory struct {
	sync.Mutex
	Data   [MEMORY_SIZE]uint8
	endian Endian
}

// Return the byte stored at
// Write len(data) number of bytes into m.data from offset and out
func (m *Memory) Write(address uint32, data []uint8) (error, int) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-len(data)) {
		return fmt.Errorf("Address out of range!"), 0
	}

	copy(m.Data[address:], data)

	return nil, len(data)
}

// Read n number of bytes from address and out
func (m *Memory) Read(address, n uint32) (error, []uint8) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data))-n {
		return fmt.Errorf("Address out of range!"), nil
	}

	bytes := make([]uint8, n)
	copy(bytes, m.Data[address:address+n])

	return nil, bytes
}

func (m *Memory) Size() uint32 {
	return MEMORY_SIZE
}

func NewMemory(endianness Endian) Memory {
	return Memory{
		endian: endianness,
	}
}
