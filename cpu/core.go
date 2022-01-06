package cpu

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	RUNNING = 2 // Core is running and executing instructions
	HALTING = 1 // Core is running and executing instructions, but will turn off in the next cycle
	HALTED  = 0 // Core is halted and will stay halted until Start() is called
)

const (
	LITTLE Endian = 0
	BIG           = 1
)

type ReservationSets struct {
	sync.Mutex
	sets   []map[uint32]bool
	lookup []*map[uint32]bool
}

func NewReservationSets(n int) (rs ReservationSets) {
	rs = ReservationSets{
		sets:   make([]map[uint32]bool, n),
		lookup: make([]*map[uint32]bool, n),
	}

	for i := 0; i < n; i++ {
		rs.sets[i] = make(map[uint32]bool)
		rs.lookup[i] = &rs.sets[i]
	}

	return
}

// A RISC-V core that runs in user mode
type Core struct {
	sync.WaitGroup
	id     int
	state  atomic.Value // core state (HALTED, HALTING, RUNNING)
	retire uint64       // number of instructions executed
	inst   uint32       // currently executing instruction
	reg    [32]uint32   // registers
	pc     uint32       // program counter
	rsets  *ReservationSets
	mc     MemoryController
}

func (c *Core) fetch() uint32 {
	err, inst := c.mc.LoadInstruction(c.pc)

	if err != nil {
		panic(err)
	}

	return inst
}

func (c *Core) UnsafeSetMemBase(base uint32) {
	c.mc.mmu.base = base
}

func (c *Core) UnsafeSetMemSize(size uint32) {
	c.mc.mmu.size = size
	c.reg[2] = c.mc.mmu.size
}

func (c *Core) run(wg *sync.WaitGroup) {
	// Start running core in loop
	if !c.state.CompareAndSwap(HALTED, RUNNING) {
		panic("Attempted to call `run()` on a core that was not in the HALTED state")
	}

	wg.Done() // core has switched to running state

	for {
		// Test if state is HALTING, swap to HALTED if so, then break
		// CompareAndSwap makes this very slow so we use Load instead
		if c.state.Load() == HALTING {
			c.state.Store(HALTED)
			break
		}

		c.UnsafeStep()
	}

	c.Done() // core is done
}

// makes sure the core has at least entered the running state before returning
// It is an error to call Start on a core that is already started
func (c *Core) StartAndWait() {
	var wg sync.WaitGroup
	wg.Add(1)
	c.Add(1) // caller can't accidentally wait on core that hasn't entered loop
	go c.run(&wg)
	wg.Wait()
}

// Leaves it up to the caller when to sync
// It is an error to call Start on a core that is already started
func (c *Core) StartAndSync(wg *sync.WaitGroup) {
	wg.Add(1)
	c.Add(1) // caller can't accidentally wait on core
	go c.run(wg)
}

func (c *Core) UnsafeReset() {
	// Reset registers
	for i := range c.reg {
		c.reg[i] = 0
	}

	// TODO: Initialize reg[2] with memory size
	c.reg[2] = c.mc.mmu.size

	c.pc = 0
	// c.state.Store(HALTED)
}

// Halts the core and waits for it to halt before returning
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndWait() {
	if c.state.Load() == HALTED {
		return
	}

	if !c.state.CompareAndSwap(RUNNING, HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	c.Wait() // Wait for core to halt
}

// Halts the core, but leaves it to the caller to sync
// TODO: It is an error to halt a core that is not running, this is a potential issue if cores can halt themselves
// Possible fix, check for halted state first, if so, return immediately
// con: it would be an error to call halt if a go run() has been performed, but not yet scheduled so state is not yet RUNNING
//   This is an acceptable solution.
func (c *Core) HaltAndSync(wg *sync.WaitGroup) {
	if !c.state.CompareAndSwap(RUNNING, HALTING) {
		panic("Attempted to halt core that was not in RUNNING state")
	}

	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

// Not really needed anymore, use core.Wait() instead
func (c *Core) SyncOnHalt(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		c.Wait()
		wg.Done()
	}()
}

func NewCore() (c Core) {
	c = Core{}

	c.UnsafeReset()

	return
}

func (c *Core) UnsafeStep() {
	// TODO: Check for interrupts

	c.retire += 1
	inst := c.fetch()

	// fmt.Printf("executing: %08X\n", inst)

	c.execute(inst)
	opcode := inst & 0x7f
	if (opcode != BRANCH) && (opcode != JAL) && (opcode != JALR) {
		c.pc += 4
	}
}

func NewCoreWithMemoryAndReservationSets(m *Memory, rs *ReservationSets, id int) (c Core) {
	c = Core{
		rsets:  rs,
		id:     id,
		retire: 0,
		mc:     NewMemoryController(m),
	}

	c.state.Store(HALTED)

	c.UnsafeReset()

	return
}

func (c *Core) InstructionsRetired() uint64 {
	return c.retire
}

func (c *Core) Misses() uint64 {
	return c.mc.misses
}

func (c *Core) Accesses() uint64 {
	return c.mc.accesses
}

func (c *Core) State() interface{} {
	return c.state.Load()
}

func (c *Core) DumpRegisters() {
	fmt.Println("Register dump")
	fmt.Printf("pc: %X\n", c.pc)
	for i, r := range c.reg {
		fmt.Printf("[%02d]: %08X\n", i, r)
	}
}
