package ram

import (
	"errors"
	"fmt"
)

const debug bool = true

//MemorySize defines the  memory size for the cpu
const MemorySize uint64 = 1024 * 1024 * 128

//RAM defining the data section which will contain the radius and height variables .
type RAM struct {
	dram [MemorySize]uint8
}

// Using the methods here
// tankDimension interface

//Load value
func (r *RAM) Load(addr uint64, size uint64) (uint64, error) {
	if debug {
		fmt.Println("RAM Load addr, size", addr, size)
	}
	if addr > MemorySize {
		return 0, errors.New("Segmentation Fault")
	}
	switch size {
	case 8:
		return r.load8(addr), nil
	case 16:
		return r.load16(addr), nil
	case 32:
		return r.load32(addr), nil
	case 64:
		return r.load64(addr), nil
	default:
		return 0, errors.New("Could not read memory size as requested")

	}

}

//Store value
func (r *RAM) Store(addr uint64, size uint64, value uint64) error {
	if debug {
		fmt.Println("RAM Store addr, MemorySize, size, value ", addr, MemorySize, size, value)
	}
	if addr > MemorySize {
		return errors.New("Segmentation Fault")
	}
	switch size {
	case 8:
		r.store8(addr, value)
		return nil
	case 16:
		r.store16(addr, value)
		return nil
	case 32:
		r.store32(addr, value)
		return nil
	case 64:
		r.store64(addr, value)
		return nil
	default:
		return errors.New("Could not write memory size as requested")
	}
}

func (r *RAM) load8(addr uint64) uint64 {
	if debug {
		fmt.Println("RAM Load8 addr", uint(addr))
	}
	var index uint = uint(addr)
	return uint64(r.dram[index])
}

func (r *RAM) load16(addr uint64) uint64 {
	if debug {
		fmt.Println("RAM Load16 addr", uint(addr))
	}
	var index uint = uint(addr)
	return (uint64(r.dram[index]) | (uint64(r.dram[index+1]) << 8))
}

func (r *RAM) load32(addr uint64) uint64 {
	if debug {
		fmt.Println("RAM Load32 addr", uint(addr))
	}
	var index uint = uint(addr)
	return (uint64(r.dram[index]) | (uint64(r.dram[index+1]) << 8) | (uint64(r.dram[index+2]) << 16) | (uint64(r.dram[index+3]) << 24))
}

func (r *RAM) load64(addr uint64) uint64 {
	if debug {
		fmt.Println("RAM Load64 addr", uint(addr))
	}
	var index uint = uint(addr)
	return (uint64(r.dram[index]) | (uint64(r.dram[index+1]) << 8) | (uint64(r.dram[index+2]) << 16) | (uint64(r.dram[index+3]) << 24) | (uint64(r.dram[index+4]) << 32) | (uint64(r.dram[index+5]) << 40) | (uint64(r.dram[index+6]) << 48) | (uint64(r.dram[index+7]) << 56))
}

func (r *RAM) store8(addr uint64, value uint64) {
	if debug {
		fmt.Println("RAM Store8 addr, value ", addr, value)
	}
	var index uint = uint(addr)
	r.dram[index] = uint8((value & 0xff))
}

func (r *RAM) store16(addr uint64, value uint64) {
	if debug {
		fmt.Println("RAM Store16 addr, value ", addr, value)
	}
	var index uint = uint(addr)
	r.dram[index] = uint8(((value) & 0xff))
	r.dram[index+1] = uint8(((value >> 8) & 0xff))
}

func (r *RAM) store32(addr uint64, value uint64) {
	if debug {
		fmt.Println("RAM Store32 addr, value ", addr, value)
	}
	var index uint = uint(addr)
	r.dram[index] = uint8(((value) & 0xff))
	r.dram[index+1] = uint8(((value >> 8) & 0xff))
	r.dram[index+2] = uint8(((value >> 16) & 0xff))
	r.dram[index+3] = uint8(((value >> 24) & 0xff))
}

func (r *RAM) store64(addr uint64, value uint64) {
	if debug {
		fmt.Println("RAM Store64 addr, value ", addr, value)
	}
	var index uint = uint(addr)
	r.dram[index] = uint8(((value) & 0xff))
	r.dram[index+1] = uint8(((value >> 8) & 0xff))
	r.dram[index+2] = uint8(((value >> 16) & 0xff))
	r.dram[index+3] = uint8(((value >> 24) & 0xff))
	r.dram[index+4] = uint8(((value >> 32) & 0xff))
	r.dram[index+5] = uint8(((value >> 40) & 0xff))
	r.dram[index+6] = uint8(((value >> 48) & 0xff))
	r.dram[index+7] = uint8(((value >> 56) & 0xff))
	r.dram[index+5] = uint8(((value >> 48) & 0xff))
}
