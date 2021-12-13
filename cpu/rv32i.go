// RV32I Base Integer Instruction Set
// https://mark.theis.site/riscv/

package cpu

func (core *Core) add(rd, rs1, rs2 uint32) {
	// add
	core.reg[rd] = core.reg[rs1] + core.reg[rs2]
}

func (core *Core) addi(rd, rs1, immi uint32) {
	// add immediate
	core.reg[rd] = core.reg[rs1] + uint32(immi)
}

func (core *Core) sub(rd, rs1, rs2 uint32) {
	// subtract
	core.reg[rd] = core.reg[rs1] - core.reg[rs2]
}

func (core *Core) ecall(inst uint32) {
	// make call
}
