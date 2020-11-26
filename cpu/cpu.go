package cpu

import (
	"errors"
	"fmt"
	"rvsim/bus"
	"rvsim/ram"
)

const debug bool = false

//Register holds all registers of the cpu
type register struct {
	//32 Bit registers
	regs [32]uint64
	//one Program Counter
	pc uint64
	//Add memory here for simplicity
	ram bus.Device
}

var cpu register
var prglength uint64

//Initialize the cpu
func Initialize(binary []uint8) {
	prglength = uint64(len(binary))
	//	var ram bus.Device
	cpu.ram = &ram.RAM{}

	//The stack pointer
	cpu.regs[2] = bus.MemoryBase + ram.MemorySize
	cpu.pc = bus.MemoryBase

	if debug {
		fmt.Println("CPU Store Program to Memory")
	}
	cpu.ram = &ram.RAM{}
	for i := 0; i <= len(binary)-1; i++ {
		if debug {
			fmt.Println("CPU Store index ", i)
		}
		err := cpu.ram.Store(uint64(i)+bus.MemoryBase, 8, uint64(binary[i]))
		if err != nil {
			fmt.Println("PANIC: could not store memory")
			break
		}
	}

}

//DumpRegisters dumps all registers x0-x31
func DumpRegisters() {
	name := [32]string{
		"zero", " ra ", " sp ", " gp ", " tp ", " t0 ", " t1 ", " t2 ",
		" s0 ", " s1 ", " a0 ", " a1 ", " a2 ", " a3 ", " a4 ", " a5 ",
		" a6 ", " a7 ", " s2 ", " s3 ", " s4 ", " s5 ", " s6 ", " s7 ",
		" s8 ", " s9 ", " s10", " s11", " t3 ", " t4 ", " t5 ", " t6 ",
	}
	for i := 0; i <= 31; i += 4 {
		for j := 0; j <= 3; j++ {
			fmt.Printf("%#.2x (%s) = %#x\t", i+j, name[i+j], cpu.regs[i+j])
		}
		fmt.Println()
	}
}

//GetPC returns the program counter
func GetPC() uint64 {
	return cpu.pc
}

//IncPC used for jump instructions
func IncPC() {
	if (cpu.pc - bus.MemoryBase) >= prglength {
		cpu.pc = 0
	} else {
		cpu.pc += 4
	}
}

//SetPC used for jump instructions
func SetPC(newPC uint64) {
	cpu.pc = newPC
}

//Fetch cycle
func Fetch() (uint64, error) {
	if debug {
		fmt.Println("CPU DEBUG Fetch cpu.pc", cpu.pc)
	}
	value, err := cpu.ram.Load(cpu.pc, 32)
	return value, err
}

//Execute executes an instruction defined by its memory address
func Execute(instruction uint64) error {
	//Simulte the zero register at x00
	cpu.regs[0] = 0

	opcode := instruction & 0x0000007f
	rd := uint((instruction & 0x00000f80) >> 7)
	rs1 := uint((instruction & 0x000f8000) >> 15)
	rs2 := uint((instruction & 0x01f00000) >> 20)

	switch opcode {
	case 0x03:
		//imm[11:0], inst[31:20]
		//imm := uint64((instruction & 0xfff00000) >> 20)

	//R-Type
	case 0x13:
		//addi
		//imm[11:0], inst[31:20]
		imm := uint64((instruction & 0xfff00000) >> 20)
		cpu.regs[rd] = cpu.regs[rs1] + imm

	//I-Type
	case 0x33:
		//add
		cpu.regs[rd] = cpu.regs[rs1] + cpu.regs[rs2]

	default:
		return errors.New("Could not execute instruction. Function not yet implementd")

	}
	return nil
	//Decode Opcode etc.
}
