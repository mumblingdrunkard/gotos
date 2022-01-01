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

	f, err = os.Open("fib/fib.text")
	if err != nil {
		panic(err)
	}
	stats, err = f.Stat()
	size = stats.Size()
	minimal := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &minimal)
	f.Close()

	mem := memory.NewMemory(memory.LITTLE) // little endian memory
	err, l := mem.Write(fib, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", l)

	err, l = mem.Write(minimal, 1024*1024*2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", l)

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
