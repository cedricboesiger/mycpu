package main

import (
	"fmt"
)

const debug bool = false

//Device for the bus need to implement load and store
type Device interface {
	// Methods signature with data types of the methods .
	Load() uint64
	Store(uint64)
}

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

// Defining the main method
func main() {
	// Here we are trying to access the elements or attributes
	// Creating a tankDimension interface
	var dev Device
	dev = &Memory{}
	fmt.Println("mem :", dev.Load())
	dev.Store(5)
	fmt.Println("mem :", dev.Load())
}
