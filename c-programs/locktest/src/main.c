#include "mutex.h"

int main() {
    struct mutex *m = (struct mutex *)(0x200);
    int *ct = (int *)(0x400);

    for (int i = 0; i < 1024 * 32; i++) {
        lock(m);
        (*ct)++;
        unlock(m);
    }

    lock(m);
    int ret = (*ct);
    unlock(m);

    return ret;
}
