package core

var (
	opcodes [0x300]func()
)

type AddressingMode uint

const (
	Carry     = 1 << iota // 0x01
	Overflow  = 1 << iota // 0x02
	Zero      = 1 << iota // 0x04
	Negative  = 1 << iota // 0x08
	IRQ_mask  = 1 << iota // 0x10
	HalfCarry = 1 << iota // 0x20
	FIRQ_mask = 1 << iota // 0x40
	Entire    = 1 << iota // 0x80
	None      = 0
)

const (
	_ = iota
	Direct
	Inherent
	Immediate
	Limmediate
	Relative
	Lrelative
	Extended
	Indexed
)

func (c *Cpu) Direct() uint16 {
	c.pc++
	return (uint16)((cpu.dp << 8) | ram.read(cpu.pc-1))
}

type Cpu struct {
	/// Accumulator register
	a Word
	/// Accumulator register
	b Word
	/// Index register
	x DWord
	/// Index register
	y DWord
	/// User stack pointer register
	u DWord
	/// Hardware stack pointer register
	s DWord
	/// Direct page register
	dp Word
	/// Condition code regsiter
	cc uint8
	/// Program counter register
	pc uint16
}

func (c *Cpu) Initialize() {
	c.Reset()
	c.initOpcodes()
}

func (c *Cpu) Reset() {
	c.a = 0
	c.b = 0
	c.x = 0
	c.y = 0
	c.u = 0
	c.s = 0
	c.dp = 0
	c.cc = 0
	c.pc = 0
}

func (c *Cpu) initOpcodes() {
	opcodes[0x00] = func() {
		c.neg(c.Direct())
	}
}

func (c *Cpu) getCarry() bool {
	return c.cc&Carry == 0x01
}

func (c *Cpu) setCarry() {
	c.cc |= Carry
}

func (c *Cpu) clearCarry() {
	c.cc &= 0xff ^ Carry
}

func (c *Cpu) getZero() bool {
	return c.cc&Zero == 0x01
}

func (c *Cpu) setZero() {
	c.cc |= Zero
}

func (c *Cpu) clearZero() {
	c.cc &= 0xff ^ Zero
}

func (c *Cpu) testAndSetZero(value Word) {
	if value == 0 {
		c.setZero()
	} else {
		c.clearZero()
	}
}

func (c *Cpu) getNegative() bool {
	return c.cc&Negative == 0x01
}

func (c *Cpu) setNegative() {
	c.cc |= Negative
}

func (c *Cpu) clearNegative() {
	c.cc &= 0xff ^ Negative
}

func (c *Cpu) testAndSetNegative(value Word) {
	if value&0x80 == 0x80 {
		c.setNegative()
	} else {
		c.clearNegative()
	}
}

func (c *Cpu) read(address uint16) Word {
	return ram[address]
}

func (c *Cpu) write(address uint16, value Word) {
	ram[address] = value
}

/** Negate - H?NxZxVxCx */
func (c *Cpu) neg(address uint16) {
	var value = c.read(address)
	var tmp = -value
	c.testAndSetZero(tmp)
	c.testAndSetNegative(tmp)
	c.write(address, tmp)
}
