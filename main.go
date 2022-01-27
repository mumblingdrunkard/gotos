package main

import (
	"encoding/binary"
	"fmt"
	"gotos/cpu"
	"gotos/system"
	"os"
	"sync"
)

func main() {
	f, err := os.Open("c-programs/locktest/locktest.text")
	if err != nil {
		panic(err)
	}
	stats, err := f.Stat()
	size := stats.Size()
	fib := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &fib)
	f.Close()

	mem := cpu.NewMemory(cpu.EndianLittle) // little endian memory
	err, l := mem.Write(0, fib)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", l)

	// All the cores share the same program, but have their stack start at different addresses.
	// By convention, the stack is placed at the end of the address space.
	// By giving all the cores different memory sizes, we essentially give them separate stack spaces.
	// It is possible that the stack for core4 could flow into the stack of core3 without causing an error.
	// For now, we just assume/hope this doesn't happen.
	// More advanced memory sharing techniques are required for this.
	rs := cpu.NewReservationSets(1)

	core0 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 0)
	core0.SetBootHandler(system.SystemStartup)
	core0.SetTrapHandler(system.TrapHandler)
	core0.UnsafeSetMemSize(1024 * 1024 * 1)

	var wg sync.WaitGroup
	core0.StartAndSync(&wg)

	core0.Wait()

	core0.DumpRegisters()
	fmt.Println("Performance statistics:")
	fmt.Printf("\ncore0: %d cycles\n", core0.InstructionsRetired())
	fmt.Printf("core0: %d misses\n", core0.Misses())
	fmt.Printf("core0: %d accesses\n", core0.Accesses())
}
