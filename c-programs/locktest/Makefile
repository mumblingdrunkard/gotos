CC = clang -nostdlib --target=riscv32 -march=rv32ima -Oz
OBJCOPY = llvm-objcopy


locktest.text: locktest
	${OBJCOPY} -O binary locktest locktest.text

locktest: obj/_start.o obj/main.o obj/sys_exit.o obj/sys_id.o obj/dbg_break.o obj/mutex.o obj/acquire.o obj/release.o obj/sync.o obj/atomic_add.o
	${CC} obj/_start.o \
		  obj/main.o \
		  obj/mutex.o \
		  obj/acquire.o \
		  obj/release.o \
		  obj/atomic_add.o \
		  obj/sync.o \
		  obj/sys_exit.o \
		  obj/sys_id.o \
		  obj/dbg_break.o \
		  -o locktest -Wl,-Ttext=0x00000000

obj/main.o: obj src/main.c
	${CC} -c src/main.c -o obj/main.o

obj/_start.o: obj src/_start.c
	${CC} -c src/_start.c -o obj/_start.o

obj/sys_exit.o: obj src/sys_exit.s
	${CC} -c src/sys_exit.s -o obj/sys_exit.o

obj/acquire.o: obj src/acquire.s
	${CC} -c src/acquire.s -o obj/acquire.o

obj/release.o: obj src/release.s
	${CC} -c src/release.s -o obj/release.o

obj/atomic_add.o: obj src/atomic_add.s
	${CC} -c src/atomic_add.s -o obj/atomic_add.o

obj/sync.o: obj src/sync.s
	${CC} -c src/sync.s -o obj/sync.o

obj/sys_id.o: obj src/sys_id.s
	${CC} -c src/sys_id.s -o obj/sys_id.o

obj/mutex.o: obj src/mutex.c
	${CC} -c src/mutex.c -o obj/mutex.o

obj/dbg_break.o: obj src/dbg_break.s
	${CC} -c src/dbg_break.s -o obj/dbg_break.o



locktest_asm: asm/main.s asm/_start.s asm/sys_exit.s

asm/main.s: asm src/main.c
	${CC} -S src/main.c -o asm/main.s

asm/_start.s: asm src/_start.c
	${CC} -S src/_start.c -o asm/_start.s

asm/sys_exit.s: asm src/sys_exit.s
	cp src/sys_exit.s asm/sys_exit.s

asm/acquire.s: asm src/acquire.s
	cp src/acquire.s asm/acquire.s

asm/sys_id.s: asm src/sys_id.s
	cp src/sys_id.s asm/sys_id.s

asm/dbg_break.s: asm src/dbg_break.s
	cp src/dbg_break.s asm/dbg_break.s



obj:
	mkdir -p obj

asm:
	mkdir -p asm



.PHONY: clean

clean:
	-@rm -rf obj asm locktest locktest.text
