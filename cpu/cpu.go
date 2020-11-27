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
	funct3 := (instruction & 0x00007000) >> 12
	funct7 := (instruction & 0xfe000000) >> 25

	switch opcode {
	case 0x03:
		//I-Type
		//imm[11:0], inst[31:20]
		imm := uint64((int64(int32(instruction))) >> 20)
		addr := cpu.regs[rs1] + imm // golang respects interger overflow on uint, see https://golang.org/ref/spec#Integer_overflow

		switch funct3 {
		case 0x0:
			//lb load byte
			val, _ := cpu.ram.Load(addr, 8)
			cpu.regs[rd] = uint64(int64(int8((val))))
		case 0x1:
			//lh load half word
			val, _ := cpu.ram.Load(addr, 16)
			cpu.regs[rd] = uint64(int64(int16(val)))
		case 0x2:
			//lw load word
			val, _ := cpu.ram.Load(addr, 32)
			cpu.regs[rd] = uint64(int64(int32(val)))
		case 0x3:
			//ld load double word
			val, _ := cpu.ram.Load(addr, 64)
			cpu.regs[rd] = uint64(val)
		case 0x4:
			//lbu load byte unsigned
			val, _ := cpu.ram.Load(addr, 8)
			cpu.regs[rd] = val
		case 0x5:
			//lhu load half word unsigned
			val, _ := cpu.ram.Load(addr, 16)
			cpu.regs[rd] = val
		case 0x6:
			//lwh load word unsigned
			val, _ := cpu.ram.Load(addr, 32)
			cpu.regs[rd] = val
		case 0x7:
			//ldu load double word unsigned
			val, _ := cpu.ram.Load(addr, 64)
			cpu.regs[rd] = val
		default:
			return errors.New("Could not execute funct3 of instruction 0x03")
		}
	case 0x13:
		//I-Type
		//imm[11:0]
		imm := uint64((int64(int32(instruction & 0xfff00000))) >> 20)
		//the shift amount is encoded in the lower 6 bits of imm for RV64I
		shamt := uint32((imm & 0x3f))

		switch funct3 {
		case 0x0:
			//addi add immediate
			cpu.regs[rd] = cpu.regs[rs1] + imm
		case 0x1:
			//slli shift left logical immediate
			cpu.regs[rd] = cpu.regs[rs1] << shamt
		case 0x2:
			//slti set if less than
			if (int64(cpu.regs[rs1])) < (int64(imm)) {
				cpu.regs[rd] = 1
			} else {
				cpu.regs[rd] = 0
			}
		case 0x3:
			//sltiu set if less than unsigned
			if cpu.regs[rs1] < imm {
				cpu.regs[rd] = 1
			} else {
				cpu.regs[rd] = 0
			}
		case 0x4:
			//xori exclusive or immediate
			cpu.regs[rd] = cpu.regs[rs1] ^ imm
		case 0x5:
			switch funct7 {
			case 0x00:
				//srli shift right logical immediate.
				cpu.regs[rd] = cpu.regs[rs1] + uint64(shamt)
			case 0x10:
				//srai shift right arithmetic immediate
				cpu.regs[rd] = uint64(int64(cpu.regs[rs1]) + int64(shamt))
			default:
				return errors.New("Coud not execute funct7 of instruction 0x13")
			}
		case 0x6:
			//ori or immediate
			cpu.regs[rd] = cpu.regs[rs1] | imm

		case 0x7:
			//andi and immediate
			cpu.regs[rd] = cpu.regs[rs1] & imm
		default:
			return errors.New("Could not execute funct3 of instruction 0x13")

		}
	case 0x33:
		//add
		cpu.regs[rd] = cpu.regs[rs1] + cpu.regs[rs2]
	case 0x17:
		// U-Type
		//imm[31:12]
		imm := uint64((int64(int32(instruction & 0xfffff000))))
		//auipc add upper immediate value to pc
		cpu.regs[rd] = cpu.pc + imm - 4
	default:
		return errors.New("Could not execute instruction. Function not yet implementd")

	}
	return nil
	//Decode Opcode etc.
}
