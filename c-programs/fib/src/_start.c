extern int main();

extern void _Noreturn sys_exit(int placeholder, int arg);
void _Noreturn exit(int res) {
        sys_exit(1, res);
}

void _Noreturn _start() {
    exit(main());
}
