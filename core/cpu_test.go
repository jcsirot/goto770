package core

import (
	"testing"

	"github.com/assertgo/assert"
)

func TestNegDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.dp = 0x20
	cpu.pc = 0x04
	ram[0x04] = 0x00 // NEG Direct
	ram[0x05] = 0x0a
	ram[0x200a] = 0x60
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0xa0)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestNegNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.dp = 0x20
	cpu.pc = 0x04
	ram[0x04] = 0x00 // NEG Direct
	ram[0x05] = 0x0a
	ram[0x200a] = 0xa0
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0x60)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestNegZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.dp = 0x20
	cpu.pc = 0x04
	ram[0x04] = 0x00 // NEG Direct
	ram[0x05] = 0x0a
	ram[0x200a] = 0x00
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0x00)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestNegOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.dp = 0x20
	cpu.pc = 0x1000
	ram[0x1000] = 0x00 // NEG Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x80
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0x80)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
}

func TestNegA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.a = 0x60
	ram[0x1000] = 0x40 // NEG A
	cpu.step()
	assert.ThatInt(int(cpu.a)).IsEqualTo(0xa0)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestNegB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.b = 0x60
	ram[0x1000] = 0x50 // COM B
	cpu.step()
	assert.ThatInt(int(cpu.b)).IsEqualTo(0xa0)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestComDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x04
	cpu.dp = 0x20
	ram[0x04] = 0x03 // COM Direct
	ram[0x05] = 0x0a
	ram[0x200a] = 0x1a
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0xe5)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestComExtended(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x04
	ram[0x04] = 0x73 // COM Extended
	ram[0x05] = 0x20
	ram[0x06] = 0x0a
	ram[0x200a] = 0x1a
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0xe5)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestComA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x43 // COM A
	cpu.a = 0x1a
	cpu.step()
	assert.ThatInt(int(cpu.a)).IsEqualTo(0xe5)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestComB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x53 // COM B
	cpu.b = 0x1a
	cpu.step()
	assert.ThatInt(int(cpu.b)).IsEqualTo(0xe5)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
}

func TestLSRDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x04 // LSR Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x66
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0x33)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
}

func TestLSRExtended(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x74 // LSR Direct
	ram[0x1001] = 0x20
	ram[0x1002] = 0x0a
	ram[0x200a] = 0x08
	cpu.step()
	assert.ThatInt(int(ram[0x200a])).IsEqualTo(0x04)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
}

func TestLSRA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x44 // LSRA
	cpu.a = 0x56
	cpu.step()
	assert.ThatInt(int(cpu.a)).IsEqualTo(0x2b)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
}

func TestLSRB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x54 // LSRB
	cpu.b = 0x56
	cpu.step()
	assert.ThatInt(int(cpu.b)).IsEqualTo(0x2b)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
}

func TestLSRZeroAndCarry(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x44 // LSRA
	cpu.a = 0x01
	cpu.step()
	assert.ThatInt(int(cpu.a)).IsEqualTo(0x00)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
}

func TestRORDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x06 // ROR Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x22
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x11)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestRORCarryAndNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x06 // ROR Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x23
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x91)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
}

func TestRORZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x46 // ROR A
	cpu.a = 0
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
}

func TestRORA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x46 // ROR A
	cpu.a = 0x22
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x11)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestRORB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x56 // ROR B
	cpu.b = 0x22
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x11)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestROLDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x09 // ROL Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x1a
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestROLA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x49 // ROLA
	cpu.a = 0x1a
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestROLB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x59 // ROLB
	cpu.b = 0x1a
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestROLZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x49 // ROLA
	cpu.a = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestROLCarryAndOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x49 // ROLA
	cpu.a = 0x81
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x03)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestROLNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x49 // ROLA
	cpu.a = 0x40
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x80)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASRDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x07 // ASR Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x02
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x01)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestASRA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x47 // ASRA
	cpu.a = 0x02
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x01)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASRB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x57 // ASRB
	cpu.b = 0x02
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x01)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASRCarry(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x47 // ASRA
	cpu.a = 0x03
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x01)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASRZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x47 // ASRA
	cpu.a = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASRNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x47 // ASRA
	cpu.a = 0x82
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0xc1)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASLDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x08 // ASL Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x1a
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestASLA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x48 // ASLA
	cpu.a = 0x1a
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASLB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x58 // ASLB
	cpu.b = 0x1a
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x34)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASLZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x48 // ASLA
	cpu.a = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASLNegativeAndOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x48 // ASLA
	cpu.a = 0x42
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x84)
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestASLCarryAndOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x48 // ASLA
	cpu.a = 0x81
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x02)
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}
