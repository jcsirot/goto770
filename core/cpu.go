package core

import "fmt"

type opcode struct {
	name   string
	f      func()
	cycles uint64
	mode   addressMode
}

var (
	opcodes [0x300]opcode
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

/* Addressing modes */

type addressMode int

const (
	inherent   addressMode = 0
	direct     addressMode = 1
	immediate  addressMode = 2
	limmediate addressMode = 3
	relative   addressMode = 4
	lrelative  addressMode = 5
	extended   addressMode = 6
	indexed    addressMode = 7
)

// Cpu structure
type CPU struct {
	/// Accumulator register
	a r8
	/// Accumulator register
	b r8
	/// Index register
	x r16
	/// Index register
	y r16
	/// User stack pointer register
	u r16
	/// Hardware stack pointer register
	s r16
	/// Direct page register
	dp r8
	/// Condition code regsiter
	cc ccr
	/// Program counter register
	pc r16
	/// Memory
	ram Memory
	///
	clock uint64
}

func (c *CPU) d() uint16 {
	return uint16(c.a.get()<<8 | c.b.get())
}

// Initialize the Cpu
func (c *CPU) Initialize(ram Memory) {
	c.ram = ram
	c.Reset()
	c.initOpcodes()
}

func (c *CPU) Reset() {
	c.a = r8{r: new(int)}
	c.b = r8{r: new(int)}
	c.x = r16{r: new(int)}
	c.y = r16{r: new(int)}
	c.u = r16{r: new(int)}
	c.s = r16{r: new(int)}
	c.dp = r8{r: new(int)}
	c.cc = ccr{r8{r: new(int)}}
	c.pc = r16{r: new(int)}
	c.clock = 0
}

func (c *CPU) initOpcodes() {
	opcodes[0x00] = opcode{"NEG", func() { c.neg(c.direct()) }, 6, direct}
	opcodes[0x03] = opcode{"COM", func() { c.com(c.direct()) }, 6, direct}
	opcodes[0x04] = opcode{"LSR", func() { c.lsr(c.direct()) }, 6, direct}
	opcodes[0x06] = opcode{"ROR", func() { c.ror(c.direct()) }, 6, direct}
	opcodes[0x07] = opcode{"ASR", func() { c.asr(c.direct()) }, 6, direct}
	opcodes[0x08] = opcode{"ASL", func() { c.asl(c.direct()) }, 6, direct}
	opcodes[0x09] = opcode{"ROL", func() { c.rol(c.direct()) }, 6, direct}
	opcodes[0x0a] = opcode{"DEC", func() { c.dec(c.direct()) }, 6, direct}
	opcodes[0x0c] = opcode{"INC", func() { c.inc(c.direct()) }, 6, direct}
	opcodes[0x0d] = opcode{"TST", func() { c.tst(c.direct()) }, 6, direct}
	opcodes[0x0e] = opcode{"JMP", func() { c.jmp(c.direct()) }, 3, direct}
	opcodes[0x0f] = opcode{"CLR", func() { c.clr(c.direct()) }, 6, direct}
	opcodes[0x12] = opcode{"NOP", func() { c.nop() }, 2, inherent}
	opcodes[0x13] = opcode{"SYNC", func() { c.sync() }, 4, inherent}
	opcodes[0x16] = opcode{"BRA", func() { c.bra(c.lrelative()) }, 5, lrelative}
	opcodes[0x17] = opcode{"BSR", func() { c.bsr(c.lrelative()) }, 9, lrelative}
	opcodes[0x19] = opcode{"DAA", func() { c.daa() }, 2, inherent}
	opcodes[0x1a] = opcode{"ORCC", func() { c.orcc(c.immediate()) }, 3, immediate}
	opcodes[0x1c] = opcode{"ANDCC", func() { c.andcc(c.immediate()) }, 3, immediate}
	opcodes[0x1d] = opcode{"SEX", func() { c.sex() }, 2, inherent}
	opcodes[0x1e] = opcode{"EXG", func() { c.exg(c.immediate()) }, 8, immediate}
	opcodes[0x1f] = opcode{"TFR", func() { c.tfr(c.immediate()) }, 6, immediate}
	opcodes[0x20] = opcode{"BRA", func() { c.bra(c.relative()) }, 3, relative}
	opcodes[0x21] = opcode{"BRN", func() { c.brn(c.relative()) }, 3, relative}
	opcodes[0x22] = opcode{"BHI", func() { c.bhi(c.relative()) }, 3, relative}
	opcodes[0x23] = opcode{"BLS", func() { c.bls(c.relative()) }, 3, relative}
	opcodes[0x24] = opcode{"BCC", func() { c.bcc(c.relative()) }, 3, relative}
	opcodes[0x25] = opcode{"BLO", func() { c.blo(c.relative()) }, 3, relative}
	opcodes[0x26] = opcode{"BNE", func() { c.bne(c.relative()) }, 3, relative}
	opcodes[0x27] = opcode{"BEQ", func() { c.beq(c.relative()) }, 3, relative}
	opcodes[0x28] = opcode{"BVC", func() { c.bvc(c.relative()) }, 3, relative}
	opcodes[0x29] = opcode{"BVS", func() { c.bvs(c.relative()) }, 3, relative}
	opcodes[0x2a] = opcode{"BPL", func() { c.bpl(c.relative()) }, 3, relative}
	opcodes[0x2b] = opcode{"BMI", func() { c.bmi(c.relative()) }, 3, relative}
	opcodes[0x2c] = opcode{"BGE", func() { c.bge(c.relative()) }, 3, relative}
	opcodes[0x2d] = opcode{"BLT", func() { c.blt(c.relative()) }, 3, relative}
	opcodes[0x2e] = opcode{"BGT", func() { c.bgt(c.relative()) }, 3, relative}
	opcodes[0x2f] = opcode{"BLE", func() { c.ble(c.relative()) }, 3, relative}
	opcodes[0x30] = opcode{"LEAX", func() { c.leax(c.indexed()) }, 4, indexed}
	opcodes[0x31] = opcode{"LEAY", func() { c.leay(c.indexed()) }, 4, indexed}
	opcodes[0x32] = opcode{"LEAS", func() { c.leas(c.indexed()) }, 4, indexed}
	opcodes[0x33] = opcode{"LEAU", func() { c.leau(c.indexed()) }, 4, indexed}
	opcodes[0x34] = opcode{"PSHS", func() { c.pshs(c.immediate()) }, 5, immediate}
	opcodes[0x35] = opcode{"PULS", func() { c.puls(c.immediate()) }, 5, immediate}
	opcodes[0x36] = opcode{"PSHU", func() { c.pshu(c.immediate()) }, 5, immediate}
	opcodes[0x37] = opcode{"PULU", func() { c.pulu(c.immediate()) }, 5, immediate}
	opcodes[0x39] = opcode{"RTS", func() { c.rts() }, 5, inherent}
	opcodes[0x3a] = opcode{"ABX", func() { c.abx() }, 3, inherent}
	opcodes[0x3b] = opcode{"RTI", func() { c.rti() }, 3, inherent}
	opcodes[0x3d] = opcode{"MUL", func() { c.mul() }, 11, inherent}
	opcodes[0x3f] = opcode{"SWI", func() { c.swi() }, 7, inherent}
	opcodes[0x40] = opcode{"NEGA", func() { c.nega() }, 2, inherent}
	opcodes[0x43] = opcode{"COMA", func() { c.coma() }, 2, inherent}
	opcodes[0x44] = opcode{"LSRA", func() { c.lsra() }, 2, inherent}
	opcodes[0x46] = opcode{"RORA", func() { c.rora() }, 2, inherent}
	opcodes[0x47] = opcode{"ASRA", func() { c.asra() }, 2, inherent}
	opcodes[0x48] = opcode{"ASLA", func() { c.asla() }, 2, inherent}
	opcodes[0x49] = opcode{"ROLA", func() { c.rola() }, 2, inherent}
	opcodes[0x4a] = opcode{"DECA", func() { c.deca() }, 2, inherent}
	opcodes[0x4c] = opcode{"INCA", func() { c.inca() }, 2, inherent}
	opcodes[0x4d] = opcode{"TSTA", func() { c.tsta() }, 2, inherent}
	opcodes[0x4f] = opcode{"CLRA", func() { c.clra() }, 2, inherent}
	opcodes[0x50] = opcode{"NEGB", func() { c.negb() }, 2, inherent}
	opcodes[0x53] = opcode{"COMB", func() { c.comb() }, 2, inherent}
	opcodes[0x54] = opcode{"LSRB", func() { c.lsrb() }, 2, inherent}
	opcodes[0x56] = opcode{"RORB", func() { c.rorb() }, 2, inherent}
	opcodes[0x57] = opcode{"ASRB", func() { c.asrb() }, 2, inherent}
	opcodes[0x58] = opcode{"ASLB", func() { c.aslb() }, 2, inherent}
	opcodes[0x59] = opcode{"ROLB", func() { c.rolb() }, 2, inherent}
	opcodes[0x5a] = opcode{"DECB", func() { c.decb() }, 2, inherent}
	opcodes[0x5c] = opcode{"INCB", func() { c.incb() }, 2, inherent}
	opcodes[0x5d] = opcode{"TSTB", func() { c.tstb() }, 2, inherent}
	opcodes[0x5f] = opcode{"CLRB", func() { c.clrb() }, 2, inherent}
	opcodes[0x60] = opcode{"NEG", func() { c.neg(c.indexed()) }, 6, indexed}
	opcodes[0x63] = opcode{"COM", func() { c.com(c.indexed()) }, 6, indexed}
	opcodes[0x64] = opcode{"LSR", func() { c.lsr(c.indexed()) }, 6, indexed}
	opcodes[0x66] = opcode{"ROR", func() { c.ror(c.indexed()) }, 6, indexed}
	opcodes[0x67] = opcode{"ASR", func() { c.asr(c.indexed()) }, 6, indexed}
	opcodes[0x68] = opcode{"ASL", func() { c.asl(c.indexed()) }, 6, indexed}
	opcodes[0x69] = opcode{"ROL", func() { c.rol(c.indexed()) }, 6, indexed}
	opcodes[0x6a] = opcode{"DEC", func() { c.dec(c.indexed()) }, 6, indexed}
	opcodes[0x6c] = opcode{"INC", func() { c.inc(c.indexed()) }, 6, indexed}
	opcodes[0x6d] = opcode{"TST", func() { c.tst(c.indexed()) }, 6, indexed}
	opcodes[0x6e] = opcode{"JMP", func() { c.jmp(c.indexed()) }, 3, indexed}
	opcodes[0x6f] = opcode{"CLR", func() { c.clr(c.indexed()) }, 6, indexed}
	opcodes[0x70] = opcode{"NEG", func() { c.neg(c.extended()) }, 7, extended}
	opcodes[0x73] = opcode{"COM", func() { c.com(c.extended()) }, 7, extended}
	opcodes[0x74] = opcode{"LSR", func() { c.lsr(c.extended()) }, 7, extended}
	opcodes[0x76] = opcode{"ROR", func() { c.ror(c.extended()) }, 7, extended}
	opcodes[0x77] = opcode{"ASR", func() { c.asr(c.extended()) }, 7, extended}
	opcodes[0x78] = opcode{"ASL", func() { c.asl(c.extended()) }, 7, extended}
	opcodes[0x79] = opcode{"ROL", func() { c.rol(c.extended()) }, 7, extended}
	opcodes[0x7a] = opcode{"DEC", func() { c.dec(c.extended()) }, 7, extended}
	opcodes[0x7c] = opcode{"INC", func() { c.inc(c.extended()) }, 7, extended}
	opcodes[0x7d] = opcode{"TST", func() { c.tst(c.extended()) }, 7, extended}
	opcodes[0x7e] = opcode{"JMP", func() { c.jmp(c.extended()) }, 4, extended}
	opcodes[0x7f] = opcode{"CLR", func() { c.clr(c.extended()) }, 7, extended}
	opcodes[0x80] = opcode{"SUBA", func() { c.suba(c.immediate()) }, 2, immediate}
	opcodes[0x81] = opcode{"CMPA", func() { c.cmpa(c.immediate()) }, 2, immediate}
	opcodes[0x82] = opcode{"SBCA", func() { c.sbca(c.immediate()) }, 2, immediate}
	opcodes[0x83] = opcode{"SUBD", func() { c.subd(c.limmediate()) }, 4, limmediate}
	opcodes[0x84] = opcode{"ANDA", func() { c.anda(c.immediate()) }, 2, immediate}
	opcodes[0x85] = opcode{"BITA", func() { c.bita(c.immediate()) }, 2, immediate}
	opcodes[0x86] = opcode{"LDA", func() { c.lda(c.immediate()) }, 2, immediate}
	opcodes[0x88] = opcode{"EORA", func() { c.eora(c.immediate()) }, 2, immediate}
	opcodes[0x89] = opcode{"ADCA", func() { c.adca(c.immediate()) }, 2, immediate}
	opcodes[0x8a] = opcode{"ORA", func() { c.ora(c.immediate()) }, 2, immediate}
	opcodes[0x8b] = opcode{"ADDA", func() { c.adda(c.immediate()) }, 2, immediate}
	opcodes[0x8c] = opcode{"CMPX", func() { c.cmpx(c.limmediate()) }, 4, limmediate}
	opcodes[0x8d] = opcode{"BSR", func() { c.bsr(c.relative()) }, 7, relative}
	opcodes[0x8e] = opcode{"LDX", func() { c.ldx(c.limmediate()) }, 3, limmediate}
	opcodes[0x90] = opcode{"SUBA", func() { c.suba(c.direct()) }, 4, direct}
	opcodes[0x91] = opcode{"CMPA", func() { c.cmpa(c.direct()) }, 4, direct}
	opcodes[0x92] = opcode{"SBCA", func() { c.sbca(c.direct()) }, 4, direct}
	opcodes[0x93] = opcode{"SUBD", func() { c.subd(c.direct()) }, 6, direct}
	opcodes[0x94] = opcode{"ANDA", func() { c.anda(c.direct()) }, 4, direct}
	opcodes[0x95] = opcode{"BITA", func() { c.bita(c.direct()) }, 4, direct}
	opcodes[0x96] = opcode{"LDA", func() { c.lda(c.direct()) }, 4, direct}
	opcodes[0x97] = opcode{"STA", func() { c.sta(c.direct()) }, 4, direct}
	opcodes[0x98] = opcode{"EORA", func() { c.eora(c.direct()) }, 4, direct}
	opcodes[0x99] = opcode{"ADCA", func() { c.adca(c.direct()) }, 4, direct}
	opcodes[0x9a] = opcode{"ORA", func() { c.ora(c.direct()) }, 4, direct}
	opcodes[0x9b] = opcode{"ADDA", func() { c.adda(c.direct()) }, 4, direct}
	opcodes[0x9c] = opcode{"CMPX", func() { c.cmpx(c.direct()) }, 7, direct}
	opcodes[0x9d] = opcode{"JSR", func() { c.jsr(c.direct()) }, 7, direct}
	opcodes[0x9e] = opcode{"LDX", func() { c.ldx(c.direct()) }, 5, direct}
	opcodes[0x9f] = opcode{"STX", func() { c.stx(c.direct()) }, 5, direct}
	opcodes[0xa0] = opcode{"SUBA", func() { c.suba(c.indexed()) }, 4, indexed}
	opcodes[0xa1] = opcode{"CMPA", func() { c.cmpa(c.indexed()) }, 4, indexed}
	opcodes[0xa2] = opcode{"SBCA", func() { c.sbca(c.indexed()) }, 4, indexed}
	opcodes[0xa3] = opcode{"SUBD", func() { c.subd(c.indexed()) }, 6, indexed}
	opcodes[0xa4] = opcode{"ANDA", func() { c.anda(c.indexed()) }, 4, indexed}
	opcodes[0xa5] = opcode{"BITA", func() { c.bita(c.indexed()) }, 4, indexed}
	opcodes[0xa6] = opcode{"LDA", func() { c.lda(c.indexed()) }, 4, indexed}
	opcodes[0xa7] = opcode{"STA", func() { c.sta(c.indexed()) }, 4, indexed}
	opcodes[0xa8] = opcode{"EORA", func() { c.eora(c.indexed()) }, 4, indexed}
	opcodes[0xa9] = opcode{"ADCA", func() { c.adca(c.indexed()) }, 4, indexed}
	opcodes[0xaa] = opcode{"ORA", func() { c.ora(c.indexed()) }, 4, indexed}
	opcodes[0xab] = opcode{"ADDA", func() { c.adda(c.indexed()) }, 4, indexed}
	opcodes[0xac] = opcode{"CMPX", func() { c.cmpx(c.indexed()) }, 7, indexed}
	opcodes[0xad] = opcode{"JSR", func() { c.jsr(c.indexed()) }, 7, indexed}
	opcodes[0xae] = opcode{"LDX", func() { c.ldx(c.indexed()) }, 5, indexed}
	opcodes[0xaf] = opcode{"STX", func() { c.stx(c.indexed()) }, 5, indexed}

}

func (c *CPU) step() uint64 {
	opcode := opcodes[c.read(c.pc.uint16())]

	instBuf := make([]uint8, 5)
	instBuf[0] = c.read(c.pc.uint16())
	instBuf[1] = c.read(c.pc.uint16() + 1)
	instBuf[2] = c.read(c.pc.uint16() + 2)
	instBuf[3] = c.read(c.pc.uint16() + 3)
	instBuf[4] = c.read(c.pc.uint16() + 4)

	// instr, len := Disassemble(opcode, instBuf)
	// format(c.pc.uint16(), instr, instBuf[0:len])

	/*
		if c.Verbose {
			Disassemble(opcode, instBuf)
		}
	*/

	c.pc.inc()

	opcode.f()
	c.clock += opcode.cycles

	return opcode.cycles

}

/***************************/
/**     Memory access     **/
/***************************/

func (c *CPU) read(address uint16) uint8 {
	return c.ram.Read(address)
}

func (c *CPU) readInt(address uint16) int {
	return int(c.read(address))
}

func (c *CPU) readw(address uint16) uint16 {
	return c.ram.Readw(address)
}

func (c *CPU) readwInt(address uint16) int {
	return int(c.ram.Readw(address))
}

func (c *CPU) write(address uint16, value uint8) {
	c.ram.Write(address, value)
}

func (c *CPU) writeInt(address uint16, value int) {
	c.ram.Write(address, uint8(value))
}

func (c *CPU) writew(address uint16, value uint16) {
	c.ram.Writew(address, value)
}

func (c *CPU) writewInt(address uint16, value int) {
	c.ram.Writew(address, uint16(value))
}

/** Negate - H?NxZxVxCx */
func (c *CPU) neg_(value int) int {
	tmp := -value
	c.updateNZVC(0, value, tmp)
	return tmp
}

/** Negate - H?NxZxVxCx */
func (c *CPU) neg(address uint16) {
	c.writeInt(address, c.neg_(c.readInt(address)))
}

/** Negate Register A - H?NxZxVxCx */
func (c *CPU) nega() {
	c.a.set(c.neg_(c.a.get()))
}

/** Negate Register B - H?NxZxVxCx */
func (c *CPU) negb() {
	c.b.set(c.neg_(c.b.get()))
}

/** Complement - H?NxZxV0C1 */
func (c *CPU) com_(value int) int {
	tmp := value ^ 0xff
	c.updateNZ(tmp)
	c.cc.clearV()
	c.cc.setC()
	return tmp
}

/** Complement - H?NxZxV0C1 */
func (c *CPU) com(address uint16) {
	c.writeInt(address, c.com_(c.readInt(address)))
}

/** Complement Register A - H?NxZxV0C1 */
func (c *CPU) coma() {
	c.a.set(c.com_(c.a.get()))
}

/** Complement Register B - H?NxZxV0C1 */
func (c *CPU) comb() {
	c.b.set(c.com_(c.b.get()))
}

/** Logical Shift Right - N0ZxCx */
func (c *CPU) lsr_(value int) int {
	tmp := value >> 1
	c.updateNZ(tmp)
	c.updateC(value&1 == 1)
	//c.testSetZN(uint8(tmp))
	//c.updateC(value&1 == 1)
	return tmp
}

/** Logical Shift Right - N0ZxCx */
func (c *CPU) lsr(address uint16) {
	c.writeInt(address, c.lsr_(c.readInt(address)))
}

/** Logical Shift Right A Register - N0ZxCx */
func (c *CPU) lsra() {
	c.a.set(c.lsr_(c.a.get()))
}

/** Logical Shift Right B Register - N0ZxCx */
func (c *CPU) lsrb() {
	c.b.set(c.lsr_(c.b.get()))
}

/** Rotate Right - NxZxCx */
func (c *CPU) ror_(value int) int {
	carry := 0
	if c.cc.getC() {
		carry = 0x80
	}
	tmp := (value >> 1) | carry
	c.updateNZ(tmp)
	c.updateC(value&1 == 1)
	//c.testSetZN(uint8(tmp))
	//c.updateC(value&1 == 1)
	return tmp
}

/** Rotate Right - NxZxCx */
func (c *CPU) ror(address uint16) {
	c.writeInt(address, c.ror_(c.readInt(address)))
}

/** Rotate Right Register A - NxZxCx */
func (c *CPU) rora() {
	c.a.set(c.ror_(c.a.get()))

}

/** Rotate Right Register B - NxZxCx */
func (c *CPU) rorb() {
	c.b.set(c.ror_(c.b.get()))
}

/** Rotate Left - NxZxVxCx */
func (c *CPU) rol_(value int) int {
	carry := 0
	if c.cc.getC() {
		carry = 1
	}
	tmp := (value << 1) | carry
	c.updateNZVC(value, value, tmp)
	//c.testSetZN(uint8(tmp))
	//c.updateC(uint8(value)>>7 == 0x01)
	//c.updateV((value>>7)^((value>>6)&0x01) == 0x01)
	return tmp
}

/** Rotate Left - NxZxVxCx */
func (c *CPU) rol(address uint16) {
	c.writeInt(address, c.rol_(c.readInt(address)))
}

/** Rotate Left Register A - NxZxVxCx */
func (c *CPU) rola() {
	c.a.set(c.rol_(c.a.get()))
}

/** Rotate Left Register B - NxZxVxCx */
func (c *CPU) rolb() {
	c.b.set(c.rol_(c.b.get()))
}

/** Arithmetic Shift Right - H?NxZxCx */
func (c *CPU) asr_(value int) int {
	tmp := (value >> 1) | (value & 0x80)
	c.updateNZ(tmp)
	//c.testSetZN(uint8(tmp))
	c.updateC(value&0x01 == 0x01)
	return tmp
}

/** Arithmetic Shift Right - H?NxZxCx */
func (c *CPU) asr(address uint16) {
	c.writeInt(address, c.asr_(c.readInt(address)))
}

/** Arithmetic Shift Right Register A - H?NxZxCx */
func (c *CPU) asra() {
	c.a.set(c.asr_(c.a.get()))
}

/** Arithmetic Shift Right Register B - H?NxZxCx */
func (c *CPU) asrb() {
	c.b.set(c.asr_(c.b.get()))
}

/** Arithmetic Shift Left / Logical Shift Left - H?NxZxVxCx */
func (c *CPU) asl_(value int) int {
	tmp := value << 1
	c.updateNZVC(value, value, tmp)
	//c.testSetZN(uint8(tmp))
	//c.updateC(value&0x80 == 0x80)
	//c.updateV((value>>7)^((value>>6)&0x01) == 0x01)
	return tmp
}

/** Arithmetic Shift Left / Logical Shift Left - H?NxZxVxCx */
func (c *CPU) asl(address uint16) {
	c.writeInt(address, c.asl_(c.readInt(address)))
}

/** Arithmetic Shift Left / Logical Shift Left Register A - H?NxZxVxCx */
func (c *CPU) asla() {
	c.a.set(c.asl_(c.a.get()))
}

/** Arithmetic Shift Left / Logical Shift Left Register B - H?NxZxVxCx */
func (c *CPU) aslb() {
	c.b.set(c.asl_(c.b.get()))
}

/** Decrement - NxZxVx */
func (c *CPU) dec_(value int) int {
	tmp := value - 1
	c.updateNZ(tmp)
	c.updateV(value == 0x80)
	return tmp
}

/** Decrement - NxZxVx */
func (c *CPU) dec(address uint16) {
	c.writeInt(address, c.dec_(c.readInt(address)))
}

/** Decrement Register A - NxZxVx */
func (c *CPU) deca() {
	c.a.set(c.dec_(c.a.get()))
}

/** Decrement Register B - NxZxVx */
func (c *CPU) decb() {
	c.b.set(c.dec_(c.b.get()))
}

/** Increment - NxZxVx */
func (c *CPU) inc_(value int) int {
	tmp := value + 1
	c.updateNZ(tmp)
	c.updateV(value == 0x7f)
	return tmp
}

/** Increment - NxZxVx */
func (c *CPU) inc(address uint16) {
	c.writeInt(address, c.inc_(c.readInt(address)))
}

/** Increment Register A - NxZxVx */
func (c *CPU) inca() {
	c.a.set(c.inc_(c.a.get()))
}

/** Increment Register B - NxZxVx */
func (c *CPU) incb() {
	c.b.set(c.inc_(c.b.get()))
}

/** Test - NxZxV0 */
func (c *CPU) tst_(value int) {
	c.updateNZ(value)
	c.cc.clearV()
}

/** Test - NxZxV0 */
func (c *CPU) tst(address uint16) {
	c.tst_(c.readInt(address))
}

/** Test Register A - NxZxV0 */
func (c *CPU) tsta() {
	c.tst_(c.a.get())
}

/** Test Register B - NxZxV0 */
func (c *CPU) tstb() {
	c.tst_(c.b.get())
}

/** Jump - NxZxV0 */
func (c *CPU) jmp(address uint16) {
	c.pc.set(address)
}

/** Clear N0Z1V0C0 */
func (c *CPU) clr(address uint16) {
	c.write(address, 0)
	c.cc.clearN()
	c.cc.setZ()
	c.cc.clearV()
	c.cc.clearC()
}

/** Clear N0Z1V0C0 */
func (c *CPU) clra() {
	c.a.set(0)
	c.cc.clearN()
	c.cc.setZ()
	c.cc.clearV()
	c.cc.clearC()
}

/** Clear N0Z1V0C0 */
func (c *CPU) clrb() {
	c.b.set(0)
	c.cc.clearN()
	c.cc.setZ()
	c.cc.clearV()
	c.cc.clearC()
}

func (c *CPU) nop() {
}

/** Synchronize to External Event */
func (c *CPU) sync() {
	// Not supported
}

/** (Long) Branch Always */
func (c *CPU) bra(address uint16) {
	c.pc.set(address)
}

/** Long Branch / Jump to Subroutine */
func (c *CPU) bsr(address uint16) {
	c.s.set(c.s.get() - 2)
	c.writew(c.s.uint16(), c.pc.uint16())
	c.pc.set(address)
}

/** Jump to Subroutine */
func (c *CPU) jsr(address uint16) {
	c.s.set(c.s.get() - 2)
	c.writew(c.s.uint16(), c.pc.uint16())
	c.pc.set(address)
}

/** Decimal Addition Adjust - NxZxV?Cx */
func (c *CPU) daa() {
	ah := c.a.uint8() & 0xf0
	al := c.a.uint8() & 0x0f
	cf := 0
	if al > 0x09 || c.cc.getH() {
		cf |= 0x06
	}
	if ah > 0x80 && al > 0x09 {
		cf |= 0x60
	}
	if ah > 0x90 || c.cc.getC() {
		cf |= 0x60
	}
	tmp := uint16(c.a.get()) + uint16(cf)
	c.a.set(int(tmp))
	carry := c.cc.getC()
	c.updateNZ(c.a.get())
	c.updateC(carry || tmp > 0xff)
}

/** Inclusive OR Memory Immediate into Condition Code Register */
func (c *CPU) orcc(address uint16) {
	value := c.readInt(address)
	c.cc.set(c.cc.get() | value)
}

/** Logical AND Immediate Memory into Condition Code Register */
func (c *CPU) andcc(address uint16) {
	value := c.readInt(address)
	c.cc.set(c.cc.get() & value)
}

/** Sign Extended - NxZx */
func (c *CPU) sex() {
	if c.b.uint8()&0x80 == 0 {
		c.a.set(0)
	} else {
		c.a.set(0xff)
	}
	if c.d() == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
	if (c.d() & 0x8000) != 0 {
		c.cc.setN()
	} else {
		c.cc.clearN()
	}
}

func (c *CPU) getRegisterFromCode(code int) uint16 {
	switch code {
	case 0:
		return uint16(c.d())
	case 1:
		return c.x.uint16()
	case 2:
		return c.y.uint16()
	case 3:
		return c.u.uint16()
	case 4:
		return c.s.uint16()
	case 5:
		return c.pc.uint16()
	case 8:
		return c.a.uint16()
	case 9:
		return c.b.uint16()
	case 10:
		return c.cc.uint16()
	case 11:
		return c.dp.uint16()
	default:
		panic(fmt.Sprintf("Invalid register code: %d", code))
	}
}

func (c *CPU) setRegisterFromCode(code int, value uint16) {
	switch code {
	case 0:
		c.a.set(int(value >> 8))
		c.b.set(int(value))
	case 1:
		c.x.set(int(value))
	case 2:
		c.y.set(int(value))
	case 3:
		c.u.set(int(value))
	case 4:
		c.s.set(int(value))
	case 5:
		c.pc.set(value)
	case 8:
		c.a.set(int(value))
	case 9:
		c.b.set(int(value))
	case 10:
		c.cc.set(value)
	case 11:
		c.dp.set(value)
	default:
		panic(fmt.Sprintf("Invalid register code: %d", code))
	}
}

/** Exchange Registers */
func (c *CPU) exg(address uint16) {
	code := c.readInt(address)
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
	code := c.readInt(address)
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
	if !c.cc.getC() && !c.cc.getZ() {
		c.pc.set(address)
	}
}

/** Branch on Lower or Same - Branch when Z = 1 || C = 1 */
func (c *CPU) bls(address uint16) {
	if c.cc.getC() || c.cc.getZ() {
		c.pc.set(address)
	}
}

/** Branch on Carry Clear - Branch when C = 0 */
func (c *CPU) bcc(address uint16) {
	if !c.cc.getC() {
		c.pc.set(address)
	}
}

/** Branch on Lower - Branch when C = 1 */
func (c *CPU) blo(address uint16) {
	if c.cc.getC() {
		c.pc.set(address)
	}
}

/** Branch on Not Equal - Branch when Z = 0 */
func (c *CPU) bne(address uint16) {
	if !c.cc.getZ() {
		c.pc.set(address)
	}
}

/** Branch on Equal - Branch when Z = 1 */
func (c *CPU) beq(address uint16) {
	if c.cc.getZ() {
		c.pc.set(address)
	}
}

/** Branch on Overflow Clear - Branch when V = 0 */
func (c *CPU) bvc(address uint16) {
	if !c.cc.getV() {
		c.pc.set(address)
	}
}

/** Branch on Overflow Set - Branch when V = 1 */
func (c *CPU) bvs(address uint16) {
	if c.cc.getV() {
		c.pc.set(address)
	}
}

/** Branch on Plus - Branch when N = 0 */
func (c *CPU) bpl(address uint16) {
	if !c.cc.getN() {
		c.pc.set(address)
	}
}

/** Branch on Minus - Branch when N = 1 */
func (c *CPU) bmi(address uint16) {
	if c.cc.getN() {
		c.pc.set(address)
	}
}

/** Branch on Greater than or Equal to Zero - Branch when N ^ V = 0 */
func (c *CPU) bge(address uint16) {
	if c.cc.getN() == c.cc.getV() {
		c.pc.set(address)
	}
}

/** Branch on Less than Zero - Branch when N ^ V = 1 */
func (c *CPU) blt(address uint16) {
	if c.cc.getN() != c.cc.getV() {
		c.pc.set(address)
	}
}

/** Branch on Greater - Branch when Z = 0 && (N ^ V) = 0 */
func (c *CPU) bgt(address uint16) {
	if !c.cc.getZ() && c.cc.getN() == c.cc.getV() {
		c.pc.set(address)
	}
}

/** Branch on Less than or Equal to Zero - Branch when Z = 1 || (N ^ V) = 1 */
func (c *CPU) ble(address uint16) {
	if c.cc.getZ() || c.cc.getN() != c.cc.getV() {
		c.pc.set(address)
	}
}

/** Load Effective Address into Register X */
func (c *CPU) leax(address uint16) {
	c.x.set(int(address))
	if address == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
}

/** Load Effective Address into Register Y */
func (c *CPU) leay(address uint16) {
	c.y.set(int(address))
	if address == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
}

/** Load Effective Address into Register S */
func (c *CPU) leas(address uint16) {
	c.s.set(address)
}

/** Load Effective Address into Register U */
func (c *CPU) leau(address uint16) {
	c.u.set(address)
}

func isBitSet(value uint8, flag uint) bool {
	return value&(1<<flag) != 0
}

func (c *CPU) pushRegister(value register, stack r16) {
	sz := value.size()
	if sz == 8 {
		stack.dec()
		c.writeInt(stack.uint16(), value.get())
		c.clock++
	} else if sz == 16 {
		stack.dec().dec()
		c.writewInt(stack.uint16(), value.get())
		c.clock += 2
	} else {
		// WTF
	}
}

func (c *CPU) pullRegister(target register, stack r16) {
	sz := target.size()
	if sz == 8 {
		target.set(uint8(c.read(stack.uint16())))
		stack.inc()
		c.clock++
	} else if sz == 16 {
		target.set(uint16(c.readw(stack.uint16())))
		stack.inc().inc()
		c.clock += 2
	} else {
		// WTF
	}
}

/** Push Registers on the Hardware Stack */
func (c *CPU) pshs(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 7) {
		c.pushRegister(c.pc, c.s)
	}
	if isBitSet(registers, 6) {
		c.pushRegister(c.u, c.s)
	}
	if isBitSet(registers, 5) {
		c.pushRegister(c.y, c.s)
	}
	if isBitSet(registers, 4) {
		c.pushRegister(c.x, c.s)
	}
	if isBitSet(registers, 3) {
		c.pushRegister(c.dp, c.s)
	}
	if isBitSet(registers, 2) {
		c.pushRegister(c.b, c.s)
	}
	if isBitSet(registers, 1) {
		c.pushRegister(c.a, c.s)
	}
	if isBitSet(registers, 0) {
		c.pushRegister(c.cc, c.s)
	}
}

