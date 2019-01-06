package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Disassembler", func() {

	BeforeEach(func() {
		var cpu CPU
		cpu.initOpcodes()
	})

	It("Should disassemble instructions with Inherent addressing mode", func() {
		op := opcodes[0x1d]
		ib := []uint8{0x1d}
		testDisassemble(op, ib, "SEX")
	})

	It("Should disassemble TFR and EXG", func() {
		op := opcodes[0x1e]
		ib := []uint8{0x1e, 0x35}
		testDisassemble(op, ib, "EXG U, PC")
		op = opcodes[0x1f]
		ib = []uint8{0x1f, 0x67}
		testDisassemble(op, ib, "TFR A, B")
	})

	It("Should disassemble instructions with Immediate addressing mode", func() {
		op := opcodes[0x8b]
		ib := []uint8{0x8b, 0x05}
		testDisassemble(op, ib, "ADDA #$05")
	})

	It("Should disassemble instructions with long Immediate addressing mode", func() {
		op := opcodes[0x8c]
		ib := []uint8{0x8c, 0xa0, 0xc4}
		testDisassemble(op, ib, "CMPX #$a0c4")
	})

	It("Should disassemble instructions with Direct addressing mode", func() {
		op := opcodes[0x00]
		ib := []uint8{0x00, 0x12}
		testDisassemble(op, ib, "NEG <$12")
	})

	It("Should disassemble instructions with Relative addressing mode", func() {
		op := opcodes[0x27]
		ib := []uint8{0x27, 0xf0}
		testDisassemble(op, ib, "BEQ *+$f0")
	})

	It("Should disassemble instructions with long Relative addressing mode", func() {
		op := opcodes[0x16]
		ib := []uint8{0x16, 0xfa, 0x50}
		testDisassemble(op, ib, "BRA *+$fa50")
	})

	It("Should disassemble instructions with Extended addressing mode", func() {
		op := opcodes[0x76]
		ib := []uint8{0x76, 0xa0, 0x18}
		testDisassemble(op, ib, "ROR $a018")
	})

	It("Should disassemble PSHS and PULS", func() {
		op := opcodes[0x34]
		ib := []uint8{0x34, 0x06}
		testDisassemble(op, ib, "PSHS B,A")
		op = opcodes[0x35]
		ib = []uint8{0x35, 0xf0}
		testDisassemble(op, ib, "PULS X,Y,U,PC")
	})

	It("Should disassemble PSHU and PULU", func() {
		op := opcodes[0x36]
		ib := []uint8{0x36, 0xff}
		testDisassemble(op, ib, "PSHU PC,S,Y,X,DP,B,A,CC")
		op = opcodes[0x37]
		ib = []uint8{0x37, 0x33}
		testDisassemble(op, ib, "PULU CC,A,X,Y")
	})

	It("Should disassemble instructions with Indexed addressing mode, 5 bits offset", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0x2b}
		testDisassemble(op, ib, "NEG 0b,Y")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0x7f}
		testDisassemble(op, ib, "NEG -01,S")
	})

	It("Should disassemble instructions with Indexed addressing mode, auto-increment", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0xc0}
		testDisassemble(op, ib, "NEG ,U+")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xc1}
		testDisassemble(op, ib, "NEG ,U++")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xd1}
		testDisassemble(op, ib, "NEG (,U++)")
	})

	It("Should disassemble instructions with Indexed addressing mode, auto-decrement", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0xe2}
		testDisassemble(op, ib, "NEG ,-S")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xe3}
		testDisassemble(op, ib, "NEG ,--S")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xf3}
		testDisassemble(op, ib, "NEG (,--S)")
	})

	It("Should disassemble instructions with Indexed addressing mode, accumulator register", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0x86}
		testDisassemble(op, ib, "NEG A,X")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xe5}
		testDisassemble(op, ib, "NEG B,S")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xab}
		testDisassemble(op, ib, "NEG D,Y")
	})

	It("Should disassemble instructions with Indexed addressing mode, 7 bits offset", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0x88, 0x6a}
		testDisassemble(op, ib, "NEG 6a,X")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xc8, 0xfc}
		testDisassemble(op, ib, "NEG -04,U")
	})

	It("Should disassemble instructions with Indexed addressing mode, 15 bits offset", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0x89, 0x6a, 0x01}
		testDisassemble(op, ib, "NEG 6a01,X")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0xc9, 0xff, 0xe9}
		testDisassemble(op, ib, "NEG -0017,U")
	})

	It("Should disassemble instructions with Indexed addressing mode, PC register with offset", func() {
		op := opcodes[0x60]
		ib := []uint8{0x60, 0x8c, 0x6a}
		testDisassemble(op, ib, "NEG 6a,PC")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0x9c, 0x6a}
		testDisassemble(op, ib, "NEG (6a,PC)")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0x8d, 0xff, 0xe9}
		testDisassemble(op, ib, "NEG -0017,PC")
		op = opcodes[0x60]
		ib = []uint8{0x60, 0x9d, 0xff, 0xe9}
		testDisassemble(op, ib, "NEG (-0017,PC)")
	})
})

func testDisassemble(op opcode, instBuf []uint8, expected string) {
	str, size := Disassemble(op, instBuf)
	//fmt.Println(str)
	//fmt.Println(size)
	ExpectWithOffset(1, str).To(Equal(expected))
	ExpectWithOffset(1, size).To(Equal(len(instBuf)))
}
