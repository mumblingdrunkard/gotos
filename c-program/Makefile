CC = clang -nostdlib --target=riscv32 -march=rv32i -Oz
OBJCOPY = llvm-objcopy


minimal.text: minimal
	${OBJCOPY} -O binary minimal minimal.text

minimal: obj/_start.o obj/main.o obj/sys_exit.o obj/dbg_break.o
	${CC} obj/_start.o obj/main.o obj/sys_exit.o obj/dbg_break.o -o minimal -Wl,-Ttext=0x00000000

obj/main.o: obj src/main.c
	${CC} -c src/main.c -o obj/main.o

obj/_start.o: obj src/_start.c
	${CC} -c src/_start.c -o obj/_start.o

obj/sys_exit.o: obj src/sys_exit.s
	${CC} -c src/sys_exit.s -o obj/sys_exit.o

obj/dbg_break.o: obj src/dbg_break.s
	${CC} -c src/dbg_break.s -o obj/dbg_break.o



minimal_asm: asm/main.s asm/_start.s asm/sys_exit.s

asm/main.s: asm src/main.c
	${CC} -S src/main.c -o asm/main.s

asm/_start.s: asm src/_start.c
	${CC} -S src/_start.c -o asm/_start.s

asm/sys_exit.s: asm src/sys_exit.s
	cp src/sys_exit.s asm/sys_exit.s

asm/dbg_break.s: asm src/dbg_break.s
	cp src/dbg_break.s asm/dbg_break.s



obj:
	mkdir -p obj

asm:
	mkdir -p asm



.PHONY: clean

clean:
	-@rm -rf obj asm minimal minimal.text
