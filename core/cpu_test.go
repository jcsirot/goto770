package core

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newCPU() *CPU {
	var cpu = CPU{}
	ram := NewRam()
	cpu.Initialize(ram)
	return &cpu
}

var _ = Describe("CPU", func() {
	var (
		cpu CPU
	)

	BeforeEach(func() {
		cpu = CPU{}
		ram := NewRam()
		cpu.Initialize(ram)
	})

	It("Register D should be the concatenation of registers A and B", func() {
		cpu.a.set(0xe5)
		cpu.b.set(0xf0)
		ExpectD(cpu, 0xe5f0)
	})

	Context("[SWI]", func() {

		It("should implement SWI", func() {
			cpu.pc.set(0x1000)
			cpu.s.set(0x2000)
			cpu.a.set(0x2f)
			cpu.dp.set(0x18)
			cpu.x.set(0xe0ff)
			cpu.writew(0xfffa, 0xc200)
			cpu.write(0x1000, 0x3f) // SWI
			cpu.step()

			ExpectWord(cpu, 0x2000-2, 0x1001)
			ExpectWord(cpu, 0x2000-4, cpu.u.get())
			ExpectWord(cpu, 0x2000-6, cpu.y.get())
			ExpectWord(cpu, 0x2000-8, cpu.x.get())
			ExpectMemory(cpu, 0x2000-9, cpu.dp.get())
			ExpectMemory(cpu, 0x2000-10, cpu.b.get())
			ExpectMemory(cpu, 0x2000-11, cpu.a.get())
			ExpectMemory(cpu, 0x2000-12, 0x80)
			ExpectPC(cpu, 0xc200)
			ExpectClock(cpu, 19)
			ExpectCCR(cpu, "EFI", "ZVNHC")
		})

		It("should implement SWI2", func() {
			cpu.pc.set(0x1000)
			cpu.s.set(0x2000)
			cpu.a.set(0x2f)
			cpu.dp.set(0x18)
			cpu.x.set(0xe0ff)
			cpu.writew(0xfff4, 0xc200)
			cpu.writew(0x1000, 0x103f) // SWI2
			cpu.cc.setC()
			cpu.step()

			ExpectWord(cpu, 0x2000-2, 0x1002)
			ExpectWord(cpu, 0x2000-4, cpu.u.get())
			ExpectWord(cpu, 0x2000-6, cpu.y.get())
			ExpectWord(cpu, 0x2000-8, cpu.x.get())
			ExpectMemory(cpu, 0x2000-9, cpu.dp.get())
			ExpectMemory(cpu, 0x2000-10, cpu.b.get())
			ExpectMemory(cpu, 0x2000-11, cpu.a.get())
			ExpectMemory(cpu, 0x2000-12, 0x81)
			ExpectPC(cpu, 0xc200)
			ExpectClock(cpu, 20)
			ExpectCCR(cpu, "EC", "ZVNHFI")
		})

		It("should implement SWI3", func() {
			cpu.pc.set(0x1000)
			cpu.s.set(0x2000)
			cpu.a.set(0x2f)
			cpu.dp.set(0x18)
			cpu.x.set(0xe0ff)
			cpu.writew(0xfff2, 0xaf10)
			cpu.writew(0x1000, 0x113f) // SWI3
			cpu.cc.setZ()
			cpu.cc.setV()
			cpu.step()

			ExpectWord(cpu, 0x2000-2, 0x1002)
			ExpectWord(cpu, 0x2000-4, cpu.u.get())
			ExpectWord(cpu, 0x2000-6, cpu.y.get())
			ExpectWord(cpu, 0x2000-8, cpu.x.get())
			ExpectMemory(cpu, 0x2000-9, cpu.dp.get())
			ExpectMemory(cpu, 0x2000-10, cpu.b.get())
			ExpectMemory(cpu, 0x2000-11, cpu.a.get())
			ExpectMemory(cpu, 0x2000-12, 0x86)
			ExpectPC(cpu, 0xaf10)
			ExpectClock(cpu, 20)
			ExpectCCR(cpu, "EZV", "CNHFI")
		})
	})

	Context("[NEG]", func() {

		It("[Direct] should implement NEG with Direct addressing mode", func() {
			cpu.dp.set(0x20)
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x00) // NEG Direct
			cpu.write(0x05, 0x0a)
			cpu.write(0x200a, 0x60)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0xa0)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Direct] should implement NEG with Direct addressing mode and negative value", func() {
			cpu.dp.set(0x20)
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x00) // NEG Direct
			cpu.write(0x05, 0x0a)
			cpu.write(0x200a, 0xa0)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x60)
			ExpectCCR(cpu, "C", "NZV")
		})

		It("[Direct] should implement NEG with Direct addressing mode and zero value", func() {
			cpu.dp.set(0x20)
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x00) // NEG Direct
			cpu.write(0x05, 0x0a)
			cpu.write(0x200a, 0x00)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x00)
			ExpectCCR(cpu, "Z", "CNV")
		})

		It("[Direct] should implement NEG with Direct addressing mode and bit V", func() {
			cpu.dp.set(0x20)
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x00) // NEG Direct
			cpu.write(0x05, 0x0a)
			cpu.write(0x200a, 0x80)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x80)
			ExpectCCR(cpu, "CNV", "Z")
		})

		It("[Inherent] should implement NEGA", func() {
			cpu.pc.set(0x1000)
			cpu.a.set(0x5d)
			cpu.write(0x1000, 0x40) // NEG A
			cpu.step()
			ExpectA(cpu, 0xa3)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Inherent] should implement NEGB", func() {
			cpu.pc.set(0x1000)
			cpu.b.set(0x60)
			cpu.write(0x1000, 0x50) // COM B
			cpu.step()
			ExpectB(cpu, 0xa0)
			ExpectCCR(cpu, "CN", "ZV")
		})
	})

	Context("[COM]", func() {

		It("[Direct] should implement COM with Direct addressing mode", func() {
			cpu.dp.set(0x20)
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x03) // COM Direct
			cpu.write(0x05, 0x0a)
			cpu.write(0x200a, 0x1a)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0xe5)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Extended] should implement COM with Extended addressing mode", func() {
			cpu.pc.set(0x04)
			cpu.write(0x04, 0x73) // COM Extended
			cpu.write(0x05, 0x20)
			cpu.write(0x06, 0x0a)
			cpu.write(0x200a, 0x1a)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0xe5)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Inherent] should implement COMA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x43) // COM A
			cpu.a.set(0x1a)
			cpu.step()
			ExpectA(cpu, 0xe5)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Inherent] should implement COMB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x53) // COM B
			cpu.b.set(0x1a)
			cpu.step()
			ExpectB(cpu, 0xe5)
			ExpectCCR(cpu, "CN", "ZV")
		})
	})

	Context("[LSR]", func() {

		It("[Direct] should implement LSR with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x04) // LSR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x66)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x33)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Direct] should implement LSR with Direct addressing mode and bits CZ", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x04) // LSR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x01)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x00)
			ExpectCCR(cpu, "CZ", "N")
		})

		It("[Extended] should implement LSR with Extended addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x74) // LSR Direct
			cpu.write(0x1001, 0x20)
			cpu.write(0x1002, 0x0a)
			cpu.write(0x200a, 0x08)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x04)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement LSRA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x44) // LSRA
			cpu.a.set(0x56)
			cpu.step()
			ExpectA(cpu, 0x2b)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement LSRB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x54) // LSRA
			cpu.b.set(0x56)
			cpu.step()
			ExpectB(cpu, 0x2b)
			ExpectCCR(cpu, "", "CNZ")
		})
	})

	Context("[ROR]", func() {

		It("[Direct] should implement ROR with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x06) // ROR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x22)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x11)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Direct] should implement ROR with Direct addressing mode and bits C", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x06) // ROR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x23)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x11)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "C", "ZN")
		})

		It("[Direct] should implement ROR with Direct addressing mode and bit N (C propagation)", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x06) // ROR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x22)
			cpu.cc.setC()
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x91)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "N", "CZ")
		})

		It("[Inherent] should implement RORA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x46) // ROR A
			cpu.a.set(0x22)
			cpu.step()

			ExpectA(cpu, 0x11)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement RORB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x56) // ROR B
			cpu.b.set(0x22)
			cpu.step()

			ExpectB(cpu, 0x11)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement RORA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x46) // ROR A
			cpu.a.set(0)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "CN")
		})
	})

	Context("[ROL]", func() {

		It("[Direct] should implement ROL with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x09) // ROL Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x1a)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x34)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ROLA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x1a)
			cpu.step()

			ExpectA(cpu, 0x34)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ROLB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x59) // ROLB
			cpu.b.set(0x1a)
			cpu.step()

			ExpectB(cpu, 0x34)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ROLA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "CNV")
		})

		It("[Inherent] should implement ROLA with bits CZV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x80)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "CZV", "N")
		})

		It("[Inherent] should implement ROLA with bits CN", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0xc1)
			cpu.step()

			ExpectA(cpu, 0x82)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Inherent] should implement ROLA with bits CV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x81)
			cpu.step()

			ExpectA(cpu, 0x02)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "CV", "NZ")
		})

		It("[Inherent] should implement ROLA with bits NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x40)
			cpu.step()

			ExpectA(cpu, 0x80)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NV", "CZ")
		})

		It("[Inherent] should implement ROLA with bit C propagation", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x49) // ROLA
			cpu.a.set(0x20)
			cpu.cc.setC()
			cpu.step()

			ExpectA(cpu, 0x41)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NCZV")
		})
	})

	Context("[ASR]", func() {

		It("[Direct] should implement ASR with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x07) // ASR Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x52)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x29)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement ASRA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x47) // ASRA
			cpu.a.set(0x02)
			cpu.step()

			ExpectA(cpu, 0x01)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement ASRB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x57) // ASRB
			cpu.b.set(0x02)
			cpu.step()

			ExpectB(cpu, 0x01)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZ")
		})

		It("[Inherent] should implement ASRA with bit C", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x47) // ASRA
			cpu.a.set(0x03)
			cpu.step()

			ExpectA(cpu, 0x01)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "C", "NZ")
		})

		It("[Inherent] should implement ASRA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x47) // ASRA
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "CN")
		})

		It("[Inherent] should implement ASRA with bit CZ", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x47) // ASRA
			cpu.a.set(0x01)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "CZ", "N")
		})

		It("[Inherent] should implement ASRA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x47) // ASRA
			cpu.a.set(0x82)
			cpu.step()

			ExpectA(cpu, 0xc1)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "CZ")
		})
	})

	Context("[CMP]", func() {

		It("[Immediate] should implement CMPA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x81) // CMPA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x3e)
			cpu.step()

			ExpectA(cpu, 0x3e)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement CMPA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x81) // CMPA
			cpu.write(0x1001, 0x06)
			cpu.a.set(0x04)
			cpu.step()

			ExpectA(cpu, 0x04)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NC", "ZV")
		})

		It("[Direct] should implement CMPA with bit NC", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x91) // CMPA
			cpu.write(0x1001, 0x40)
			cpu.write(0x2040, 0x02)
			cpu.a.set(0x01)

			cpu.step()

			ExpectA(cpu, 0x01)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "NC", "ZV")
		})

		It("[Immediate] should implement CMPB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc1) // CMPB
			cpu.write(0x1001, 0x04)
			cpu.b.set(0x3e)
			cpu.step()

			ExpectB(cpu, 0x3e)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement CMPX", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8c) // CMPX
			cpu.writew(0x1001, 0x0104)
			cpu.x.set(0x07f9)
			cpu.step()

			ExpectX(cpu, 0x07f9)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement CMPX with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8c) // CMPX
			cpu.writew(0x1001, 0x0104)
			cpu.x.set(0x0104)
			cpu.step()

			ExpectX(cpu, 0x0104)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "Z", "NVC")
		})

		It("[Immediate] should implement CMPX with bit NC", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8c) // CMPX
			cpu.writew(0x1001, 0x0104)
			cpu.x.set(0x0102)
			cpu.step()

			ExpectX(cpu, 0x0102)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "NC", "ZV")
		})
	})

	Context("[SUB]", func() {

		It("[Immediate] should implement SUBA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x80) // SUBA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x3e)
			cpu.step()

			ExpectA(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SUBB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc0) // SUBB
			cpu.write(0x1001, 0x04)
			cpu.b.set(0x3e)
			cpu.step()

			ExpectB(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Direct] should implement SUBA with Direct addessing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x90) // SUBA
			cpu.write(0x1001, 0x40)
			cpu.write(0x2040, 0x04)
			cpu.a.set(0x3e)
			cpu.step()

			ExpectA(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SUBA with Immediate addessing mode with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x80) // SUBA
			cpu.write(0x1001, 0x2b)
			cpu.a.set(0x2b)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NVC")
		})

		It("[Immediate] should implement SUBA with Immediate addessing mode with bit NC", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x80) // SUBA
			cpu.write(0x1001, 0x06)
			cpu.a.set(0x04)
			cpu.step()

			ExpectA(cpu, 0xfe)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NC", "ZV")
		})

		It("[Immediate] should implement SUBD with Immediate addessing mode", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x83) // SUBD
			cpu.writew(0x1001, 0x08f3)
			cpu.a.set(0x3e)
			cpu.b.set(0xa0)
			cpu.step()

			ExpectD(cpu, 0x35ad)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SUBD with Immediate addessing mode and bit CN", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x83) // SUBD
			cpu.writew(0x1001, 0x48f3)
			cpu.a.set(0x3e)
			cpu.b.set(0xa0)
			cpu.step()

			ExpectD(cpu, 0xf5ad)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "CN", "ZV")
		})

		It("[Extended] should implement SUBA with Extended addessing mode", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xb0) // SUBA
			cpu.writew(0x1001, 0x30a0)
			cpu.write(0x30a0, 0x0e)
			cpu.a.set(0x3e)
			cpu.step()

			ExpectA(cpu, 0x30)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Extended] should implement SUBD with Extended addessing mode", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xb3) // SUBD
			cpu.writew(0x1001, 0x20a0)
			cpu.writew(0x20a0, 0x1401)
			cpu.a.set(0x16)
			cpu.b.set(0x50)
			cpu.step()

			ExpectD(cpu, 0x024f)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 7)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Extended] should implement SUBD with Extended addessing mode and bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xb3) // SUBD
			cpu.writew(0x1001, 0x20a0)
			cpu.writew(0x20a0, 0x18b8)
			cpu.a.set(0x18)
			cpu.b.set(0xb8)
			cpu.step()

			ExpectD(cpu, 0)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 7)
			ExpectCCR(cpu, "Z", "CNV")
		})
	})

	Context("[SBC]", func() {

		It("[Immediate] should implement SBCA with Immediate addessing mode", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x82) // SBCA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x3e)

			cpu.step()

			ExpectA(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SBCA with Immediate addessing mode with bit C", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x82) // SBCA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x3e)
			cpu.cc.setC()

			cpu.step()

			ExpectA(cpu, 0x39)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SBCA with Immediate addessing mode with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x82) // SBCA
			cpu.write(0x1001, 0x3e)
			cpu.a.set(0x3e)

			cpu.step()

			ExpectA(cpu, 0)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NVC")
		})

		It("[Immediate] should implement SBCA with Immediate addessing mode with bit NC", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x82) // SBCA
			cpu.write(0x1001, 0x0b)
			cpu.a.set(0x09)

			cpu.step()

			ExpectA(cpu, 0xfe)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NC", "ZV")
		})

		It("[Immediate] should implement SBCA with Immediate addessing mode with bit V", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x82) // SBCA
			cpu.write(0x1001, 0x03)
			cpu.a.set(0x82)

			cpu.step()

			ExpectA(cpu, 0x7f)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "V", "NZC")
		})

		It("[Direct] should implement SBCA", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x30)
			cpu.write(0x1000, 0x92) // SBCA
			cpu.write(0x1001, 0xa0)
			cpu.write(0x30a0, 0x04)
			cpu.a.set(0x3e)

			cpu.step()

			ExpectA(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement SBCB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc2) // SBCB
			cpu.write(0x1001, 0x04)
			cpu.b.set(0x3e)

			cpu.step()

			ExpectB(cpu, 0x3a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVC")
		})
	})

	Context("[ASL]", func() {

		It("[Direct] should implement ASL with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x08) // ASL Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x1a)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x34)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ASLA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x48) // ASLA
			cpu.a.set(0x1a)
			cpu.step()

			ExpectA(cpu, 0x34)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ASLB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x58) // ASLB
			cpu.b.set(0x1a)
			cpu.step()

			ExpectB(cpu, 0x34)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "CNZV")
		})

		It("[Inherent] should implement ASLA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x48) // ASLA
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "CNV")
		})

		It("[Inherent] should implement ASLA with bit NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x48) // ASLA
			cpu.a.set(0x42)
			cpu.step()

			ExpectA(cpu, 0x84)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NV", "CZ")
		})

		It("[Inherent] should implement ASLA with bit CV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x48) // ASLA
			cpu.a.set(0x81)
			cpu.step()

			ExpectA(cpu, 0x02)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "CV", "NZ")
		})
	})

	Context("[DEC]", func() {

		It("[Direct] should implement DEC with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x0a) // DEC Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x2b)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x2a)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement DECA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4a) // DECA
			cpu.a.set(0x2b)
			cpu.step()
			ExpectA(cpu, 0x2a)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement DECA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x5a) // DECB
			cpu.b.set(0x64)
			cpu.step()
			ExpectB(cpu, 0x63)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement DECA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4a) // DECA
			cpu.a.set(0x01)
			cpu.step()
			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Inherent] should implement DECA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4a) // DECA
			cpu.a.set(0x00)
			cpu.step()
			ExpectA(cpu, 0xff)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Inherent] should implement DECA with bit V", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4a) // DECA
			cpu.a.set(0x80)
			cpu.step()
			ExpectA(cpu, 0x7f)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "V", "NZ")
		})
	})

	Context("[AND]", func() {

		It("[Immediate] should implement ANDA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x84) // ANDA
			cpu.write(0x1001, 0x91)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0x11)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement ANDA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x84) // ANDA
			cpu.write(0x1001, 0x80)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement ANDA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x84) // ANDA
			cpu.write(0x1001, 0x83)
			cpu.a.set(0x81)
			cpu.step()

			ExpectA(cpu, 0x81)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement ANDB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc4) // ANDA
			cpu.write(0x1001, 0x91)
			cpu.b.set(0x55)
			cpu.step()

			ExpectB(cpu, 0x11)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})
	})

	Context("[BIT]", func() {

		It("[Immediate] should implement BITA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x85) // BITA
			cpu.write(0x1001, 0x80)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0x55)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement BITA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x85) // BITA
			cpu.write(0x1001, 0x8a)
			cpu.a.set(0x86)
			cpu.step()

			ExpectA(cpu, 0x86)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement BITB with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc5) // BITB
			cpu.write(0x1001, 0x80)
			cpu.b.set(0x55)
			cpu.step()

			ExpectB(cpu, 0x55)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement BITB with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc5) // BITB
			cpu.write(0x1001, 0x8a)
			cpu.b.set(0x86)
			cpu.step()

			ExpectB(cpu, 0x86)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})
	})

	Context("[LD]", func() {

		It("[Immediate] should implement LDA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x86) // LDA
			cpu.write(0x1001, 0x67)
			cpu.step()

			ExpectA(cpu, 0x67)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x86) // LDA
			cpu.write(0x1001, 0x00)
			cpu.a.set(0xf5)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement LDA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x86) // LDA
			cpu.write(0x1001, 0xa1)
			cpu.step()

			ExpectA(cpu, 0xa1)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement LDB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc6) // LDB
			cpu.write(0x1001, 0x67)
			cpu.step()

			ExpectB(cpu, 0x67)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDB with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc6) // LDB
			cpu.write(0x1001, 0x00)
			cpu.a.set(0xf5)
			cpu.step()

			ExpectB(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement LDB with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc6) // LDB
			cpu.write(0x1001, 0xa1)
			cpu.step()

			ExpectB(cpu, 0xa1)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement LDD", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xcc) // LDD
			cpu.writew(0x1001, 0x2a89)
			cpu.step()

			ExpectD(cpu, 0x2a89)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDD with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xcc) // LDD
			cpu.writew(0x1001, 0x00)
			cpu.a.set(0xf5)
			cpu.b.set(0x89)
			cpu.step()

			ExpectD(cpu, 0x0000)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement LDD with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xcc) // LDD
			cpu.writew(0x1001, 0xa189)
			cpu.step()

			ExpectD(cpu, 0xa189)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Direct] should implement LDD", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x3f)
			cpu.write(0x1000, 0xdc) // LDD
			cpu.write(0x1001, 0x80)
			cpu.writew(0x3f80, 0x45ab)
			cpu.step()

			ExpectD(cpu, 0x45ab)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDU", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xce) // LDU
			cpu.writew(0x1001, 0x2a89)
			cpu.step()

			ExpectU(cpu, 0x2a89)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDU with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xce) // LDU
			cpu.writew(0x1001, 0x00)
			cpu.u.set(0xf589)
			cpu.step()

			ExpectU(cpu, 0x0000)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement LDU with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xce) // LDU
			cpu.writew(0x1001, 0xa189)
			cpu.step()

			ExpectU(cpu, 0xa189)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Direct] should implement LDU", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x55)
			cpu.write(0x1000, 0xde) // LDU
			cpu.write(0x1001, 0x20)
			cpu.writew(0x5520, 0x20b0)
			cpu.step()

			ExpectU(cpu, 0x20b0)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDX", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8e) // LDX
			cpu.writew(0x1001, 0x2a89)
			cpu.step()

			ExpectX(cpu, 0x2a89)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 3)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement LDY", func() {
			cpu.pc.set(0x1000)
			cpu.writew(0x1000, 0x108e) // LDY
			cpu.writew(0x1002, 0x2a89)
			cpu.step()

			ExpectY(cpu, 0x2a89)
			ExpectPC(cpu, 0x1004)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZV")
		})
	})

	Context("[ST]", func() {

		It("[Direct] should implement STA", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x97) // STA
			cpu.write(0x1001, 0x60)
			cpu.a.set(0x76)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0x76)
			ExpectA(cpu, 0x76)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Direct] should implement STA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x97) // STA
			cpu.write(0x1001, 0x60)
			cpu.a.set(0x00)
			cpu.writew(0x2060, 0xff)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Direct] should implement STA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x97) // STA
			cpu.write(0x1001, 0x60)
			cpu.a.set(0xba)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0xba)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Direct] should implement STB", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0xd7) // STB
			cpu.write(0x1001, 0x60)
			cpu.b.set(0x76)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0x76)
			ExpectB(cpu, 0x76)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Direct] should implement STB with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0xd7) // STB
			cpu.write(0x1001, 0x60)
			cpu.b.set(0x00)
			cpu.writew(0x2060, 0xff)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Direct] should implement STB with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0xd7) // STB
			cpu.write(0x1001, 0x60)
			cpu.b.set(0xba)
			cpu.step()

			ExpectMemory(cpu, 0x2060, 0xba)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Direct] should implement STD", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0xdd) // STD
			cpu.write(0x1001, 0x60)
			cpu.a.set(0x58)
			cpu.b.set(0xb0)
			cpu.step()

			ExpectWord(cpu, 0x2060, 0x58b0)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Direct] should implement STU", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0xb0)
			cpu.write(0x1000, 0xdf) // STU
			cpu.write(0x1001, 0xa0)
			cpu.u.set(0x19f0)
			cpu.step()

			ExpectWord(cpu, 0xb0a0, 0x19f0)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Direct] should implement STU with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0xb0)
			cpu.write(0x1000, 0xdf) // STU
			cpu.write(0x1001, 0xa0)
			cpu.u.set(0xb851)
			cpu.step()

			ExpectWord(cpu, 0xb0a0, 0xb851)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Direct] should implement STU with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0xb0)
			cpu.write(0x1000, 0xdf) // STU
			cpu.write(0x1001, 0xa0)
			cpu.u.set(0x0000)
			cpu.step()

			ExpectWord(cpu, 0xb0a0, 0x0000)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Direct] should implement STX", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x9f) // STX
			cpu.write(0x1001, 0x60)
			cpu.x.set(0x6b08)
			cpu.step()

			ExpectWord(cpu, 0x2060, 0x6b08)
			ExpectX(cpu, 0x6b08)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 5)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Direct] should implement STY", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.writew(0x1000, 0x109f) // STY
			cpu.write(0x1002, 0x60)
			cpu.y.set(0x6b08)
			cpu.step()

			ExpectWord(cpu, 0x2060, 0x6b08)
			ExpectY(cpu, 0x6b08)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "NZV")
		})
	})

	Context("[EOR][XOR]", func() {

		It("[Immediate] should implement EORA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x88) // EORA
			cpu.write(0x1001, 0x55)
			cpu.a.set(0x31)
			cpu.step()

			ExpectA(cpu, 0x64)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement EORA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x88) // EORA
			cpu.write(0x1001, 0x55)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement EORA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x88) // EORA
			cpu.write(0x1001, 0x84)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0xd1)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement EORB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc8) // EORB
			cpu.write(0x1001, 0x7a)
			cpu.b.set(0x22)
			cpu.step()

			ExpectB(cpu, 0x58)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})
	})

	Context("[ADC]", func() {

		It("[Immediate] should implement ADCA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x50)
			cpu.step()

			ExpectA(cpu, 0x54)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVHC")
		})

		It("[Immediate] should implement ADCA with Carry", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x50)
			cpu.cc.setC()
			cpu.step()

			ExpectA(cpu, 0x55)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVHC")
		})

		It("[Immediate] should implement ADCA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0x01)
			cpu.a.set(0xf0)
			cpu.step()

			ExpectA(cpu, 0xf1)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZVHC")
		})

		It("[Immediate] should implement ADCA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0x00)
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NVHC")
		})

		It("[Immediate] should implement ADCA with bit NHC", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0xff)
			cpu.a.set(0xff)
			cpu.step()

			ExpectA(cpu, 0xfe)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NCH", "ZV")
		})

		It("[Immediate] should implement ADCA with bit NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x89) // ADCA
			cpu.write(0x1001, 0x40)
			cpu.a.set(0x42)
			cpu.step()

			ExpectA(cpu, 0x82)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NV", "ZHC")
		})

		It("[Immediate] should implement ADCB with Carry", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc9) // ADCB
			cpu.write(0x1001, 0x0a)
			cpu.b.set(0x62)
			cpu.cc.setC()
			cpu.step()

			ExpectB(cpu, 0x6d)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVHC")
		})
	})

	Context("[INC]", func() {

		It("[Direct] should implement INC with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x0c) // INC Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x2b)
			cpu.step()

			ExpectMemory(cpu, 0x200a, 0x2c)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement INCA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4c) // INCA
			cpu.a.set(0x2b)
			cpu.step()

			ExpectA(cpu, 0x2c)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement INCB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x5c) // INCB
			cpu.b.set(0x2b)
			cpu.step()

			ExpectB(cpu, 0x2c)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implement INCA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4c) // INCA
			cpu.a.set(0xff)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Inherent] should implement INCA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4c) // INCA
			cpu.a.set(0xfb)
			cpu.step()

			ExpectA(cpu, 0xfc)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Inherent] should implement INCA with bit NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4c) // INCA
			cpu.a.set(0x7f)
			cpu.step()

			ExpectA(cpu, 0x80)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NV", "Z")
		})
	})

	Context("[OR]", func() {

		It("[Immediate] should implement ORA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8a) // ORA
			cpu.write(0x1001, 0x55)
			cpu.a.set(0x31)
			cpu.step()

			ExpectA(cpu, 0x75)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Immediate] should implement ORA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8a) // ORA
			cpu.write(0x1001, 0x00)
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})

		It("[Immediate] should implement ORA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8a) // ORA
			cpu.write(0x1001, 0x84)
			cpu.a.set(0x55)
			cpu.step()

			ExpectA(cpu, 0xd5)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Immediate] should implement ORB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xca) // ORB
			cpu.write(0x1001, 0x42)
			cpu.b.set(0x79)
			cpu.step()

			ExpectB(cpu, 0x7b)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})
	})

	Context("[ADD]", func() {

		It("[Immediate] should implement ADDA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0x04)
			cpu.a.set(0x50)
			cpu.step()

			ExpectA(cpu, 0x54)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVHC")
		})

		It("[Immediate] should implement ADDA with bit H", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0x0f)
			cpu.a.set(0x01)
			cpu.step()

			ExpectA(cpu, 0x10)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "H", "NZVC")
		})

		It("[Immediate] should implement ADDA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0x01)
			cpu.a.set(0xf0)
			cpu.step()

			ExpectA(cpu, 0xf1)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZVHC")
		})

		It("[Immediate] should implement ADDA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0x00)
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NVHC")
		})

		It("[Immediate] should implement ADDA with bit NHC", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0xff)
			cpu.a.set(0xff)
			cpu.step()

			ExpectA(cpu, 0xfe)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NCH", "ZV")
		})

		It("[Immediate] should implement ADDA with bit NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x8b) // ADDA
			cpu.write(0x1001, 0x40)
			cpu.a.set(0x42)
			cpu.step()

			ExpectA(cpu, 0x82)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "NV", "ZHC")
		})

		It("[Immediate] should implement ADDB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xcb) // ADDB
			cpu.write(0x1001, 0x14)
			cpu.b.set(0x23)
			cpu.step()

			ExpectB(cpu, 0x37)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZVHC")
		})

		It("[Immediate] should implement ADDD", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc3) // ADDD
			cpu.writew(0x1001, 0x2c80)
			cpu.a.set(0x50)
			cpu.b.set(0x50)
			cpu.step()

			ExpectD(cpu, 0x7cd0)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "", "NZVC")
		})

		It("[Immediate] should implement ADDD with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc3) // ADDD
			cpu.writew(0x1001, 0x0004)
			cpu.a.set(0x8f)
			cpu.b.set(0x76)
			cpu.step()

			ExpectD(cpu, 0x8f7a)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "N", "ZVC")
		})

		It("[Immediate] should implement ADDD with bit NV", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc3) // ADDD
			cpu.writew(0x1001, 0x0003)
			cpu.a.set(0x7f)
			cpu.b.set(0xfe)
			cpu.step()

			ExpectD(cpu, 0x8001)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "NV", "ZC")
		})

		It("[Immediate] should implement ADDD with bit C", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0xc3) // ADDD
			cpu.writew(0x1001, 0x0003)
			cpu.a.set(0xff)
			cpu.b.set(0xfe)
			cpu.step()

			ExpectD(cpu, 0x0001)
			ExpectPC(cpu, 0x1003)
			ExpectClock(cpu, 4)
			ExpectCCR(cpu, "C", "NZV")
		})
	})

	Context("[TST]", func() {

		It("[Direct] should implement TST with Direct addressing mode", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.write(0x1000, 0x0d) // TST Direct
			cpu.write(0x1001, 0x0a)
			cpu.write(0x200a, 0x32)
			cpu.step()
			ExpectMemory(cpu, 0x200a, 0x32)
			ExpectPC(cpu, 0x1002)
			ExpectClock(cpu, 6)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implment TSTA", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4d) // TSTA
			cpu.a.set(0x32)
			cpu.step()

			ExpectA(cpu, 0x32)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implment TSTB", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x5d) // TSTA
			cpu.b.set(0x32)
			cpu.step()

			ExpectB(cpu, 0x32)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "", "NZV")
		})

		It("[Inherent] should implment TSTA with bit N", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4d) // TSTA
			cpu.a.set(0xd8)
			cpu.step()

			ExpectA(cpu, 0xd8)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "N", "ZV")
		})

		It("[Inherent] should implment TSTA with bit Z", func() {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, 0x4d) // TSTA
			cpu.a.set(0x00)
			cpu.step()

			ExpectA(cpu, 0x00)
			ExpectPC(cpu, 0x1001)
			ExpectClock(cpu, 2)
			ExpectCCR(cpu, "Z", "NV")
		})
	})

	Context("[BSR]", func() {

		It("[Long Relative] should implement LBSR", func() {
			cpu.pc.set(0x1000)
			cpu.s.set(0x400)
			cpu.write(0x1000, 0x17) // LBSR
			cpu.writew(0x1001, 0x0180)
			cpu.step()

			ExpectPC(cpu, 0x1183)
			ExpectS(cpu, 0x3fe)
			ExpectWord(cpu, cpu.s.uint16(), 0x1003)
			ExpectClock(cpu, 9)
		})

		It("[Relative] should implement BSR", func() {
			cpu.pc.set(0x1000)
			cpu.s.set(0x400)
			cpu.write(0x1000, 0x8d) // BSR
			cpu.write(0x1001, 0x46)
			cpu.step()

			ExpectPC(cpu, 0x1048)
			ExpectS(cpu, 0x3fe)
			ExpectWord(cpu, cpu.s.uint16(), 0x1002)
			ExpectClock(cpu, 7)
		})
	})

	Context("[JSR]", func() {

		It("[Direct] should implement JSR", func() {
			cpu.pc.set(0x1000)
			cpu.dp.set(0x20)
			cpu.s.set(0x400)
			cpu.write(0x1000, 0x9d) // JSR
			cpu.write(0x1001, 0x60)
			cpu.step()

			ExpectPC(cpu, 0x2060)
			ExpectS(cpu, 0x3fe)
			ExpectWord(cpu, cpu.s.uint16(), 0x1002)
			ExpectClock(cpu, 7)
		})
	})

	Context("[BR]", func() {
		ccflags := []struct {
			flag  string
			set   func()
			clear func()
		}{
			{"H", cpu.cc.setH, cpu.cc.clearH},
			{"C", cpu.cc.setC, cpu.cc.clearC},
			{"N", cpu.cc.setN, cpu.cc.clearN},
			{"Z", cpu.cc.setZ, cpu.cc.clearZ},
			{"V", cpu.cc.setV, cpu.cc.clearV},
		}

		branchingOpcodeTest8 := func(opcode uint8, flags string, branch bool, cycles int) {
			cpu.pc.set(0x1000)
			cpu.write(0x1000, opcode)
			cpu.write(0x1001, 0x10)
			for _, ccf := range ccflags {
				if strings.Contains(flags, ccf.flag) {
					ccf.set()
				} else {
					ccf.clear()
				}
			}
			cpu.step()
			offset := 0
			if branch {
				offset = 0x10
			}
			ExpectPC(cpu, 0x1002+offset)
			ExpectClock(cpu, cycles)
		}

		branchingOpcodeTest16 := func(opcode uint16, flags string, branch bool, cycles int) {
			cpu.pc.set(0x1000)
			cpu.writew(0x1000, opcode)
			cpu.writew(0x1002, 0x1000)
			for _, ccf := range ccflags {
				if strings.Contains(flags, ccf.flag) {
					ccf.set()
				} else {
					ccf.clear()
				}
			}
			cpu.step()
			offset := 0
			if branch {
				offset = 0x1000
			}
			ExpectPC(cpu, 0x1004+offset)
			ExpectClock(cpu, cycles)
		}

		It("[Relative] should implement BHI C=0 Z=0 (o)", func() {
			branchingOpcodeTest8(0x22, "", true, 3)
		})

		It("[Relative] should implement BHI C=1 Z=0 (x)", func() {
			branchingOpcodeTest8(0x22, "C", false, 3)
		})

		It("[Relative] should implement BHI C=0 Z=1 (x)", func() {
			branchingOpcodeTest8(0x22, "Z", false, 3)
		})

		It("[Relative] should implement BHI C=1 Z=1 (x)", func() {
			branchingOpcodeTest8(0x22, "CZ", false, 3)
		})

		It("[Long Relative] should implement LBHI C=0 Z=0 (o)", func() {
			branchingOpcodeTest16(0x1022, "", true, 6)
		})

		It("[Long Relative] should implement LBHI C=1 Z=0 (x)", func() {
			branchingOpcodeTest16(0x1022, "C", false, 5)
		})

		It("[Long Relative] should implement LBHI C=0 Z=1 (x)", func() {
			branchingOpcodeTest16(0x1022, "Z", false, 5)
		})

		It("[Long Relative] should implement LBHI C=1 Z=1 (x)", func() {
			branchingOpcodeTest16(0x1022, "CZ", false, 5)
		})

		It("[Relative] should implement BLS C=0 Z=0 (o)", func() {
			branchingOpcodeTest8(0x23, "", false, 3)
		})

		It("[Relative] should implement BLS C=1 Z=0 (x)", func() {
			branchingOpcodeTest8(0x23, "C", true, 3)
		})

		It("[Relative] should implement BLS C=0 Z=1 (x)", func() {
			branchingOpcodeTest8(0x23, "Z", true, 3)
		})

		It("[Relative] should implement BLS C=1 Z=1 (x)", func() {
			branchingOpcodeTest8(0x23, "CZ", true, 3)
		})

		It("[Long Relative] should implement LBLS C=0 Z=0 (o)", func() {
			branchingOpcodeTest16(0x1023, "", false, 5)
		})

		It("[Long Relative] should implement LBLS C=1 Z=0 (x)", func() {
			branchingOpcodeTest16(0x1023, "C", true, 6)
		})

		It("[Long Relative] should implement LBLS C=0 Z=1 (x)", func() {
			branchingOpcodeTest16(0x1023, "Z", true, 6)
		})

		It("[Long Relative] should implement LBLS C=1 Z=1 (x)", func() {
			branchingOpcodeTest16(0x1023, "CZ", true, 6)
		})

		It("[Relative] should implement BCC C=0 (o)", func() {
			branchingOpcodeTest8(0x24, "", true, 3)
		})

		It("[Relative] should implement BCC C=1 (x)", func() {
			branchingOpcodeTest8(0x24, "C", false, 3)
		})

		It("[Relative] should implement LBCC C=0 (o)", func() {
			branchingOpcodeTest16(0x1024, "", true, 6)
		})

		It("[Relative] should implement LBCC C=1 (x)", func() {
			branchingOpcodeTest16(0x1024, "C", false, 5)
		})

		It("[Relative] should implement BLO C=0 (x)", func() {
			branchingOpcodeTest8(0x25, "", false, 3)
		})

		It("[Relative] should implement BLO C=1 (o)", func() {
			branchingOpcodeTest8(0x25, "C", true, 3)
		})

		It("[Long Relative] should implement LBLO C=0 (x)", func() {
			branchingOpcodeTest16(0x1025, "", false, 5)
		})

		It("[Long Relative] should implement LBLO C=1 (o)", func() {
			branchingOpcodeTest16(0x1025, "C", true, 6)
		})

		It("[Relative] should implement BNE Z=0 (o)", func() {
			branchingOpcodeTest8(0x26, "", true, 3)
		})

		It("[Relative] should implement BNE Z=1 (x)", func() {
			branchingOpcodeTest8(0x26, "Z", false, 3)
		})

		It("[Long Relative] should implement LBNE Z=0 (o)", func() {
			branchingOpcodeTest16(0x1026, "", true, 6)
		})

		It("[Long Relative] should implement LBNE Z=1 (x)", func() {
			branchingOpcodeTest16(0x1026, "Z", false, 5)
		})

		It("[Relative] should implement BEQ Z=0 (x)", func() {
			branchingOpcodeTest8(0x27, "", false, 3)
		})

		It("[Relative] should implement BEQ Z=1 (o)", func() {
			branchingOpcodeTest8(0x27, "Z", true, 3)
		})

		It("[Long Relative] should implement LBEQ Z=0 (x)", func() {
			branchingOpcodeTest16(0x1027, "", false, 5)
		})

		It("[Long Relative] should implement LBEQ Z=1 (o)", func() {
			branchingOpcodeTest16(0x1027, "Z", true, 6)
		})

		It("[Relative] should implement BVC V=0 (o)", func() {
			branchingOpcodeTest8(0x28, "", true, 3)
		})

		It("[Relative] should implement BVC V=1 (x)", func() {
			branchingOpcodeTest8(0x28, "V", false, 3)
		})

		It("[Long Relative] should implement LBVC V=0 (o)", func() {
			branchingOpcodeTest16(0x1028, "", true, 6)
		})

		It("[Long Relative] should implement LBVC V=1 (x)", func() {
			branchingOpcodeTest16(0x1028, "V", false, 5)
		})

		It("[Relative] should implement BVS V=0 (x)", func() {
			branchingOpcodeTest8(0x29, "", false, 3)
		})

		It("[Relative] should implement BVS V=1 (o)", func() {
			branchingOpcodeTest8(0x29, "V", true, 3)
		})

		It("[Long Relative] should implement LBVS V=0 (x)", func() {
			branchingOpcodeTest16(0x1029, "", false, 5)
		})

		It("[Long Relative] should implement LBVS V=1 (o)", func() {
			branchingOpcodeTest16(0x1029, "V", true, 6)
		})

		It("[Relative] should implement BPL N=0 (o)", func() {
			branchingOpcodeTest8(0x2a, "", true, 3)
		})

		It("[Relative] should implement BPL N=1 (x)", func() {
			branchingOpcodeTest8(0x2a, "N", false, 3)
		})

		It("[Long Relative] should implement LBPL N=0 (o)", func() {
			branchingOpcodeTest16(0x102a, "", true, 6)
		})

		It("[Long Relative] should implement LBPL N=1 (x)", func() {
			branchingOpcodeTest16(0x102a, "N", false, 5)
		})

		It("[Relative] should implement BMI N=0 (x)", func() {
			branchingOpcodeTest8(0x2b, "", false, 3)
		})

		It("[Relative] should implement BMI N=1 (o)", func() {
			branchingOpcodeTest8(0x2b, "N", true, 3)
		})

		It("[Long Relative] should implement LBMI N=0 (x)", func() {
			branchingOpcodeTest16(0x102b, "", false, 5)
		})

		It("[Long Relative] should implement LBMI N=1 (o)", func() {
			branchingOpcodeTest16(0x102b, "N", true, 6)
		})

		It("[Relative] should implement BGE N=0 V=0 (o)", func() {
			branchingOpcodeTest8(0x2c, "", true, 3)
		})

		It("[Relative] should implement BGE N=0 V=1 (x)", func() {
			branchingOpcodeTest8(0x2c, "V", false, 3)
		})

		It("[Relative] should implement BGE N=1 V=0 (x)", func() {
			branchingOpcodeTest8(0x2c, "N", false, 3)
		})

		It("[Relative] should implement BGE N=1 V=1 (o)", func() {
			branchingOpcodeTest8(0x2c, "NV", true, 3)
		})

		It("[Long Relative] should implement LBGE N=0 V=0 (o)", func() {
			branchingOpcodeTest16(0x102c, "", true, 6)
		})

		It("[Long Relative] should implement LBGE N=0 V=1 (x)", func() {
			branchingOpcodeTest16(0x102c, "V", false, 5)
		})

		It("[Long Relative] should implement LBGE N=1 V=0 (x)", func() {
			branchingOpcodeTest16(0x102c, "N", false, 5)
		})

		It("[Long Relative] should implement LBGE N=1 V=1 (o)", func() {
			branchingOpcodeTest16(0x102c, "NV", true, 6)
		})

		It("[Relative] should implement BLT N=0 V=0 (x)", func() {
			branchingOpcodeTest8(0x2d, "", false, 3)
		})

		It("[Relative] should implement BLT N=0 V=1 (o)", func() {
			branchingOpcodeTest8(0x2d, "V", true, 3)
		})

		It("[Relative] should implement BLT N=1 V=0 (o)", func() {
			branchingOpcodeTest8(0x2d, "N", true, 3)
		})

		It("[Relative] should implement BLT N=1 V=1 (x)", func() {
			branchingOpcodeTest8(0x2d, "NV", false, 3)
		})

		It("[Long Relative] should implement LBLT N=0 V=0 (x)", func() {
			branchingOpcodeTest16(0x102d, "", false, 5)
		})

		It("[Long Relative] should implement LBLT N=0 V=1 (o)", func() {
			branchingOpcodeTest16(0x102d, "V", true, 6)
		})

		It("[Long Relative] should implement LBLT N=1 V=0 (o)", func() {
			branchingOpcodeTest16(0x102d, "N", true, 6)
		})

		It("[Long Relative] should implement LBLT N=1 V=1 (x)", func() {
			branchingOpcodeTest16(0x102d, "NV", false, 5)
		})

		It("[Relative] should implement BGT Z=0 N=0 V=0 (o)", func() {
			branchingOpcodeTest8(0x2e, "", true, 3)
		})

		It("[Relative] should implement BGT Z=0 N=0 V=1 (x)", func() {
			branchingOpcodeTest8(0x2e, "V", false, 3)
		})

		It("[Relative] should implement BGT Z=0 N=1 V=0 (x)", func() {
			branchingOpcodeTest8(0x2e, "N", false, 3)
		})

		It("[Relative] should implement BGT Z=0 N=1 V=1 (o)", func() {
			branchingOpcodeTest8(0x2e, "NV", true, 3)
		})

		It("[Relative] should implement BGT Z=1 N=0 V=0 (x)", func() {
			branchingOpcodeTest8(0x2e, "Z", false, 3)
		})

		It("[Relative] should implement BGT Z=1 N=0 V=1 (x)", func() {
			branchingOpcodeTest8(0x2e, "ZV", false, 3)
		})

		It("[Relative] should implement BGT Z=1 N=1 V=0 (x)", func() {
			branchingOpcodeTest8(0x2e, "ZN", false, 3)
		})

		It("[Relative] should implement BGT Z=1 N=1 V=1 (x)", func() {
			branchingOpcodeTest8(0x2e, "ZNV", false, 3)
		})

		It("[Long Relative] should implement LBGT Z=0 N=0 V=0 (o)", func() {
			branchingOpcodeTest16(0x102e, "", true, 6)
		})

		It("[Long Relative] should implement LBGT Z=0 N=0 V=1 (x)", func() {
			branchingOpcodeTest16(0x102e, "V", false, 5)
		})

		It("[Long Relative] should implement LBGT Z=0 N=1 V=0 (x)", func() {
			branchingOpcodeTest16(0x102e, "N", false, 5)
		})

		It("[Long Relative] should implement LBGT Z=0 N=1 V=1 (o)", func() {
			branchingOpcodeTest16(0x102e, "NV", true, 6)
		})

		It("[Long Relative] should implement LBGT Z=1 N=0 V=0 (x)", func() {
			branchingOpcodeTest16(0x102e, "Z", false, 5)
		})

		It("[Long Relative] should implement LBGT Z=1 N=0 V=1 (x)", func() {
			branchingOpcodeTest16(0x102e, "ZV", false, 5)
		})

		It("[Long Relative] should implement LBGT Z=1 N=1 V=0 (x)", func() {
			branchingOpcodeTest16(0x102e, "ZN", false, 5)
		})

		It("[Long Relative] should implement LBGT Z=1 N=1 V=1 (x)", func() {
			branchingOpcodeTest16(0x102e, "ZNV", false, 5)
		})

		It("[Relative] should implement BLE Z=0 N=0 V=0 (x)", func() {
			branchingOpcodeTest8(0x2f, "", false, 3)
		})

		It("[Relative] should implement BLE Z=0 N=0 V=1 (o)", func() {
			branchingOpcodeTest8(0x2f, "V", true, 3)
		})

		It("[Relative] should implement BLE Z=0 N=1 V=0 (o)", func() {
			branchingOpcodeTest8(0x2f, "N", true, 3)
		})

		It("[Relative] should implement BLE Z=0 N=1 V=1 (x)", func() {
			branchingOpcodeTest8(0x2f, "NV", false, 3)
		})

		It("[Relative] should implement BLE Z=1 N=0 V=0 (o)", func() {
			branchingOpcodeTest8(0x2f, "Z", true, 3)
		})

		It("[Relative] should implement BLE Z=1 N=0 V=1 (o)", func() {
			branchingOpcodeTest8(0x2f, "ZV", true, 3)
		})

		It("[Relative] should implement BLE Z=1 N=1 V=0 (o)", func() {
			branchingOpcodeTest8(0x2f, "ZN", true, 3)
		})

		It("[Relative] should implement BLE Z=1 N=1 V=1 (o)", func() {
			branchingOpcodeTest8(0x2f, "ZNV", true, 3)
		})

		It("[Long Relative] should implement LBLE Z=0 N=0 V=0 (x)", func() {
			branchingOpcodeTest16(0x102f, "", false, 5)
		})

		It("[Long Relative] should implement LBLE Z=0 N=0 V=1 (o)", func() {
			branchingOpcodeTest16(0x102f, "V", true, 6)
		})

		It("[Long Relative] should implement LBLE Z=0 N=1 V=0 (o)", func() {
			branchingOpcodeTest16(0x102f, "N", true, 6)
		})

		It("[Long Relative] should implement LBLE Z=0 N=1 V=1 (x)", func() {
			branchingOpcodeTest16(0x102f, "NV", false, 5)
		})

		It("[Long Relative] should implement LBLE Z=1 N=0 V=0 (o)", func() {
			branchingOpcodeTest16(0x102f, "Z", true, 6)
		})

		It("[Long Relative] should implement LBLE Z=1 N=0 V=1 (o)", func() {
			branchingOpcodeTest16(0x102f, "ZV", true, 6)
		})

		It("[Long Relative] should implement LBLE Z=1 N=1 V=0 (o)", func() {
			branchingOpcodeTest16(0x102f, "ZN", true, 6)
		})

		It("[Long Relative] should implement LBLE Z=1 N=1 V=1 (o)", func() {
			branchingOpcodeTest16(0x102f, "ZNV", true, 6)
		})
	})
})

