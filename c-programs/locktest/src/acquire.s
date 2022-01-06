.section .text
.globl acquire
.type acq, @function

acquire:
	li           t0, 1
again:
	lr.w         t1, (a0)
	bnez         t1, again
	amoswap.w.aq t1, t0, (a0)
	bnez         t1, again
	ret
