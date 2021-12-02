package cpu

import (
	"encoding/binary"
	"sync"
)

const (
	MEMORY_SIZE           = 1024 * 1024 * 32 // 32 MiB of memory
	INSTRUCTIONS_PER_LOCK = 10
)

type CPU struct {
	sync.Mutex
	running bool
	reg     [32]uint64
	pc      uint64
	mem     []uint8
}

type State struct {
	reg [32]uint64
	pc  uint64
}

func (state *State) Reg() [32]uint64 {
	return state.reg
}

func (state *State) Pc() uint64 {
	return state.pc
}

func (cpu *CPU) fetch() uint32 {
	pc := cpu.pc
	cpu.pc += 4
	return binary.LittleEndian.Uint32(cpu.mem[pc : pc+4])
}

func (cpu *CPU) execute(inst uint32) {
	// Register 0 is hardwired with all 0s
	cpu.reg[0] = 0
	cpu.reg[3]++

	op := inst & 0x7f            // 7 bits
	funct3 := (inst >> 12) & 0x7 // 3 bits
	rd := (inst >> 7) & 0x1f     // 5 bits
	rs1 := (inst >> 15) & 0x1f   // 5 bits
	rs2 := (inst >> 20) & 0x1f   // 5 bits
	immi := (inst >> 20) & 0xfff // 12 bits
	// imms := (inst & 0x1f) | ((inst >> 20) & 0x7f) // 12 bits
	// immu := (inst >> 12) & 0xfffff // 20 bits

	switch op {
	case 0x03:
		cpu.load(rd, funct3, rs1, immi)
	case 0x33: // add
		cpu.add(rd, rs1, rs2)
	case 0x13: // addi
		cpu.addi(rd, rs1, immi)
	default:
		cpu.nop()
	}
}

func (cpu *CPU) nop() {
	cpu.pc = 0
}

func (cpu *CPU) load(rd, funct3, rs1, immi uint32) {
	switch funct3 {
	case 0x0:
		cpu.lb(rd, rs1, immi)
	case 0x1:
		cpu.lh(rd, rs1, immi)
	case 0x2:
	case 0x3:
	case 0x4:
	case 0x5:
	case 0x6:
	}
}
func (cpu *CPU) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start running cpu in loop
outer_loop:
	for {
		cpu.Lock()

		for i := 0; i < INSTRUCTIONS_PER_LOCK; i++ {
			// Check if cpu should stop running
			if !cpu.running {
				cpu.Unlock()
				break outer_loop
			}

			inst := cpu.fetch()
			cpu.execute(inst)
		}

		cpu.Unlock()
	}
}

func (cpu *CPU) LoadMemory(data []uint8, offset uint64) {
	cpu.Lock()
	for i := range data {
		cpu.mem[offset+uint64(i)] = data[i]
	}
	cpu.Unlock()
}

func (cpu *CPU) GetState() State {
	cpu.Lock()
	defer cpu.Unlock()
	return State{
		reg: cpu.reg,
		pc:  cpu.pc,
	}
}

func (cpu *CPU) RestoreState(state State) {
	cpu.Lock()
	cpu.reg = state.reg
	cpu.pc = state.pc
	cpu.Unlock()
}

func (cpu *CPU) Reset() {
	cpu.Lock()
	cpu.running = true

	// Reset registers
	for i := range cpu.reg {
		cpu.reg[i] = 0
	}

	// Initialize reg[2] with memory size
	cpu.reg[2] = MEMORY_SIZE

	cpu.pc = 0

	// TODO reset memory
	cpu.Unlock()
}

func (cpu *CPU) Stop() {
	cpu.Lock()
	cpu.running = false
	cpu.Unlock()
}

func NewCPU() CPU {
	return CPU{
		mem: make([]uint8, MEMORY_SIZE),
	}
}
