package system

import (
	"encoding/binary"
	"gotos/cpu"
	"os"
)

func (s *System) Load(fname string, pc, sp, pid uint32) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	stats, err := f.Stat()
	size := stats.Size()
	program := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &program)
	f.Close()

	err, _ = s.memory.Write(pc, program)
	pcb := PCB{
		PC:  pc,
		PID: pid,
	}
	pcb.IReg[cpu.Reg_SP] = sp

	s.Scheduler.Push(&pcb)
}
