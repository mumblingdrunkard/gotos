extern void dbg_break();
extern int sys_id();

#include "mutex.h"

int main() {
  struct mutex *m = (struct mutex *)4096;
  int *ct = (int *)2048;

  for (int i = 0; i < 1024; i++) {
    lock(m);
    (*ct)++;
    unlock(m);
  }

  lock(m);
  int ret = (*ct);
  unlock(m);

  return ret;
}
