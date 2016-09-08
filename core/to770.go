package core

import (
	"fmt"
)

var (
	cpu Cpu
	ram Ram
)

func Start() {
	ram = NewRam()
	cpu.Initialize()

	cpu.neg(10)
	fmt.Println("cpu = %s", cpu)
}
