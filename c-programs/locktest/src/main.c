#include "mutex.h"

int main() {
    struct mutex *m = (struct mutex *)(0x20000);
    int *ct = (int *)(0x21000);

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
