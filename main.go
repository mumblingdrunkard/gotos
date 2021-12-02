package main

import (
	"fmt"
	"gotos/cpu"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	core := cpu.NewCPU()

	file, err := os.Open("add-addi.bin")

	if err != nil {
		panic(err)
	}

	program := make([]byte, 128)
	_, err = file.Read(program)

	if err != nil {
		panic(err)
	}

	core.Reset()
	core.LoadMemory(program, 0)

	fmt.Println("Starting core...")

	wg.Add(1)
	go core.Start(&wg)

	fmt.Println("Running...")

	time.Sleep(1 * time.Millisecond)
	core.Stop()

	fmt.Println("Waiting for core to finish...")
	wg.Wait()

	fmt.Println("Done!")

	state := core.GetState()

	fmt.Println()
	fmt.Println("Register dump:")
	for i, val := range state.Reg() {
		fmt.Printf("[%d]:\t 0x%x\n", i, val)
	}
	fmt.Println()
	fmt.Printf("pc: 0x%x\n", state.Pc())
}
