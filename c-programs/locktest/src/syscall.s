.section .text
.globl syscall
.type syscall, @function

syscall:
	mv a7, a0
	ecall
	ret
