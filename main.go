package main

import (
	"encoding/binary"
	"fmt"
	"gotos/cpu"
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

	mem := cpu.NewMemory(cpu.ENDIAN_LITTLE) // little endian memory
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

	rs := cpu.NewReservationSets(4)

	core0 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 0)
	core0.UnsafeSetMemSize(1024 * 1024 * 1)

	core1 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 1)
	core1.UnsafeSetMemSize(1024 * 1024 * 2)

	core2 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 2)
	core2.UnsafeSetMemSize(1024 * 1024 * 3)

	core3 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 3)
	core3.UnsafeSetMemSize(1024 * 1024 * 4)

	var wg sync.WaitGroup
	core0.StartAndSync(&wg)
	core1.StartAndSync(&wg)
	core2.StartAndSync(&wg)
	core3.StartAndSync(&wg)
	// don't need to wait on wg since we're waiting on cores

	core0.Wait()
	core1.Wait()
	core2.Wait()
	core3.Wait()

	fmt.Printf("\ncore0: %d cycles\n", core0.InstructionsRetired())
	fmt.Printf("core0: %d misses\n", core0.Misses())
	fmt.Printf("core0: %d accesses\n", core0.Accesses())
	core0.DumpRegisters()

	fmt.Printf("\ncore1: %d cycles\n", core1.InstructionsRetired())
	fmt.Printf("core1: %d misses\n", core1.Misses())
	fmt.Printf("core1: %d accesses\n", core1.Accesses())
	core1.DumpRegisters()

	fmt.Printf("\ncore2: %d cycles\n", core2.InstructionsRetired())
	fmt.Printf("core2: %d misses\n", core2.Misses())
	fmt.Printf("core2: %d accesses\n", core2.Accesses())
	core2.DumpRegisters()

	fmt.Printf("\ncore3: %d cycles\n", core3.InstructionsRetired())
	fmt.Printf("core3: %d misses\n", core3.Misses())
	fmt.Printf("core3: %d accesses\n", core3.Accesses())
	core3.DumpRegisters()
}
