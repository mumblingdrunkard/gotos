extern void dbg_break();
extern int sys_id();
extern void atomic_add(int *addr, int val);

#include "mutex.h"

int main() {
  struct mutex *m = (struct mutex *)4096;
  int *ct = (int *)2048;

  for (int i = 0; i < 1024 * 64; i++) {
    //   lock(m);
    //   (*ct)++;
    //   unlock(m);
    atomic_add(ct, 1);
  }

  lock(m);
  int ret = (*ct);
  unlock(m);

  return ret;
}
