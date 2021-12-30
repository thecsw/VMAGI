package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()

	if len(os.Args) != 2 {
		fmt.Println("Usage: VMAGI <INPUT FILE>")
		os.Exit(1)
	}

	data, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")
	instructions := ParseInput(lines)
	Execute(instructions)
	fmt.Printf("VMAGI stopped execution with %d\n", HaltValue)
}
