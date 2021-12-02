// RV64I Base Integer Instruction Set (addition to RV32I)

package cpu

//
// Load
//

func (cpu *CPU) lwu(rd, rs1, immi uint32) {
	// TODO
	// load word unsigned
}

func (cpu *CPU) ld(rd, rs1, immi uint32) {
	// TODO
	// load double
}

//
// Store
//

func (cpu *CPU) sd(rs1, rs2, imms uint32) {
	// TODO
	// store double
}

//
// Integer Arithmetic
//

func (cpu *CPU) addw(rd, rs1, rs2 uint32) {
	// add word
	cpu.reg[rd] = cpu.reg[rs1] + cpu.reg[rs2]
}

func (cpu *CPU) addiw(rd, rs1, immi uint32) {
	// add immediate word
	cpu.reg[rd] = cpu.reg[rs1] + uint64(immi)
}

func (cpu *CPU) subw(rd, rs1, rs2 uint32) {
	// subtract word
	cpu.reg[rd] = cpu.reg[rs1] - cpu.reg[rs2]
}

//
// Bitwise
//

func (cpu *CPU) slli_s(rd, rs1, immi uint32) {
	// TODO
	// shift left logical immediate signed
}

func (cpu *CPU) srli_s(rd, rs1, immi uint32) {
	// TODO
	// shift right logical immediate signed
}

func (cpu *CPU) srai_s(rd, rs1, immi uint32) {
	// TODO
	// shift right arithmetic immediate signed
}

func (cpu *CPU) slliw(rd, rs1, immi uint32) {
	// TODO
	// shift left logical immediate word
}

func (cpu *CPU) srliw(rd, rs1, immi uint32) {
	// TODO
	// shift right logical immediate word
}

func (cpu *CPU) sraiw(rd, rs1, immi uint32) {
	// TODO
	// shift right arithmetic immediate word
}

func (cpu *CPU) sllw(rd, rs1, rs2 uint32) {
	// TODO
	// shift left logical word
}

func (cpu *CPU) srlw(rd, rs1, rs2 uint32) {
	// TODO
	// shift right logical word
}

func (cpu *CPU) sraw(rd, rs1, rs2 uint32) {
	// TODO
	// shift right arithmetic word
}
