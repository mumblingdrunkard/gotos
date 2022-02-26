.section .text
.globl sys_exit
.type sys_ext, @function

sys_exit:
	li    a0, 1
	ecall