/*

func TestJMPDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.dp = 0x20
	cpu.write(0x1000, 0x0e) // JMP Direct
	cpu.write(0x1001, 0x0a)
	cpu.step()
	assert.That(cpu.pc).AsInt().IsEqualTo(0x200a)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestCLRDirect(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.dp = 0x20
	cpu.write(0x1000, 0x0f) // CLR Direct
	cpu.write(0x1001, 0x0a)
	cpu.write(0x200a, 0x4d)
	cpu.step()
	assert.That(cpu.read(0x200a)).AsInt().IsEqualTo(0x00)
	assert.ThatBool(cpu.getN()).IsFalse()
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getV()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestCLRA(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x4f) // CLRA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x5f) // CLRB
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x20) // BRA
	cpu.write(0x1001, 0x10)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1012)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestLBRAPositive(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x16) // LBRA
	cpu.writew(0x1001, 0x1000)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2003)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestBRANegative(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x20) // LBRA
	cpu.write(0x1001, 0xfe)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1000)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestLBRANegative(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x16) // LBRA
	cpu.writew(0x1001, 0xff00)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x0f03)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestDAA(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x19) // DAA
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.cc = carry | negative
	cpu.write(0x1000, 0x1a) // ORCC
	cpu.write(0x1001, 0x82)
	cpu.step()
	assert.That(cpu.cc).AsInt().IsEqualTo(0x8b)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestANDCC(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.cc = carry | negative
	cpu.write(0x1000, 0x1c) // ANDCC
	cpu.write(0x1001, 0xfe)
	cpu.step()
	assert.That(cpu.cc).AsInt().IsEqualTo(0x08)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestSEX(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1d) // SEX
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1d) // SEX
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1d) // SEX
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1e) // EXG
	cpu.write(0x1001, 0x89)
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1e) // EXG
	cpu.write(0x1001, 0x12)
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
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1e) // EXG
	cpu.write(0x1001, 0x85)
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Try to exchange 8-bit with 16-bits registers")
	}()
	cpu.step()
}

func TestEXGInvalidCode(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1e) // EXG
	cpu.write(0x1001, 0x67)
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Invalid register code")
	}()
	cpu.step()
}

func TestTFRAtoB(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1f) // TFR
	cpu.write(0x1001, 0x89)
	cpu.a = 0x27
	cpu.step()
	assert.That(cpu.a).AsInt().IsEqualTo(0x27)
	assert.That(cpu.b).AsInt().IsEqualTo(0x27)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestTFRPCtoX(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1f) // TFR
	cpu.write(0x1001, 0x51)
	cpu.step()
	assert.That(cpu.x).AsInt().IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestTFRInvalidRegisterCombination(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1f) // TFR
	cpu.write(0x1001, 0x85)
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Try to transfer 8-bit and 16-bits registers")
	}()
	cpu.step()
}

func TestTFRInvalidRegisterCode(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x1f) // TFR
	cpu.write(0x1001, 0x8f)
	defer func() {
		r := recover()
		assert.That(r).AsString().Contains("Invalid register code")
	}()
	cpu.step()
}

func TestLEAX(t *testing.T) {
	assert := assert.New(t)
	var cpu CPU
	ram := NewRam()
	cpu.Initialize(ram)
	cpu.pc = 0x1000
	cpu.y = 0xd000
	cpu.a = 0x5a
	cpu.write(0x1000, 0x30) // LEAX
	cpu.write(0x1001, 0xa6) // EA = Y + ACCA
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.That(cpu.x).AsInt().IsEqualTo(0xd05a)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestLEAXZero(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.y = 0x100
	cpu.a = 0xff
	cpu.b = 0x00 // D = 0xff00
	cpu.x = 0x0100
	cpu.write(0x1000, 0x30) // LEAX
	cpu.write(0x1001, 0xab) // EA = Y + ACCD
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.That(cpu.x).AsInt().IsEqualTo(0)
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(8)
}

func TestLEAY(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, 0x31) // LEAY
	cpu.write(0x1001, 0x8c) // EA = PC + 8 bits offset
	cpu.write(0x1002, 0x0a)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1003)
	assert.That(cpu.y).AsInt().IsEqualTo(0x100d)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestLEAYZero(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.y = 0x0100
	cpu.u = 0
	cpu.write(0x1000, 0x31) // LEAY
	cpu.write(0x1001, 0xc4) // EA = U
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.That(cpu.y).AsInt().IsEqualTo(0)
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(4)
}

func TestLEAS(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.x = 0x800
	cpu.write(0x1000, 0x32) // LEAS
	cpu.write(0x1001, 0x94) // EA = [X]
	cpu.writew(0x800, 0x1f40)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.That(cpu.s).AsInt().IsEqualTo(0x1f40)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(7)
}

func TestLEAU(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.y = 0x2000
	cpu.write(0x1000, 0x33) // LEAU
	cpu.write(0x1001, 0xb1) // EA = [Y++]
	cpu.writew(0x2000, 0x1f40)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.That(cpu.u).AsInt().IsEqualTo(0x1f40)
	assert.That(cpu.y).AsInt().IsEqualTo(0x2002)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(10)
}

func TestPSHS(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.u = 0x7fff
	cpu.y = 0xa200
	cpu.x = 0xa100
	cpu.dp = 0x04
	cpu.b = 0x4f
	cpu.a = 0x05
	cpu.cc = 0x03
	cpu.s = 0xd000
	cpu.write(0x1000, 0x34) // PSHS
	cpu.write(0x1001, 0xff) // Push All
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.readw(0xcffe))).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.readw(0xcffc))).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.readw(0xcffa))).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.readw(0xcff8))).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.read(0xcff7))).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.read(0xcff6))).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0xcff4)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(17)
}

func TestPULS(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.s = 0xcff4
	cpu.write(0x1000, 0x35) // PULS
	cpu.write(0x1001, 0xff) // All registers
	cpu.writew(0xcffe, 0x2000)
	cpu.writew(0xcffc, 0x7fff)
	cpu.writew(0xcffa, 0xa200)
	cpu.writew(0xcff8, 0xa100)
	cpu.write(0xcff7, 0x04)
	cpu.write(0xcff6, 0x4f)
	cpu.write(0xcff5, 0x05)
	cpu.write(0xcff4, 0x03)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2000)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0xd000)
	assert.ThatInt(int(cpu.cc)).IsEqualTo(0x03)
	assert.ThatInt(int(cpu.a)).IsEqualTo(0x05)
	assert.ThatInt(int(cpu.b)).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.dp)).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.x)).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.y)).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.u)).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(17)
}

func TestPSHU(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.s = 0x7fff
	cpu.y = 0xa200
	cpu.x = 0xa100
	cpu.dp = 0x04
	cpu.b = 0x4f
	cpu.a = 0x05
	cpu.cc = 0x03
	cpu.u = 0xd000
	cpu.write(0x1000, 0x36) // PSHU
	cpu.write(0x1001, 0xff) // Push All
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.readw(0xcffe))).IsEqualTo(0x1002)
	assert.ThatInt(int(cpu.readw(0xcffc))).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.readw(0xcffa))).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.readw(0xcff8))).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.read(0xcff7))).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.read(0xcff6))).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.u)).IsEqualTo(0xcff4)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(17)
}

func TestPULU(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.u = 0xcff4
	cpu.write(0x1000, 0x37) // PULU
	cpu.write(0x1001, 0xff) // All registers
	cpu.writew(0xcffe, 0x2000)
	cpu.writew(0xcffc, 0x7fff)
	cpu.writew(0xcffa, 0xa200)
	cpu.writew(0xcff8, 0xa100)
	cpu.write(0xcff7, 0x04)
	cpu.write(0xcff6, 0x4f)
	cpu.write(0xcff5, 0x05)
	cpu.write(0xcff4, 0x03)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2000)
	assert.ThatInt(int(cpu.u)).IsEqualTo(0xd000)
	assert.ThatInt(int(cpu.cc)).IsEqualTo(0x03)
	assert.ThatInt(int(cpu.a)).IsEqualTo(0x05)
	assert.ThatInt(int(cpu.b)).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.dp)).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.x)).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.y)).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(17)
}

func TestRTS(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.s = 0xcffe
	cpu.write(0x1000, 0x39) // RTS
	cpu.writew(cpu.s, 0x3000)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x3000)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0xd000)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(5)
}

func TestABX(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.x = 0x2f00
	cpu.b = 0x50
	cpu.write(0x1000, 0x3a) // ABX
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.x)).IsEqualTo(0x2f50)
	assert.ThatInt(int(cpu.b)).IsEqualTo(0x50)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(3)
}

func TestRTIPartial(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.s = 0xbffd
	cpu.write(0x1000, 0x3b) // RTI
	cpu.write(0xbffd, 0x05)
	cpu.writew(0xbffe, 0x2000)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2000)
	assert.ThatInt(int(cpu.cc)).IsEqualTo(0x05)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0xc000)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(6)
}

func TestRTIEntire(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.s = 0xbff4
	cpu.write(0x1000, 0x3b) // RTI
	cpu.writew(0xbffe, 0x2000)
	cpu.writew(0xbffc, 0x7fff)
	cpu.writew(0xbffa, 0xa200)
	cpu.writew(0xbff8, 0xa100)
	cpu.write(0xbff7, 0x04)
	cpu.write(0xbff6, 0x4f)
	cpu.write(0xbff5, 0x05)
	cpu.write(0xbff4, 0x83)
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2000)
	assert.ThatInt(int(cpu.cc)).IsEqualTo(0x83)
	assert.ThatInt(int(cpu.a)).IsEqualTo(0x05)
	assert.ThatInt(int(cpu.b)).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.dp)).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.x)).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.y)).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.u)).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x2000)
	assert.ThatInt(int(cpu.s)).IsEqualTo(0xc000)
	assert.ThatInt(int(cpu.clock)).IsEqualTo(15)
}

func TestMUL(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.a = 0x25
	cpu.b = 0xd0
	cpu.write(0x1000, 0x3d) // MUL
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.d())).IsEqualTo(0x1e10)
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(11)
}

func TestMULCarry(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.a = 0x65
	cpu.b = 0xdf
	cpu.write(0x1000, 0x3d) // MUL
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.d())).IsEqualTo(0x57fb)
	assert.ThatBool(cpu.getZ()).IsFalse()
	assert.ThatBool(cpu.getC()).IsTrue()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(11)
}

func TestMULZero(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.a = 0x65
	cpu.b = 0x00
	cpu.write(0x1000, 0x3d) // MUL
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.d())).IsEqualTo(0)
	assert.ThatBool(cpu.getZ()).IsTrue()
	assert.ThatBool(cpu.getC()).IsFalse()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(11)
}

func TestSWI(t *testing.T) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.u = 0x7fff
	cpu.y = 0xa200
	cpu.x = 0xa100
	cpu.dp = 0x04
	cpu.b = 0x4f
	cpu.a = 0x05
	cpu.cc = 0x03
	cpu.s = 0xd000
	cpu.write(0x1000, 0x3f)    // SWI
	cpu.writew(0xfffa, 0xe000) // Interupt Vector Address
	cpu.step()
	assert.ThatInt(int(cpu.pc)).IsEqualTo(0xe000)
	assert.ThatInt(int(cpu.readw(0xcffe))).IsEqualTo(0x1001)
	assert.ThatInt(int(cpu.readw(0xcffc))).IsEqualTo(0x7fff)
	assert.ThatInt(int(cpu.readw(0xcffa))).IsEqualTo(0xa200)
	assert.ThatInt(int(cpu.readw(0xcff8))).IsEqualTo(0xa100)
	assert.ThatInt(int(cpu.read(0xcff7))).IsEqualTo(0x04)
	assert.ThatInt(int(cpu.read(0xcff6))).IsEqualTo(0x4f)
	assert.ThatInt(int(cpu.read(0xcff5))).IsEqualTo(0x05)
	assert.ThatInt(int(cpu.read(0xcff4))).IsEqualTo(0x83)
	assert.ThatBool(cpu.getE()).IsTrue()
	assert.ThatBool(cpu.getI()).IsTrue()
	assert.ThatBool(cpu.getF()).IsTrue()
	assert.ThatInt(int(cpu.clock)).IsEqualTo(19)
}
*/

