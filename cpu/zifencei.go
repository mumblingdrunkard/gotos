// This file contains implementations of the instructions specified in
// the Zifencei extension of the RISC-V unprivileged specification.
//   Refer to the specification for instruction documentation.

package cpu

func (c *Core) fence_i(inst uint32) {
	c.FENCE_I()
}
