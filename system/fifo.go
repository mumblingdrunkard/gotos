package system

import "sync"

type FIFO struct {
	sync.Mutex
	queue []*PCB
}

func (f *FIFO) Push(pcb *PCB) {
	f.Lock()
	defer f.Unlock()
	f.queue = append(f.queue, pcb)
}

func (f *FIFO) Pop() *PCB {
	f.Lock()
	defer f.Unlock()

	if len(f.queue) == 0 {
		return nil
	}

	pcb := f.queue[0]
	f.queue = f.queue[1:]
	return pcb
}
