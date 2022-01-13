package cpu

// Allocated Unprivileged CSR addresses
const (
	// Unprivileged Floating-Point CSRs
	csr_FFLAGS uint32 = 0x001
	csr_FRM           = 0x002
	csr_FCSR          = 0x003

	// Unprivileged Counters/Timers
	csr_CYCLE        = 0xC00
	csr_TIME         = 0xC01
	csr_INSTRET      = 0xC02
	csr_HPMCOUNTER3  = 0xC03
	csr_HPMCOUNTER4  = 0xC04
	csr_HPMCOUNTER5  = 0xC05
	csr_HPMCOUNTER6  = 0xC06
	csr_HPMCOUNTER7  = 0xC07
	csr_HPMCOUNTER8  = 0xC08
	csr_HPMCOUNTER9  = 0xC09
	csr_HPMCOUNTER10 = 0xC0A
	csr_HPMCOUNTER11 = 0xC0B
	csr_HPMCOUNTER12 = 0xC0C
	csr_HPMCOUNTER13 = 0xC0D
	csr_HPMCOUNTER14 = 0xC0E
	csr_HPMCOUNTER15 = 0xC0F
	csr_HPMCOUNTER16 = 0xC10
	csr_HPMCOUNTER17 = 0xC11
	csr_HPMCOUNTER18 = 0xC12
	csr_HPMCOUNTER19 = 0xC13
	csr_HPMCOUNTER20 = 0xC14
	csr_HPMCOUNTER21 = 0xC15
	csr_HPMCOUNTER22 = 0xC16
	csr_HPMCOUNTER23 = 0xC17
	csr_HPMCOUNTER24 = 0xC18
	csr_HPMCOUNTER25 = 0xC19
	csr_HPMCOUNTER26 = 0xC1A
	csr_HPMCOUNTER27 = 0xC1B
	csr_HPMCOUNTER28 = 0xC1C
	csr_HPMCOUNTER29 = 0xC1D
	csr_HPMCOUNTER30 = 0xC1E
	csr_HPMCOUNTER31 = 0xC1F
	// Upper bits
	csr_CYCLEH        = 0xC80
	csr_TIMEH         = 0xC81
	csr_INSTRETH      = 0xC82
	csr_HPMCOUNTER3H  = 0xC83
	csr_HPMCOUNTER4H  = 0xC84
	csr_HPMCOUNTER5H  = 0xC85
	csr_HPMCOUNTER6H  = 0xC86
	csr_HPMCOUNTER7H  = 0xC87
	csr_HPMCOUNTER8H  = 0xC88
	csr_HPMCOUNTER9H  = 0xC89
	csr_HPMCOUNTER10H = 0xC8A
	csr_HPMCOUNTER11H = 0xC8B
	csr_HPMCOUNTER12H = 0xC8C
	csr_HPMCOUNTER13H = 0xC8D
	csr_HPMCOUNTER14H = 0xC8E
	csr_HPMCOUNTER15H = 0xC8F
	csr_HPMCOUNTER16H = 0xC90
	csr_HPMCOUNTER17H = 0xC91
	csr_HPMCOUNTER18H = 0xC92
	csr_HPMCOUNTER19H = 0xC93
	csr_HPMCOUNTER20H = 0xC94
	csr_HPMCOUNTER21H = 0xC95
	csr_HPMCOUNTER22H = 0xC96
	csr_HPMCOUNTER23H = 0xC97
	csr_HPMCOUNTER24H = 0xC98
	csr_HPMCOUNTER25H = 0xC99
	csr_HPMCOUNTER26H = 0xC9A
	csr_HPMCOUNTER27H = 0xC9B
	csr_HPMCOUNTER28H = 0xC9C
	csr_HPMCOUNTER29H = 0xC9D
	csr_HPMCOUNTER30H = 0xC9E
	csr_HPMCOUNTER31H = 0xC9F
)

func (c *Core) csrrw(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	// Don't need to check rd for 0
	// We just reset the zero register to 0 in the next instruction anyway

	// Check permissions
	if (csr&0xC00 != 0) && (csr&0xC00 != 0xC00) { // verify read and write permissions
		old := c.csr[csr]
		c.csr[csr] = c.reg[rs1]
		c.reg[rd] = old
	} else {
		c.trap(TrapIllegalInstruction)
	}
}

func (c *Core) csrrs(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
	}

	if rs1 != RegZero {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= c.reg[rs1]
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrc(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
	}

	if rs1 != RegZero {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (c.reg[rs1] ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrwi(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	// Don't need to check rd for 0
	// We just reset the zero register to 0 in the next instruction anyway

	// Check permissions
	if (csr&0xC00 != 0) && (csr&0xC00 != 0xC00) { // verify read and write permissions
		old := c.csr[csr]
		c.csr[csr] = imm4_0
		c.reg[rd] = old
	} else {
		// illegal instruction?
		c.trap(TrapIllegalInstruction)
	}
}

func (c *Core) csrrsi(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= imm4_0
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrci(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (imm4_0 ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}
