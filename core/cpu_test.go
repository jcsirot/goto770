package core

import (
	"testing"

	"github.com/assertgo/assert"
)

func TestRegisterD(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.a = 0xe5
	cpu.b = 0xf0
	assert.That(cpu.d()).AsInt().IsEqualTo(0xe5f0)
}

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

func TestDECDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x0a // DEC Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x2b
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x2a)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestDECA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4a // DECA
	cpu.a = 0x2b
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x2a)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDECB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x5a // DECB
	cpu.b = 0x2b
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x2a)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDECZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4a // DECA
	cpu.a = 0x01
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDECNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4a // DECA
	cpu.a = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0xff)
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDECOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4a // DECA
	cpu.a = 0x80
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x7f)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestINCDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x0c // DEC Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x2b
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x2c)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestINCA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4c // DEC Direct
	cpu.a = 0x2b
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x2c)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestINCB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x5c // DEC Direct
	cpu.b = 0x2b
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x2c)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestINCZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4c // DEC Direct
	cpu.a = 0xff
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestINCNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4c // DEC Direct
	cpu.a = 0xfb
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0xfc)
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestINCOverflow(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4c // DEC Direct
	cpu.a = 0x7f
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x80)
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestTSTDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x0d // TST Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x32
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x32)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestTSTA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4d // TSTA
	cpu.a = 0x32
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x32)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestTSTB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x5d // TSTA
	cpu.b = 0x32
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x32)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestTSTNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4d // TSTA
	cpu.a = 0xd8
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0xd8)
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestTSTZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4d // TSTA
	cpu.a = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestJMPDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x0e // JMP Direct
	ram[0x1001] = 0x0a
	cpu.step()
	assert.That(cpu.pc).AsInt().IsEqualTo(0x200a)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestCLRDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.dp = 0x20
	ram[0x1000] = 0x0f // CLR Direct
	ram[0x1001] = 0x0a
	ram[0x200a] = 0x4d
	cpu.step()
	assert.That(ram[0x200a]).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestCLRA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x4f // CLRA
	cpu.a = 0x4d
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestCLRB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x5f // CLRB
	cpu.b = 0x4d
	cpu.step()
	assert.That(cpu.b).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestBRAPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x20 // BRA
	ram[0x1001] = 0x10
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1012)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestLBRAPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x16 // LBRA
	ram[0x1001] = 0x10
	ram[0x1002] = 0x00
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2003)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestBRANegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x20 // LBRA
	ram[0x1001] = 0xfe
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1000)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestLBRANegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x16 // LBRA
	ram[0x1001] = 0xff
	ram[0x1002] = 0x00
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x0f03)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestLBSR(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.s = 0x400
	ram[0x1000] = 0x17 // LBSR
	ram[0x1001] = 0x01
	ram[0x1002] = 0x80
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1183)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0x3fe)
	assert.ThatInt(int(cpu.readw(cpu.s))).IsEqualTo(0x1003)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(9)
}

func TestDAA(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0x62
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x62)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDAALsb(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0x4a
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x50)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDAALsbWithHalfCarry(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0x22
	cpu.setH()
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x28)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDAAMsb(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0xb7
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x17)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDAAMsbWithCarry(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0x19
	cpu.setC()
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x79)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestDAALsbAndMsb(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x19 // DAA
	cpu.a = 0x9a
	cpu.setC()
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestORCC(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.cc = carry | negative
	ram[0x1000] = 0x1a // ORCC
	ram[0x1001] = 0x82
	cpu.step()
	assert.That(cpu.cc).AsInt().IsEqualTo(0x8b)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestANDCC(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.cc = carry | negative
	ram[0x1000] = 0x1c // ANDCC
	ram[0x1001] = 0xfe
	cpu.step()
	assert.That(cpu.cc).AsInt().IsEqualTo(0x08)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestSEX(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1d // SEX
	cpu.b = 0x16
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.That(cpu.b).AsInt().IsEqualTo(0x16)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestSEXNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1d // SEX
	cpu.b = 0xa5
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0xff)
	assert.That(cpu.b).AsInt().IsEqualTo(0xa5)
	assert.ThatBool(cpu.getN()).IsTrue()
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestSEXZero(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1d // SEX
	cpu.b = 0x00
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x00)
	assert.That(cpu.b).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(2)
}

func TestEXGAB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1e // EXG
	ram[0x1001] = 0x89
	cpu.a = 0x27
	cpu.b = 0x0b
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x0b)
	assert.That(cpu.b).AsInt().IsEqualTo(0x27)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(8)
}

func TestEXGXY(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1e // EXG
	ram[0x1001] = 0x12
	cpu.x = 0x1f00
	cpu.y = 0x4000
	cpu.step()
	assert.That(cpu.x).AsInt().IsEqualTo(0x4000)
	assert.That(cpu.y).AsInt().IsEqualTo(0x1f00)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(8)
}

