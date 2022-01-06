package cpu

// TODO: Exceptions

// TODO: Fix decoding for OP
// TODO: Add decoding for RV32M extension

const (
	OP_IMM   uint32 = 0b0010011
	LUI             = 0b0110111
	AUIPC           = 0b0010111
	OP              = 0b0110011
	JAL             = 0b1101111
	JALR            = 0b1100111
	BRANCH          = 0b1100011
	LOAD            = 0b0000011
	STORE           = 0b0100011
	MISC_MEM        = 0b0001111
	SYSTEM          = 0b1110011
	AMO             = 0b0101111
)

func (c *Core) execute(inst uint32) {
	// Register 0 is hardwired with all 0s have to reset to 0 for every
	// cycle because some instructions may use this as their /dev/null
	c.reg[0] = 0

	opcode := inst & 0x7f

	switch opcode {
	case OP:
		c.op(inst)
	case OP_IMM:
		c.op_imm(inst)
	case LUI:
		c.lui(inst)
	case AUIPC:
		c.auipc(inst)
	case JAL:
		c.jal(inst)
	case JALR:
		c.jalr(inst)
	case BRANCH:
		c.branch(inst)
	case LOAD:
		c.load(inst)
	case STORE:
		c.store(inst)
	case MISC_MEM:
		c.misc_mem(inst)
	case SYSTEM:
		c.system(inst)
	case AMO:
		c.amo(inst)
	default:
		panic("Unknown instruction")
	}
}

// OP funct7
const (
	OP_A   uint32 = 0b0000000
	OP_B          = 0b0100000
	MULDIV        = 0b0000001
)

func (c *Core) op(inst uint32) {
	funct7 := (inst >> 25) & 0x7f
	switch funct7 {
	case OP_A:
		c.op_a(inst)
	case OP_B:
		c.op_b(inst)
	case MULDIV:
		c.muldiv(inst)
	default:
		panic("Illegal instruction format")
	}
}

// OP_A funct3
const (
	ADD  uint32 = 0b000
	SLL         = 0b001
	SLT         = 0b010
	SLTU        = 0b011
	XOR         = 0b100
	SRL         = 0b101
	OR          = 0b110
	AND         = 0b111
)

func (c *Core) op_a(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case ADD:
		c.add(inst)
	case SLL:
		c.sll(inst)
	case SLT:
		c.slt(inst)
	case SLTU:
		c.sltu(inst)
	case XOR:
		c.xor(inst)
	case SRL:
		c.srl(inst)
	case OR:
		c.or(inst)
	case AND:
		c.and(inst)
	default:
		panic("Unknown instruction")
	}
}

const (
	SUB uint32 = 0b000
	SRA        = 0b101
)

func (c *Core) op_b(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case SUB:
		c.sub(inst)
	case SRA:
		c.sra(inst)
	default:
		panic("Unknown instruction")
	}
}

const (
	MUL    uint32 = 0b000
	MULH          = 0b001
	MULHSU        = 0b010
	MULHU         = 0b011
	DIV           = 0b100
	DIVU          = 0b101
	REM           = 0b110
	REMU          = 0b111
)

func (c *Core) muldiv(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case MUL:
		c.mul(inst)
	case MULH:
		c.mulh(inst)
	case MULHSU:
		c.mulhsu(inst)
	case MULHU:
		c.mulhu(inst)
	case DIV:
		c.div(inst)
	case DIVU:
		c.divu(inst)
	case REM:
		c.rem(inst)
	case REMU:
		c.rem(inst)
	default:
		panic("Unknown instruction")
	}
}

// op-imm funct3
const (
	ADDI  uint32 = 0b000
	SLTI         = 0b010
	SLTIU        = 0b011
	XORI         = 0b100
	ORI          = 0b110
	ANDI         = 0b111
	SLLI         = 0b001
	SRLI         = 0b101
)

func (c *Core) op_imm(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case ADDI:
		c.addi(inst)
	case SLTI:
		c.slti(inst)
	case SLTIU:
		c.sltiu(inst)
	case XORI:
		c.xori(inst)
	case ORI:
		c.ori(inst)
	case ANDI:
		c.andi(inst)
	case SLLI:
		c.slli(inst)
	case SRLI:
		c.srli(inst)
	default:
		panic("Illegal instruction format")
	}
}

// branch funct3
const (
	BEQ  uint32 = 0b000
	BNE         = 0b001
	BLT         = 0b100
	BGE         = 0b101
	BLTU        = 0b110
	BGEU        = 0b111
)

func (c *Core) branch(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case BEQ:
		c.beq(inst)
	case BNE:
		c.bne(inst)
	case BLT:
		c.blt(inst)
	case BGE:
		c.bge(inst)
	case BLTU:
		c.bltu(inst)
	case BGEU:
		c.bgeu(inst)
	default:
		panic("Illegal/unimplemented instruction format")
	}
}

// load funct3
const (
	LB  uint32 = 0b000
	LH         = 0b001
	LW         = 0b010
	LBU        = 0b100
	LHU        = 0b101
)

func (c *Core) load(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case LB:
		c.lb(inst)
	case LH:
		c.lh(inst)
	case LW:
		c.lw(inst)
	case LBU:
		c.lbu(inst)
	case LHU:
		c.lhu(inst)
	default:
		panic("Illegal/unimplemented instruction format")
	}
}

// store funct3
const (
	SB uint32 = 0b000
	SH        = 0b001
	SW        = 0b010
)

func (c *Core) store(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case SB:
		c.sb(inst)
	case SH:
		c.sh(inst)
	case SW:
		c.sw(inst)
	default:
		panic("Illegal/unimplemented instruction format")
	}
}

// misc-mem funct3
const (
	FENCE uint32 = 0b000
)

func (c *Core) misc_mem(inst uint32) {
	funct3 := (inst >> 12) & 0x7
	switch funct3 {
	case FENCE:
		c.fence(inst)
	default:
		panic("Illegal instruction format")
	}
}

// system funct12
const (
	ECALL  uint32 = 0b000000000000
	EBREAK        = 0b000000000001
)

func (c *Core) system(inst uint32) {
	funct12 := (inst >> 20) & 0xfff
	switch funct12 {
	case ECALL:
		c.ecall(inst)
	case EBREAK:
		c.ebreak(inst)
	default:
		panic("Illegal/unimplemented instruction format")
	}
}

// AMO funct5
const (
	LR      uint32 = 0b00010
	SC             = 0b00011
	AMOSWAP        = 0b00001
)

func (c *Core) amo(inst uint32) {
	funct5 := inst >> 27
	switch funct5 {
	case LR:
		c.lr_w(inst)
	case SC:
		c.sc_w(inst)
	case AMOSWAP:
		c.amoswap_w(inst)
	}
}
