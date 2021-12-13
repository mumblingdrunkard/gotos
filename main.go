package main

import (
	"fmt"
	"gotos/cpu"
	"sync"
	"time"
)

func main() {
	core1 := cpu.NewCore()
	core2 := cpu.NewCore()
	core3 := cpu.NewCore()
	core4 := cpu.NewCore()

	fmt.Println("Starting core...")

	var wg sync.WaitGroup

	core1.StartAndSync(&wg)
	core2.StartAndSync(&wg)
	core3.StartAndSync(&wg)
	core4.StartAndSync(&wg)

	fmt.Println("Running...")

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("Waiting for core to finish...")

	core1.HaltAndSync(&wg)
	core2.HaltAndSync(&wg)
	core3.HaltAndSync(&wg)
	core4.HaltAndSync(&wg)

	fmt.Println("Done!")

	state1 := core1.UnsafeGetState()
	state2 := core2.UnsafeGetState()
	state3 := core3.UnsafeGetState()
	state4 := core4.UnsafeGetState()

	fmt.Println()
	fmt.Println("Register 3 on core 1:")
	fmt.Printf("[%d]:\t %x₁₆ = %d₁₀\n", 3, state1.Reg()[3], state1.Reg()[3])
	fmt.Println()

	fmt.Println()
	fmt.Println("Register 3 on core 2:")
	fmt.Printf("[%d]:\t %x₁₆ = %d₁₀\n", 3, state2.Reg()[3], state2.Reg()[3])
	fmt.Println()

	fmt.Println()
	fmt.Println("Register 3 on core 3:")
	fmt.Printf("[%d]:\t %x₁₆ = %d₁₀\n", 3, state3.Reg()[3], state3.Reg()[3])
	fmt.Println()

	fmt.Println()
	fmt.Println("Register 3 on core 4:")
	fmt.Printf("[%d]:\t %x₁₆ = %d₁₀\n", 3, state4.Reg()[3], state4.Reg()[3])
	fmt.Println()
}
