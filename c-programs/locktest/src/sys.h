#ifndef SYS_H
#define SYS_H

extern int syscall(int, ...);

extern void _Noreturn sys_exit(int);

// system calls
int getpid();
int putint(unsigned int);
int print(char*, unsigned int);

#endif
