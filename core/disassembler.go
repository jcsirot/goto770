package core

import (
	"fmt"
	"strings"
)

var registers = []string{"D", "X", "Y", "U", "S", "PC", "A", "B", "CC", "DP"}
var stackRegisters = [][]int{
	{5, 3, 2, 1, 9, 7, 6, 8},
	{5, 4, 2, 1, 9, 7, 6, 8},
}
var indexRegisters = []int{1, 2, 3, 4}

// Disassemble an instruction and return the string representation and the size of the instruction
func Disassemble(op opcode, instBuf []uint8) (string, int) {
	var sb strings.Builder

	sb.WriteString(op.name)
	var size int

	switch op.mode {
	case inherent:
		// else no EA to decode
		size = 1
	case direct:
		sb.WriteString(fmt.Sprintf(" <$%02x", instBuf[1]))
		size = 2
	case relative:
		sb.WriteString(fmt.Sprintf(" *+$%02x", instBuf[1]))
		size = 2
	case lrelative:
		sb.WriteString(fmt.Sprintf(" *+$%02x%02x", instBuf[1], instBuf[2]))
		size = 3
	case immediate:
		if op.name == "TFR" || op.name == "EXG" {
			sb.WriteString(fmt.Sprintf(" %s, %s", registers[(instBuf[1]>>4)&0xf], registers[instBuf[1]&0xf]))
		} else if instBuf[0]&0xfc == 0x34 { // PSHS PULS PSHU PULU
			regs := make([]string, 0)
			if instBuf[0]&0x01 == 0 { // Push
				for i := 0; i <= 7; i++ {
					if (0x80>>uint(i))&instBuf[1] != 0 {
						regs = append(regs, registers[stackRegisters[(instBuf[0]&2)>>1][i]])
					}
				}
			} else { // Pull
				for i := 7; i >= 0; i-- {
					if (0x80>>uint(i))&instBuf[1] != 0 {
						regs = append(regs, registers[stackRegisters[(instBuf[0]&2)>>1][i]])
					}
				}
			}
			sb.WriteString(fmt.Sprintf(" %s", strings.Join(regs, ",")))
		} else {
			sb.WriteString(fmt.Sprintf(" #$%02x", instBuf[1]))
		}
		size = 2
	case limmediate:
		sb.WriteString(fmt.Sprintf(" #$%02x%02x", instBuf[1], instBuf[2]))
		size = 3
	case extended:
		sb.WriteString(fmt.Sprintf(" $%02x%02x", instBuf[1], instBuf[2]))
		size = 3
	case indexed:
		postbyte := instBuf[0x01]
		if postbyte&0x80 == 0 {
			if postbyte&0x10 == 0 {
				offset := postbyte & 0x0f
				sb.WriteString(fmt.Sprintf(" %02x,%s", offset, registers[indexRegisters[(postbyte&0x60)>>5]]))
			} else {
				offset := ((postbyte & 0x0f) ^ 0x0f) + 1
				sb.WriteString(fmt.Sprintf(" -%02x,%s", offset, registers[indexRegisters[(postbyte&0x60)>>5]]))
			}
			size = 2
		} else {
			sb.WriteString(" ")
			if postbyte&0x10 == 0x10 { // Indirect mode
				sb.WriteString("(")
			}

			switch postbyte & 0x0f {
			case 0x00:
				sb.WriteString(fmt.Sprintf(",%s+", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x01:
				sb.WriteString(fmt.Sprintf(",%s++", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x02:
				sb.WriteString(fmt.Sprintf(",-%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x03:
				sb.WriteString(fmt.Sprintf(",--%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x04:
				sb.WriteString(fmt.Sprintf(",%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x05:
				sb.WriteString(fmt.Sprintf("B,%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x06:
				sb.WriteString(fmt.Sprintf("A,%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x08:
				offset := int8(instBuf[2])
				if offset < 0 {
					sb.WriteString("-")
					offset = -offset
				}
				sb.WriteString(fmt.Sprintf("%02x,%s", offset, registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 3

			case 0x09:
				offset := int16(instBuf[2])<<8 | int16(instBuf[3])
				if offset < 0 {
					sb.WriteString("-")
					offset = -offset
				}
				sb.WriteString(fmt.Sprintf("%04x,%s", offset, registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 4

			case 0x0b:
				sb.WriteString(fmt.Sprintf("D,%s", registers[indexRegisters[(postbyte&0x60)>>5]]))
				size = 2

			case 0x0c:
				offset := int8(instBuf[2])
				if offset < 0 {
					sb.WriteString("-")
					offset = -offset
				}
				sb.WriteString(fmt.Sprintf("%02x,PC", offset))
				size = 3

			case 0x0d:
				offset := int16(instBuf[2])<<8 | int16(instBuf[3])
				if offset < 0 {
					sb.WriteString("-")
					offset = -offset
				}
				sb.WriteString(fmt.Sprintf("%04x,PC", offset))
				size = 4
			}

			if postbyte&0x10 == 0x10 { // Indirect mode
				sb.WriteString(")")
			}
		}
	default:
		sb.WriteString(" ??? (NYE)")
	}

	return sb.String(), size
}

func format(pc uint16, instruction string, binary []uint8) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%04x", pc))
	sb.WriteString(" | ")
	sb.WriteString(instruction)
	sb.WriteString(" (")
	hexa := []string{}
	for _, x := range binary {
		hexa = append(hexa, fmt.Sprintf("%02x", x))
	}

	fmt.Printf("%04x | %s | (%s)\n", pc, padRight(instruction, " ", 12), strings.Join(hexa, " "))
}

func padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
