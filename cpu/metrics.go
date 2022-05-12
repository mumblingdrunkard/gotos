package cpu

import "fmt"

// DumpRegisters will print all registers, and the program counter as
// hexadecimal values.
func (c *Core) DumpRegisters() {
	fmt.Printf("\n=== Register dump for core %d ===\n", c.csr[Csr_MHARTID])
	fmt.Printf("pc: %X\n", c.pc)
	fmt.Println("Integer registers")
	for i, r := range c.reg {
		if r != 0 {
			fmt.Printf("[%02d]: %08X\n", i, r)
		}
	}
}
