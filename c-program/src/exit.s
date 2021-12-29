.section .text
.globl sys_exit
.type sys_ext, @function

sys_exit:
	li    a7, 1 // syscall number
	ecall
