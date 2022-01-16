package cpu

// Allocated Unprivileged CSR addresses
const (
	// Unprivileged Floating-Point CSRs
	Csr_FFLAGS uint32 = 0x001
	Csr_FRM           = 0x002
	Csr_FCSR          = 0x003

	// --- Unprivileged counters/timers ---
	Csr_CYCLE   = 0xC00
	Csr_TIME    = 0xC01
	Csr_INSTRET = 0xC02
	// Csr_HPMCOUNTER3  = 0xC03
	// Csr_HPMCOUNTER4  = 0xC04
	// Csr_HPMCOUNTER5  = 0xC05
	// Csr_HPMCOUNTER6  = 0xC06
	// Csr_HPMCOUNTER7  = 0xC07
	// Csr_HPMCOUNTER8  = 0xC08
	// Csr_HPMCOUNTER9  = 0xC09
	// Csr_HPMCOUNTER10 = 0xC0A
	// Csr_HPMCOUNTER11 = 0xC0B
	// Csr_HPMCOUNTER12 = 0xC0C
	// Csr_HPMCOUNTER13 = 0xC0D
	// Csr_HPMCOUNTER14 = 0xC0E
	// Csr_HPMCOUNTER15 = 0xC0F
	// Csr_HPMCOUNTER16 = 0xC10
	// Csr_HPMCOUNTER17 = 0xC11
	// Csr_HPMCOUNTER18 = 0xC12
	// Csr_HPMCOUNTER19 = 0xC13
	// Csr_HPMCOUNTER20 = 0xC14
	// Csr_HPMCOUNTER21 = 0xC15
	// Csr_HPMCOUNTER22 = 0xC16
	// Csr_HPMCOUNTER23 = 0xC17
	// Csr_HPMCOUNTER24 = 0xC18
	// Csr_HPMCOUNTER25 = 0xC19
	// Csr_HPMCOUNTER26 = 0xC1A
	// Csr_HPMCOUNTER27 = 0xC1B
	// Csr_HPMCOUNTER28 = 0xC1C
	// Csr_HPMCOUNTER29 = 0xC1D
	// Csr_HPMCOUNTER30 = 0xC1E
	// Csr_HPMCOUNTER31 = 0xC1F

	// --- Upper bits of unprivileged counters/timers ---
	Csr_CYCLEH   = 0xC80
	Csr_TIMEH    = 0xC81
	Csr_INSTRETH = 0xC82
	// Csr_HPMCOUNTER3H  = 0xC83
	// Csr_HPMCOUNTER4H  = 0xC84
	// Csr_HPMCOUNTER5H  = 0xC85
	// Csr_HPMCOUNTER6H  = 0xC86
	// Csr_HPMCOUNTER7H  = 0xC87
	// Csr_HPMCOUNTER8H  = 0xC88
	// Csr_HPMCOUNTER9H  = 0xC89
	// Csr_HPMCOUNTER10H = 0xC8A
	// Csr_HPMCOUNTER11H = 0xC8B
	// Csr_HPMCOUNTER12H = 0xC8C
	// Csr_HPMCOUNTER13H = 0xC8D
	// Csr_HPMCOUNTER14H = 0xC8E
	// Csr_HPMCOUNTER15H = 0xC8F
	// Csr_HPMCOUNTER16H = 0xC90
	// Csr_HPMCOUNTER17H = 0xC91
	// Csr_HPMCOUNTER18H = 0xC92
	// Csr_HPMCOUNTER19H = 0xC93
	// Csr_HPMCOUNTER20H = 0xC94
	// Csr_HPMCOUNTER21H = 0xC95
	// Csr_HPMCOUNTER22H = 0xC96
	// Csr_HPMCOUNTER23H = 0xC97
	// Csr_HPMCOUNTER24H = 0xC98
	// Csr_HPMCOUNTER25H = 0xC99
	// Csr_HPMCOUNTER26H = 0xC9A
	// Csr_HPMCOUNTER27H = 0xC9B
	// Csr_HPMCOUNTER28H = 0xC9C
	// Csr_HPMCOUNTER29H = 0xC9D
	// Csr_HPMCOUNTER30H = 0xC9E
	// Csr_HPMCOUNTER31H = 0xC9F

	// --- Machine information registers ---
	// Csr_MVENDORID  = 0xF11
	// Csr_MARCHID    = 0xF12
	// Csr_MIMPID     = 0xF13
	Csr_MHARTID = 0xF14
	// Csr_MCONFIGPTR = 0xF15

	// --- Machine trap setup ---
	// Csr_MSTATUS    = 0x300
	// Csr_MISA       = 0x301
	// Csr_MEDELEG    = 0x302
	// Csr_MIDELEG    = 0x303
	// Csr_MIE        = 0x304
	// Csr_MTVEC      = 0x305
	// Csr_MCOUNTEREN = 0x306
	// Csr_MSTATUSH   = 0x310

	// --- Machine trap handling ---
	// Csr_MSCRATCH = 0x340
	Csr_MEPC   = 0x341
	Csr_MCAUSE = 0x342
	Csr_MTVAL  = 0x343
	// Csr_MIP      = 0x344
	// Csr_MTINST   = 0x34A
	// Csr_MTVAL2   = 0x34B

	// --- Machine configuration---
	// -- unused --

	// --- Machine memory protection ---
	// -- unused --
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

	if rs1 != Reg_ZERO {
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

	if rs1 != Reg_ZERO {
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