/** Pull Registers from the Hardware Stack */
func (c *CPU) puls(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 0) {
		c.pullRegister(c.cc, c.s)
	}
	if isBitSet(registers, 1) {
		c.pullRegister(c.a, c.s)
	}
	if isBitSet(registers, 2) {
		c.pullRegister(c.b, c.s)
	}
	if isBitSet(registers, 3) {
		c.pullRegister(c.dp, c.s)
	}
	if isBitSet(registers, 4) {
		c.pullRegister(c.x, c.s)
	}
	if isBitSet(registers, 5) {
		c.pullRegister(c.y, c.s)
	}
	if isBitSet(registers, 6) {
		c.pullRegister(c.u, c.s)
	}
	if isBitSet(registers, 7) {
		c.pullRegister(c.pc, c.s)
	}
}

/** Push Registers on the User Stack */
func (c *CPU) pshu(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 7) {
		c.pushRegister(c.pc, c.u)
	}
	if isBitSet(registers, 6) {
		c.pushRegister(c.s, c.u)
	}
	if isBitSet(registers, 5) {
		c.pushRegister(c.y, c.u)
	}
	if isBitSet(registers, 4) {
		c.pushRegister(c.x, c.u)
	}
	if isBitSet(registers, 3) {
		c.pushRegister(c.dp, c.u)
	}
	if isBitSet(registers, 2) {
		c.pushRegister(c.b, c.u)
	}
	if isBitSet(registers, 1) {
		c.pushRegister(c.a, c.u)
	}
	if isBitSet(registers, 0) {
		c.pushRegister(c.cc, c.u)
	}
}

