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
	// TODO create queue of PCBs with programs loaded in memory, then pass this
	// off to the system to handle instead of starting the core at some random address.

	system.Load("c-programs/locktest/locktest")

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

	err, l := mem.Write(0x4000, fib)
	if err != nil {
		panic(err)
	}

	err, l = mem.Write(0x14000, fib)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", l)

	rs := cpu.NewReservationSets(2)

	core0 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 0)
	core1 := cpu.NewCoreWithMemoryAndReservationSets(&mem, &rs, 1)

	core0.SetPC(0x4000) // TODO remove this
	core0.SetIRegister(cpu.Reg_SP, 0x8000)

	core1.SetPC(0x14000) // TODO remove this
	core1.SetIRegister(cpu.Reg_SP, 0x18000)

	core0.SetBootHandler(system.SystemStartup)
	core0.SetTrapHandler(system.TrapHandler)

	core1.SetBootHandler(system.SystemStartup)
	core1.SetTrapHandler(system.TrapHandler)

	var wg sync.WaitGroup
	core0.StartAndSync(&wg)
	core1.StartAndSync(&wg)

	core0.Wait()
	core1.Wait()

	core0.DumpRegisters()
	fmt.Println("Performance statistics:")
	fmt.Printf("\ncore0: %d cycles\n", core0.InstructionsRetired())
	fmt.Printf("core0: %d misses\n", core0.Misses())
	fmt.Printf("core0: %d accesses\n", core0.Accesses())

	core1.DumpRegisters()
	fmt.Println("Performance statistics:")
	fmt.Printf("\ncore1: %d cycles\n", core1.InstructionsRetired())
	fmt.Printf("core1: %d misses\n", core1.Misses())
	fmt.Printf("core1: %d accesses\n", core1.Accesses())
}
