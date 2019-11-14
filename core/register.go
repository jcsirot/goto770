package core

import log "github.com/sirupsen/logrus"

type r8 struct {
	n string
	r *int
}

type r16 struct {
	n string
	r *int
}

type register interface {
	name() string
	size() int
	set(value interface{})
	get() int
	uint8() uint8
	int8() int8
	uint16() uint16
	int16() int16
}

// 8-bits register

func (r r8) name() string {
	return r.n
}

func (r r8) size() int {
	return 8
}

func (r r8) checkAndConvert(value interface{}) int {
	v1, ok := value.(int)
	if ok {
		return int(v1 & 0xff)
	}
	v2, ok := value.(uint8)
	if ok {
		return int(v2 & 0xff)
	}
	v3, ok := value.(int8)
	if ok {
		return int(v3) & 0xff
	}
	log.Fatalf("Type conversion error: %T is not an integer type", value)
	return 0
}

func (r r8) set(value interface{}) {
	*r.r = r.checkAndConvert(value)
}

func (r r8) get() int {
	return int(*r.r)
}

func (r r8) uint8() uint8 {
	return uint8(*r.r)
}

func (r r8) int8() int8 {
	return int8(*r.r)
}

func (r r8) uint16() uint16 {
	return uint16(*r.r)
}

func (r r8) int16() int16 {
	return int16(*r.r)
}

// 16-bits register

func (r r16) name() string {
	return r.n
}

func (r r16) size() int {
	return 16
}

func (r r16) checkAndConvert(value interface{}) int {
	v1, ok := value.(int)
	if ok {
		return v1 & 0xffff
	}
	v2, ok := value.(uint8)
	if ok {
		return int(v2) & 0xffff
	}
	v3, ok := value.(int8)
	if ok {
		return int(v3) & 0xffff
	}
	v4, ok := value.(uint16)
	if ok {
		return int(v4) & 0xffff
	}
	v5, ok := value.(int16)
	if ok {
		return int(v5) & 0xffff
	}
	log.Fatalf("Type conversion error: %T is not an integer type", value)
	return 0
}

func (r r16) set(value interface{}) {
	*r.r = r.checkAndConvert(value)
}

func (r r16) get() int {
	return int(*r.r)
}

func (r r16) uint8() uint8 {
	return uint8(*r.r)
}

func (r r16) int8() int8 {
	return int8(*r.r)
}

func (r r16) uint16() uint16 {
	return uint16(*r.r)
}

func (r r16) int16() int16 {
	return int16(*r.r)
}

func (r r16) inc() r16 {
	*r.r++
	return r
}

func (r r16) dec() r16 {
	*r.r--
	return r
}
