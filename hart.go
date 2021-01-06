package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"rvsim/cpu"
	"time"
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

	//Figure execution Hz
	start := time.Now()
	cycle := 0
	//the fetch/decode/execute cycles
	for {

		//Fetch
		inst, err := cpu.Fetch()
		if err != nil {
			fmt.Println("Error Fetch inst from ram")
			break
		}

		cpu.IncPC()

		if debug {
			fmt.Println("HART cpu.pc, ", cpu.GetPC())
		}

		//Workaround to abort loop
		if cpu.GetPC() == 4 && cycle > 0 {
			if debug {
				fmt.Println("HLT DUE TO END OF PROGAM")
			}
			break
		}
		//Decode / Execute
		err = cpu.Execute(inst)

		if err != nil {
			fmt.Println("PANIC: ", err)
			break
		}
		cycle++
	}
	fmt.Printf("CPU speed %.1f kHz", float64(cycle)/time.Since(start).Seconds()/1000)
	fmt.Println()
	//Show all registers
	cpu.DumpRegisters()
}
