package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Adressing mode", func() {
	var (
		cpu CPU
		ram Memory
	)

	BeforeEach(func() {
		ram = NewRam()
		cpu.Initialize(ram)
	})

	It("should work with Indexed 5-bits Offset Positive", func() {
		cpu.pc.set(0x1001)
		cpu.y.set(0x2000)
		cpu.write(0x1001, 0x25) // EA = Y + 5-bit offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x2005))
		Expect(cpu.pc.get()).To(BeEquivalentTo(0x1002))
		Expect(cpu.clock).To(BeEquivalentTo(1))
	})

	It("should work with Indexed 5-bits Offset Negative", func() {
		cpu.pc.set(0x1001)
		cpu.y.set(0x2000)
		cpu.write(0x1001, 0x3e) // EA = Y + 5-bit offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x1ffe))
		Expect(cpu.pc.get()).To(BeEquivalentTo(0x1002))
		Expect(cpu.clock).To(BeEquivalentTo(1))
	})

	It("should work with Indexed autoincrement by 1", func() {
		cpu.pc.set(0x1001)
		cpu.u.set(0x2000)
		cpu.write(0x1001, 0xc0) // EA = U+
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x2000))
		Expect(cpu.u.get()).To(BeEquivalentTo(0x2001))
		Expect(cpu.pc.get()).To(BeEquivalentTo(0x1002))
	})

	It("should work with Indexed autodecrement by 1", func() {
		cpu.pc.set(0x1001)
		cpu.s.set(0x2000)
		cpu.write(0x1001, 0xe2) // EA = -S
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x1fff))
		Expect(cpu.s.get()).To(BeEquivalentTo(0x1fff))
		Expect(cpu.pc.get()).To(BeEquivalentTo(0x1002))
	})

	It("should work with Indexed Idx B Positive", func() {
		cpu.pc.set(0x1001)
		cpu.x.set(0x2000)
		cpu.b.set(0x05)
		cpu.write(0x1001, 0x85) // EA = ,X ± ACCB offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x2005))
	})

	It("should work with Indexed Idx B Negative", func() {
		cpu.pc.set(0x1001)
		cpu.x.set(0x2000)
		cpu.b.set(0xf0)
		cpu.write(0x1001, 0x85) // EA = ,X ± ACCB offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x1FF0))
	})

	It("should work with Indexed Idx A Positive", func() {
		cpu.pc.set(0x1001)
		cpu.x.set(0x2000)
		cpu.a.set(0x32)
		cpu.write(0x1001, 0x86) // EA = ,X ± ACCA offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x2032))
		Expect(cpu.clock).To(BeEquivalentTo(1))

	})

	It("should work with Indexed Idx A Negative", func() {
		cpu.pc.set(0x1001)
		cpu.x.set(0x2000)
		cpu.a.set(0x32)
		cpu.write(0x1001, 0x86) // EA = ,X ± ACCA offset
		address := cpu.indexed()
		Expect(address).To(BeEquivalentTo(0x2032))
		Expect(cpu.clock).To(BeEquivalentTo(1))
	})
})
