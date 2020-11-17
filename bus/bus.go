package main

import "fmt"

const debug bool = false

type Device interface {
	// Methods signature with data types of the methods .
	Load() uint64
	Store(uint64)
}

//Defining the data section which will contain the radius and height variables .
type memory struct {
	dram uint64
}

// Using the methods here
// tankDimension interface

func (m memory) Load() uint64 {
	return m.dram
}
func (m *memory) Store(binary uint64) {
	fmt.Println("store:", binary)
	m.dram = binary
}

// Defining the main method
func main() {
	// Here we are trying to access the elements or attributes
	// Creating a tankDimension interface
	var dev Device
	dev = &memory{}
	fmt.Println("mem :", dev.Load())
	dev.Store(5)
	fmt.Println("mem :", dev.Load())
}
