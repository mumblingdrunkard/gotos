CC = clang -nostdlib --target=riscv32 -march=rv32g -Oz
OBJCOPY = llvm-objcopy


main.text: main
	${OBJCOPY} -O binary main main.text

main:
	${CC} src/_start.s     \
		  -o main -Wl,-Ttext=0x00004000

.PHONY: clean

clean:
	-@rm -rf obj asm main main.text
