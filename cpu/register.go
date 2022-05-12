// this file declares the type `Reg` and gives names (and aliases) for the 32
// RISC-V registers.

package cpu

type Reg int

// Register mnemonics
const (
	// Unnamed
	Reg_X0  Reg = 0
	Reg_X1      = 1
	Reg_X2      = 2
	Reg_X3      = 3
	Reg_X4      = 4
	Reg_X5      = 5
	Reg_X6      = 6
	Reg_X7      = 7
	Reg_X8      = 8
	Reg_X9      = 9
	Reg_X10     = 10
	Reg_X11     = 11
	Reg_X12     = 12
	Reg_X13     = 13
	Reg_X14     = 14
	Reg_X15     = 15
	Reg_X16     = 16
	Reg_X17     = 17
	Reg_X18     = 18
	Reg_X19     = 19
	Reg_X20     = 20
	Reg_X21     = 21
	Reg_X22     = 22
	Reg_X23     = 23
	Reg_X24     = 24
	Reg_X25     = 25
	Reg_X26     = 26
	Reg_X27     = 27
	Reg_X28     = 28
	Reg_X29     = 29
	Reg_X30     = 30
	Reg_X31     = 31

	// ABI Names
	Reg_ZERO = 0  // Hard-wired zero
	Reg_RA   = 1  // Return address
	Reg_SP   = 2  // Stack pointer
	Reg_GP   = 3  // Global pointer
	Reg_TP   = 4  // Thread pointer
	Reg_T0   = 5  // Temporary/alternate link register
	Reg_T1   = 6  // Temporaries
	Reg_T2   = 7  //
	Reg_S0   = 8  // Saved register/frame pointer
	Reg_FP   = 8  //
	Reg_S1   = 9  // Saved register
	Reg_A0   = 10 // Function arguments/return values
	Reg_A1   = 11 //
	Reg_A2   = 12 //
	Reg_A3   = 13 //
	Reg_A4   = 14 //
	Reg_A5   = 15 //
	Reg_A6   = 16 //
	Reg_A7   = 17 //
	Reg_S2   = 18 // Saved registers
	Reg_S3   = 19 //
	Reg_S4   = 20 //
	Reg_S5   = 21 //
	Reg_S6   = 22 //
	Reg_S7   = 23 //
	Reg_S8   = 24 //
	Reg_S9   = 25 //
	Reg_S10  = 26 //
	Reg_S11  = 27 //
	Reg_T3   = 28 // Temporaries
	Reg_T4   = 29 //
	Reg_T5   = 30 //
	Reg_T6   = 31 //
)
