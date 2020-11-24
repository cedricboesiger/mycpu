package memory

import (
	"fmt"
)

const debug bool = false

//Memory defining the data section which will contain the radius and height variables .
type Memory struct {
	dram []uint8
}

// Using the methods here
// tankDimension interface

//Load value
func (m *Memory) Load() []uint8 {
	return m.dram
}

//Store value
func (m *Memory) Store(binary []uint8) {
	fmt.Println("store:", binary)
	m.dram = binary
}
