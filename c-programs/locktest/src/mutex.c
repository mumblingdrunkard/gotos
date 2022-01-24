extern void acquire(unsigned int *);
extern void release(unsigned int *);

// Becomes a `FENCE` instruction.
// Because of the non-discriminatory nature of the FENCE call when it comes to cache interactions, it is best to avoid this as much as possible.
// For updates on large data structures, this is perfectly fine as more of the cache is likely to be touched.
// It is less fine for things like counters, in which case you should probably use the instructions from the A extension.
extern void sync();

#include "mutex.h"

void lock(struct mutex *m) {
    acquire(&m->lock);

    // Invalidates and flushes the cache of this core.
    sync();
    // This is important to get all updates that other cores may have performed when holding the mutex.
}

void unlock(struct mutex *m) {
    // Invalidates and flushes the cache of this core.
    sync();
    // This is important to let all cores get updates that have been held while holding the mutex.

    release(&m->lock);
}
