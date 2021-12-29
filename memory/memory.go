package memory

import (
	"encoding/binary"
	"fmt"
)

const (
	// 128 MiB of memory ought to be enough
	MEMORY_SIZE = 1024 * 1024 * 128
)

type Endian uint8

const (
	LITTLE Endian = 0
	BIG           = 1
)

type Memory struct {
	data   [MEMORY_SIZE]uint8
	endian Endian
}

// Return the byte stored at
func (m *Memory) LoadByte(address int) (error, uint8) {
	if address > len(m.data)-1 {
		return fmt.Errorf("Address out of range!"), 0
	}

	return nil, m.data[address]
}

func (m *Memory) LoadHalfWord(address int) (error, uint16) {
	if address > len(m.data)-2 {
		return fmt.Errorf("Address out of range!"), 0
	}

	if m.endian == BIG {
		return nil, binary.BigEndian.Uint16(m.data[address : address+2])
	} else {
		return nil, binary.LittleEndian.Uint16(m.data[address : address+2])
	}
}

func (m *Memory) LoadWord(address int) (error, uint32) {
	if address > len(m.data)-4 {
		return fmt.Errorf("Address out of range!"), 0
	}

	if m.endian == BIG {
		return nil, binary.BigEndian.Uint32(m.data[address : address+4])
	} else {
		return nil, binary.LittleEndian.Uint32(m.data[address : address+4])
	}
}

// Return the byte stored at
func (m *Memory) StoreByte(address int, b uint8) error {
	if address > len(m.data)-1 {
		return fmt.Errorf("Address out of range!")
	}

	m.data[address] = b

	return nil
}

func (m *Memory) StoreHalfWord(address int, hw uint16) error {
	if address > len(m.data)-2 {
		return fmt.Errorf("Address out of range!")
	}

	var bytes []uint8

	if m.endian == BIG {
		binary.BigEndian.PutUint16(bytes, hw)
	} else {
		binary.LittleEndian.PutUint16(bytes, hw)
	}

	for i, b := range bytes {
		m.data[address+i] = b
	}

	return nil
}

func (m *Memory) StoreWord(address int, w uint32) error {
	if address > len(m.data)-4 {
		return fmt.Errorf("Address out of range!")
	}

	var bytes []uint8

	if m.endian == BIG {
		binary.BigEndian.PutUint32(bytes, w)
	} else {
		binary.LittleEndian.PutUint32(bytes, w)
	}

	for i, b := range bytes {
		m.data[address+i] = b
	}

	return nil
}

// Write len(data) number of bytes into m.data from offset and out
func (m *Memory) Write(data []uint8, address int) error {
	if address > len(m.data)-len(data) {
		return fmt.Errorf("Address out of range!")
	}

	for i, b := range data {
		m.data[address+i] = b
	}

	return nil
}

// Read n number of bytes from address and out
func (m *Memory) Read(address, n int) (error, []uint8) {
	if address > len(m.data)-n {
		return fmt.Errorf("Address out of range!"), nil
	}

	bytes := make([]uint8, n)
	for i := range bytes {
		bytes[i] = m.data[address+i]
	}

	return nil, bytes
}

func NewMemory(endianness Endian) Memory {
	return Memory{
		endian: endianness,
	}
}
