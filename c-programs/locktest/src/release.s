.section .text
.globl release
.type rel, @function

release:
	amoswap.w.rl x0, x0, (a0)
	ret
