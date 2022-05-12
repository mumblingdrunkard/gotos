int main() {
    int *const restrict a = (int *const restrict)0x00040000;
    const int len = 5 * 0x1000;
    for (int i = 0; i < 1000; i++) {
        for (int j = 0; j < len; j++) {
            a[j] += i;
        }
    }

    int sum = 0;
    for (int i = 0; i < len; i++) {
        sum += a[i];
    }

    return sum;
}
