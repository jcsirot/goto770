package core

type ccr struct {
	r8
}

/*******************************/
/**     CC register flags     **/
/*******************************/

func (cc *ccr) getE() bool {
	cc.get()
	return *cc.r8.r&entire == entire
}

func (cc *ccr) setE() {
	*cc.r8.r |= entire
}

func (cc *ccr) getH() bool {
	return *cc.r8.r&halfCarry == halfCarry
}

func (cc *ccr) setH() {
	*cc.r8.r |= halfCarry
}

func (cc *ccr) clearH() {
	*cc.r8.r &= 0xff ^ halfCarry
}

func (cc *ccr) getC() bool {
	return *cc.r8.r&carry == carry
}

func (cc *ccr) setC() {
	*cc.r8.r |= carry
}

func (cc *ccr) clearC() {
	*cc.r8.r &= 0xff ^ carry
}

func (cc *ccr) getZ() bool {
	return *cc.r8.r&zero == zero
}

func (cc *ccr) setZ() {
	*cc.r8.r |= zero
}

func (cc *ccr) clearZ() {
	*cc.r8.r &= 0xff ^ zero
}

func (cc *ccr) getN() bool {
	return *cc.r8.r&negative == negative
}

func (cc *ccr) setN() {
	*cc.r8.r |= negative
}

func (cc *ccr) clearN() {
	*cc.r8.r &= 0xff ^ negative
}

func (cc *ccr) getV() bool {
	return *cc.r8.r&overflow == overflow
}

func (cc *ccr) setV() {
	*cc.r8.r |= overflow
}

func (cc *ccr) clearV() {
	*cc.r8.r &= 0xff ^ overflow
}

func (cc *ccr) getF() bool {
	return *cc.r8.r&firqmask == firqmask
}

func (cc *ccr) setF() {
	*cc.r8.r |= firqmask
}

func (cc *ccr) getI() bool {
	return *cc.r8.r&irqmask == irqmask
}

func (cc *ccr) setI() {
	*cc.r8.r |= irqmask
}

func (c *CPU) updateNZVC(a, b, r int) {
	c.updateZ(r)
	c.updateN(r)
	c.updateC(r&0x100 != 0)
	c.updateV(((a ^ b ^ r ^ (r >> 1)) & 0x80) != 0)
}

func (c *CPU) updateHNZVC(a, b, r int) {
	c.updateH(((a ^ b ^ r) & 0x10) != 0)
	c.updateZ(r)
	c.updateN(r)
	c.updateC(r&0x100 != 0)
	c.updateV(((a ^ b ^ r ^ (r >> 1)) & 0x80) != 0)
}

func (c *CPU) updateNZVC16(a, b, r int) {
	c.updateZ16(r)
	c.updateN16(r)
	c.updateC(r&0x10000 != 0)
	c.updateV(((a ^ b ^ r ^ (r >> 1)) & 0x8000) != 0)
}

func (c *CPU) updateNZ(r int) {
	c.updateZ(r)
	c.updateN(r)
}

func (c *CPU) updateNZ16(r int) {
	c.updateZ16(r)
	c.updateN16(r)
}

func (c *CPU) updateV(value bool) {
	if value {
		c.cc.setV()
	} else {
		c.cc.clearV()
	}
}

func (c *CPU) updateZ(value int) {
	if value&0xff == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
}

func (c *CPU) updateZ16(value int) {
	if value&0xffff == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
}

func (c *CPU) updateN(value int) {
	if value&0x80 != 0 {
		c.cc.setN()
	} else {
		c.cc.clearN()
	}
}

func (c *CPU) updateN16(value int) {
	if value&0x8000 != 0 {
		c.cc.setN()
	} else {
		c.cc.clearN()
	}
}

func (c *CPU) updateC(value bool) {
	if value {
		c.cc.setC()
	} else {
		c.cc.clearC()
	}
}

func (c *CPU) updateH(value bool) {
	if value {
		c.cc.setH()
	} else {
		c.cc.clearH()
	}
}
