CC = clang -nostdlib --target=riscv32 -march=rv32i -O0
OBJCOPY = llvm-objcopy


minimal.text: minimal
	${OBJCOPY} -O binary minimal minimal.text

minimal: obj/_start.o obj/main.o obj/exit.o
	${CC} obj/_start.o obj/main.o obj/exit.o -o minimal -Wl,-Ttext=0x80000000

obj/main.o: obj src/main.c
	${CC} -c src/main.c -o obj/main.o

obj/_start.o: obj src/_start.c
	${CC} -c src/_start.c -o obj/_start.o

obj/exit.o: obj src/exit.s
	${CC} -c src/exit.s -o obj/exit.o



minimal_asm: asm/main.s asm/_start.s asm/exit.s

asm/main.s: asm src/main.c
	${CC} -S src/main.c -o asm/main.s

asm/_start.s: asm src/_start.c
	${CC} -S src/_start.c -o asm/_start.s

asm/exit.s: asm src/exit.s
	cp src/exit.s asm/exit.s




obj:
	mkdir -p obj

asm:
	mkdir -p asm



.PHONY: clean

clean:
	-@rm -rf obj asm minimal minimal.text