// ExpectMemory asserts a value at a memory address
func ExpectMemory(cpu CPU, address uint16, expected interface{}) {
	ExpectWithOffset(1, cpu.read(address)).To(BeEquivalentTo(expected), "Expected value at address 0x%04x to be 0x%02x but is 0x%02x", address, expected, cpu.read(address))
}

// ExpectWord asserts a value at a memory address
func ExpectWord(cpu CPU, address uint16, expected interface{}) {
	ExpectWithOffset(1, cpu.readw(address)).To(BeEquivalentTo(expected), "Expected value at address 0x%04x to be 0x%04x but is 0x%04x", address, expected, cpu.readw(address))
}

// ExpectA asserts a value in A registry
func ExpectA(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.a.get()).To(BeEquivalentTo(expected), "Expected A register to be 0x%x but is 0x%x", expected, cpu.a.get())
}

// ExpectB asserts a value in B registry
func ExpectB(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.b.get()).To(BeEquivalentTo(expected), "Expected B register to be 0x%x but is 0x%x", expected, cpu.b.get())
}

// ExpectD asserts a value in D registry
func ExpectD(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.d()).To(BeEquivalentTo(expected), "Expected D register to be 0x%04x but is 0x%04x", expected, cpu.d())
}

// ExpectX asserts a value in X registry
func ExpectX(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.x.get()).To(BeEquivalentTo(expected), "Expected X register to be 0x%04x but is 0x%04x", expected, cpu.x.get())
}