func TestEXGInvalidRegisterCombination(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1e // EXG
	ram[0x1001] = 0x85
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Try to exchange 8-bit with 16-bits registers")
	}()
	cpu.step()
}

func TestEXGInvalidCode(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1e // EXG
	ram[0x1001] = 0x67
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Invalid register code")
	}()
	cpu.step()
}

func TestTFRAtoB(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1f // TFR
	ram[0x1001] = 0x89
	cpu.a = 0x27
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x27)
	assert.That(cpu.b).AsInt().IsEqualTo(0x27)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestTFRPCtoX(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1f // TFR
	ram[0x1001] = 0x51
	cpu.step()
	assert.That(cpu.x).AsInt().IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestTFRInvalidRegisterCombination(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1f // TFR
	ram[0x1001] = 0x85
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Try to transfer 8-bit and 16-bits registers")
	}()
	cpu.step()
}

func TestTFRInvalidRegisterCode(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = 0x1f // TFR
	ram[0x1001] = 0x8f
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Invalid register code")
	}()
	cpu.step()
}

func branchingOpcodeTest(t *testing.T, opcode Word, flags []uint8, branch bool) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	ram[0x1000] = opcode
	ram[0x1001] = 0x10
	for _, flag := range flags {
		cpu.cc |= flag
	}
	cpu.step()
	offset := 0
	if branch {
		offset = 0x10
	}
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002 + offset)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestBHIC0Z0(t *testing.T) {
	branchingOpcodeTest(t, 0x22, []uint8{}, true)
}

func TestBHIC1Z0(t *testing.T) {
	branchingOpcodeTest(t, 0x22, []uint8{carry}, false)
}

func TestBHIC0Z1(t *testing.T) {
	branchingOpcodeTest(t, 0x22, []uint8{zero}, false)
}

func TestBHIC1Z1(t *testing.T) {
	branchingOpcodeTest(t, 0x22, []uint8{carry, zero}, false)
}

func TestBLSC0Z0(t *testing.T) {
	branchingOpcodeTest(t, 0x23, []uint8{}, false)
}

func TestBLSC1Z0(t *testing.T) {
	branchingOpcodeTest(t, 0x23, []uint8{carry}, true)
}

func TestBLSC0Z1(t *testing.T) {
	branchingOpcodeTest(t, 0x23, []uint8{zero}, true)
}

func TestBLSC1Z1(t *testing.T) {
	branchingOpcodeTest(t, 0x23, []uint8{carry, zero}, true)
}

func TestBCCC0(t *testing.T) {
	branchingOpcodeTest(t, 0x24, []uint8{}, true)
}

func TestBCCC1(t *testing.T) {
	branchingOpcodeTest(t, 0x24, []uint8{carry}, false)
}

func TestBLOC0(t *testing.T) {
	branchingOpcodeTest(t, 0x25, []uint8{}, false)
}

func TestBLOC1(t *testing.T) {
	branchingOpcodeTest(t, 0x25, []uint8{carry}, true)
}

func TestBNEZ0(t *testing.T) {
	branchingOpcodeTest(t, 0x26, []uint8{}, true)
}

func TestBNEZ1(t *testing.T) {
	branchingOpcodeTest(t, 0x26, []uint8{zero}, false)
}

func TestBEQZ0(t *testing.T) {
	branchingOpcodeTest(t, 0x27, []uint8{}, false)
}

func TestBEQZ1(t *testing.T) {
	branchingOpcodeTest(t, 0x27, []uint8{zero}, true)
}

func TestBVCV0(t *testing.T) {
	branchingOpcodeTest(t, 0x28, []uint8{}, true)
}

func TestBVCV1(t *testing.T) {
	branchingOpcodeTest(t, 0x28, []uint8{overflow}, false)
}

func TestBVSV0(t *testing.T) {
	branchingOpcodeTest(t, 0x29, []uint8{}, false)
}

func TestBVSV1(t *testing.T) {
	branchingOpcodeTest(t, 0x29, []uint8{overflow}, true)
}

func TestBPLN0(t *testing.T) {
	branchingOpcodeTest(t, 0x2a, []uint8{}, true)
}

func TestBPLN1(t *testing.T) {
	branchingOpcodeTest(t, 0x2a, []uint8{negative}, false)
}

func TestBMIN0(t *testing.T) {
	branchingOpcodeTest(t, 0x2b, []uint8{}, false)
}

func TestBMIN1(t *testing.T) {
	branchingOpcodeTest(t, 0x2b, []uint8{negative}, true)
}

func TestBGEN0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2c, []uint8{}, true)
}

func TestBGEN1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2c, []uint8{negative}, false)
}

func TestBGEN0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2c, []uint8{overflow}, false)
}

func TestBGEN1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2c, []uint8{negative, overflow}, true)
}

func TestBLTN0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2d, []uint8{}, false)
}

func TestBLTN1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2d, []uint8{negative}, true)
}

func TestBLTN0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2d, []uint8{overflow}, true)
}

func TestBLTN1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2d, []uint8{negative, overflow}, false)
}
