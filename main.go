package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sanity-io/litter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: VMAGI <INPUT FILE>")
		os.Exit(1)
	}

	data, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(string(data), "\n")

	instructions := ParseInput(lines)

	//litter.Dump(instructions)

	Execute(instructions)

	fmt.Println("========= LABELS")
	litter.Dump(Labels)
	fmt.Println("========= STACK")
	litter.Dump(StackPointer)
	litter.Dump(Stack[0:10])
	fmt.Println("========= MEMORY")
	litter.Dump(MemoryPointer)
	litter.Dump(Memory[MemoryPointer])
	fmt.Println("========= RETURNS")
	litter.Dump(ReturnStackPointer)
	litter.Dump(ReturnStack[0:10])
}
