#include "sys.h"

extern int main();
void _Noreturn _start() { sys_exit(main()); }