// ExpectY asserts a value in Y registry
func ExpectY(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.y.get()).To(BeEquivalentTo(expected), "Expected Y register to be 0x%04x but is 0x%04x", expected, cpu.y.get())
}

// ExpectS asserts a value in S registry
func ExpectS(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.s.get()).To(BeEquivalentTo(expected), "Expected S register to be 0x%04x but is 0x%04x", expected, cpu.s.get())
}

// ExpectU asserts a value in U registry
func ExpectU(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.u.get()).To(BeEquivalentTo(expected), "Expected U register to be 0x%04x but is 0x%04x", expected, cpu.u.get())
}

// ExpectPC asserts a value in PC registry
func ExpectPC(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.pc.get()).To(BeEquivalentTo(expected), "Expected PC register to be 0x%x but is 0x%x", expected, cpu.pc.get())
}

// ExpectClock asserts a value in PC registry
func ExpectClock(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.clock).To(BeEquivalentTo(expected), "Expected clock to be %d but is %d", expected, cpu.clock)
}

// ExpectCCR asserts set and clear flags of the CCR registry
func ExpectCCR(cpu CPU, set string, clear string) {
	type ccr struct {
		flag string
		call func() bool
	}
	flags := []ccr{
		ccr{"E", cpu.cc.getE},
		ccr{"F", cpu.cc.getF},
		ccr{"H", cpu.cc.getH},
		ccr{"I", cpu.cc.getI},
		ccr{"N", cpu.cc.getN},
		ccr{"Z", cpu.cc.getZ},
		ccr{"V", cpu.cc.getV},
		ccr{"C", cpu.cc.getC},
	}
	for _, f := range flags {
		if strings.Contains(set, f.flag) {
			ExpectWithOffset(1, f.call()).To(BeTrue(), "Expected CCR[%s] to be set but is clear", f.flag)
		}
		if strings.Contains(clear, f.flag) {
			ExpectWithOffset(1, f.call()).To(BeFalse(), "Expected CCR[%s] to be clear but is set", f.flag)
		}
	}
}
