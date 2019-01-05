package core

var (
	Cpu CPU
	Ram Memory
)

func Start() {
	Ram = NewRam()
	Cpu.Initialize(Ram)
}
