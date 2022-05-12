// This file contains implementations of the instructions specified in
// the Zicsr extension of the RISC-V unprivileged specification.
//   Refer to the specification for instruction documentation.

package cpu

type Csr int

// Allocated Unprivileged CSR addresses
const (
	// Unprivileged Floating-Point CSRs
	csr_FFLAGS Csr = 0x001
	csr_FRM        = 0x002
	csr_FCSR       = 0x003

	// --- Machine information registers ---
	// Csr_MVENDORID  = 0xF11
	// Csr_MARCHID    = 0xF12
	// Csr_MIMPID     = 0xF13
	Csr_MHARTID = 0xF14
	// Csr_MCONFIGPTR = 0xF15

	// --- Machine trap handling ---
	// Csr_MSCRATCH = 0x340
	Csr_MEPC   = 0x341
	Csr_MCAUSE = 0x342
	Csr_MTVAL  = 0x343
	// Csr_MIP      = 0x344
	// Csr_MTINST   = 0x34A
	// Csr_MTVAL2   = 0x34B

	// --- Supervisor Protection and Translation
	Csr_SATP = 0x180
)

func (c *Core) csrrw(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
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
		return
	}
}

func (c *Core) csrrs(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
		return
	}

	if rs1 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= c.reg[rs1]
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
			return
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrc(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
		return
	}

	if rs1 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (c.reg[rs1] ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
			return
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrwi(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
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
		return
	}
}

func (c *Core) csrrsi(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
		return
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] |= imm4_0
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
			return
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}

func (c *Core) csrrci(inst uint32) {
	if !xZicsrEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	imm4_0 := (inst >> 15) & 0x1f
	csr := (inst >> 20) & 0xfff

	if csr&0xC00 != 0 { // check if user mode register
		c.trap(TrapIllegalInstruction)
		return
	}

	if imm4_0 != 0 {
		if csr&0xC00 != 0xC00 { // verify read and write permissions
			old := c.csr[csr]
			c.csr[csr] &= (imm4_0 ^ 0xFFFFFFFF) // AND with inverse of bit-pattern to unset select bits
			c.reg[rd] = old
		} else {
			c.trap(TrapIllegalInstruction)
			return
		}
	} else { // don't write csr, just read
		c.reg[rd] = c.csr[csr]
	}
}
