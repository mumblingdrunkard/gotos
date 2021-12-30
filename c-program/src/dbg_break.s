.section .text
.globl dbg_break
.type dbg_brk, @function

// Just an environment break.
dbg_break:
	ebreak
	ret
