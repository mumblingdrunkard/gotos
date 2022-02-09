.section .text
.globl _start
.type _strt, @function

_start:
	li		t0, 0
	li		t1, 42
loop:
	addi    t0, t0, 1
	blt     t0, t1, loop
	ecall
