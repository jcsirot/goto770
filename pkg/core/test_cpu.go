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

	Context("[SUB]", func() {

		It("[Immediate] should implement SUBA with Immediate addessing mode", func() {
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

		It("[Immediate] should implement SUBA with Immediate addessing mode with bit Z", func() {
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

func branchingOpcodeTest(t *testing.T, opcode uint8, flags []uint8, branch bool) {
	assert := assert.New(t)
	var cpu = newCPU()
	cpu.pc = 0x1000
	cpu.write(0x1000, opcode)
	cpu.write(0x1001, 0x10)
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

func TestBGTZ0N0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{}, true)
}

func TestBGTZ0N0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{overflow}, false)
}

func TestBGTZ0N1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{negative}, false)
}

func TestBGTZ0N1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{overflow, negative}, true)
}

func TestBGTZ1N0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{zero}, false)
}

func TestBGTZ1N0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{zero, overflow}, false)
}

func TestBGTZ1N1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{zero, negative}, false)
}

func TestBGTZ1N1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2e, []uint8{zero, overflow, negative}, false)
}

func TestBLEZ0N0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{}, false)
}

func TestBLEZ0N0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{overflow}, true)
}

func TestBLEZ0N1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{negative}, true)
}

func TestBLEZ0N1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{overflow, negative}, false)
}

func TestBLEZ1N0V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{zero}, true)
}

func TestBLEZ1N0V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{zero, overflow}, true)
}

func TestBLEZ1N1V0(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{zero, negative}, true)
}

func TestBLEZ1N1V1(t *testing.T) {
	branchingOpcodeTest(t, 0x2f, []uint8{zero, overflow, negative}, true)
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
	ExpectWithOffset(1, cpu.read(address)).To(BeEquivalentTo(expected))
}

func ExpectWord(cpu CPU, address uint16, expected interface{}) {
	ExpectWithOffset(1, cpu.readw(address)).To(BeEquivalentTo(expected))
}

// ExpectA asserts a value in A registry
func ExpectA(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.a.get()).To(BeEquivalentTo(expected), "Expected A register to be 0x%x but is 0x%x", expected, cpu.a.get())
}

// ExpectB asserts a value in B registry
func ExpectB(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.b.get()).To(BeEquivalentTo(expected))
}

// ExpectD asserts a value in D registry
func ExpectD(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.d()).To(BeEquivalentTo(expected), "Expected D register to be 0x%x but is 0x%x", expected, cpu.d())
}

// ExpectS asserts a value in S registry
func ExpectS(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.s.get()).To(BeEquivalentTo(expected), "Expected S register to be 0x%x but is 0x%x", expected, cpu.s.get())
}

// ExpectS asserts a value in S registry
func ExpectU(cpu CPU, expected interface{}) {
	ExpectWithOffset(1, cpu.u.get()).To(BeEquivalentTo(expected), "Expected U register to be 0x%x but is 0x%x", expected, cpu.u.get())
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