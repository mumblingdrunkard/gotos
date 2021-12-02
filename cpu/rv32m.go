// RV32M Standard Extension for Integer Multiply and Divide

package cpu

//
// Multiply
//

func (cpu *CPU) mul(rd, rs1, rs2 uint32) {
	// TODO
	// multiply
}

func (cpu *CPU) mulh(rd, rs1, rs2 uint32) {
	// TODO
	// multiply high signed signed
}

func (cpu *CPU) mulhsu(rd, rs1, rs2 uint32) {
	// TODO
	// multiply high signed unsigned
}

func (cpu *CPU) mulhu(rd, rs1, rs2 uint32) {
	// TODO
	// multiply high unsigned unsigned
}

//
// Divide
//

func (cpu *CPU) div(rd, rs1, rs2 uint32) {
	// TODO
	// divide signed
}

func (cpu *CPU) divu(rd, rs1, rs2 uint32) {
	// TODO
	// divide unsigned
}

//
// Remainder
//

func (cpu *CPU) rem(rd, rs1, rs2 uint32) {
	// TODO
	// remainder signed
}

func (cpu *CPU) remu(rd, rs1, rs2 uint32) {
	// TODO
	// remainder unsigned
}
