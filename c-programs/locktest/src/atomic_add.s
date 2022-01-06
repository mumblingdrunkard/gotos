.section .text
.globl atomic_add
.type am_add, @function

atomic_add:
	amoadd.w x0, a1, (a0)
	ret