/** Pull Registers from the User Stack */
func (c *CPU) pulu(address uint16) {
	registers := uint8(c.read(address))
	if isBitSet(registers, 0) {
		c.pullRegister(c.cc, c.u)
	}
	if isBitSet(registers, 1) {
		c.pullRegister(c.a, c.u)
	}
	if isBitSet(registers, 2) {
		c.pullRegister(c.b, c.u)
	}
	if isBitSet(registers, 3) {
		c.pullRegister(c.dp, c.u)
	}
	if isBitSet(registers, 4) {
		c.pullRegister(c.x, c.u)
	}
	if isBitSet(registers, 5) {
		c.pullRegister(c.y, c.u)
	}
	if isBitSet(registers, 6) {
		c.pullRegister(c.s, c.u)
	}
	if isBitSet(registers, 7) {
		c.pullRegister(c.pc, c.u)
	}
}

/** Return from Subroutine */
func (c *CPU) rts() {
	c.pc.set(c.readw(c.s.uint16()))
	c.s.inc().inc()
}

/** Add Accumulator B into Index Register X */
func (c *CPU) abx() {
	c.x.set(c.x.get() + c.b.get())
}

/** Return from Interrupt */
func (c *CPU) rti() {
	c.pullRegister(c.cc, c.s)
	if c.cc.getE() {
		c.pullRegister(c.a, c.s)
		c.pullRegister(c.b, c.s)
		c.pullRegister(c.dp, c.s)
		c.pullRegister(c.x, c.s)
		c.pullRegister(c.y, c.s)
		c.pullRegister(c.u, c.s)
	}
	c.pullRegister(c.pc, c.s)
}

