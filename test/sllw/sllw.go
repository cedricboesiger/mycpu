package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// const debug bool = true

func main() {

	prg := "hart.go"
	inst := "test/sllw/sllw.bin"
	// cmd := exec.Command("pwd")
	cmd := exec.Command("go", "run", prg, "-f", inst)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Printf("failed exec: %s ", inst)
		return
	}

	if strings.Contains(string(stdout), "0x1e ( t5 ) = 0x15	0x1f ( t6 ) = 0x17") {
		fmt.Printf("passed test: %s", inst)
		fmt.Println()
		return
	}

	fmt.Printf("failed test: %s", inst)
	fmt.Println()
}
