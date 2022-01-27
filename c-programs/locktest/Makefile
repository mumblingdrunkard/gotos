CC = clang -nostdlib --target=riscv32 -march=rv32ima -Oz
OBJCOPY = llvm-objcopy


locktest.text: locktest
	${OBJCOPY} -O binary locktest locktest.text

locktest:
	${CC} src/_start.c     \
		  src/main.c       \
		  src/mutex.c      \
		  src/acquire.s    \
		  src/release.s    \
		  src/atomic_add.s \
		  src/sync.s       \
		  src/sys_exit.s   \
		  src/dbg_break.s  \
		  src/syscall.s    \
		  -o locktest -Wl,-Ttext=0x00004000

.PHONY: clean

clean:
	-@rm -rf obj asm locktest locktest.text
