.section .text
.globl sys_exit
.type sys_ext, @function

sys_exit:
	mv    a1, a0
	li    a0, 1
	ecall
