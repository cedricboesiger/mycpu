package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"rvsim/cpu"
)

const debug bool = true

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

		//Fetch
		inst := cpu.Fetch()
		cpu.IncPC()
		if debug {
			fmt.Println("HART cpu.pc, ", cpu.GetPC())
		}

		//Workaround to abort loop
		if cpu.GetPC() == 0 {
			if debug {
				fmt.Println("HLT DUE TO END OF PROGAM")
			}
			break
		}
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
