// RV64M Standard Extension for Integer Multiply and Divide (addition to RV32M)

package cpu

//
// Multiply
//

func (cpu *CPU) mulw(rd, rs1, rs2 uint32) {
	// TODO
	// multiply word
}

//
// Divide
//

func (cpu *CPU) divw(rd, rs1, rs2 uint32) {
	// TODO
	// divide signed word
}

func (cpu *CPU) divuw(rd, rs1, rs2 uint32) {
	// TODO
	// divide unsigned word
}

//
// Remainder
//

func (cpu *CPU) remw(rd, rs1, rs2 uint32) {
	// TODO
	// remainder signed word
}

func (cpu *CPU) remuw(rd, rs1, rs2 uint32) {
	// TODO
	// remainder unsigned word
}
