// This file contains implementations of the instructions specified in
// the F extension of the RISC-V unprivileged specification.
//   Refer to the specification for instruction documentation.

// WARNING This implementation is _NOT_ compliant with the RISC-V
// specification.
//   A compliant implementation can likely be made, but it is not a
// priority.
//   User/application mode programs should not depend on IEEE compliant
// floating point numbers.
//   A compliant implementation would likely also be much, much slower.
//   Programs should also not depend on the state of the FCSR containing
// correct exception flags.

package cpu

import (
	"math"
)

type FReg int

// Mnemonics for Floating-Point registers
const (
	// Unnamed
	FReg_F0  FReg = 0
	FReg_F1       = 1
	FReg_F2       = 2
	FReg_F3       = 3
	FReg_F4       = 4
	FReg_F5       = 5
	FReg_F6       = 6
	FReg_F7       = 7
	FReg_F8       = 8
	FReg_F9       = 9
	FReg_F10      = 10
	FReg_F11      = 11
	FReg_F12      = 12
	FReg_F13      = 13
	FReg_F14      = 14
	FReg_F15      = 15
	FReg_F16      = 16
	FReg_F17      = 17
	FReg_F18      = 18
	FReg_F19      = 19
	FReg_F20      = 20
	FReg_F21      = 21
	FReg_F22      = 22
	FReg_F23      = 23
	FReg_F24      = 24
	FReg_F25      = 25
	FReg_F26      = 26
	FReg_F27      = 27
	FReg_F28      = 28
	FReg_F29      = 29
	FReg_F30      = 30
	FReg_F31      = 31

	// ABI Names
	FReg_FT0  = 0  // FP temporaries
	FReg_FT1  = 1  //
	FReg_FT2  = 2  //
	FReg_FT3  = 3  //
	FReg_FT4  = 4  //
	FReg_FT5  = 5  //
	FReg_FT6  = 6  //
	FReg_FT7  = 7  //
	FReg_FS0  = 8  // FP saved registers
	FReg_FS1  = 9  //
	FReg_FA0  = 10 // FP arguments/return values
	FReg_FA1  = 11 //
	FReg_FA2  = 12 // FP arguments
	FReg_FA3  = 13 //
	FReg_FA4  = 14 //
	FReg_FA5  = 15 //
	FReg_FA6  = 16 //
	FReg_FA7  = 17 //
	FReg_FS2  = 18 // FP saved registers
	FReg_FS3  = 19 //
	FReg_FS4  = 20 //
	FReg_FS5  = 21 //
	FReg_FS6  = 22 //
	FReg_FS7  = 23 //
	FReg_FS8  = 24 //
	FReg_FS9  = 25 //
	FReg_FS10 = 26 //
	FReg_FS11 = 27 //
	FReg_FT8  = 28 // FP temporaries
	FReg_FT9  = 29 //
	FReg_FT10 = 30 //
	FReg_FT11 = 31 //
)

const (
	fcsrFlagNV uint32 = 0b10000
	fcsrFlagDZ        = 0b01000
	fcsrFlagOF        = 0b00100
	fcsrFlagUF        = 0b00010
	fcsrFlagNX        = 0b00001
)

func (c *Core) flw(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f             // destination fpregister
	rs1 := (inst >> 15) & 0x1f           // base register
	imm11_0 := uint32(int32(inst) >> 20) // sign extended

	address := c.reg[rs1] + imm11_0

	if success, word := c.loadWord(address); success {
		c.freg[rd] = 0xFFFFFFFF00000000 | uint64(word)
	}
}

func (c *Core) fsw(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rs1 := (inst >> 15) & 0x1f           // base register
	rs2 := (inst >> 20) & 0x1f           // source fp register
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f

	offset := (imm11_5 << 5) | imm4_0

	address := c.reg[rs1] + offset

	c.storeWord(address, uint32(c.freg[rs2]))
}

