.section .text
.globl syscall
.type syscall, @function

syscall:
	ecall
	ret
