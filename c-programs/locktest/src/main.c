#include "sys.h"
#include "mutex.h"

extern void atomic_add(int *, int);

int main() {
    char s[] = "Tello, World!";
    s[0] = 'H';

    struct mutex *m = (struct mutex *)(0x20000);
    int *ct = (int *)(0x21000);

    putint((unsigned int)s);
    print(s, 13);

    int id = getpid();
    putint(id);

    for (int i = 0; i < 1024 * 32; i++) {
        lock(m);
        (*ct)++;
        unlock(m);
    }

    putint(id);

    lock(m);
    int ret = (*ct);
    unlock(m);

    return ret;
}
