// This file should contain utilities to load programs as "ELF" sections.
//
// Ideas:
// - Load program into memory (or at least store the segments somewhere), then
//	 return a PCB that is ready to be scheduled. This requires locking the
//   memory.
//
// Read https://wiki.osdev.org/ELF for a short introduction.

package system

import (
	"debug/elf"
	"fmt"
)

type Section struct {
	Offset uint32
	Data   []uint8
}

type ELF struct {
	Program Section
	Bss     Section
	Data    Section
}

// TODO should load the different sections of an ELF, place it into
// memory/storage, create a pcb and return the pcb.
// TODO in the future, it will also need to create a page table for the process
func Load(fname string) (error, *pcb) {
	f, err := elf.Open(fname)
	if err != nil {
		return err, nil
	}

	for _, s := range f.Sections {
		fmt.Println(s.Name)
		fmt.Println(s.Addr)
		fmt.Println(s.Offset)
		fmt.Println(s.Size)
		fmt.Println(s.Addralign)
		fmt.Println(s.Flags)
	}

	return nil, &pcb{}
}
