package cpu

import (
	"errors"
	"fmt"
)

const debug bool = false

//MemorySize defines the  memory size for the cpu
const memorySize uint64 = 1024 * 1024 * 128

//Register holds all registers of the cpu
type register struct {
	//32 Bit registers
	regs [32]uint64
	//one Program Counter
	pc uint64
	//Add memory here for simplicity
	memory []uint8
}

var cpu register

//Initialize the cpu
func Initialize(binary []uint8) {
	cpu.regs[2] = memorySize
	cpu.pc = 0
	cpu.memory = binary

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
	cpu.pc += 4
}

//SetPC used for jump instructions
func SetPC(newPC uint64) {
	cpu.pc = newPC
}

//GetMemory return memory as array
func GetMemory() []uint8 {
	return cpu.memory
}

//Fetch cycle
func Fetch() uint64 {
	if debug {
		fmt.Println("DEBUG Fetch cpu.pc", cpu.pc)
	}
	return uint64(read32(cpu.pc))
}

//read32 for 32 bit code
func read32(addr uint64) uint64 {
	index := uint(addr)

	//Shift bits to little-endian order
	return (uint64(cpu.memory[index]) | (uint64(cpu.memory[index+1]) << 8) | (uint64(cpu.memory[index+2]) << 16) | (uint64(cpu.memory[index+3]) << 24))
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