/** Multiply - ZxCx */
func (c *CPU) mul() {
	value := c.a.get() * c.b.get()
	c.a.set(value >> 8)
	c.b.set(value & 0xff)
	if c.d() == 0 {
		c.cc.setZ()
	} else {
		c.cc.clearZ()
	}
	c.updateC(c.b.get()&0x80 != 0)
}

/** Software Interrupt */
func (c *CPU) swi() {
	c.cc.setE()
	c.pushRegister(c.pc, c.s)
	c.pushRegister(c.u, c.s)
	c.pushRegister(c.y, c.s)
	c.pushRegister(c.x, c.s)
	c.pushRegister(c.dp, c.s)
	c.pushRegister(c.b, c.s)
	c.pushRegister(c.a, c.s)
	c.pushRegister(c.cc, c.s)
	c.cc.setF()
	c.cc.setI()
	c.pc.set(c.readw(0xfffa))
}

/** Subtract Memory - H?NxZxVxCx */
func (c *CPU) sub_(reg int, value int) int {
	tmp := reg - value
	c.updateNZVC(reg, value, tmp)
	return tmp
}

/** Subtract Memory (16 bits) - H?NxZxVxCx */
func (c *CPU) sub16_(reg int, value int) int {
	tmp := reg - value
	c.updateNZVC16(reg, value, tmp)
	return tmp
}

