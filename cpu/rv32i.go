// RV32I Base Integer Instruction Set
// https://mark.theis.site/riscv/

// TODO: Write tests for instructions

package cpu

import "fmt"

// add immediate
func (c *Core) addi(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	// interpret as int32 to get arithmetic right shift
	imm11_0 := uint32(int32(inst) >> 20)
	c.reg[rd] = c.reg[rs1] + imm11_0
}

// set less than immediate
func (c *Core) slti(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := int32(inst) >> 20
	if int32(c.reg[rs1]) < imm11_0 {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

// set less than immediate unsigned
func (c *Core) sltiu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := uint32(int32(inst) >> 20)
	if c.reg[rs1] < imm11_0 {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

// xor immediate
func (c *Core) xori(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := uint32(int32(inst) >> 20)
	c.reg[rd] = c.reg[rs1] ^ imm11_0
}

// or immediate
func (c *Core) ori(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := uint32(int32(inst) >> 20)
	c.reg[rd] = c.reg[rs1] | imm11_0
}

// and immediate
func (c *Core) andi(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := uint32(int32(inst) >> 20)
	c.reg[rd] = c.reg[rs1] & imm11_0
}

// shift left logical immediate
func (c *Core) slli(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm4_0 := (inst >> 20) & 0x1f // shamt, shift amount
	c.reg[rd] = c.reg[rs1] << imm4_0
}

// shift right logical immediate
func (c *Core) srli(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm4_0 := (inst >> 20) & 0x1f    // shamt, shift amount
	c.reg[rd] = c.reg[rs1] >> imm4_0 // logical shift
}

// shift right arithmetic immediate
func (c *Core) srai(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm4_0 := (inst >> 20) & 0x1f                   // shamt, shift amount
	c.reg[rd] = uint32(int32(c.reg[rs1]) >> imm4_0) // logical shift
}

// load upper immediate
func (c *Core) lui(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm31_12 := inst & 0xfffff000
	c.reg[rd] = imm31_12
}

// add upper immediate pc
func (c *Core) auipc(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm31_12 := inst & 0xfffff000
	c.reg[rd] = c.pc + imm31_12
}

// add
func (c *Core) add(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	c.reg[rd] = c.reg[rs1] + c.reg[rs2]
}

// set less than (signed)
func (c *Core) slt(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	if int32(c.reg[rs1]) < int32(c.reg[rs2]) {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

// set less than unsigned
func (c *Core) sltu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	if c.reg[rs1] < c.reg[rs2] {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

// bitwise and
func (c *Core) and(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	c.reg[rd] = c.reg[rs1] & c.reg[rs2]
}

// bitwise or
func (c *Core) or(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	c.reg[rd] = c.reg[rs1] | c.reg[rs2]
}

// exclusive bitwise or
func (c *Core) xor(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	c.reg[rd] = c.reg[rs1] ^ c.reg[rs2]
}

// shift left logical
func (c *Core) sll(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	shift_amount := c.reg[rs2] & 0x1f
	c.reg[rd] = c.reg[rs1] << shift_amount
}

// shift right logical
func (c *Core) srl(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	shift_amount := c.reg[rs2] & 0x1f
	c.reg[rd] = c.reg[rs1] >> shift_amount
}

// sub
func (c *Core) sub(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	c.reg[rd] = c.reg[rs1] - c.reg[rs2]
}

// shift right arithmetic
func (c *Core) sra(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	shift_amount := c.reg[rs2] & 0x1f
	c.reg[rd] = uint32(int32(c.reg[rs1]) >> shift_amount)
}

// TODO: figure out how to raise exceptions when jumps are misaligned
// jump and link
func (c *Core) jal(inst uint32) {
	rd := (inst >> 7) & 0x1f
	// What the fuck is this format?
	imm19_12 := (inst >> 12) & 0xff
	imm11 := (inst >> 20) & 1
	imm10_1 := (inst >> 21) & 0x3ff
	imm20 := uint32(int32(inst) >> 31) // for sign extension
	// Why couldn't this just be imm[20:1] ?
	// I think this is how it's supposed to work?
	offset := (imm10_1 << 1) | (imm11 << 11) | (imm19_12 << 12) | (imm20 << 20)
	c.reg[rd] = c.pc + 4
	c.pc = c.pc + offset
}

// jump and link register
func (c *Core) jalr(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs1_val := c.reg[rs1]
	imm11_0 := uint32(int32(inst) >> 20)
	c.reg[rd] = c.pc + 4
	c.pc = (rs1_val + imm11_0) & 0xfffffffe
}

// branch equal
func (c *Core) beq(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if c.reg[rs1] == c.reg[rs2] {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// branch not equal
func (c *Core) bne(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if c.reg[rs1] != c.reg[rs2] {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// branch less than
func (c *Core) blt(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if int32(c.reg[rs1]) < int32(c.reg[rs2]) {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// branch less than unsigned
func (c *Core) bltu(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if c.reg[rs1] < c.reg[rs2] {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// branch greater than or equal to
func (c *Core) bge(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if int32(c.reg[rs1]) >= int32(c.reg[rs2]) {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// branch greater than or equal to unsigned
func (c *Core) bgeu(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x1f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	if c.reg[rs1] >= c.reg[rs2] {
		c.pc += offset
	} else {
		c.pc += 4
	}
}

// TODO: Caching and address translation for load and store instructions

// load byte (signed)
func (c *Core) lb(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	err, b := c.mem.LoadByte(int(address))

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}

	signed := int8(b)
	extended := int32(signed)
	converted := uint32(extended)
	c.reg[rd] = converted
}

// load half (signed)
func (c *Core) lh(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	err, hw := c.mem.LoadHalfWord(int(address))

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}

	signed := int16(hw)
	extended := int32(signed)
	converted := uint32(extended)
	c.reg[rd] = converted
}

// load word
func (c *Core) lw(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	err, w := c.mem.LoadWord(int(address))

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}

	c.reg[rd] = w
}

// load byte unsigned
func (c *Core) lbu(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	err, b := c.mem.LoadByte(int(address))

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}

	c.reg[rd] = uint32(b)
}

// load half unsigned
func (c *Core) lhu(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	err, hw := c.mem.LoadHalfWord(int(address))

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}

	c.reg[rd] = uint32(hw)
}

// store byte
func (c *Core) sb(inst uint32) {
	rs1 := (inst >> 15) & 0x1f           // base
	rs2 := (inst >> 20) & 0x1f           // src
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f
	offset := (imm11_5 << 5) | imm4_0

	address := offset + c.reg[rs1]
	b := uint8(c.reg[rs2] & 0xff)

	err := c.mem.StoreByte(int(address), b)

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}
}

// store half
func (c *Core) sh(inst uint32) {
	rs1 := (inst >> 15) & 0x1f           // base
	rs2 := (inst >> 20) & 0x1f           // src
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f
	offset := (imm11_5 << 5) | imm4_0

	address := offset + c.reg[rs1]
	hw := uint16(c.reg[rs2] & 0xffff)

	err := c.mem.StoreHalfWord(int(address), hw)

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}
}

// store word
func (c *Core) sw(inst uint32) {
	rs1 := (inst >> 15) & 0x1f           // base
	rs2 := (inst >> 20) & 0x1f           // src
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f
	offset := (imm11_5 << 5) | imm4_0

	address := offset + c.reg[rs1]

	err := c.mem.StoreWord(int(address), c.reg[rs2])

	if err != nil {
		fmt.Println("Register dump")
		fmt.Printf("pc: %x\n", c.pc)
		for i, r := range c.reg {
			fmt.Printf("[%d]: %x\n", i, r)
		}
		panic(err)
	}
}

// fence
func (c *Core) fence(inst uint32) {
	// TODO
	// no-op for now
}

// environment call
func (c *Core) ecall(inst uint32) {
	// TODO just halts for now
	// c.state.Store(HALTING)
	c.UnsafeReset()
}

// environment break
func (c *Core) ebreak(inst uint32) {
	// TODO no-op for now
}
