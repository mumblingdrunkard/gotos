package memory

import (
	"encoding/binary"
	"fmt"
	"sync"
)

const (
	// 4 MiB of memory ought to be enough
	MEMORY_SIZE = 1024 * 1024 * 4
)

type Endian uint8

const (
	LITTLE Endian = 0
	BIG           = 1
)

type Memory struct {
	sync.Mutex
	Data   [MEMORY_SIZE]uint8
	endian Endian
}

// Return the byte stored at
func (m *Memory) LoadByte(address uint32) (error, uint8) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-1) {
		return fmt.Errorf("Address `%x` out of range!", address), 0
	}

	return nil, m.Data[address]
}

func (m *Memory) LoadHalfWord(address uint32) (error, uint16) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-2) {
		return fmt.Errorf("Address `%x` out of range!", address), 0
	}

	if m.endian == BIG {
		return nil, binary.BigEndian.Uint16(m.Data[address : address+2])
	} else {
		return nil, binary.LittleEndian.Uint16(m.Data[address : address+2])
	}
}

func (m *Memory) LoadWord(address uint32) (error, uint32) {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-4) {
		return fmt.Errorf("Address `%x` out of range!", address), 0
	}

	if m.endian == BIG {
		return nil, binary.BigEndian.Uint32(m.Data[address : address+4])
	} else {
		return nil, binary.LittleEndian.Uint32(m.Data[address : address+4])
	}
}

// Return the byte stored at
func (m *Memory) StoreByte(address uint32, b uint8) error {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-1) {
		return fmt.Errorf("Address `%x` out of range!", address)
	}

	m.Data[address] = b

	return nil
}

func (m *Memory) StoreHalfWord(address uint32, hw uint16) error {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-2) {
		return fmt.Errorf("Address `%x` out of range!", address)
	}

	bytes := make([]uint8, 2)

	if m.endian == BIG {
		binary.BigEndian.PutUint16(bytes, hw)
	} else {
		binary.LittleEndian.PutUint16(bytes, hw)
	}

	copy(m.Data[address:address+2], bytes)

	return nil
}

func (m *Memory) StoreWord(address uint32, w uint32) error {
	m.Lock()
	defer m.Unlock()
	if address > uint32(len(m.Data)-4) {
		return fmt.Errorf("Address `%x` out of range!", address)
	}

	bytes := make([]uint8, 4)

	if m.endian == BIG {
		binary.BigEndian.PutUint32(bytes, w)
	} else {
		binary.LittleEndian.PutUint32(bytes, w)
	}

	copy(m.Data[address:address+4], bytes)

	return nil
}

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

func (m *Memory) Dump() {
	fmt.Println("Memory dump:")
	for i, v := range m.Data {
		fmt.Printf("[%06X] : %02X \n", i, v)
	}
}

func NewMemory(endianness Endian) Memory {
	return Memory{
		endian: endianness,
	}
}
