package system

import (
	"fmt"
	"gotos/cpu"
)

//
// syscalls
//

func syscall(c *cpu.Core, number uint32) {
	const (
		sys_exit   = 1
		sys_id     = 2
		sys_getpid = 6
		sys_putint = 8
		sys_print  = 12
	)

	switch number {
	case sys_exit:
		sysExit(c)
	case sys_id:
		sysId(c)
	case sys_getpid:
		sysGetPID(c)
	case sys_putint:
		sysPutInt(c)
	case sys_print:
		sysPrint(c)
	}
}

func getArgs(c *cpu.Core) [6]uint32 {
	var args [6]uint32
	args[0] = c.GetIRegister(cpu.Reg_A1)
	args[1] = c.GetIRegister(cpu.Reg_A2)
	args[2] = c.GetIRegister(cpu.Reg_A3)
	args[3] = c.GetIRegister(cpu.Reg_A4)
	args[4] = c.GetIRegister(cpu.Reg_A5)
	args[5] = c.GetIRegister(cpu.Reg_A6)
	return args
}

func returnValue(c *cpu.Core, v uint32) {
	c.SetIRegister(cpu.Reg_A0, v)
}

func sysExit(c *cpu.Core) {
	// TODO sys_exit just halts the core
	// Instead, extract return value and store somewhere safe.
	// Start the next process in the queue
	// Halt if the queue is empty
	c.HaltIfRunning()
}

func sysId(c *cpu.Core) {
	c.SetIRegister(cpu.Reg_A0, c.GetCSR(cpu.Csr_MHARTID))
	// return to program
	// This is done by setting the pc to the address *after* the one that caused the trap
	// c.MEPC() will contain the address of the instruction that caused the trap
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4)
}

func sysGetPID(c *cpu.Core) {
	c.SetIRegister(cpu.Reg_A0, c.GetCSR(cpu.Csr_MHARTID))
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4)
}

func sysPutInt(c *cpu.Core) {
	fmt.Println(c.GetIRegister(cpu.Reg_A1))
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4)
}

func sysPrint(c *cpu.Core) {
	args := getArgs(c)

	// virtual pointer to string will be first argument
	vAddr := args[0]

	// length will be second argument
	length := args[1]

	// Ensure all data is available in memory
	c.CacheWriteback()

	_, pAddr, _ := c.Translate(vAddr)

	_, data := c.ReadMemory(pAddr, length)

	s := string(data)
	fmt.Println(s)

	returnValue(c, length)
	c.SetPC(c.GetCSR(cpu.Csr_MEPC) + 4)
}
