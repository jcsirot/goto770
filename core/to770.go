package core

import (
	"fmt"
)

var (
	Cpu CPU
	Ram Memory
)

func Start() {
	Ram = NewRam()
	Cpu.Initialize(Ram)

	Cpu.neg(10)
	fmt.Println("cpu = %s", Cpu)
}
