package main

import (
	"encoding/binary"
	"fmt"
	"gotos/cpu"
	"gotos/memory"
	"os"
	"sync"
)

func main() {
	f, err := os.Open("c-program/minimal.text")
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

	var wg sync.WaitGroup

	core1.StartAndSync(&wg)

	wg.Wait()

	core1.SyncOnHalt(&wg)

	wg.Wait()

	state1 := core1.UnsafeGetState()

	fmt.Println("Register dump:")
	fmt.Printf("pc: %x\n", state1.Pc())
	for i, r := range state1.Reg() {
		fmt.Printf("[%d]: %x\n", i, r)
	}
}
