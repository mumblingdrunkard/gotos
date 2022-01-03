extern void dbg_break();

int main() {
  // volatile to avoid optimizing away the branch
  volatile int a = 1;
  volatile int b = 2;
  // dbg_break();
  if (a < b) {
    return 105;
  } else {
    return 1056;
  }
}
