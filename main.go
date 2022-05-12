package main

import (
	"encoding/binary"
	"gotos/system"
)

const (
	pageFlagValid    uint32 = 0x01 // the virtual address is valid
	pageFlagRead            = 0x02 // indicates that the processor is allowed to read data from this address
	pageFlagWrite           = 0x04 // indicates that the processor is allowed to write data to this address
	pageFlagExec            = 0x08 // indicates that the processor is allowed to fetch instructions from this address
	pageFlagUser            = 0x10 // indicates that the processor can access this page in user mode
	pageFlagGlobal          = 0x20 // whether this page is globally mapped into all address spaces (probably unused here)
	pageFlagAccessed        = 0x40 // whether this page has been accessed since the access bit was last cleared
	pageFlagDirty           = 0x80 // whether this page has been written to since the dirty bit was last cleared
)

func main() {
	// create a system with 1 core
	sys := system.NewSystem(4)

	// create a simple batch scheduler queue (FIFO)
	fifo := &system.FIFO{}
	sys.Scheduler = fifo

	scratchTable := [1024]uint32{}
	page := tableToPage(scratchTable[:])

	// ---- PROCESS 0 ----
	// top level page table
	scratchTable[0] = pageFlagValid | 0x002<<10
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x00001000, page[:])

	// second level page table
	scratchTable[0x000] = (0x000 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	// program  u v a       x
	scratchTable[0x004] = (0x004 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagExec
	// stack    u v a d r w
	scratchTable[0x005] = (0x005 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite

	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x00002000, page[:])

	// load the program
	// file, pc, sp, pid, addr
	sys.Load("c-programs/fib/main.text", 0x00004000, 0x00006000, 0, 0x00004000, 0x00001000)

	// ---- PROCESS 1 ----
	// top level page table
	scratchTable[0] = 0x007<<10 | pageFlagValid
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x00006000, page[:])

	// second level page table
	scratchTable[0x000] = (0x000 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	// program  u v a       x
	scratchTable[0x004] = (0x004 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagExec
	// stack    u v a d r w
	scratchTable[0x005] = (0x008 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x00007000, page[:])

	sys.Load("c-programs/fib/main.text", 0x00004000, 0x00006000, 1, 0x00004000, 0x00006000)

	// ---- PROCESS 2 ----
	// top level page table
	scratchTable[0] = 0x00A<<10 | pageFlagValid
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x00009000, page[:])

	// second level page table
	scratchTable[0x000] = (0x000 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	// program  u v a       x
	scratchTable[0x004] = (0x004 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagExec
	// stack    u v a d r w
	scratchTable[0x005] = (0x00B << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x0000A000, page[:])

	sys.Load("c-programs/fib/main.text", 0x00004000, 0x00006000, 2, 0x00004000, 0x00009000)

	// ---- PROCESS 3 ----
	// top level page table
	scratchTable[0] = 0x00D<<10 | pageFlagValid
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x0000C000, page[:])

	// second level page table
	scratchTable[0x000] = (0x000 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	// program  u v a       x
	scratchTable[0x004] = (0x004 << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagExec
	// stack    u v a d r w
	scratchTable[0x005] = (0x00E << 10) | pageFlagUser | pageFlagValid | pageFlagAccessed | pageFlagDirty | pageFlagRead | pageFlagWrite
	page = tableToPage(scratchTable[:])
	sys.Memory().WriteRaw(0x0000D000, page[:])

	sys.Load("c-programs/fib/main.text", 0x00004000, 0x00006000, 3, 0x00004000, 0x0000C000)

	// run the system
	sys.Run()
}

func tableToPage(table []uint32) [4096]uint8 {
	var page [4096]uint8
	for i, w := range table {
		var bytes [4]uint8
		binary.LittleEndian.PutUint32(bytes[:], w)
		for j, b := range bytes {
			page[i*4+j] = b
		}
		table[i] = 0
	}
	return page
}
