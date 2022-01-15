// This file contains logic for decoding 32 bit RISC-V instructions.
// Decoding is implemented for the I, M, A, F, D, Zicsr, and Zifencei extensions.

// RISC-V privileged specification says:
//
// ---
//
// If **mtval** is written with a nonzero value when an illegal-instruction exception occurs, then **mtval** will contain the shortest of:
//
// 1. the actual faulting instruction
// 2. the first ILEN bits of the faulting instruction
// 3. the first MXLEN bits of the faulting instruction.
//
// ---
//
// I have opted for the first alternative.

package cpu

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
		LOAD_FP         = 0b0000111
		STORE_FP        = 0b0100111
		OP_FP           = 0b1010011
		FMADD           = 0b1000011
		FMSUB           = 0b1000111
		FNMSUB          = 0b1001011
		FNMADD          = 0b1001111
	)

	// Register 0 is hardwired with all 0s have to reset to 0 for every cycle because some instructions may use this as their /dev/null
	c.reg[RegZero] = 0

	// From the RISC-V privileged spec:
	//
	// ---
	//
	// For other traps, **mtval** is set to zero, but a future standard may redefine **mtval**'s setting for other traps.
	//
	// ---
	//
	// This ensures that mtval is always 0 before entering an instruction.
	// Before a trap is raised, mtval may be set.
	c.mtval = 0

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
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
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
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
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
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case LUI:
		c.lui(inst)
	case AUIPC:
		c.auipc(inst)
	case JAL:
		c.jal(inst)
	case JALR:
		c.jalr(inst)
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case SYSTEM:
		// system funct3
		const (
			ECALL_EBREAK uint32 = 0b000
			CSRRW               = 0b001
			CSRRS               = 0b010
			CSRRC               = 0b011
			CSRRWI              = 0b101
			CSRRSI              = 0b110
			CSRRCI              = 0b111
		)

		// system funct12
		const (
			ECALL  uint32 = 0b000000000000
			EBREAK        = 0b000000000001
		)

		funct3 := (inst >> 12) & 0x7
		switch funct3 {
		case ECALL_EBREAK:
			funct12 := (inst >> 20) & 0xfff
			switch funct12 {
			case ECALL:
				c.ecall(inst)
			case EBREAK:
				c.ebreak(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case CSRRW:
			c.csrrw(inst)
		case CSRRS:
			c.csrrs(inst)
		case CSRRC:
			c.csrrc(inst)
		case CSRRWI:
			c.csrrwi(inst)
		case CSRRSI:
			c.csrrsi(inst)
		case CSRRCI:
			c.csrrci(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
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
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case LOAD_FP:
		const (
			W uint32 = 0b010
			D        = 0b011
		)

		funct3 := (inst >> 12) & 0x7
		switch funct3 {
		case W:
			c.flw(inst)
		case D:
			c.fld(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case STORE_FP:
		const (
			W uint32 = 0b010
			D        = 0b011
		)

		funct3 := (inst >> 12) & 0x7
		switch funct3 {
		case W:
			c.fsw(inst)
		case D:
			c.fsd(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case OP_FP:
		const (
			// F extension
			FADD_S              uint32 = 0b0000000
			FSUB_S                     = 0b0000100
			FMUL_S                     = 0b0001000
			FDIV_S                     = 0b0001100
			FSQRT_S                    = 0b0101100
			FSGNJZ_S                   = 0b0010000 // FSGNJ_S, FSGNJN_S, FSGNJX_S
			FMNX_S                     = 0b0010100
			FCVT_WX_S                  = 0b1100000 // FCVT_W_S, FCVT_WU_S
			FMV_X_W_OR_FCLASS_S        = 0b1110000 // FMV_X_W, FCLASS_S
			FCMP_S                     = 0b1010000 // FEQ_S, FLT_S, FLE_S
			FCVT_S_WX                  = 0b1101000 // FCVT_S_W, FCVT_S_WU
			FMV_W_X                    = 0b1111000
			// D extension
			FADD_D    uint32 = 0b0000001
			FSUB_D           = 0b0000101
			FMUL_D           = 0b0001001
			FDIV_D           = 0b0001101
			FSQRT_D          = 0b0101101
			FSGNJZ_D         = 0b0010001 // FDGNJ_D, FDGNJN_D, FDGNJX_D
			FMNX_D           = 0b0010101 // FMIN_D, FMAX_D
			FCVT_S_D         = 0b0100000
			FCVT_D_S         = 0b0100001
			FCMP_D           = 0b1010001 // FEQ_D, FLT_D, FLE_D
			FCLASS_D         = 0b1110001
			FCVT_WX_D        = 0b1100001 // FCVT_W_D, FCVT_WU_D
			FCVT_D_WX        = 0b1101001 // FCVT_D_W, FCVT_D_WU
		)

		funct7 := (inst >> 25) & 0x7f

		switch funct7 {
		// F extension
		case FADD_S:
			c.fadd_s(inst)
		case FSUB_S:
			c.fsub_s(inst)
		case FMUL_S:
			c.fmul_s(inst)
		case FDIV_S:
			c.fdiv_s(inst)
		case FSQRT_S:
			c.fsqrt_s(inst)
		case FSGNJZ_S:
			const (
				FSGNJ_S  uint32 = 0b000
				FSGNJN_S        = 0b001
				FSGNJX_S        = 0b010
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FSGNJ_S:
				c.fsgnj_s(inst)
			case FSGNJN_S:
				c.fsgnjn_s(inst)
			case FSGNJX_S:
				c.fsgnjx_s(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FMNX_S:
			const (
				FMIN_S uint32 = 0b000
				FMAX_S        = 0b001
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FMIN_S:
				c.fmin_s(inst)
			case FMAX_S:
				c.fmax_s(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCVT_WX_S:
			const (
				FCVT_W_S  uint32 = 0b00000
				FCVT_WU_S        = 0b00001
			)
			funct5 := (inst >> 20) & 0x1f
			switch funct5 {
			case FCVT_W_S:
				c.fcvt_w_s(inst)
			case FCVT_WU_S:
				c.fcvt_wu_s(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FMV_X_W_OR_FCLASS_S:
			const (
				FMV_X_W  uint32 = 0b000
				FCLASS_S        = 0b001
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FMV_X_W:
				c.fmv_x_w(inst)
			case FCLASS_S:
				c.fclass_s(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCMP_S:
			const (
				FEQ_S uint32 = 0b010
				FLT_S        = 0b001
				FLE_S        = 0b000
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FEQ_S:
				c.feq_s(inst)
			case FLT_S:
				c.flt_s(inst)
			case FLE_S:
				c.fle_s(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCVT_S_WX:
			const (
				FCVT_S_W  uint32 = 0b00000
				FCVT_S_WU        = 0b00001
			)
			funct5 := (inst >> 20) & 0x1f
			switch funct5 {
			case FCVT_S_W:
				c.fcvt_s_w(inst)
			case FCVT_S_WU:
				c.fcvt_s_wu(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FMV_W_X:
			c.fmv_w_x(inst)
		// D extension
		case FADD_D:
			c.fadd_d(inst)
		case FSUB_D:
			c.fsub_d(inst)
		case FMUL_D:
			c.fmul_d(inst)
		case FDIV_D:
			c.fdiv_d(inst)
		case FSQRT_D:
			c.fsqrt_d(inst)
		case FSGNJZ_D:
			const (
				FSGNJ_D  uint32 = 0b000
				FSGNJN_D        = 0b001
				FSGNJX_D        = 0b010
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FSGNJ_D:
				c.fsgnj_d(inst)
			case FSGNJN_D:
				c.fsgnjn_d(inst)
			case FSGNJX_D:
				c.fsgnjn_d(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FMNX_D:
			const (
				FMIN_D uint32 = 0b000
				FMAX_D        = 0b001
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FMIN_D:
				c.fmin_d(inst)
			case FMAX_D:
				c.fmax_d(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCVT_S_D:
			c.fcvt_s_d(inst)
		case FCVT_D_S:
			c.fcvt_d_s(inst)
		case FCMP_D:
			const (
				FEQ_D uint32 = 0b010
				FLT_D        = 0b001
				FLE_D        = 0b000
			)
			funct3 := (inst >> 12) & 0x7
			switch funct3 {
			case FEQ_D:
				c.feq_d(inst)
			case FLT_D:
				c.flt_d(inst)
			case FLE_D:
				c.fle_d(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCLASS_D:
			c.fclass_d(inst)
		case FCVT_WX_D:
			const (
				FCVT_W_D  uint32 = 0b00000
				FCVT_WU_D        = 0b00001
			)
			funct5 := (inst >> 20) & 0x1f
			switch funct5 {
			case FCVT_W_D:
				c.fcvt_w_d(inst)
			case FCVT_WU_D:
				c.fcvt_wu_d(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		case FCVT_D_WX:
			const (
				FCVT_D_W  uint32 = 0b00000
				FCVT_D_WU        = 0b00001
			)
			funct5 := (inst >> 20) & 0x1f
			switch funct5 {
			case FCVT_D_W:
				c.fcvt_d_w(inst)
			case FCVT_D_WU:
				c.fcvt_d_wu(inst)
			default:
				c.mtval = inst
				c.trap(TrapIllegalInstruction)
			}
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case FMADD:
		const (
			S uint32 = 0b00
			D        = 0b01
		)
		format2 := (inst >> 25) & 0x3
		switch format2 {
		case S:
			c.fmadd_s(inst)
		case D:
			c.fmadd_d(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case FMSUB:
		const (
			S uint32 = 0b00
			D        = 0b01
		)
		format2 := (inst >> 25) & 0x3
		switch format2 {
		case S:
			c.fmsub_s(inst)
		case D:
			c.fmsub_d(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case FNMSUB:
		const (
			S uint32 = 0b00
			D        = 0b01
		)
		format2 := (inst >> 25) & 0x3
		switch format2 {
		case S:
			c.fnmadd_s(inst)
		case D:
			c.fnmadd_d(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	case FNMADD:
		const (
			S uint32 = 0b00
			D        = 0b01
		)
		format2 := (inst >> 25) & 0x3
		switch format2 {
		case S:
			c.fnmsub_s(inst)
		case D:
			c.fnmsub_d(inst)
		default:
			c.mtval = inst
			c.trap(TrapIllegalInstruction)
		}
	default:
		c.mtval = inst
		c.trap(TrapIllegalInstruction)
	}
}
