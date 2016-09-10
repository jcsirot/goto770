package core

type Opcode struct {
	f      func()
	cycles uint64
}

var (
	opcodes [0x300]Opcode
)

const (
	carry     = 1 << iota // 0x01
	overflow  = 1 << iota // 0x02
	zero      = 1 << iota // 0x04
	negative  = 1 << iota // 0x08
	irqmask   = 1 << iota // 0x10
	halfCarry = 1 << iota // 0x20
	firqmask  = 1 << iota // 0x40
	entire    = 1 << iota // 0x80
	none      = 0
)

/*
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
) */

/* Addressing modes */

func (c *CPU) direct() uint16 {
	ea := uint16(c.dp)<<8 | uint16(c.read(c.pc))
	c.pc++
	return ea
}

func (c *CPU) extended() uint16 {
	ea := uint16(c.readw(c.pc))
	c.pc += 2
	return ea
}

// Cpu structure
type CPU struct {
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
	/// Memory
	ram Memory
	///
	clock uint64
}

// Initialize the Cpu
func (c *CPU) Initialize(ram Memory) {
	c.ram = ram
	c.Reset()
	c.initOpcodes()
}

func (c *CPU) Reset() {
	c.a = 0
	c.b = 0
	c.x = 0
	c.y = 0
	c.u = 0
	c.s = 0
	c.dp = 0
	c.cc = 0
	c.pc = 0
	c.clock = 0
}

func (c *CPU) initOpcodes() {
	opcodes[0x00] = Opcode{func() { c.neg(c.direct()) }, 6}
	opcodes[0x03] = Opcode{func() { c.com(c.direct()) }, 6}
	opcodes[0x04] = Opcode{func() { c.lsr(c.direct()) }, 6}
	opcodes[0x06] = Opcode{func() { c.ror(c.direct()) }, 6}
	opcodes[0x07] = Opcode{func() { c.asr(c.direct()) }, 6}
	opcodes[0x08] = Opcode{func() { c.asl(c.direct()) }, 6}
	opcodes[0x09] = Opcode{func() { c.rol(c.direct()) }, 6}
	//opcodes[0x0a] = Opcode{func() { c.dec(c.direct()) }, 6}
	opcodes[0x40] = Opcode{func() { c.nega() }, 2}
	opcodes[0x43] = Opcode{func() { c.coma() }, 2}
	opcodes[0x44] = Opcode{func() { c.lsra() }, 2}
	opcodes[0x46] = Opcode{func() { c.rora() }, 2}
	opcodes[0x47] = Opcode{func() { c.asra() }, 2}
	opcodes[0x48] = Opcode{func() { c.asla() }, 2}
	opcodes[0x49] = Opcode{func() { c.rola() }, 2}
	opcodes[0x50] = Opcode{func() { c.negb() }, 2}
	opcodes[0x53] = Opcode{func() { c.comb() }, 2}
	opcodes[0x54] = Opcode{func() { c.lsrb() }, 2}
	opcodes[0x56] = Opcode{func() { c.rorb() }, 2}
	opcodes[0x57] = Opcode{func() { c.asrb() }, 2}
	opcodes[0x58] = Opcode{func() { c.aslb() }, 2}
	opcodes[0x59] = Opcode{func() { c.rolb() }, 2}
	//opcodes[0x60] = Opcode{func() { c.com(c.indexed()) }, 6}
	//opcodes[0x63] = Opcode{func() { c.com(c.indexed()) }, 6}
	//opcodes[0x64] = Opcode{func() { c.lsr(c.extended()) }, 6}
	opcodes[0x70] = Opcode{func() { c.neg(c.extended()) }, 7}
	opcodes[0x73] = Opcode{func() { c.com(c.extended()) }, 7}
	opcodes[0x74] = Opcode{func() { c.lsr(c.extended()) }, 7}
}

func (c *CPU) step() uint64 {
	opcode := opcodes[c.read(c.pc)]
	c.pc++

	/*
		if c.Verbose {
			Disassemble(opcode, c, c.ProgramCounter)
		}
	*/

	opcode.f()
	c.clock += opcode.cycles

	return opcode.cycles

}

func (c *CPU) getC() bool {
	return c.cc&carry == carry
}

func (c *CPU) setC() {
	c.cc |= carry
}

func (c *CPU) clearC() {
	c.cc &= 0xff ^ carry
}

func (c *CPU) updateC(value bool) {
	if value {
		c.setC()
	} else {
		c.clearC()
	}
}

func (c *CPU) getZ() bool {
	return c.cc&zero == zero
}

func (c *CPU) setZ() {
	c.cc |= zero
}

func (c *CPU) clearZ() {
	c.cc &= 0xff ^ zero
}

func (c *CPU) testSetZ(value Word) {
	if value == 0 {
		c.setZ()
	} else {
		c.clearZ()
	}
}

func (c *CPU) getN() bool {
	return c.cc&negative == negative
}

func (c *CPU) setN() {
	c.cc |= negative
}

func (c *CPU) clearN() {
	c.cc &= 0xff ^ negative
}

func (c *CPU) testSetN(value Word) {
	if value&0x80 == 0x80 {
		c.setN()
	} else {
		c.clearN()
	}
}

func (c *CPU) testSetZN(value Word) {
	c.testSetZ(value)
	c.testSetN(value)
}

func (c *CPU) getV() bool {
	return c.cc&overflow == overflow
}

func (c *CPU) setV() {
	c.cc |= overflow
}