/** Subtract Memory with borrow - H?NxZxVxCx */
func (c *CPU) sbc_(reg int, value int) int {
	borrow := 0
	if c.cc.getC() {
		borrow = 1
	}
	tmp := reg - value - borrow
	c.updateNZVC(reg, value, tmp)
	return tmp
}

/** Subtract Memory from Register A - H?NxZxVxCx */
func (c *CPU) suba(address uint16) {
	value := c.readInt(address)
	c.a.set(c.sub_(c.a.get(), value))
}

/** Subtract Memory from Register D - H?NxZxVxCx */
func (c *CPU) subd(address uint16) {
	value := c.readwInt(address)
	res := c.sub16_(int(c.d()), value)
	c.a.set(res >> 8)
	c.b.set(res & 0xff)
}

/** Compare Memory from Register A - H?NxZxVxCx */
func (c *CPU) cmpa(address uint16) {
	value := c.readInt(address)
	c.sub_(c.a.get(), value)
}

/** Compare Memory from Register X - NxZxVxCx */
func (c *CPU) cmpx(address uint16) {
	value := c.readwInt(address)
	c.sub16_(c.x.get(), value)
}

/** Compare Memory from Register A - H?NxZxVxCx */
func (c *CPU) sbca(address uint16) {
	value := c.readInt(address)
	c.a.set(c.sbc_(c.a.get(), value))
}

