extern void dbg_break();
extern int sys_id();

#include "mutex.h"

int main() {
  struct mutex *m = (struct mutex *)4096;
  int *ct = (int *)2048;

  lock(m);
  (*ct)++;
  unlock(m);

  return *ct;
}
