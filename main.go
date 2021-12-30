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
		return
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("You gave me a bad file:", err.Error())
		return
	}

	// Execute the input by splitting by newlines and parsing
	Execute(ParseInput(strings.Split(string(data), "\n")))

	// Show the halt value
	fmt.Printf("VMAGI stopped execution with %d\n", HaltValue)
}
