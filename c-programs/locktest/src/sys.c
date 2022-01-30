#include "sys.h"

#define SYS_GETPID 6
#define SYS_PUTINT 8
#define SYS_PRINT 12

int getpid() { return syscall(SYS_GETPID); }

int putint(unsigned int c) { return syscall(SYS_PUTINT, c); }

int print(char* s, unsigned int length) { return syscall(SYS_PRINT, s, length); };
