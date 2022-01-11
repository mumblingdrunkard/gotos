package cpu

// WARNING This implementation is _NOT_ compliant with the RISC-V specification.
// A compliant implementation can likely be made, but it is not a priority.
// Programs should not depend on IEEE compliant floating point numbers.
// A compliant implementation would likely also be much, much slower.
// Programs should also not depend on the state of the FCSR containing correct exception flags.

import (
	"math"
)

type FRegisterNumber uint8

// Mnemonics for Floating-Point registers
const (
	FRegFT0  FRegisterNumber = 0  // FP temporaries
	FRegFT1                  = 1  //
	FRegFT2                  = 2  //
	FRegFT3                  = 3  //
	FRegFT4                  = 4  //
	FRegFT5                  = 5  //
	FRegFT6                  = 6  //
	FRegFT7                  = 7  //
	FRegFS0                  = 8  // FP saved registers
	FRegFS1                  = 9  //
	FRegFA0                  = 10 // FP arguments/return values
	FRegFA1                  = 11 //
	FRegFA2                  = 12 // FP arguments
	FRegFA3                  = 13 //
	FRegFA4                  = 14 //
	FRegFA5                  = 15 //
	FRegFA6                  = 16 //
	FRegFA7                  = 17 //
	FRegFS2                  = 18 // FP saved registers
	FRegFS3                  = 19 //
	FRegFS4                  = 20 //
	FRegFS5                  = 21 //
	FRegFS6                  = 22 //
	FRegFS7                  = 23 //
	FRegFS8                  = 24 //
	FRegFS9                  = 25 //
	FRegFS10                 = 26 //
	FRegFS11                 = 27 //
	FRegFT8                  = 28 // FP temporaries
	FRegFT9                  = 29 //
	FRegFT10                 = 30 //
	FRegFT11                 = 31 //
)

const (
	fcsrFlagNV uint32 = 0b10000
	fcsrFlagDZ        = 0b01000
	fcsrFlagOF        = 0b00100
	fcsrFlagUF        = 0b00010
	fcsrFlagNX        = 0b00001
)

func (c *Core) flw(inst uint32) {
	rd := (inst >> 7) & 0x1f             // destination fpregister
	rs1 := (inst >> 15) & 0x1f           // base register
	imm11_0 := uint32(int32(inst) >> 20) // sign extended

	address := c.reg[rs1] + imm11_0

	success, word := c.loadWord(address)

	if !success {
		c.DumpRegisters()
		panic("Failed")
	}

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(word)
}

func (c *Core) fsw(inst uint32) {
	rs1 := (inst >> 15) & 0x1f           // base register
	rs2 := (inst >> 20) & 0x1f           // source fp register
	imm11_5 := uint32(int32(inst) >> 25) // sign extended
	imm4_0 := (inst >> 7) & 0x1f

	offset := (imm11_5 << 5) | imm4_0

	address := c.reg[rs1] + offset

	success := c.storeWord(address, uint32(c.freg[rs2]))

	if !success {
		c.DumpRegisters()
		panic("Failed")
	}
}

// float multiply and add
func (c *Core) fmadd_s(inst uint32) {
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
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 + f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsub_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 - f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmul_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 * f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fdiv_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := f1 / f2

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsqrt_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	if rs2 != 0 {
		panic("Illegal instruction")
	}

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	res := float32(math.Sqrt(float64(f1)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fsgnj_s(inst uint32) {
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
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := uint32(c.freg[rs1])
	src2 := uint32(c.freg[rs2])

	res := (src1 & 0x7FFFFFFF) | (src2 & 0x80000000) ^ 0x80000000 // flip the sign bit
	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(res)                 // NaN boxing
}

func (c *Core) fsgnjx_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := uint32(c.freg[rs1])
	src2 := uint32(c.freg[rs2])

	res := (src1 & 0x7FFFFFFF) | (src2 & 0x80000000) ^ (src1 & 0x80000000) // xor the sign bits
	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(res)                          // NaN boxing
}

func (c *Core) fmin_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := float32(math.Min(float64(f1), float64(f2)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmax_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	f2 := math.Float32frombits(uint32(c.freg[rs2]))

	res := float32(math.Max(float64(f1), float64(f2)))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fcvt_w_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	res := int32(f1)

	c.reg[rd] = uint32(res)
}

func (c *Core) fcvt_wu_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))

	c.reg[rd] = uint32(f1)
}

func (c *Core) fmv_x_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	rs2 := (inst >> 20) & 0x1f // source fp register
	if rs2 != 0 {
		panic("Illegal instruction")
	}

	rm := (inst >> 12) & 0x7
	if rm != 0 {
		panic("Illegal instruction")
	}

	c.reg[rd] = uint32(c.freg[rs1])
}

func (c *Core) feq_s(inst uint32) {
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
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	if rs2 != 0 {
		panic("Illegal instruction")
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
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float32(int32(c.reg[rs1]))

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fcvt_s_wu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float32(c.reg[rs1])

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(res))
}

func (c *Core) fmv_w_x(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	rs2 := (inst >> 20) & 0x1f // source fp register
	if rs2 != 0 {
		panic("Illegal instruction")
	}

	rm := (inst >> 12) & 0x7
	if rm != 0 {
		panic("Illegal instruction")
	}

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(c.reg[rs1])
}
