package core

import (
	"testing"

	"github.com/assertgo/assert"
)

func TestIndexed5bitsOffsetPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.y = 0x2000
	ram[0x1001] = 0x25 // EA = Y + 5-bit offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x2005)
	assert.That(cpu.pc).AsInt().IsEqualTo(0x1002)
	assert.That(int(cpu.clock)).AsInt().IsEqualTo(1)
}

func TestIndexed5bitsOffsetNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.y = 0x2000
	ram[0x1001] = 0x3e // EA = Y + 5-bit offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x1ffe)
	assert.That(cpu.pc).AsInt().IsEqualTo(0x1002)
	assert.That(int(cpu.clock)).AsInt().IsEqualTo(1)
}

func TestIndexedAutoIncrement1(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.u = 0x2000
	ram[0x1001] = 0xc0 // EA = U+
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x2000)
	assert.That(cpu.u).AsInt().IsEqualTo(0x2001)
	assert.That(cpu.pc).AsInt().IsEqualTo(0x1002)
}

func TestIndexedAutoDecrement1(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.s = 0x2000
	ram[0x1001] = 0xe2 // EA = -S
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x1fff)
	assert.That(cpu.s).AsInt().IsEqualTo(0x1fff)
	assert.That(cpu.pc).AsInt().IsEqualTo(0x1002)
}

func TestIndexedIdxbPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.x = 0x2000
	cpu.b = 0x05
	ram[0x1001] = 0x85 // EA = ,X ± ACCB offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x2005)
}

func TestIndexedIdxbNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.x = 0x2000
	cpu.b = 0xf0
	ram[0x1001] = 0x85 // EA = ,X ± ACCB offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x1FF0)
}

func TestIndexedIdxaPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.x = 0x2000
	cpu.a = 0x32
	ram[0x1001] = 0x86 // EA = ,X ± ACCA offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x2032)
	assert.That(int(cpu.clock)).AsInt().IsEqualTo(1)
}

func TestIndexedIdxaNegative(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1001
	cpu.x = 0x2000
	cpu.a = 0xa5
	ram[0x1001] = 0x86 // EA = ,X ± ACCA offset
	address := cpu.indexed()
	assert.That(address).AsInt().IsEqualTo(0x1FA5)
	assert.That(int(cpu.clock)).AsInt().IsEqualTo(1)
}