func (c *CPU) clearV() {
	c.cc &= 0xff ^ overflow
}

func (c *CPU) updateV(value bool) {
	if value {
		c.setV()
	} else {
		c.clearV()
	}
}

func (c *CPU) read(address uint16) Word {
	return c.ram[address]
}

func (c *CPU) write(address uint16, value Word) {
	c.ram[address] = value
}

func (c *CPU) readw(address uint16) DWord {
	hi := c.read(address)
	lo := c.read(address + 1)
	return (DWord)(uint16(hi)<<8 | uint16(lo))
}

/** Negate - H?NxZxVxCx */
func (c *CPU) neg_(value Word) Word {
	tmp := -value
	c.testSetZN(tmp)
	c.updateC(value != 0)
	c.updateV(value == 0x80)
	return tmp
}

/** Negate - H?NxZxVxCx */
func (c *CPU) neg(address uint16) {
	c.write(address, c.neg_(c.read(address)))
}

/** Negate Register A - H?NxZxVxCx */
func (c *CPU) nega() {
	c.a = c.neg_(c.a)
}

/** Negate Register B - H?NxZxVxCx */
func (c *CPU) negb() {
	c.b = c.neg_(c.b)
}

/** Complement - H?NxZxV0C1 */
func (c *CPU) com_(value Word) Word {
	tmp := value ^ 0xff
	c.testSetZN(tmp)
	c.setC()
	c.clearV()
	return tmp
}

/** Complement - H?NxZxV0C1 */
func (c *CPU) com(address uint16) {
	c.write(address, c.com_(c.read(address)))
}

/** Complement Register A - H?NxZxV0C1 */
func (c *CPU) coma() {
	c.a = c.com_(c.a)
}

/** Complement Register B - H?NxZxV0C1 */
func (c *CPU) comb() {
	c.b = c.com_(c.b)
}

/** Logical Shift Right - N0ZxCx */
func (c *CPU) lsr_(value Word) Word {
	tmp := value >> 1
	c.testSetZN(tmp)
	c.updateC(value&1 == 1)
	return tmp
}

/** Logical Shift Right - N0ZxCx */
func (c *CPU) lsr(address uint16) {
	c.write(address, c.lsr_(c.read(address)))
}

/** Logical Shift Right A Register - N0ZxCx */
func (c *CPU) lsra() {
	c.a = c.lsr_(c.a)
}

/** Logical Shift Right B Register - N0ZxCx */
func (c *CPU) lsrb() {
	c.b = c.lsr_(c.b)
}

/** Rotate Right - NxZxCx */
func (c *CPU) ror_(value Word) Word {
	tmp := (value >> 1) | (value << 7)
	c.testSetZN(tmp)
	c.updateC(value&1 == 1)
	return tmp
}

/** Rotate Right - NxZxCx */
func (c *CPU) ror(address uint16) {
	c.write(address, c.ror_(c.read(address)))
}

/** Rotate Right Register A - NxZxCx */
func (c *CPU) rora() {
	c.a = c.ror_(c.a)
}

/** Rotate Right Register B - NxZxCx */
func (c *CPU) rorb() {
	c.b = c.ror_(c.b)
}

/** Rotate Left - NxZxVxCx */
func (c *CPU) rol_(value Word) Word {
	tmp := (value << 1) | ((value >> 7) & 0x01)
	c.testSetZN(tmp)
	c.updateC(value>>7 == 0x01)
	c.updateV((value>>7)^((value>>6)&0x01) == 0x01)
	return tmp
}

/** Rotate Left - NxZxVxCx */
func (c *CPU) rol(address uint16) {
	c.write(address, c.rol_(c.read(address)))
}

/** Rotate Left Register A - NxZxVxCx */
func (c *CPU) rola() {
	c.a = c.rol_(c.a)
}

/** Rotate Left Register B - NxZxVxCx */
func (c *CPU) rolb() {
	c.b = c.rol_(c.b)
}

/** Arithmetic Shift Right - H?NxZxCx */
func (c *CPU) asr_(value Word) Word {
	tmp := (value >> 1) | (value & 0x80)
	c.testSetZN(tmp)
	c.updateC(value&0x01 == 0x01)
	return tmp
}

/** Arithmetic Shift Right - H?NxZxCx */
func (c *CPU) asr(address uint16) {
	c.write(address, c.asr_(c.read(address)))
}

/** Arithmetic Shift Right Register A - H?NxZxCx */
func (c *CPU) asra() {
	c.a = c.asr_(c.a)
}

/** Arithmetic Shift Right Register B - H?NxZxCx */
func (c *CPU) asrb() {
	c.b = c.asr_(c.b)
}

/** Arithmetic Shift Left / Logical Shift Left - H?NxZxVxCx */
func (c *CPU) asl_(value Word) Word {
	tmp := value << 1
	c.testSetZN(tmp)
	c.updateC(value&0x80 == 0x80)
	c.updateV((value>>7)^((value>>6)&0x01) == 0x01)
	return tmp
}

/** Arithmetic Shift Left / Logical Shift Left - H?NxZxVxCx */
func (c *CPU) asl(address uint16) {
	c.write(address, c.asl_(c.read(address)))
}

/** Arithmetic Shift Left / Logical Shift Left Register A - H?NxZxVxCx */
func (c *CPU) asla() {
	c.a = c.asl_(c.a)
}

/** Arithmetic Shift Left / Logical Shift Left Register B - H?NxZxVxCx */
func (c *CPU) aslb() {
	c.b = c.asl_(c.b)
}
