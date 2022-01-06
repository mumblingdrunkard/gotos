.section .text
.globl sys_id
.type sysid, @function

sys_id:
	li    a7, 2
	ecall
	ret
