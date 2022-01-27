package system

// The paging system is "inspired" by the Sv32 _mode_ from the RISC-V
// privileged specification, chapter 4.3. This is a two-level translation
// scheme.

// From the RISC-V privileged specification:
//
// >           31  20   19  10   9 8   7   6   5   4   3   2   1   0
// >         | PPN[1] | PPN[0] | RSW | D | A | G | U | X | W | R | V |
// >             12       10      2    1   1   1   1   1   1   1   1
// >
// >                    Figure 4.18: Sv32 page table entry
// >
// > The PTE format for Sv32 is shown in Figure 4.18. The V bit indicates
// > whether the PTE is valid; if it is 0, all other bits in the PTE are
// > don't-cares and may be used freely by software. The permission bits, R, W,
// > and X, indicate whether the page is readable, writable, and executable,
// > respectively. When all three are zero, the PTE is a pointer to the next
// > level of the page table; otherwise, it is a leaf PTE. Writable pages must
// > also be marked readable; the contrary combinations are reserved for future
// > use. Table 4.5 summarizes the encoding of the permission bits.
// >
// >           | X | W | R || Meaning                              |
// >           |---|---|---||--------------------------------------|
// >           | 0 | 0 | 0 || Pointer to next level of page table. |
// >           | 0 | 0 | 1 || Read-only page.                      |
// >           | 0 | 1 | 0 || _Reserved for future use_.           |
// >           | 0 | 1 | 1 || Read-write page.                     |
// >           | 1 | 0 | 0 || Execute-only page.                   |
// >           | 1 | 0 | 1 || Read-execute page.                   |
// >           | 1 | 1 | 0 || _Reserved for future use_.           |
// >           | 1 | 1 | 1 || Read-write-execute page.             |
// >
// >                  Table 4.5: Encoding of PTE R/W/X fields.
//
// The "Pointer to the next level of the page table." is interesting to us.
// We won't be storing pages in memory that is generally accessible to the core.
// I suppose we instead store an index to a page table in a table of tables.
//
// > The U bit indicates whether the page is accessible to user mode. U-mode
// > software may only access the page when U=1. If the SUM bit in the sstatus
// > register is set, supervisor mode software may also access pages with U=1.
// > However, supervisor code normally operates with the SUM bit clear, in
// > which case, supervisor code will fault on accesses to user-mode pages.
// > Irrespective of SUM, the supervisor may not execute code on pages with
// > U=1.

// I expect we mark all pages as U-mode pages.
// Higher privilege levels wouldn't use this memory.
// The U bit could perhaps be used for something else.

// > The G bit designates a global mapping. Global mappings are those that
// > exist in all address spaces. For non-leaf PTEs, the global setting implies
// > that all mappings in the subsequent levels of the page table are global.
// > Note that failing to mark a global mapping as global merely reduces
// > performance, whereas marking a non-global mapping as global is a software
// > bug that, after switching to an address space with a different non-global
// > mapping for that address range, can unpredictably result in either mapping
// > being used.

// I don't think we'll use the G bit either, but it would be interesting.
// It could be used to map things like IO into the address space of all processes,
// but I think this is better done dynamically.

// > The RSW field is reserved for use by supervisor software; the
// > implementation shall ignore this field.
// >
// > Each leaf PTE contains an accessed (A) and dirty (D) bit. The A bit
// > indicates the virtual page has been read, written, or fetched from since
// > the last time the A bit was cleared. The D bit indicates the virtual page
// > has been written since the last time the D bit was cleared.
// >
// > Two schemes to manage the A and D bits are permitted:
// >
// >   * When a virtual page is accessed and the A bit is clear, or is written
// >   and the D bit is clear, a page-fault exception is raised.
// >
// >   * When a virtual page is accessed and the A bit is clear, or is written
// >   and the D bit is clear, the implementation sets the corresponding bit(s)
// >   in the PTE. The PTE update must be atomic with respect to other accesses
// >   to the PTE, and must atomically check that the PTE is valid and grants
// >   sufficient permissions. Updates of the A bit may be performed as a
// >   result of speculation, but updates to the D bit must be exact (i.e., not
// >   speculative), and observed in program order by the local hart.
// >   Furthermore, the PTE update must appear in the global memory order no
// >   later than the explicit memory access, or any subsequent explicit memory
// >   access to that virtual page by the local hart. The ordering on loads and
// >   stores provided by FENCE instructions and the acquire/release bits on
// >   atomic instructions also orders the PTE updates associated with those
// >   loads and stores as observed by remote harts. The PTE update is not
// >   required to be atomic with respect to the explicit memory access that
// >   caused the update, and the sequence is interruptible. However, the hart
// >   must not perform the explicit memory access before the PTE update is
// >   globally visible.
// >
// > All harts in a system must employ the same PTE-update scheme as each other.
// >

