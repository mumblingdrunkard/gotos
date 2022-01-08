package cpu

// TODO: Exceptions

// TODO: Fix decoding for OP
// TODO: Add decoding for RV32M extension

// TODO: Add decoding for FENCE_I

func (c *Core) execute(inst uint32) {
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
	// Register 0 is hardwired with all 0s have to reset to 0 for every
	// cycle because some instructions may use this as their /dev/null
	c.reg[0] = 0

	opcode := inst & 0x7f

	switch opcode {
	case OP:
		// OP funct7
		const (
			OP_A   uint32 = 0b0000000
			OP_B          = 0b0100000
			MULDIV        = 0b0000001
		)

		funct7 := (inst >> 25) & 0x7f
		switch funct7 {
		case OP_A:
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
		case OP_B:
			const (
				SUB uint32 = 0b000
				SRA        = 0b101
			)

			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case SUB:
				c.sub(inst)
			case SRA:
				c.sra(inst)
			default:
				panic("Unknown instruction")
			}
		case MULDIV:
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
		default:
			panic("Illegal instruction format")
		}
	case OP_IMM:
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
	case LUI:
		c.lui(inst)
	case AUIPC:
		c.auipc(inst)
	case JAL:
		c.jal(inst)
		c.pc -= 4 // decrement pc as it will be incremented straight after
	case JALR:
		c.jalr(inst)
		c.pc -= 4 // decrement pc as it will be incremented straight after
	case BRANCH:
		// branch funct3
		const (
			BEQ  uint32 = 0b000
			BNE         = 0b001
			BLT         = 0b100
			BGE         = 0b101
			BLTU        = 0b110
			BGEU        = 0b111
		)

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
		c.pc -= 4 // decrement pc as it will be incremented straight after
	case LOAD:
		// load funct3
		const (
			LB  uint32 = 0b000
			LH         = 0b001
			LW         = 0b010
			LBU        = 0b100
			LHU        = 0b101
		)

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
	case STORE:
		// store funct3
		const (
			SB uint32 = 0b000
			SH        = 0b001
			SW        = 0b010
		)

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
	case MISC_MEM:
		// misc-mem funct3
		const (
			FENCE   uint32 = 0b000
			FENCE_I        = 0b001
		)

		funct3 := (inst >> 12) & 0x7
		switch funct3 {
		case FENCE:
			c.fence(inst)
		case FENCE_I:
			c.fence_i(inst)
		default:
			panic("Illegal instruction format")
		}
	case SYSTEM:
		// system funct12
		const (
			ECALL  uint32 = 0b000000000000
			EBREAK        = 0b000000000001
		)

		funct12 := (inst >> 20) & 0xfff
		switch funct12 {
		case ECALL:
			c.ecall(inst)
		case EBREAK:
			c.ebreak(inst)
		default:
			panic("Illegal/unimplemented instruction format")
		}
	case AMO:
		// AMO funct5
		const (
			LR      uint32 = 0b00010
			SC             = 0b00011
			AMOSWAP        = 0b00001
			AMOADD         = 0b00000
			AMOXOR         = 0b00100
			AMOAND         = 0b01100
			AMOOR          = 0b01000
			AMOMIN         = 0b10000
			AMOMAX         = 0b10100
			AMOMINU        = 0b11000
			AMOMAXU        = 0b11100
		)

		funct5 := inst >> 27
		switch funct5 {
		case LR:
			c.lr_w(inst)
		case SC:
			c.sc_w(inst)
		case AMOSWAP:
			c.amoswap_w(inst)
		case AMOADD:
			c.amoadd_w(inst)
		case AMOXOR:
			c.amoxor_w(inst)
		case AMOAND:
			c.amoand_w(inst)
		case AMOOR:
			c.amoor_w(inst)
		case AMOMIN:
			c.amomin_w(inst)
		case AMOMAX:
			c.amomax_w(inst)
		case AMOMINU:
			c.amominu_w(inst)
		case AMOMAXU:
			c.amomaxu_w(inst)
		default:
			panic("Unknown instruction")
		}
	default:
		panic("Unknown instruction")
	}
}
