package core

type Word uint8
type DWord uint16
type Address uint16

type Memory []Word

func NewRam() Memory {
	return make([]Word, 0x10000)
}

func (r *Memory) read(address uint16) Word {
	return 0
}
