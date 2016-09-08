package core

type Word uint8
type DWord uint16

type Ram []Word

func NewRam() Ram {
	return make([]Word, 0x10000)
}

func (r *Ram) read(address uint16) Word {
	return 0
}
