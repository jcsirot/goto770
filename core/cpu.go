package core

import "fmt"

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
	u uint16
	/// Hardware stack pointer register
	s uint16
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

func (c *CPU) d() DWord {
	return DWord(int(c.a)<<8 | int(c.b))
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
	opcodes[0x0a] = Opcode{func() { c.dec(c.direct()) }, 6}
	opcodes[0x0c] = Opcode{func() { c.inc(c.direct()) }, 6}
	opcodes[0x0d] = Opcode{func() { c.tst(c.direct()) }, 6}
	opcodes[0x0e] = Opcode{func() { c.jmp(c.direct()) }, 3}
	opcodes[0x0f] = Opcode{func() { c.clr(c.direct()) }, 6}
	opcodes[0x12] = Opcode{func() { c.nop() }, 2}
	opcodes[0x13] = Opcode{func() { c.sync() }, 4}
	opcodes[0x16] = Opcode{func() { c.bra(c.lrelative()) }, 5}
	opcodes[0x17] = Opcode{func() { c.bsr(c.lrelative()) }, 9}
	opcodes[0x19] = Opcode{func() { c.daa() }, 2}
	opcodes[0x1a] = Opcode{func() { c.orcc(c.immediate()) }, 3}
	opcodes[0x1c] = Opcode{func() { c.andcc(c.immediate()) }, 3}
	opcodes[0x1d] = Opcode{func() { c.sex() }, 2}
	opcodes[0x1e] = Opcode{func() { c.exg(c.immediate()) }, 8}
	opcodes[0x1f] = Opcode{func() { c.tfr(c.immediate()) }, 6}
	opcodes[0x20] = Opcode{func() { c.bra(c.relative()) }, 3}
	opcodes[0x21] = Opcode{func() { c.brn(c.relative()) }, 3}
	opcodes[0x22] = Opcode{func() { c.bhi(c.relative()) }, 3}
	opcodes[0x23] = Opcode{func() { c.bls(c.relative()) }, 3}
	opcodes[0x24] = Opcode{func() { c.bcc(c.relative()) }, 3}
	opcodes[0x25] = Opcode{func() { c.blo(c.relative()) }, 3}
	opcodes[0x26] = Opcode{func() { c.bne(c.relative()) }, 3}
	opcodes[0x27] = Opcode{func() { c.beq(c.relative()) }, 3}
	opcodes[0x28] = Opcode{func() { c.bvc(c.relative()) }, 3}
	opcodes[0x29] = Opcode{func() { c.bvs(c.relative()) }, 3}
	opcodes[0x2a] = Opcode{func() { c.bpl(c.relative()) }, 3}
	opcodes[0x2b] = Opcode{func() { c.bmi(c.relative()) }, 3}
	opcodes[0x2c] = Opcode{func() { c.bge(c.relative()) }, 3}
	opcodes[0x2d] = Opcode{func() { c.blt(c.relative()) }, 3}
	opcodes[0x2e] = Opcode{func() { c.bgt(c.relative()) }, 3}
	opcodes[0x2f] = Opcode{func() { c.ble(c.relative()) }, 3}
	opcodes[0x30] = Opcode{func() { c.leax(c.indexed()) }, 4}
	opcodes[0x31] = Opcode{func() { c.leay(c.indexed()) }, 4}
	opcodes[0x32] = Opcode{func() { c.leas(c.indexed()) }, 4}
	opcodes[0x33] = Opcode{func() { c.leau(c.indexed()) }, 4}
	opcodes[0x34] = Opcode{func() { c.pshs(c.immediate()) }, 5}
	opcodes[0x35] = Opcode{func() { c.puls(c.immediate()) }, 5}
	opcodes[0x36] = Opcode{func() { c.pshu(c.immediate()) }, 5}
	opcodes[0x37] = Opcode{func() { c.pulu(c.immediate()) }, 5}
	opcodes[0x39] = Opcode{func() { c.rts() }, 5}
	opcodes[0x3a] = Opcode{func() { c.abx() }, 3}
	opcodes[0x3b] = Opcode{func() { c.rti() }, 3}
	opcodes[0x3d] = Opcode{func() { c.mul() }, 11}
	opcodes[0x3f] = Opcode{func() { c.swi() }, 7}
	opcodes[0x40] = Opcode{func() { c.nega() }, 2}
	opcodes[0x43] = Opcode{func() { c.coma() }, 2}
	opcodes[0x44] = Opcode{func() { c.lsra() }, 2}
	opcodes[0x46] = Opcode{func() { c.rora() }, 2}
	opcodes[0x47] = Opcode{func() { c.asra() }, 2}
	opcodes[0x48] = Opcode{func() { c.asla() }, 2}
	opcodes[0x49] = Opcode{func() { c.rola() }, 2}
	opcodes[0x4a] = Opcode{func() { c.deca() }, 2}
	opcodes[0x4c] = Opcode{func() { c.inca() }, 2}
	opcodes[0x4d] = Opcode{func() { c.tsta() }, 2}
	opcodes[0x4f] = Opcode{func() { c.clra() }, 2}
	opcodes[0x50] = Opcode{func() { c.negb() }, 2}
	opcodes[0x53] = Opcode{func() { c.comb() }, 2}
	opcodes[0x54] = Opcode{func() { c.lsrb() }, 2}
	opcodes[0x56] = Opcode{func() { c.rorb() }, 2}
	opcodes[0x57] = Opcode{func() { c.asrb() }, 2}
	opcodes[0x58] = Opcode{func() { c.aslb() }, 2}
	opcodes[0x59] = Opcode{func() { c.rolb() }, 2}
	opcodes[0x5a] = Opcode{func() { c.decb() }, 2}
	opcodes[0x5c] = Opcode{func() { c.incb() }, 2}
	opcodes[0x5d] = Opcode{func() { c.tstb() }, 2}
	opcodes[0x5f] = Opcode{func() { c.clrb() }, 2}
	opcodes[0x60] = Opcode{func() { c.com(c.indexed()) }, 6}
	opcodes[0x63] = Opcode{func() { c.com(c.indexed()) }, 6}
	opcodes[0x64] = Opcode{func() { c.lsr(c.extended()) }, 6}
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

/*******************************/
/**     CC register flags     **/
/*******************************/

func (c *CPU) getE() bool {
	return c.cc&entire == entire
}

func (c *CPU) setE() {
	c.cc |= entire
}

func (c *CPU) getH() bool {
	return c.cc&halfCarry == halfCarry
}

func (c *CPU) setH() {
	c.cc |= halfCarry
}

func (c *CPU) clearH() {
	c.cc &= 0xff ^ halfCarry
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

func (c *CPU) getF() bool {
	return c.cc&firqmask == firqmask
}

func (c *CPU) setF() {
	c.cc |= firqmask
}

func (c *CPU) getI() bool {
	return c.cc&irqmask == irqmask
}

func (c *CPU) setI() {
	c.cc |= irqmask
}

/***************************/
/**     Memory access     **/
/***************************/

func (c *CPU) read(address uint16) Word {
	return c.ram[address]
}

func (c *CPU) readw(address uint16) DWord {
	hi := c.read(address)
	lo := c.read(address + 1)
	return (DWord)(uint16(hi)<<8 | uint16(lo))
}

func (c *CPU) write(address uint16, value Word) {
	c.ram[address] = value
}

func (c *CPU) writew(address uint16, value DWord) {
	c.write(address+1, Word(value&0xff))
	c.write(address, Word((value >> 8)))
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

/** Decrement - NxZxVx */
func (c *CPU) dec_(value Word) Word {
	tmp := value - 1
	c.testSetZN(tmp)
	c.updateV(value == 0x80)
	return tmp
}

/** Decrement - NxZxVx */
func (c *CPU) dec(address uint16) {
	c.write(address, c.dec_(c.read(address)))
}

/** Decrement Register A - NxZxVx */
func (c *CPU) deca() {
	c.a = c.dec_(c.a)
}

/** Decrement Register B - NxZxVx */
func (c *CPU) decb() {
	c.b = c.dec_(c.b)
}

/** Increment - NxZxVx */
func (c *CPU) inc_(value Word) Word {
	tmp := value + 1
	c.testSetZN(tmp)
	c.updateV(value == 0x7f)
	return tmp
}

/** Increment - NxZxVx */
func (c *CPU) inc(address uint16) {
	c.write(address, c.inc_(c.read(address)))
}

/** Increment Register A - NxZxVx */
func (c *CPU) inca() {
	c.a = c.inc_(c.a)
}

/** Increment Register B - NxZxVx */
func (c *CPU) incb() {
	c.b = c.inc_(c.b)
}

/** Test - NxZxV0 */
func (c *CPU) tst_(value Word) {
	c.testSetZN(value)
	c.clearV()
}

/** Test - NxZxV0 */
func (c *CPU) tst(address uint16) {
	c.tst_(c.read(address))
}

/** Test Register A - NxZxV0 */
func (c *CPU) tsta() {
	c.tst_(c.a)
}

/** Test Register B - NxZxV0 */
func (c *CPU) tstb() {
	c.tst_(c.b)
}

/** Jump - NxZxV0 */
func (c *CPU) jmp(address uint16) {
	c.pc = address
}

/** Clear N0Z1V0C0 */
func (c *CPU) clr(address uint16) {
	c.write(address, 0)
	c.clearN()
	c.setZ()
	c.clearV()
	c.clearC()
}

/** Clear N0Z1V0C0 */
func (c *CPU) clra() {
	c.a = 0
	c.clearN()
	c.setZ()
	c.clearV()
	c.clearC()
}

/** Clear N0Z1V0C0 */
func (c *CPU) clrb() {
	c.b = 0
	c.clearN()
	c.setZ()
	c.clearV()
	c.clearC()
}

func (c *CPU) nop() {
}

/** Synchronize to External Event */
func (c *CPU) sync() {
	// Not supported
}

/** (Long) Branch Always */
func (c *CPU) bra(address uint16) {
	c.pc = address
}

/** Long Branch to Subroutine */
func (c *CPU) bsr(address uint16) {
	c.s -= 2
	c.writew(c.s, DWord(c.pc))
	c.pc = address
}

/** Decimal Addition Adjust - NxZxV?Cx */
func (c *CPU) daa() {
	ah := c.a & 0xf0
	al := c.a & 0x0f
	cf := 0
	if al > 0x09 || c.getH() {
		cf |= 0x06
	}
	if ah > 0x80 && al > 0x09 {
		cf |= 0x60
	}
	if ah > 0x90 || c.getC() {
		cf |= 0x60
	}
	tmp := uint16(c.a) + uint16(cf)
	c.a = Word(tmp)
	carry := c.getC()
	c.testSetZN(c.a)
	c.updateC(carry || tmp > 0xff)
}

/** Inclusive OR Memory Immediate into Condition Code Register */
func (c *CPU) orcc(address uint16) {
	value := c.read(address)
	c.cc |= uint8(value)
}

/** Logical AND Immediate Memory into Condition Code Register */
func (c *CPU) andcc(address uint16) {
	value := c.read(address)
	c.cc &= uint8(value)
}

/** Sign Extended - NxZx */
func (c *CPU) sex() {
	if c.b&0x80 == 0 {
		c.a = 0
	} else {
		c.a = 0xff
	}
	if c.d() == 0 {
		c.setZ()
	} else {
		c.clearZ()
	}
	if (c.d() & 0x8000) != 0 {
		c.setN()
	} else {
		c.clearN()
	}
}

func (c *CPU) getRegisterFromCode(code int) uint16 {
	switch code {
	case 0:
		return uint16(c.d())
	case 1:
		return uint16(c.x)
	case 2:
		return uint16(c.y)
	case 3:
		return uint16(c.u)
	case 4:
		return uint16(c.s)
	case 5:
		return uint16(c.pc)
	case 8:
		return uint16(c.a)
	case 9:
		return uint16(c.b)
	case 10:
		return uint16(c.cc)
	case 11:
		return uint16(c.dp)
	default:
		panic(fmt.Sprintf("Invalid register code: %d", code))
	}
}

func (c *CPU) setRegisterFromCode(code int, value uint16) {
	switch code {
	case 0:
		c.a = Word(value >> 8)
		c.b = Word(value)
	case 1:
		c.x = DWord(value)
	case 2:
		c.y = DWord(value)
	case 3:
		c.u = value
	case 4:
		c.s = value
	case 5:
		c.pc = value
	case 8:
		c.a = Word(value)
	case 9:
		c.b = Word(value)
	case 10:
		c.cc = uint8(value)
	case 11:
		c.dp = Word(value)
	default:
		panic(fmt.Sprintf("Invalid register code: %d", code))
	}
}

/** Exchange Registers */
func (c *CPU) exg(address uint16) {
	code := int(c.read(address))
	if ((code&0x80)>>7)^((code&0x08)>>3) == 1 {
		panic("Try to exchange 8-bit with 16-bits registers")
	}
	value1 := c.getRegisterFromCode(code >> 4)
	value2 := c.getRegisterFromCode(code & 0x0f)
	c.setRegisterFromCode(code>>4, value2)
	c.setRegisterFromCode(code&0x0f, value1)
}

/** Transfer Register to Register */
func (c *CPU) tfr(address uint16) {
	code := int(c.read(address))
	if ((code&0x80)>>7)^((code&0x08)>>3) == 1 {
		panic("Try to transfer 8-bit and 16-bits registers")
	}
	value := c.getRegisterFromCode(code >> 4)
	c.setRegisterFromCode(code&0x0f, value)
}

/** Branch Never */
func (c *CPU) brn(address uint16) {
	// NOP
}

/** Branch if Higher - Branch when Z = 0 && C = 0 */
func (c *CPU) bhi(address uint16) {
	if !c.getC() && !c.getZ() {
		c.pc = address
	}
}

/** Branch on Lower or Same - Branch when Z = 1 || C = 1 */
func (c *CPU) bls(address uint16) {
	if c.getC() || c.getZ() {
		c.pc = address
	}
}

/** Branch on Carry Clear - Branch when C = 0 */
func (c *CPU) bcc(address uint16) {
	if !c.getC() {
		c.pc = address
	}
}

/** Branch on Lower - Branch when C = 1 */
func (c *CPU) blo(address uint16) {
	if c.getC() {
		c.pc = address
	}
}

/** Branch on Not Equal - Branch when Z = 0 */
func (c *CPU) bne(address uint16) {
	if !c.getZ() {
		c.pc = address
	}
}

/** Branch on Equal - Branch when Z = 1 */
func (c *CPU) beq(address uint16) {
	if c.getZ() {
		c.pc = address
	}
}

/** Branch on Overflow Clear - Branch when V = 0 */
func (c *CPU) bvc(address uint16) {
	if !c.getV() {
		c.pc = address
	}
}

/** Branch on Overflow Set - Branch when V = 1 */
func (c *CPU) bvs(address uint16) {
	if c.getV() {
		c.pc = address
	}
}

/** Branch on Plus - Branch when N = 0 */
func (c *CPU) bpl(address uint16) {
	if !c.getN() {
		c.pc = address
	}
}

/** Branch on Minus - Branch when N = 1 */
func (c *CPU) bmi(address uint16) {
	if c.getN() {
		c.pc = address
	}
}

/** Branch on Greater than or Equal to Zero - Branch when N ^ V = 0 */
func (c *CPU) bge(address uint16) {
	if c.getN() == c.getV() {
		c.pc = address
	}
}

/** Branch on Less than Zero - Branch when N ^ V = 1 */
func (c *CPU) blt(address uint16) {
	if c.getN() != c.getV() {
		c.pc = address
	}
}

/** Branch on Greater - Branch when Z = 0 && (N ^ V) = 0 */
func (c *CPU) bgt(address uint16) {
	if !c.getZ() && c.getN() == c.getV() {
		c.pc = address
	}
}

/** Branch on Less than or Equal to Zero - Branch when Z = 1 || (N ^ V) = 1 */
func (c *CPU) ble(address uint16) {
	if c.getZ() || c.getN() != c.getV() {
		c.pc = address
	}
}

/** Load Effective Address into Register X */
func (c *CPU) leax(address uint16) {
	c.x = DWord(address)
	if address == 0 {
		c.setZ()
	} else {
		c.clearZ()
	}
}

/** Load Effective Address into Register Y */
func (c *CPU) leay(address uint16) {
	c.y = DWord(address)
	if address == 0 {
		c.setZ()
	} else {
		c.clearZ()
	}
}

/** Load Effective Address into Register S */
func (c *CPU) leas(address uint16) {
	c.s = address
}

/** Load Effective Address into Register U */
func (c *CPU) leau(address uint16) {
	c.u = address
}

func isBitSet(value uint8, flag uint) bool {
	return value&(1<<flag) != 0
}

func (c *CPU) pushRegister(value interface{}, stack *uint16) {
	switch t := value.(type) {
	case uint8:
		*stack--
		c.write(*stack, Word(t))
		c.clock++
	case Word:
		*stack--
		c.write(*stack, Word(t))
		c.clock++
	case uint16:
		*stack -= 2
		c.writew(*stack, DWord(t))
		c.clock += 2
	case DWord:
		*stack -= 2
		c.writew(*stack, DWord(t))
		c.clock += 2
	default:
		// WTF
	}
}

func (c *CPU) pullRegister(target interface{}, stack *uint16) {
	switch t := target.(type) {
	case *uint8:
		*t = uint8(c.read(*stack))
		*stack++
		c.clock++
	case *Word:
		*t = Word(c.read(*stack))
		*stack++
		c.clock++
	case *uint16:
		*t = uint16(c.readw(*stack))
		*stack += 2
		c.clock += 2
	case *DWord:
		*t = DWord(c.readw(*stack))
		*stack += 2
		c.clock += 2
	default:
		// WTF
	}
}

/** Push Registers on the Hardware Stack */
func (c *CPU) pshs(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 7) {
		c.pushRegister(c.pc, &c.s)
	}
	if isBitSet(registers, 6) {
		c.pushRegister(c.u, &c.s)
	}
	if isBitSet(registers, 5) {
		c.pushRegister(c.y, &c.s)
	}
	if isBitSet(registers, 4) {
		c.pushRegister(c.x, &c.s)
	}
	if isBitSet(registers, 3) {
		c.pushRegister(c.dp, &c.s)
	}
	if isBitSet(registers, 2) {
		c.pushRegister(c.b, &c.s)
	}
	if isBitSet(registers, 1) {
		c.pushRegister(c.a, &c.s)
	}
	if isBitSet(registers, 0) {
		c.pushRegister(c.cc, &c.s)
	}
}

/** Pull Registers from the Hardware Stack */
func (c *CPU) puls(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 0) {
		c.pullRegister(&c.cc, &c.s)
	}
	if isBitSet(registers, 1) {
		c.pullRegister(&c.a, &c.s)
	}
	if isBitSet(registers, 2) {
		c.pullRegister(&c.b, &c.s)
	}
	if isBitSet(registers, 3) {
		c.pullRegister(&c.dp, &c.s)
	}
	if isBitSet(registers, 4) {
		c.pullRegister(&c.x, &c.s)
	}
	if isBitSet(registers, 5) {
		c.pullRegister(&c.y, &c.s)
	}
	if isBitSet(registers, 6) {
		c.pullRegister(&c.u, &c.s)
	}
	if isBitSet(registers, 7) {
		c.pullRegister(&c.pc, &c.s)
	}
}

/** Push Registers on the User Stack */
func (c *CPU) pshu(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 7) {
		c.pushRegister(c.pc, &c.u)
	}
	if isBitSet(registers, 6) {
		c.pushRegister(c.s, &c.u)
	}
	if isBitSet(registers, 5) {
		c.pushRegister(c.y, &c.u)
	}
	if isBitSet(registers, 4) {
		c.pushRegister(c.x, &c.u)
	}
	if isBitSet(registers, 3) {
		c.pushRegister(c.dp, &c.u)
	}
	if isBitSet(registers, 2) {
		c.pushRegister(c.b, &c.u)
	}
	if isBitSet(registers, 1) {
		c.pushRegister(c.a, &c.u)
	}
	if isBitSet(registers, 0) {
		c.pushRegister(c.cc, &c.u)
	}
}

