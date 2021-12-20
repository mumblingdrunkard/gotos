// RV32I Base Integer Instruction Set
// https://mark.theis.site/riscv/

package cpu

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
	imm10_1 := (inst >> 20) & 0x3ff
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
	imm11_0 := uint32(int32(inst) >> 20)

	c.reg[rd] = c.pc + 4
	c.pc = (c.pc + imm11_0) & 0xfffffffe
}

func (c *Core) beq(inst uint32) {
	// imm[12] | imm[10:5] | rs2 | rs1 | funct3 | imm[4:1] | imm[11] | opcode
}
