struct mutex {
    unsigned int lock;
};

void lock(struct mutex *m);
void unlock(struct mutex *m);
