package core

type Memory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
	Readw(address uint16) uint16
	Writew(address uint16, value uint16)
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
