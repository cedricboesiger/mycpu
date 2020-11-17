package memory

import (
	"fmt"
)

const debug bool = false

//Memory defining the data section which will contain the radius and height variables .
type Memory struct {
	dram uint64
}

// Using the methods here
// tankDimension interface

//Load value
func (m Memory) Load() uint64 {
	return m.dram
}

//Store value
func (m *Memory) Store(binary uint64) {
	fmt.Println("store:", binary)
	m.dram = binary
}
