package core

import (
	"fmt"
	"strings"
)

type Memory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
	Readw(address uint16) uint16
	Writew(address uint16, value uint16)

	Dump()
}

type memoryImpl struct {
	RAM []uint8
}

func NewRam() Memory {
	return &memoryImpl{RAM: make([]uint8, 0x10000)}
}

func (mem *memoryImpl) Read(address uint16) uint8 {
	return (*mem).RAM[address]
}

func (mem *memoryImpl) Write(address uint16, value uint8) {
	(*mem).RAM[address] = value
}

func (mem *memoryImpl) Readw(address uint16) uint16 {
	hi := mem.Read(address)
	lo := mem.Read(address + 1)
	return (uint16(hi)<<8 | uint16(lo))
}

func (mem *memoryImpl) Writew(address uint16, value uint16) {
	mem.Write(address+1, uint8(value&0xff))
	mem.Write(address, uint8(value>>8))
}

func (mem *memoryImpl) Dump() {
	var sb strings.Builder
	for a := 0; a < 0x10000; a += 16 {
		sb.WriteString(fmt.Sprintf("%04x | ", a))
		hexa := []string{}
		for _, x := range mem.RAM[a : a+16] {
			hexa = append(hexa, fmt.Sprintf("%02x", x))
		}
		sb.WriteString(strings.Join(hexa[0:8], " "))
		sb.WriteString("  ")
		sb.WriteString(strings.Join(hexa[8:16], " "))
		sb.WriteString("\n")
	}
	fmt.Printf(sb.String())
}
