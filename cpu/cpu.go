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
	if debug {
		fmt.Println("CPU Initialize")
	}
	prglength = uint64(len(binary))
	//	var ram bus.Device
	cpu.ram = &ram.RAM{}

	//The stack pointer
	cpu.regs[2] = ram.MemorySize
	cpu.pc = 0

	if debug {
		fmt.Println("CPU Store Program to Memory")
	}
	cpu.ram = &ram.RAM{}
	for i := 0; i <= len(binary)-1; i++ {
		if debug {
			fmt.Println("CPU Store index ", i)
		}
		err := cpu.ram.Store(uint64(i), 8, uint64(binary[i]))
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
	if (cpu.pc) >= prglength {
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
	value, err := cpu.ram.Load(cpu.pc, 32)

	return value, err
}

//Execute executes an instruction defined by its memory address
func Execute(instruction uint64) error {
	//Simulte the zero register at x00
	cpu.regs[0] = 0

	opcode := instruction & 0x7f
	rd := uint((instruction >> 7) & 0x1f)
	rs1 := uint((instruction >> 15) & 0x1f)
	rs2 := uint((instruction >> 20) & 0x1f)
	funct3 := (instruction >> 12) & 0x7
	funct7 := (instruction >> 25) & 0x7f
	if debug {
		fmt.Printf("CPU DEBUG pc %d instruction %d opcode %d funct3 %d funct7 %d", cpu.pc, instruction, opcode, funct3, funct7)
		fmt.Println()
	}

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
			//lwu load word unsigned
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
				cpu.regs[rd] = cpu.regs[rs1] >> uint64(shamt)
			case 0x10:
				//srai shift right arithmetic immediate
				cpu.regs[rd] = uint64(int64(cpu.regs[rs1]) >> int64(shamt))
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
	case 0x17:
		// U-Type
		//imm[31:12]
		imm := uint64((int64(int32(instruction & 0xfffff000))))
		//auipc add upper immediate value to pc
		cpu.regs[rd] = cpu.pc + imm - 4
	case 0x1b:
		//I-Type
		//imm[11:0], inst[31,20]
		imm := uint64(int64(int32(instruction)) >> 20)
		shamt := uint32(imm & 0x1f)
		switch funct3 {
		case 0x0:
			//addiw add word immediate
			cpu.regs[rd] = uint64(int64(int32(cpu.regs[rs1] + imm)))
		case 0x1:
			//slliw shift left logical word immediate
			cpu.regs[rd] = uint64(int64(int32(cpu.regs[rs1])) + int64(shamt))
		case 0x5:
			switch funct7 {
			case 0x00:
				//srliw shift right logical word immediate
				cpu.regs[rd] = uint64(int64(int32(uint32(cpu.regs[rs1]) >> shamt)))
			case 0x20:
				//sraiw shift right arithmetic word immediate
				cpu.regs[rd] = uint64(int64(int32(cpu.regs[rs1] >> shamt)))
			default:
				return errors.New("Could not execute funct7 of instruction 0x1b")
			}
		default:
			return errors.New("Could not execute funct3 of instruction 0x1b")
		}
	case 0x23:
		// S-Type
		// imm[11:5]|4:0], inst[31:24|11:7]
		imm := uint64((int64(int32(instruction&0xfe000000)))>>20) | ((instruction >> 7) & 0x1f)
		addr := cpu.regs[rs1] + imm // golang respects interger overflow on uint, see https://golang.org/ref/spec#Integer_overflow
		switch funct3 {
		case 0x0:
			//sb
			err := cpu.ram.Store(addr, 8, cpu.regs[rs2])
			if err != nil {
				return err
			}
		case 0x1:
			//sh
			err := cpu.ram.Store(addr, 16, cpu.regs[rs2])
			if err != nil {
				return err
			}
		case 0x2:
			//sw
			err := cpu.ram.Store(addr, 32, cpu.regs[rs2])
			if err != nil {
				return err
			}
		case 0x3:
			//sd
			err := cpu.ram.Store(addr, 64, cpu.regs[rs2])
			if err != nil {
				return err
			}
		default:
			return errors.New("Could not execute funct3 of instruction 0x23")
		}
	case 0x33:
		shamt := uint32(uint64(cpu.regs[rs2] & 0x3f))
		switch funct3 {
		case 0x0:
			switch funct7 {
			case 0x00:
				//add
				cpu.regs[rd] = cpu.regs[rs1] + cpu.regs[rs2]
			case 0x01:
				//mul
				cpu.regs[rd] = cpu.regs[rs1] * cpu.regs[rs2]
			case 0x02:
				//sub
				cpu.regs[rd] = cpu.regs[rs1] - cpu.regs[rs2]
			default:
				return errors.New("Could not execute funct7 of funct3 0x0 of instruction 0x33")
			}
		case 0x1:
			switch funct7 {
			case 0x00:
				//sll
				cpu.regs[rd] = cpu.regs[rs1] << (shamt)
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x1 instruction 0x33")
			}
		case 0x2:
			switch funct7 {
			case 0x00:
				//slt
				if int64(cpu.regs[rs1]) < int64(cpu.regs[rs2]) {
					cpu.regs[rd] = 1
				} else {
					cpu.regs[rd] = 0
				}
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x2 instruction 0x33")
			}
		case 0x3:
			switch funct7 {
			case 0x00:
				//sltu
				if cpu.regs[rs1] < cpu.regs[rs2] {
					cpu.regs[rd] = 1
				} else {
					cpu.regs[rd] = 0
				}
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x3 instruction 0x33")
			}
		case 0x4:
			switch funct7 {
			case 0x00:
				//xor
				cpu.regs[rd] = cpu.regs[rs1] ^ cpu.regs[rs2]
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x4 instruction 0x33")
			}
		case 0x5:
			switch funct7 {
			case 0x00:
				//srl
				cpu.regs[rd] = cpu.regs[rs1] >> (shamt)
			case 0x20:
				//sra
				cpu.regs[rd] = uint64(int64(cpu.regs[rs1]) >> (shamt))
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x5 instruction 0x33")
			}
		case 0x6:
			switch funct7 {
			case 0x00:
				//or
				cpu.regs[rd] = cpu.regs[rs1] | cpu.regs[rs2]
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x6 instruction 0x33")
			}
		case 0x7:
			switch funct7 {
			case 0x00:
				//and
				cpu.regs[rd] = cpu.regs[rs1] & cpu.regs[rs2]
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x7 instruction 0x33")
			}
		default:
			return errors.New("Could not execute funct3 of instruction 0x33")
		}
	case 0x37:
		//lui
		cpu.regs[rd] = uint64(int64(int32((instruction & 0xfffff000))))
	case 0x3b:
		shamt := uint32(cpu.regs[rs2] & 0x1f)
		switch funct3 {
		case 0x0:
			switch funct7 {
			case 0x00:
				//addw
				cpu.regs[rd] = uint64(int64(int32(cpu.regs[rs1] + cpu.regs[rs2])))
			case 0x20:
				//subw
				cpu.regs[rd] = uint64(int32(cpu.regs[rs1] - cpu.regs[rs2]))
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x0 instruction 0x3b")
			}
		case 0x1:
			switch funct7 {
			case 0x00:
				//sllw
				cpu.regs[rd] = uint64(int32(uint32(cpu.regs[rs1]) << (shamt)))
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x1 instruction 0x3b")
			}
		case 0x5:
			switch funct7 {
			case 0x00:
				//slrw
				cpu.regs[rd] = uint64(int32(uint32(cpu.regs[rs1]) >> (shamt)))
			case 0x20:
				//sraw
				cpu.regs[rd] = uint64(uint32(cpu.regs[rs1]) >> int32((shamt)))
			default:
				return errors.New("Could not execute funct7 of funct3 of 0x5 instruction 0x3b")
			}
		default:
			return errors.New("Could not execute funct3 of instruction 0x3b")
		}
	case 0x63:
		// B-Type
		// imm[12|10:5|4:1|11]
		imm := uint64((int64(int32(instruction&0x80000000)) >> 19)) | ((instruction & 0x80) << 4) | ((instruction >> 20) & 0x7e0) | ((instruction >> 7) & 0x1e)
		switch funct3 {
		case 0x0:
			//beq
			if debug {
				fmt.Printf("DEBUG BEQ rs1 %d rs2 %d cpu %d imm %d", cpu.regs[rs1], cpu.regs[rs2], cpu.pc, imm)
				fmt.Println()
			}
			if cpu.regs[rs1] == cpu.regs[rs2] {
				cpu.pc = cpu.pc + imm - 4
			}
		case 0x1:
			//bne
			if cpu.regs[rs1] != cpu.regs[rs2] {
				cpu.pc = cpu.pc + imm - 4
			}
		case 0x4:
			//blt
			if int64(cpu.regs[rs1]) < int64(cpu.regs[rs2]) {
				cpu.pc = cpu.pc + imm - 4
			}
		case 0x5:
			//bge
			if int64(cpu.regs[rs1]) >= int64(cpu.regs[rs2]) {
				cpu.pc = cpu.pc + imm - 4
			}
		case 0x6:
			//bltu
			if cpu.regs[rs1] < cpu.regs[rs2] {
				cpu.pc = cpu.pc + imm - 4
			}
		case 0x7:
			//bgeu
			if cpu.regs[rs1] >= cpu.regs[rs2] {
				cpu.pc = cpu.pc + imm - 4
			}
		default:
			return errors.New("Could not execute funct3 of instruction 0x63")
		}
	case 0x67:
		//jal
		// Don'd add 4 to t because pc already moved
		t := cpu.pc
		imm := uint64(int64(int32((instruction & 0xfff00000))) >> 20)
		cpu.pc = cpu.regs[rs1] + imm
		cpu.regs[rd] = t
	case 0x6f:
		//jal
		cpu.regs[rd] = cpu.pc

		// imm[20|10:1|11|19:12]
		imm := uint64((int64(int32(instruction&0x80000000)))>>11) | (instruction & 0xff000) | ((instruction >> 9) & 0x800) | ((instruction >> 20) & 0x7fe)
		cpu.pc = cpu.pc + imm - 4
	default:
		return errors.New("Could not execute instruction. Function not yet implementd")
	}
	return nil
	//Decode Opcode etc.
}