// I propose we use the first scheme.
// This puts the responsibility on the system to mark the pages for certain operations.
// Some notes are provided about the A and D bits:

// >   -------------------------------------------------------------------------
// >   Prior versions of this specification required PTE A bit updates to be
// >   exact, but allowing the A bit to be updated as a result of speculation
// >   simplifies the implementation of address translation prefetchers. System
// >   software typically uses the A bit as a page replacement policy hint, but
// >   does not require exactness for functional correctness. On the other
// >   hand, D bit updates are still required to be exact and performed in
// >   program order, as the D bit affects the functional correctness of page
// >   eviction.
//
// >   Implementations are of course still permitted to perform both A and D
// >   bit updates only in an exact manner.
//
// >   In both cases, requiring atomicity ensures that the PTE update will not
// >   be interrupted by other intervening writes to the page table, as such
// >   interruptions could lead to A/D bits being set on PTEs that have been
// >   reused for other purposes, on memory that has been reclaimed for other
// >   purposes, and so on. Simple implementations may instead generate
// >   page-fault exceptions.
// >
// >   The A and D bits are never cleared by the implementation. If the
// >   supervisor software does not rely on accessed and/or dirty bits, e.g. if
// >   it does not swap memory pages to secondary storage or if the pages are
// >   being used to map I/O space, it should always set them to 1 in the PTE
// >   to improve performance.
// >   -------------------------------------------------------------------------

// It is likely a good idea to always set the A and D bits.
// I don't plan on supporting page swapping in the first go.

// > Any level of PTE may be a leaf PTE, so in addition to 4 KiB pages, Sv32
// > supports 4 MiB megapages. A megapage must be virtually and physically
// > aligned to a 4 MiB boundary; a page-fault exception is raised if the
// > physical address is insufficiently aligned.

// May need to have 2 TLB's per core here, one for megapages and one for normal pages.
// Check the megapage TLB first, then the normal TLB if the megapage TLB does not contain the translation.

// > For non-leaf PTEs, the D, A, and U bits are reserved for future standard
// > use. Until their use is defined by a standard extension, they must be
// > cleared by software for forward compatibility.
// >
// > For implementations with both page-based virtual memory and the “A”
// > standard extension, the LR/SC reservation set must lie completely within a
// > single base page (i.e., a naturally aligned 4 KiB region).
// >
// >                                  --- The RISC-V Instruction Set Manual,
// >                                        Volume II: Privileged Architecture,
// >                                        Version 20211203,
// >                                        Section 4.3.1.

// The specification also gives a short guide for address translation using
// this scheme:

const (
	pageFlagValid    uint32 = 0x01 // the virtual address is valid
	pageFlagRead            = 0x02 // indicates that the processor is allowed to read data from this address
	pageFlagWrite           = 0x04 // indicates that the processor is allowed to write data to this address
	pageFlagExec            = 0x08 // indicates that the processor is allowed to fetch instructions from this address
	pageFlagUser            = 0x10
	pageFlagGlobal          = 0x20
	pageFlagAccessed        = 0x40
	pageFlagDirty           = 0x80
)

type pte uint32 // page-table entry

type ptable struct {
	ptes [1024]pte
}

func testmain() {
	// example for page table management
}