/** Pull Registers from the User Stack */
func (c *CPU) pulu(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 0) {
		c.pullRegister(&c.cc, &c.u)
	}
	if isBitSet(registers, 1) {
		c.pullRegister(&c.a, &c.u)
	}
	if isBitSet(registers, 2) {
		c.pullRegister(&c.b, &c.u)
	}
	if isBitSet(registers, 3) {
		c.pullRegister(&c.dp, &c.u)
	}
	if isBitSet(registers, 4) {
		c.pullRegister(&c.x, &c.u)
	}
	if isBitSet(registers, 5) {
		c.pullRegister(&c.y, &c.u)
	}
	if isBitSet(registers, 6) {
		c.pullRegister(&c.s, &c.u)
	}
	if isBitSet(registers, 7) {
		c.pullRegister(&c.pc, &c.u)
	}
}

/** Return from Subroutine */
func (c *CPU) rts() {
	c.pc = uint16(c.readw(c.s))
	c.s += 2
}

/** Add Accumulator B into Index Register X */
func (c *CPU) abx() {
	c.x += DWord(c.b)
}

/** Return from Interrupt */
func (c *CPU) rti() {
	c.pullRegister(&c.cc, &c.s)
	if c.getE() {
		c.pullRegister(&c.a, &c.s)
		c.pullRegister(&c.b, &c.s)
		c.pullRegister(&c.dp, &c.s)
		c.pullRegister(&c.x, &c.s)
		c.pullRegister(&c.y, &c.s)
		c.pullRegister(&c.u, &c.s)
	}
	c.pullRegister(&c.pc, &c.s)
}

/** Multiply - ZxCx */
func (c *CPU) mul() {
	value := uint16(c.a) * uint16(c.b)
	c.a = Word(value >> 8)
	c.b = Word(value & 0xff)
	if c.d() == 0 {
		c.setZ()
	} else {
		c.clearZ()
	}
	c.updateC(c.b&0x80 != 0)
}

/** Software Interrupt */
func (c *CPU) swi() {
	c.setE()
	c.pushRegister(c.pc, &c.s)
	c.pushRegister(c.u, &c.s)
	c.pushRegister(c.y, &c.s)
	c.pushRegister(c.x, &c.s)
	c.pushRegister(c.dp, &c.s)
	c.pushRegister(c.b, &c.s)
	c.pushRegister(c.a, &c.s)
	c.pushRegister(c.cc, &c.s)
	c.setF()
	c.setI()
	c.pc = uint16(c.readw(0xfffa))
}
