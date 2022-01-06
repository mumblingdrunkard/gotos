.section .text
.globl sys_exit
.type sys_ext, @function

// exit is special. It doesn't fill the return address as there should be no returing from the method
sys_exit:
	li    a7, 1 // syscall number
	ecall
