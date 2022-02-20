package cpu

import "fmt"

// TODO Should mmu be moved into its own package perhaps?

const (
	pageFlagValid    uint32 = 0x01 // the virtual address is valid
	pageFlagRead            = 0x02 // indicates that the processor is allowed to read data from this address
	pageFlagWrite           = 0x04 // indicates that the processor is allowed to write data to this address
	pageFlagExec            = 0x08 // indicates that the processor is allowed to fetch instructions from this address
	pageFlagUser            = 0x10 // indicates that the processor can access this page in user mode
	pageFlagGlobal          = 0x20 // whether this page is globally mapped into all address spaces (probably unused here)
	pageFlagAccessed        = 0x40 // whether this page has been accessed since the access bit was last cleared
	pageFlagDirty           = 0x80 // whether this page has been written to since the dirty bit was last cleared
)

// =============================================================================
// >           31  20   19  10   9 8   7   6   5   4   3   2   1   0
// >         | PPN[1] | PPN[0] | RSW | D | A | G | U | X | W | R | V |
// >             12       10      2    1   1   1   1   1   1   1   1
// >
// >                    Figure 4.18: Sv32 page table entry

// Two schemes to manage the A and D bits are permitted:
// - When a virtual page is accessed and the A bit is clear, or is written and
//   the D bit is clear, a page-fault exception is raised.
// - ...
//
// ...
//
//     ---
//         The A and D bits are never cleared by the implementation.
//     If the supervisor software does not rely on accessed and/or dirty bits,
//     e.g. if it does not swap memory pages to secondary storage or if the
//     pages are being used to map I/O space, it should always set them to 1 in
//     the PTE to improve performance.
//     ---

type accessType int

const (
	accessTypeLoad             accessType = 0
	accessTypeStore                       = 1
	accessTypeInstructionFetch            = 2
)

const (
	pagesize = 4096
)

// Translates the address and returns the flags that apply if the address is valid
func (c *Core) translate(vAddr uint32, aType accessType) (success bool, pAddr uint32) {
	// get the satp register
	satp := c.csr[Csr_SATP]

	success = true
	if satp&0x80000000 == 0 { // bare mode, no translation or protection
		return true, vAddr
	}

	vpn0 := vAddr >> 12

	pte := uint32(0)
	i := 0

	if present, p := c.mc.tlb0.load(vpn0); present {
		// normal page
		pte = p
		i = 0
	} else {
		fmt.Println("Missed")
		// not in tlb, walk table
		j, p := c.walkTable(vpn0)
		i = j

		if j < 0 {
			success = false
		}

		pte = p

		// perform checks after walking the page table.
		// This ensures that only valid entries are stored in the TLB, ensuring
		// the hot-path stays largely without much checking.

		// 5. A leaf PTE has been found. Determine if the requested memory access
		// is allowed by the pte.r, pte.w, pte.x, and pte.u bits, given the current
		// privilege mode and the value of the SUM and MXR fields of the mstatus
		// register. If not, stop and raise a page-fault exception corresponding to
		// the original access type.

		// 7. If pte.a = 0, or if the original memory access is a store and pte.d = 0,
		// either raise a page-fault exception corresponding to the original access
		// type...

		success = success && pte&pageFlagUser != 0

		if aType == accessTypeInstructionFetch && success {
			success = pte&pageFlagExec != 0 && pte&pageFlagAccessed != 0
		} else if aType == accessTypeLoad && success {
			success = pte&pageFlagRead != 0 && pte&pageFlagAccessed != 0
		} else if aType == accessTypeStore && success {
			success = pte&pageFlagWrite != 0 && pte&pageFlagAccessed != 0 && pte&pageFlagDirty != 0
		}

		if !success {
			c.csr[Csr_MTVAL] = vAddr
			switch aType {
			case accessTypeLoad:
				c.trap(TrapLoadPageFault)
			case accessTypeStore:
				c.trap(TrapStorePageFault)
			case accessTypeInstructionFetch:
				c.trap(TrapInstructionPageFault)
			default:
				panic("Invalid access type.")
			}
			return false, 0
		}

		// store the pte in table as well
		if j == 0 {
			c.mc.tlb0.store(vpn0, pte)
		} else if j == 1 {
			// treat superpage entries as several entries of normal pages instead
			pte |= (vAddr & 0x003FF000) >> 2 // edit the PTE
			c.mc.tlb0.store(vpn0, pte)
		}

		pAddr = ((pte & 0xFFFFFC00) << 2) | (vAddr & 0x00000FFF)
	}

	// 8. The translation is successful. The translated physical address is given
	// as follows: pa.pgoff = va.pgoff. If i > 0, then this is a superpage
	// translation and
	// pa.ppn[i − 1 : 0] = va.vpn[i − 1 : 0]. pa.ppn[LEVELS − 1 : i] = pte.ppn[LEVELS − 1 : i].

	pAddr = ((pte & 0xFFFFFC00) << 2) | (vAddr & (0xFFFFFFFF >> (20 - 10*i)))

	return true, pAddr
}

// Walks the table that satp is currently pointing to
func (c *Core) walkTable(vpn uint32) (int, uint32) {
	// fmt.Println("walking table!")
	// fmt.Printf("%08x\n", vpn)
	satp := c.csr[Csr_SATP]

	// 1. Let a be satp.ppn × PAGESIZE, and let i = LEVELS − 1. (For Sv32,
	// PAGESIZE=2¹² and LEVELS=2.) The satp register must be active, i.e., the
	// effective privilege mode must be S-mode or U-mode.
	a := (satp & 0x003FFFFF) * pagesize
	i := 1
	for {
		// fmt.Println(a)
		// 2. Let pte be the value of the PTE at address a+va.vpn[i]×PTESIZE.
		// (For Sv32, PTESIZE=4.) If accessing pte violates a PMA or PMP check,
		// raise an access-fault exception corresponding to the original access
		// type.
		vpni := (vpn >> (10 * i)) & 0x3FF
		success, pte := c.loadWordPhysical(a + vpni*4)
		if !success {
			panic("Invalid")
		}

		// 3. If pte.v = 0, or if pte.r = 0 and pte.w = 1, or if any bits or
		// encodings that are reserved for future standard use are set within
		// pte, stop and raise a page-fault exception corresponding to the
		// original access type.
		if pte&pageFlagValid == 0 || (pte&pageFlagRead == 0 && pte&pageFlagWrite == 1) {
			return i, 0
		}

		// fmt.Println("Page valid")

		// 4. Otherwise, the PTE is valid. If pte.r = 1 or pte.x = 1, go to
		// step 5.  Otherwise, this PTE is a pointer to the next level of the
		// page table. Let i = i − 1. If i < 0, stop and raise a page-fault
		// exception corresponding to the original access type. Otherwise, let
		// a = pte.ppn × PAGESIZE and go to step 2.
		if pte&pageFlagRead == 0 && pte&pageFlagExec == 0 {
			i = i - 1
			if i < 0 {
				return i, 0
			}
			a = (pte >> 10) * pagesize
			// fmt.Println("Next level")
			continue
		}

		// 6. If i > 0 and pte.ppn[i − 1 : 0] ̸= 0, this is a misaligned
		// superpage; stop and raise a page-fault exception corresponding to
		// the original access type.
		if i > 0 && pte&0x000FFC00 != 0 {
			return i, 0
		}

		// fmt.Printf("%08x\n", pte)
		return i, pte
		// do checking in translate
	}
}
