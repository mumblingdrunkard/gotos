#ifndef SYS_H
#define SYS_H

extern int syscall(int, ...);

#define SYS_GETPID 6

extern void dbg_break();
extern void _Noreturn sys_exit(int);

// system calls
int getpid() { return syscall(SYS_GETPID); };

#endif
