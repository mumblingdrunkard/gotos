extern void acquire(unsigned int *);
extern void release(unsigned int *);
extern void sync();

#include "mutex.h"

void lock(struct mutex* m) {
    acquire(&m->lock);
    sync();
}

void unlock(struct mutex* m) {
    sync();
    release(&m->lock);
}