/** Logical AND Memory into Register - NxZxV0 */
func (c *CPU) and_(reg int, value int) int {
	tmp := reg & value
	c.updateNZ(tmp)
	c.cc.clearV()
	return tmp
}

/** Logical AND Memory into Register - NxZxV0 */
func (c *CPU) anda(address uint16) {
	value := c.readInt(address)
	c.a.set(c.and_(c.a.get(), value))
}

/** Logical AND Memory and Register - NxZxV0 */
func (c *CPU) bita(address uint16) {
	value := c.readInt(address)
	c.and_(c.a.get(), value)
}

/** Load Register from Memory - NxZxV0 */
func (c *CPU) lda(address uint16) {
	value := c.readInt(address)
	c.updateNZ(value)
	c.cc.clearV()
	c.a.set(value)
}

/** Load Register from Memory - NxZxV0 */
func (c *CPU) ldx(address uint16) {
	value := c.readwInt(address)
	c.updateNZ16(value)
	c.cc.clearV()
	c.x.set(value)
}

/** Exclusive OR - NxZxV0 */
func (c *CPU) eor_(reg int, value int) int {
	tmp := reg ^ value
	c.updateNZ(tmp)
	c.cc.clearV()
	return tmp
}

/** Exclusive OR - NxZxV0 */
func (c *CPU) eora(address uint16) {
	value := c.readInt(address)
	c.a.set(c.eor_(c.a.get(), value))
}

