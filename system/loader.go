package system

import (
	"encoding/binary"
	"gotos/cpu"
	"os"
)

// Load loads a raw binary from file `fname` and places it at `addr` in system
// memory, and creates a process with `pc`, `sp`, `pid`, and a pointer to a page
// table `ptableAddr`.
//   `addr` has to be aligned on an INSTRUCTION_WIDTH byte boundary (4 bytes).
func (s *System) Load(fname string, pc, sp, pid, addr, ptableAddr uint32) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	stats, err := f.Stat()
	size := stats.Size()
	program := make([]uint8, size)
	binary.Read(f, binary.BigEndian, &program)
	f.Close()

	err, _ = s.memory.WriteRaw(addr, program)
	pcb := PCB{
		PC:     pc,
		PID:    pid,
		PTable: ptableAddr,
	}
	pcb.IReg[cpu.Reg_SP] = sp
	s.Scheduler.Push(&pcb)
}
