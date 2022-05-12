.section .text
.globl acquire
.type acq, @function

acquire:
	li           t0, 1
	mv           a1, a0
	li           a0, 10
again:
	amoswap.w.aq t1, t0, (a1)
	bnez         t1, yield
	ret
yield:
	ecall
	j            again
