package cpu

// Allocated Unprivileged CSR addresses
const (
	// Unprivileged Floating-Point CSRs
	CSR_FFLAGS uint32 = 0x001
	CSR_FRM           = 0x002
	CSR_FCSR          = 0x003

	// Unprivileged Counters/Timers
	CSR_CYCLE        = 0xC00
	CSR_TIME         = 0xC01
	CSR_INSTRET      = 0xC02
	CSR_HPMCOUNTER3  = 0xC03
	CSR_HPMCOUNTER4  = 0xC04
	CSR_HPMCOUNTER5  = 0xC05
	CSR_HPMCOUNTER6  = 0xC06
	CSR_HPMCOUNTER7  = 0xC07
	CSR_HPMCOUNTER8  = 0xC08
	CSR_HPMCOUNTER9  = 0xC09
	CSR_HPMCOUNTER10 = 0xC0A
	CSR_HPMCOUNTER11 = 0xC0B
	CSR_HPMCOUNTER12 = 0xC0C
	CSR_HPMCOUNTER13 = 0xC0D
	CSR_HPMCOUNTER14 = 0xC0E
	CSR_HPMCOUNTER15 = 0xC0F
	CSR_HPMCOUNTER16 = 0xC10
	CSR_HPMCOUNTER17 = 0xC11
	CSR_HPMCOUNTER18 = 0xC12
	CSR_HPMCOUNTER19 = 0xC13
	CSR_HPMCOUNTER20 = 0xC14
	CSR_HPMCOUNTER21 = 0xC15
	CSR_HPMCOUNTER22 = 0xC16
	CSR_HPMCOUNTER23 = 0xC17
	CSR_HPMCOUNTER24 = 0xC18
	CSR_HPMCOUNTER25 = 0xC19
	CSR_HPMCOUNTER26 = 0xC1A
	CSR_HPMCOUNTER27 = 0xC1B
	CSR_HPMCOUNTER28 = 0xC1C
	CSR_HPMCOUNTER29 = 0xC1D
	CSR_HPMCOUNTER30 = 0xC1E
	CSR_HPMCOUNTER31 = 0xC1F
	// Upper bits
	CSR_CYCLEH        = 0xC80
	CSR_TIMEH         = 0xC81
	CSR_INSTRETH      = 0xC82
	CSR_HPMCOUNTER3H  = 0xC83
	CSR_HPMCOUNTER4H  = 0xC84
	CSR_HPMCOUNTER5H  = 0xC85
	CSR_HPMCOUNTER6H  = 0xC86
	CSR_HPMCOUNTER7H  = 0xC87
	CSR_HPMCOUNTER8H  = 0xC88
	CSR_HPMCOUNTER9H  = 0xC89
	CSR_HPMCOUNTER10H = 0xC8A
	CSR_HPMCOUNTER11H = 0xC8B
	CSR_HPMCOUNTER12H = 0xC8C
	CSR_HPMCOUNTER13H = 0xC8D
	CSR_HPMCOUNTER14H = 0xC8E
	CSR_HPMCOUNTER15H = 0xC8F
	CSR_HPMCOUNTER16H = 0xC90
	CSR_HPMCOUNTER17H = 0xC91
	CSR_HPMCOUNTER18H = 0xC92
	CSR_HPMCOUNTER19H = 0xC93
	CSR_HPMCOUNTER20H = 0xC94
	CSR_HPMCOUNTER21H = 0xC95
	CSR_HPMCOUNTER22H = 0xC96
	CSR_HPMCOUNTER23H = 0xC97
	CSR_HPMCOUNTER24H = 0xC98
	CSR_HPMCOUNTER25H = 0xC99
	CSR_HPMCOUNTER26H = 0xC9A
	CSR_HPMCOUNTER27H = 0xC9B
	CSR_HPMCOUNTER28H = 0xC9C
	CSR_HPMCOUNTER29H = 0xC9D
	CSR_HPMCOUNTER30H = 0xC9E
	CSR_HPMCOUNTER31H = 0xC9F
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
		// illegal instruction?
		panic("Illegal instruction (tried to access csr that was not user mode)")
	}
}

func (c *Core) csrrs(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		panic("Illegal access")
	}

	if rs1 != REG_ZERO {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= c.reg[rs1]
			c.reg[rd] = old
		} else {
			panic("Illegal instruction (tried to access csr that was not user mode)")
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
		panic("Illegal access")
	}

	if rs1 != REG_ZERO {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (c.reg[rs1] ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			panic("Illegal instruction (tried to access csr that was not user mode)")
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
		panic("Illegal instruction (tried to access csr that was not user mode)")
	}
}

func (c *Core) csrrsi(inst uint32) {
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		panic("Illegal access")
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= imm4_0
			c.reg[rd] = old
		} else {
			panic("Illegal instruction (tried to access csr that was not user mode)")
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
		panic("Illegal access")
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (imm4_0 ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			panic("Illegal instruction (tried to access csr that was not user mode)")
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}
