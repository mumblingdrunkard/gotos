.section .text
.globl _start
.type _strt, @function

_start:
	li	t0, 17
	li	t1, 25
	add	a0, t0, t1
exit:
	li	a0, 1
	ecall
