extern int main();
extern void _Noreturn sys_exit(int);

void _Noreturn _start(void) {
    sys_exit(main());
}
