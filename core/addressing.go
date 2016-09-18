package core

import "fmt"

func (c *CPU) direct() uint16 {
	ea := uint16(c.dp)<<8 | uint16(c.read(c.pc))
	c.pc++
	return ea
}

func (c *CPU) immediate() uint16 {
	ea := c.pc
	c.pc++
	return ea
}

func (c *CPU) extended() uint16 {
	ea := uint16(c.readw(c.pc))
	c.pc += 2
	return ea
}

func (c *CPU) relative() uint16 {
	offset := int8(c.read(c.pc))
	c.pc++
	address := c.pc
	if offset < 0 {
		address -= uint16(-offset)
	} else {
		address += uint16(offset)
	}
	return address
}

func (c *CPU) lrelative() uint16 {
	offset := int16(c.readw(c.pc))
	c.pc += 2
	address := c.pc
	if offset < 0 {
		address -= uint16(-offset)
	} else {
		address += uint16(offset)
	}
	return address
}

func (c *CPU) indexed() uint16 {
	postbyte := uint8(c.read(c.pc))
	c.pc++
	ea := c.getIndexedAddress(postbyte)
	if uint8(postbyte)&0x90 == 0x90 { // Indirect mode?
		c.clock += 3
		ea = uint16(c.readw(ea))
	}
	return ea
}

func (c *CPU) readIndexedRegister(postbyte uint8) uint16 {
	code := (postbyte & 0x60) >> 5
	switch code {
	case 0:
		return uint16(c.x)
	case 1:
		return uint16(c.y)
	case 2:
		return uint16(c.u)
	case 3:
		return uint16(c.s)
	default:
		panic(fmt.Sprintf("Undefined indexed addressing mode register code %d at pc=%x", code, c.pc))
	}
}

func (c *CPU) writeIndexedRegister(postbyte uint8, value uint16) {
	code := (postbyte & 0x60) >> 5
	switch code {
	case 0:
		c.x = value
	case 1:
		c.y = value
	case 2:
		c.u = value
	case 3:
		c.s = value
	default:
		panic(fmt.Sprintf("Undefined indexed addressing mode register code %d at pc=%x", code, c.pc))
	}
}

func (c *CPU) getIndexedAddress(postbyte uint8) uint16 {
	var address uint16
	if postbyte&0x80 == 0 {
		/* idx5off - 5 bits offset from Register */
		address = c.readIndexedRegister(postbyte)
		c.clock++
		offset := postbyte & 0x1f
		if offset > 0x0f {
			address -= uint16(32 - offset)
		} else {
			address += uint16(offset)
		}
	} else if postbyte&0x0f == 0x00 {
		/* idxinc1 - Autoincrement by 1 from Register */
		address = c.readIndexedRegister(postbyte)
		c.writeIndexedRegister(postbyte, address+1)
		c.clock += 2
	} else if postbyte&0x0f == 0x01 {
		/* idxinc2 - Autoincrement by 2 from Register */
		address = c.readIndexedRegister(postbyte)
		c.writeIndexedRegister(postbyte, address+2)
		c.clock += 3
	} else if postbyte&0x0f == 0x02 {
		/* idxdec1 - Autodecrement by 1 from Register */
		address = c.readIndexedRegister(postbyte) - 1
		c.writeIndexedRegister(postbyte, address)
		c.clock += 2
	} else if postbyte&0x0f == 0x03 {
		/* Autodecrement by 2 from Register */
		address = c.readIndexedRegister(postbyte) - 2
		c.writeIndexedRegister(postbyte, address)
		c.clock += 3
	} else if postbyte&0x0f == 0x04 {
		/* No Offset from Register */
		address = c.readIndexedRegister(postbyte)
	} else if postbyte&0x0f == 0x05 {
		/* idxb - B Accumulator Offset from Register */
		address = c.readIndexedRegister(postbyte)
		offset := int8(c.b)
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock++
	} else if postbyte&0x0f == 0x06 {
		/* idxa - A Accumulator Offset from Register */
		address = c.readIndexedRegister(postbyte)
		offset := int8(c.a)
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock++
	} else if postbyte&0x0f == 0x08 {
		/* idx8off - 8 bits offset from Register */
		address = c.readIndexedRegister(postbyte)
		offset := int8(c.read(c.pc))
		c.pc++
		c.clock++
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
	} else if postbyte&0x0f == 0x09 {
		/* idx8off - 16 bits offset from Register */
		address = c.readIndexedRegister(postbyte)
		offset := int16(c.readw(c.pc))
		c.pc += 2
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock += 4
	} else if postbyte&0x0f == 0x0b {
		/* idxd - D Accumulator Offset from Register */
		address = c.readIndexedRegister(postbyte)
		offset := int16(c.d())
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock += 4
	} else if postbyte&0x0f == 0x0c {
		/* idxpc8 - 8 bits Offset from Program Counter */
		offset := int8(c.read(c.pc))
		c.pc++
		address = c.pc
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock++
	} else if postbyte&0x0f == 0x0d {
		/* idxpc16 - 16 bits Offset from Program Counter */
		offset := int16(c.read(c.pc))
		c.pc += 2
		address = c.pc
		if offset >= 0 {
			address += uint16(offset)
		} else {
			address -= uint16(-offset)
		}
		c.clock += 5
	} else if postbyte&0x0f == 0x0f {
		/* idxext - Extended Indirect */
		address = uint16(c.readw(c.pc))
		c.pc += 2
		c.clock += 2
	} else {
		panic(fmt.Sprintf("Undefined indexed submode code %d at pc=%x", postbyte, c.pc))
	}
	return address
}
