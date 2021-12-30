package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
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

	// fmt.Println("========= LABELS")
	// litter.Dump(Labels)
	// fmt.Println("========= STACK")
	// litter.Dump(StackPointer)
	// litter.Dump(Stack[0:10])
	// fmt.Println("========= MEMORY")
	// litter.Dump(MemoryPointer)
	// litter.Dump(Memory[MemoryPointer])
	// fmt.Println("========= RETURNS")
	// litter.Dump(ReturnStackPointer)
	// litter.Dump(ReturnStack[0:10])
}