// float multiply and add
func (c *Core) fmadd_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))
	f3 := math.Float32frombits(uint32(c.freg[rs3]))

	res := f1*f2 + f3

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmsub_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))
	f3 := math.Float32frombits(uint32(c.freg[rs3]))

	res := f1*f2 - f3

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fnmsub_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))
	f3 := math.Float32frombits(uint32(c.freg[rs3]))

	res := -f1*f2 - f3

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fnmadd_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))
	f3 := math.Float32frombits(uint32(c.freg[rs3]))

	res := -f1*f2 + f3

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fadd_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 + f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsub_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 - f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmul_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 * f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fdiv_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 / f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsqrt_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	if rs2 != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	res := float32(math.Sqrt(float64(f1)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsgnj_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := uint32(c.freg[rs1])
	src2 := uint32(c.freg[rs2])

	// NaN boxing
	res := (src1 & 0x7FFFFFFF) | (src2 & 0x80000000)
	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(res)
}

func (c *Core) fsgnjn_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := uint32(c.freg[rs1])
	src2 := uint32(c.freg[rs2])

	res := (src1 & 0x7FFFFFFF) | (src2 & 0x80000000) ^ 0x80000000 // flip the sign bit
	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(res)                 // NaN boxing
}

func (c *Core) fsgnjx_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := uint32(c.freg[rs1])
	src2 := uint32(c.freg[rs2])

	res := (src1 & 0x7FFFFFFF) | (src2 & 0x80000000) ^ (src1 & 0x80000000) // xor the sign bits
	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(res)                          // NaN boxing
}

func (c *Core) fmin_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := float32(math.Min(float64(f1), float64(f2)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmax_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := float32(math.Max(float64(f1), float64(f2)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fcvt_w_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	res := int32(f1)

	c.reg[rd] = uint32(res)
}

func (c *Core) fcvt_wu_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	c.reg[rd] = uint32(f1)
}

func (c *Core) fmv_x_w(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	rs2 := (inst >> 20) & 0x1f // source fp register
	if rs2 != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	rm := (inst >> 12) & 0x7
	if rm != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	c.reg[rd] = uint32(c.freg[rs1])
}

func (c *Core) feq_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	// quiet comparison, only signal if either input is signaling
	if math.IsNaN(float64(f1)) {
		if c.freg[rs1]&0x00400000 != 0x00400000 { // it's a signaling NaN
			c.csr[csr_FCSR] |= fcsrFlagNV
		}
		c.reg[rd] = 0
		return
	}

	if math.IsNaN(float64(f2)) {
		if c.freg[rs2]&0x00400000 != 0x00400000 { // it's a signaling NaN
			c.csr[csr_FCSR] |= fcsrFlagNV
		}
		c.reg[rd] = 0
		return
	}

	if f1 == f2 {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

func (c *Core) flt_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	if math.IsNaN(float64(f1)) {
		c.csr[csr_FCSR] |= fcsrFlagNV
		c.reg[rd] = 0
		return
	}

	if math.IsNaN(float64(f2)) {
		c.csr[csr_FCSR] |= fcsrFlagNV
		c.reg[rd] = 0
		return
	}

	if f1 < f2 {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

func (c *Core) fle_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	if math.IsNaN(float64(f1)) {
		c.csr[csr_FCSR] |= fcsrFlagNV
		c.reg[rd] = 0
		return
	}

	if math.IsNaN(float64(f2)) {
		c.csr[csr_FCSR] |= fcsrFlagNV
		c.reg[rd] = 0
		return
	}

	if f1 <= f2 {
		c.reg[rd] = 1
	} else {
		c.reg[rd] = 0
	}
}

func (c *Core) fclass_s(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	if rs2 != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	c.reg[rd] = 0

	if math.IsNaN(float64(f1)) {
		if c.freg[rs1]&0x00400000 != 0x00400000 { // it's a signaling NaN
			c.reg[rd] |= 0b0100000000 // signaling NaN
		} else {
			c.reg[rd] |= 0b1000000000 // quiet NaN
		}
	}

	if math.IsInf(float64(f1), 1) {
		c.reg[rd] |= 0b0010000000 // positive infinity
	}

	if math.IsInf(float64(f1), -1) {
		c.reg[rd] |= 0b0000000001 // positive infinity
	}

	// See table 11.5 in RISC-V spec
	// TODO subnormal detection     - 4 bits
	// TODO positive and negative 0 - 2 bits
}

func (c *Core) fcvt_s_w(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float32(int32(c.reg[rs1]))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fcvt_s_wu(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float32(c.reg[rs1])

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmv_w_x(inst uint32) {
	if !xFEnable {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	rs2 := (inst >> 20) & 0x1f // source fp register
	if rs2 != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	rm := (inst >> 12) & 0x7
	if rm != 0 {
		c.csr[Csr_MTVAL] = inst
		c.trap(TrapIllegalInstruction)
		return
	}

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(c.reg[rs1])
}
