.section .text
.globl _start
.type _strt, @function

_start:
    li    t0, 0
    li    t1, 1
    li    t2, 0
    li    t3, 0
loop:
    add   t2, t0, t1
    mv    t0, t1
    mv    t1, t2
    addi  t3, t3, 1
    blt   t3, a0, loop
    mv    a1, t2
exit:
    li    a0, 1
    ecall