/** Inclusive OR Memory into Register - NxZxV0 */
func (c *CPU) or_(reg int, value int) int {
	tmp := reg | value
	c.updateNZ(tmp)
	c.cc.clearV()
	return tmp
}

/** Inclusive OR Memory into Register - NxZxV0 */
func (c *CPU) ora(address uint16) {
	value := c.readInt(address)
	c.a.set(c.or_(c.a.get(), value))
}

/** Add with Carry into Register - HxNxZxVxCx */
func (c *CPU) adc_(reg int, value int) int {
	carry := 0
	if c.cc.getC() {
		carry = 1
	}
	tmp := reg + value + carry
	c.updateHNZVC(reg, value, tmp)
	return tmp
}

/** Add with Carry into Register - HxNxZxVxCx */
func (c *CPU) adca(address uint16) {
	value := c.readInt(address)
	c.a.set(c.adc_(c.a.get(), value))
}

/** Add Memory into Register - HxNxZxVxCx */
func (c *CPU) add_(reg int, value int) int {
	tmp := reg + value
	c.updateHNZVC(reg, value, tmp)
	return tmp
}

/** Add Memory into Register - HxNxZxVxCx */
func (c *CPU) adda(address uint16) {
	value := c.readInt(address)
	c.a.set(c.add_(c.a.get(), value))
}

/** Store Register A into Memory - NxZxV0 */
func (c *CPU) sta(address uint16) {
	tmp := c.a.get()
	c.writeInt(address, tmp)
	c.updateNZ(tmp)
	c.cc.clearV()
}

/** Store Register X into Memory - NxZxV0 */
func (c *CPU) stx(address uint16) {
	tmp := c.x.get()
	c.writewInt(address, tmp)
	c.updateNZ16(tmp)
	c.cc.clearV()
}
