// This file contains implementations of the instructions specified in
// the M extension of the RISC-V unprivileged specification.
//   Refer to the specification for instruction documentation.

package cpu

func (c *Core) mul(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := int64(int32(c.reg[rs1])) // sign extended to 64 bits
	b := int64(int32(c.reg[rs2])) // sign extended to 64 bits

	c.reg[rd] = uint32(uint64(a * b))
}

func (c *Core) mulh(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := int64(int32(c.reg[rs1])) // sign extended to 64 bits
	b := int64(int32(c.reg[rs2])) // sign extended to 64 bits

	c.reg[rd] = uint32(uint64(a*b) >> 32)
}

func (c *Core) mulhu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := uint64(c.reg[rs1])
	b := uint64(c.reg[rs2])

	c.reg[rd] = uint32((a * b) >> 32)
}

func (c *Core) mulhsu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := int64(int32(c.reg[rs1])) // sign extended to 64 bits
	b := int64(c.reg[rs2])        // sign extended to 64 bits

	c.reg[rd] = uint32((a * b) >> 32)
}

func (c *Core) div(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := int32(c.reg[rs1])
	b := int32(c.reg[rs2])

	if b == 0 {
		c.reg[rd] = 0xFFFFFFFF
		return
	}

	c.reg[rd] = uint32(a / b)
}

func (c *Core) divu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := c.reg[rs1]
	b := c.reg[rs2]

	if b == 0 {
		c.reg[rd] = 0xFFFFFFFF
		return
	}

	c.reg[rd] = a / b
}

func (c *Core) rem(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := int32(c.reg[rs1])
	b := int32(c.reg[rs2])

	if b == 0 {
		c.reg[rd] = uint32(a)
		return
	}

	c.reg[rd] = uint32(a % b)
}

func (c *Core) remu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	a := c.reg[rs1]
	b := c.reg[rs2]

	if b == 0 {
		c.reg[rd] = uint32(a)
		return
	}

	c.reg[rd] = a % b
}
