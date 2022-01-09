package cpu

// TODO how to support the different rounding modes?
// Go has math/big/float which supports the IEEE rounding modes.
// Find out how this works

const (
	FREG_FT0  = 0  // FP temporaries
	FREG_FT1  = 1  //
	FREG_FT2  = 2  //
	FREG_FT3  = 3  //
	FREG_FT4  = 4  //
	FREG_FT5  = 5  //
	FREG_FT6  = 6  //
	FREG_FT7  = 7  //
	FREG_FS0  = 8  // FP saved registers
	FREG_FS1  = 9  //
	FREG_FA0  = 10 // FP arguments/return values
	FREG_FA1  = 11 //
	FREG_FA2  = 12 // FP arguments
	FREG_FA3  = 13 //
	FREG_FA4  = 14 //
	FREG_FA5  = 15 //
	FREG_FA6  = 16 //
	FREG_FA7  = 17 //
	FREG_FS2  = 18 // FP saved registers
	FREG_FS3  = 19 //
	FREG_FS4  = 20 //
	FREG_FS5  = 21 //
	FREG_FS6  = 22 //
	FREG_FS7  = 23 //
	FREG_FS8  = 24 //
	FREG_FS9  = 25 //
	FREG_FS10 = 26 //
	FREG_FS11 = 27 //
	FREG_FT8  = 28 // FP temporaries
	FREG_FT9  = 29 //
	FREG_FT10 = 30 //
	FREG_FT11 = 31 //
)

func (c *Core) flw(inst uint32) {
	// TODO
}

func (c *Core) fsw(inst uint32) {
	// TODO
}

func (c *Core) fmadd_s(inst uint32) {
	// TODO
}

func (c *Core) fmsub_s(inst uint32) {
	// TODO
}

func (c *Core) fnmsub_s(inst uint32) {
	// TODO
}

func (c *Core) fnmadd_s(inst uint32) {
	// TODO
}

func (c *Core) fadd_s(inst uint32) {
	// TODO
}

func (c *Core) fsub_s(inst uint32) {
	// TODO
}

func (c *Core) fmul_s(inst uint32) {
	// TODO
}

func (c *Core) fdiv_s(inst uint32) {
	// TODO
}

func (c *Core) fsqrt_s(inst uint32) {
	// TODO
}

func (c *Core) fsgnj_s(inst uint32) {
	// TODO
}

func (c *Core) fsgnjn_s(inst uint32) {
	// TODO
}

func (c *Core) fsgnjx_s(inst uint32) {
	// TODO
}

func (c *Core) fmin_s(inst uint32) {
	// TODO
}

func (c *Core) fmax_s(inst uint32) {
	// TODO
}

func (c *Core) fcvt_w_s(inst uint32) {
	// TODO
}

func (c *Core) fcvt_wu_s(inst uint32) {
	// TODO
}

func (c *Core) fmv_x_w(inst uint32) {
	// TODO
}

func (c *Core) feq_s(inst uint32) {
	// TODO
}

func (c *Core) flt_s(inst uint32) {
	// TODO
}

func (c *Core) fle_s(inst uint32) {
	// TODO
}

func (c *Core) fclass_s(inst uint32) {
	// TODO
}

func (c *Core) fcvt_s_w(inst uint32) {
	// TODO
}

func (c *Core) fcvt_s_wu(inst uint32) {
	// TODO
}

func (c *Core) fmv_w_x(inst uint32) {
	// TODO
}
