#include "mutex.h"

extern void atomic_add(int *, int);

int main() {
    struct mutex *m = (struct mutex *)4096;
    int *ct = (int *)2048;

    for (int i = 0; i < 1024 * 32; i++) {
        // This really sucks for cache and we're only updating one counter so
        // the performance is pure garbage. An improvement would be to use one
        // of the instructions from the A extension such as AMOADD. This
        // instruction is available as the function `atomic_add(...)` which
        // takes a pointer to an int and a value to increment the number that
        // this pointer points to.
        //
        // Example code
        //
        //     ```c
        //     atomic_add(ct, 1);
        //     ```
        //
        //     Code example: Using `atomic_add` to increment a counter `ct`.
        //
        // is the equivalent, thread-safe method to increment the counter ct
        // by 1, but has much higher performance because it doesn't thrash the
        // cache and doesn't require locking and unlocking.

        lock(m);
        (*ct)++;
        unlock(m);
    }

    lock(m);
    int ret = (*ct);
    unlock(m);

    return ret;
}
