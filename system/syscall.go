package system

import "gotos/cpu"

//
// syscalls
//

func syscall(c *cpu.Core, number uint32) {
	const (
		sys_exit = 1
		sys_id   = 2
	)

	switch number {
	case sys_exit:
		sysExit(c)
	case sys_id:
		coreID(c)
	}
}

func sysExit(c *cpu.Core) {
	// TODO sys_exit just halts the core
	// Instead, extract return value and store somewhere safe.
	// Start the next process in the queue
	// Halt if the queue is empty
	c.HaltIfRunning()
}

func coreID(c *cpu.Core) {
	c.SetIRegister(cpu.Reg_A0, c.GetCSR(cpu.Csr_MHARTID))
	// return to program
	// This is done by setting the pc to the address *after* the one that caused the trap
	// c.MEPC() will contain the address of the instruction that caused the trap
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4)
}
