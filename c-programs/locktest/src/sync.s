.section .text
.globl sync
.type syn, @function

sync:
	fence
	ret
