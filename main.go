package main

import (
	"encoding/binary"
	"fmt"
	"gotos/cpu"
	"gotos/memory"
	"os"
)

func main() {
	f, err := os.Open("fib/fib.text")
	if err != nil {
		panic(err)
	}
	stats, err := f.Stat()
	size := stats.Size()
	program := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &program)
	f.Close()

	mem := memory.NewMemory(memory.LITTLE) // little endian memory
	err, l := mem.Write(program, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", l)

	core1 := cpu.NewCoreWithMemory(&mem)

	for core1.State() != cpu.HALTING {
		// core1.DumpRegisters()
		core1.Step()
	}

	core1.DumpRegisters()
}
