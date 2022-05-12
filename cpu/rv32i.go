// This file contains implementations of the instructions specified in
// the I extension of the RISC-V unprivileged specification.
//   Refer to the specification for instruction documentation.

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

	targetAddress := c.pc + offset
	if targetAddress&0x3 != 0 {
		c.csr[Csr_MTVAL] = targetAddress
		c.trap(TrapInstructionAddressMisaligned)
		return
	} else {
		c.reg[rd] = c.pc + 4
		c.pc = targetAddress
		c.jumped = true
	}
}

// jump and link register
func (c *Core) jalr(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs1_val := c.reg[rs1]
	imm11_0 := uint32(int32(inst) >> 20)

	targetAddress := (imm11_0 + rs1_val) & 0xfffffffe

	if targetAddress&0x3 != 0 {
		c.csr[Csr_MTVAL] = targetAddress
		c.trap(TrapInstructionAddressMisaligned)
		return
	} else {
		c.reg[rd] = c.pc + 4
		c.pc = targetAddress
		c.jumped = true
	}
}

// branch equal
func (c *Core) beq(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if c.reg[rs1] == c.reg[rs2] {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}

}

// branch not equal
func (c *Core) bne(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if c.reg[rs1] != c.reg[rs2] {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}
}

// branch less than
func (c *Core) blt(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if int32(c.reg[rs1]) < int32(c.reg[rs2]) {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}
}

// branch less than unsigned
func (c *Core) bltu(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if c.reg[rs1] < c.reg[rs2] {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}
}

// branch greater than or equal to
func (c *Core) bge(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if int32(c.reg[rs1]) >= int32(c.reg[rs2]) {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}
}

// branch greater than or equal to unsigned
func (c *Core) bgeu(inst uint32) {
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm4_1 := (inst >> 8) & 0xf
	imm11 := (inst >> 7) & 0x1
	imm10_5 := (inst >> 25) & 0x3f
	imm12 := (int32(inst) >> 31) // sign extended
	offset := (uint32(imm12) << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)

	targetAddress := c.pc + offset

	if c.reg[rs1] >= c.reg[rs2] {
		if targetAddress&0x3 != 0 {
			c.csr[Csr_MTVAL] = targetAddress
			c.trap(TrapInstructionAddressMisaligned)
			return
		} else {
			c.pc = targetAddress
			c.jumped = true
		}
	}
}

// load byte (signed)
func (c *Core) lb(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	if success, b := c.loadByte(address); success {
		signed := int8(b)
		extended := int32(signed)
		converted := uint32(extended)
		c.reg[rd] = converted
	}
}

// load half (signed)
func (c *Core) lh(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	if success, hw := c.loadHalfWord(address); success {
		signed := int16(hw)
		extended := int32(signed)
		converted := uint32(extended)
		c.reg[rd] = converted
	}
}

// load word
func (c *Core) lw(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	if success, w := c.loadWord(address); success {
		c.reg[rd] = w
	}
}

// load byte unsigned
func (c *Core) lbu(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	if success, b := c.loadByte(address); success {
		c.reg[rd] = uint32(b)
	}
}

// load half unsigned
func (c *Core) lhu(inst uint32) {
	rd := (inst >> 7) & 0x1f             // dest
	rs1 := (inst >> 15) & 0x1f           // base
	imm11_0 := uint32(int32(inst) >> 20) // offset
	address := imm11_0 + c.reg[rs1]

	if success, hw := c.loadHalfWord(address); success {
		c.reg[rd] = uint32(hw)
	}
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

	c.storeByte(address, b)
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

	c.storeHalfWord(address, hw)
}

// store word
func (c *Core) sw(inst uint32) {
	rs1 := (inst >> 15) & 0x1f           // base
	rs2 := (inst >> 20) & 0x1f           // src
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f
	offset := (imm11_5 << 5) | imm4_0

	address := offset + c.reg[rs1]

	c.storeWord(address, c.reg[rs2])
}

// fence
func (c *Core) fence(inst uint32) {
	// I think this should flush the cache?
	// Will need some input from an expert or something.
	c.FENCE()

	// This ensures no memory operations from this hart can be observed before any memory operation that comes after the FENCE.

	// Any FENCE to invalidate the cache?

	// Do I have to invalidate all the other caches?
	// Seems pretty wasteful...
	// Could I invalidate the cache for this hart with a fence?

	// Locking will need some cache invalidation to ensure coherence

	// Lock
	//   Atomically acquire lock
	//   Invalidate my cache (new data may be available from other harts)

	// Unlock
	//   Flush my cache
	//   Atomically release lock

	// The only way to guarantee memory consistency would then be by ensuring the lock is acquired before trying to access memory
}

// environment call
func (c *Core) ecall(inst uint32) {
	c.trap(TrapEcallUMode)
	return
}

// environment break
func (c *Core) ebreak(inst uint32) {
	c.trap(TrapBreakpoint)
	return
}
