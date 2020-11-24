package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"rvsim/cpu"
	"rvsim/ram"
)

const debug bool = false

func main() {

	filePtr := flag.String("f", "test/a.out", "<binary-filename> of executable")
	flag.Parse()

	data, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		fmt.Println("Error reading binary file: ", err)
	}

	if debug {
		fmt.Println("HART DUMP, ", data)
	}
	//Initialize an fresh and empty cpu
	cpu.Initialize(data)

	//the fetch/decode/execute cycles
	for {
		if cpu.GetPC() > uint64(ram.MemorySize) {
			if debug {
				fmt.Println("HART BREAK, MemorySize cpu.pc", uint64(ram.MemorySize), cpu.GetPC())
			}
			break
		}
		if debug {
			fmt.Println("HART DEBUG, MemorySize cpu.pc", uint64(ram.MemorySize), cpu.GetPC())
		}

		//Fetch
		inst := cpu.Fetch()
		cpu.IncPC()

		//Decode / Execute
		err := cpu.Execute(inst)

		if err != nil {
			fmt.Println("PANIC: ", err)
			break
		}
	}

	//Show all registers
	cpu.DumpRegisters()
}
