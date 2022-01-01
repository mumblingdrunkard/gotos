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
	fib := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &fib)
	f.Close()

	mem := memory.NewMemory(memory.LITTLE) // little endian memory
	err, l := mem.Write(fib, 0)
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
	core1 := cpu.NewCoreWithMemory(&mem)
	core1.UnsafeSetMemSize(1024 * 1024 * 1)

	core2 := cpu.NewCoreWithMemory(&mem)
	core2.UnsafeSetMemSize(1024 * 1024 * 2)

	core3 := cpu.NewCoreWithMemory(&mem)
	core3.UnsafeSetMemSize(1024 * 1024 * 3)

	core4 := cpu.NewCoreWithMemory(&mem)
	core4.UnsafeSetMemSize(1024 * 1024 * 4)

	core1.StartAndWait()
	core2.StartAndWait()
	core3.StartAndWait()
	core4.StartAndWait()

	core1.Wait()
	core2.Wait()
	core3.Wait()
	core4.Wait()

	fmt.Printf("core1: %d cycles\n", core1.Cycles())
	core1.DumpRegisters()

	fmt.Printf("core2: %d cycles\n", core2.Cycles())
	core2.DumpRegisters()

	fmt.Printf("core3: %d cycles\n", core3.Cycles())
	core3.DumpRegisters()

	fmt.Printf("core3: %d cycles\n", core3.Cycles())
	core4.DumpRegisters()
}
