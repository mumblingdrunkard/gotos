CC = clang -nostdlib --target=riscv32 -march=rv32g -Oz
OBJCOPY = llvm-objcopy

fib.text: fib
	${OBJCOPY} -O binary fib fib.text

fib:
	${CC} src/_start.c   \
		  src/main.c     \
		  src/sys_exit.s \
		  -o fib -Wl,-Ttext=0x00004000

.PHONY: clean

clean:
	-@rm -rf fib fib.text
