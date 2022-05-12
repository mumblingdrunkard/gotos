.section .text
.globl sys_exit
.type sys_ext, @function

sys_exit:
	mv    a1, a0 # sys_exit is called with the argument
	li    a0, 1  # syscall number
	ecall
