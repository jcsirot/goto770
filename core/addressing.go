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
		c.x = DWord(value)
	case 1:
		c.y = DWord(value)
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
			address += uint16(32 - offset)
		} else {
			address += uint16(offset)
		}
	} else if postbyte&0x0f == 0x01 {
		/* idxinc1 - Autoincrement by 1 from Register */
		address = c.readIndexedRegister(postbyte)
		c.writeIndexedRegister(postbyte, address+1)
		c.clock += 2
	} else if postbyte&0x0f == 0x02 {
		/* idxinc2 - Autoincrement by 2 from Register */
		address = c.readIndexedRegister(postbyte)
		c.writeIndexedRegister(postbyte, address+2)
		c.clock += 3
	}
	return address
}