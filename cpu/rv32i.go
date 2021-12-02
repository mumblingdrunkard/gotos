// RV32I Base Integer Instruction Set
// https://mark.theis.site/riscv/

package cpu

//
// Branching
//

func (cpu *CPU) auipc(rd, unknown uint32) {
	// TODO
	// add upper immediate to pc
}

func (cpu *CPU) jal(rd, unknown uint32) {
	// TODO
	// jump and link
}

func (cpu *CPU) jalr(rd, rs1, unknown uint32) {
	// TODO
	// jump and link register
}

func (cpu *CPU) beq(rs1, rs2, unknown uint32) {
	// TODO
	// branch equal
}

func (cpu *CPU) blt(rs1, rs2, unknown uint32) {
	// TODO
	// branch equal
}

func (cpu *CPU) bge(rs1, rs2, unknown uint32) {
	// TODO
	// branch equal
}

func (cpu *CPU) bltu(rs1, rs2, unknown uint32) {
	// TODO
	// branch equal
}

func (cpu *CPU) bgeu(rs1, rs2, unknown uint32) {
	// TODO
	// branch equal
}

//
// Load
//

func (cpu *CPU) lui(rd, unknown uint32) {
	// TODO
	// load upper immediate
}

func (cpu *CPU) lb(rd, rs1, immi uint32) {
	// TODO
	// load byte
}

func (cpu *CPU) lh(rd, rs1, immi uint32) {
	// TODO
	// load half
}

func (cpu *CPU) lw(rd, rs1, immi uint32) {
	// TODO
	// load word
}

func (cpu *CPU) lbu(rd, rs1, immi uint32) {
	// TODO
	// load byte unsigned
}

func (cpu *CPU) lhu(rd, rs1, immi uint32) {
	// TODO
	// load half unsigned
}

//
// Store
//

func (cpu *CPU) sb(rs1, rs2, imms uint32) {
	// TODO
	// store byte
}

func (cpu *CPU) sh(rs1, rs2, imms uint32) {
	// TODO
	// store half
}

func (cpu *CPU) sw(rs1, rs2, imms uint32) {
	// TODO
	// store word
}

//
// Integer arithmetic
//

func (cpu *CPU) add(rd, rs1, rs2 uint32) {
	// add
	cpu.reg[rd] = cpu.reg[rs1] + cpu.reg[rs2]
}

func (cpu *CPU) addi(rd, rs1, immi uint32) {
	// add immediate
	cpu.reg[rd] = cpu.reg[rs1] + uint64(immi)
}

func (cpu *CPU) sub(rd, rs1, rs2 uint32) {
	// subtract
	cpu.reg[rd] = cpu.reg[rs1] - cpu.reg[rs2]
}

//
// Bitwise
//

func (cpu *CPU) xori(rd, rs1, immi uint32) {
	// TODO
	// xor immediate
}

func (cpu *CPU) ori(rd, rs1, immi uint32) {
	// TODO
	// or immediate
}

func (cpu *CPU) andi(rd, rs1, immi uint32) {
	// TODO
	// and immediate
}

func (cpu *CPU) slli(rd, rs1, immi uint32) {
	// TODO
	// shift left logical immediate
}

func (cpu *CPU) srli(rd, rs1, immi uint32) {
	// TODO
	// shift right logical immediate
}

func (cpu *CPU) srai(rd, rs1, immi uint32) {
	// TODO
	// shift right arithmetic immediate
}

func (cpu *CPU) sll(rd, rs1, rs2 uint32) {
	// TODO
	// shift left logical
}

func (cpu *CPU) srl(rd, rs1, rs2 uint32) {
	// TODO
	// shift right logical
}

func (cpu *CPU) sra(rd, rs1, rs2 uint32) {
	// TODO
	// shift right arithmetic
}

func (cpu *CPU) xor(rd, rs1, rs2 uint32) {
	// TODO
	// xor
}

func (cpu *CPU) or(rd, rs1, rs2 uint32) {
	// TODO
	// or
}

func (cpu *CPU) and(rd, rs1, rs2 uint32) {
	// TODO
	// and
}

//
// Unplaced
//

func (cpu *CPU) slti(rd, rs1, unknown uint32) {
	// TODO
	// set less than immediate
}

func (cpu *CPU) sltiu(rd, rs1, unknown uint32) {
	// TODO
	// set less than immediate unsigned
}

func (cpu *CPU) slt(rd, rs1, rs2 uint32) {
	// TODO
	// set less than
}

func (cpu *CPU) sltu(rd, rs1, rs2 uint32) {
	// TODO
	// set less than unsigned
}

func (cpu *CPU) fence(unknown1, unknown2 uint32) {
	// TODO
	// fence
}

func (cpu *CPU) fence_i() {
	// TODO
	// fence instruction
}
