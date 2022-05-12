#include "sys.h"

#define SYS_GETPID 6
#define SYS_PUTINT 8

int getpid() { return syscall(SYS_GETPID); }

int putint(unsigned int c) { return syscall(SYS_PUTINT, c); }
