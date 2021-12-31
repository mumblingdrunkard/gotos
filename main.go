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

	core1.StartAndWait()

	core1.Wait()

	fmt.Printf("%d cycles\n", core1.Cycles())

	core1.DumpRegisters()
}
