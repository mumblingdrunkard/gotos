package cpu

import "math"

// WARNING This implementation is _NOT_ compliant with the RISC-V specification.
// A compliant implementation can likely be made, but it is not a priority.
// User/application mode programs should not depend on IEEE compliant floating point numbers.
// A compliant implementation would likely also be much, much slower.
// Programs should also not depend on the state of the FCSR containing correct exception flags.

func (c *Core) fld(inst uint32) {
	// TODO
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	imm11_0 := uint32(int32(inst) >> 20)

	err, val := c.mc.LoadDoubleWord(c.reg[rs1] + imm11_0)

	if err != nil {
		panic(err)
	}

	c.freg[rd] = val
}

func (c *Core) fsd(inst uint32) {
	// TODO
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	imm11_5 := uint32(int32(inst) >> 25)
	imm4_0 := (inst >> 7) & 0x1f

	addr := c.reg[rs1] + ((imm11_5 << 5) | imm4_0)

	err := c.mc.StoreDoubleWord(addr, c.freg[rs2])

	if err != nil {
		panic(err)
	}
}

func (c *Core) fmadd_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])
	f3 := math.Float64frombits(c.freg[rs3])

	res := f1*f2 + f3

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fmsub_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])
	f3 := math.Float64frombits(c.freg[rs3])

	res := f1*f2 - f3

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fnmsub_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])
	f3 := math.Float64frombits(c.freg[rs3])

	res := -f1*f2 - f3

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fnmadd_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f
	rs3 := (inst >> 27) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])
	f3 := math.Float64frombits(c.freg[rs3])

	res := -f1*f2 + f3

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fadd_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := f1 + f2

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fsub_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := f1 - f2

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fmul_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := f1 * f2

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fdiv_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := f1 / f2

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fsqrt_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	if rs2 != 0 {
		panic("Illegal instruction")
	}

	f1 := math.Float64frombits(c.freg[rs1])

	res := f1

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fsgnj_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := c.freg[rs1]
	src2 := c.freg[rs2]

	res := (src1 & 0x7FFFFFFFFFFFFFFF) | (src2 & 0x8000000000000000)
	c.freg[rd] = res

}

func (c *Core) fsgnjn_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := c.freg[rs1]
	src2 := c.freg[rs2]

	res := (src1 & 0x7FFFFFFFFFFFFFFF) | (src2 & 0x8000000000000000) ^ 0x8000000000000000 // inverts the sign bit
	c.freg[rd] = res
}

func (c *Core) fsgnjx_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	src1 := c.freg[rs1]
	src2 := c.freg[rs2]

	res := (src1 & 0x7FFFFFFFFFFFFFFF) | (src2 & 0x8000000000000000) ^ (src1 & 0x8000000000000000) // xor the sign bits
	c.freg[rd] = res
}

func (c *Core) fmin_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := math.Min(f1, f2)

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fmax_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	res := math.Max(f1, f2)

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fcvt_s_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	d1 := math.Float64frombits(c.freg[rs1])
	f1 := float32(d1)

	c.freg[rd] = 0xFFFFFFFF00000000 | uint64(math.Float32bits(f1))
}

func (c *Core) fcvt_d_s(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float32frombits(uint32(c.freg[rs1]))
	d1 := float64(f1)

	c.freg[rd] = math.Float64bits(d1)
}

func (c *Core) feq_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	c.reg[rd] = 0

	if math.IsNaN(f1) {
		c.csr[CSR_FCSR] |= FCSR_F_NV // fuck it, always signaling
		return
	}

	if math.IsNaN(f2) {
		c.csr[CSR_FCSR] |= FCSR_F_NV
		return
	}

	if f1 == f2 {
		c.reg[rd] = 1
	}
}

func (c *Core) flt_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	c.reg[rd] = 0

	if math.IsNaN(f1) {
		c.csr[CSR_FCSR] |= FCSR_F_NV // fuck it, always signaling
		return
	}

	if math.IsNaN(f2) {
		c.csr[CSR_FCSR] |= FCSR_F_NV
		return
	}

	if f1 < f2 {
		c.reg[rd] = 1
	}
}

func (c *Core) fle_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	f1 := math.Float64frombits(c.freg[rs1])
	f2 := math.Float64frombits(c.freg[rs2])

	c.reg[rd] = 0

	if math.IsNaN(f1) {
		c.csr[CSR_FCSR] |= FCSR_F_NV // fuck it, always signaling
		return
	}

	if math.IsNaN(f2) {
		c.csr[CSR_FCSR] |= FCSR_F_NV
		return
	}

	if f1 <= f2 {
		c.reg[rd] = 1
	}
}

func (c *Core) fclass_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f // source fp register

	if rs2 != 0 {
		panic("Illegal instruction")
	}

	f1 := math.Float64frombits(c.freg[rs1])

	c.reg[rd] = 0

	if math.IsNaN(float64(f1)) {
		// TODO find out how to check for signaling nan in f64
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

func (c *Core) fcvt_w_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])

	res := int32(f1)

	c.reg[rd] = uint32(res)
}

func (c *Core) fcvt_wu_d(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	f1 := math.Float64frombits(c.freg[rs1])

	res := uint32(f1)

	c.reg[rd] = res
}

func (c *Core) fcvt_d_w(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float64(int32(c.reg[rs1]))

	c.freg[rd] = math.Float64bits(res)
}

func (c *Core) fcvt_d_wu(inst uint32) {
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f

	res := float64(c.reg[rs1])

	c.freg[rd] = math.Float64bits(res)
}
