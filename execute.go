package main

import (
	"fmt"
	"os"
)

const (
	STACK_DEPTH        = 10000
	RETURN_STACK_DEPTH = 10000
)

var (
	Memory        = map[uint32](map[uint16]int64){}
	MemoryPointer = uint32(0)

	Registers = map[uint8]int32{}
	Labels    = map[string]uint32{}

	Stack        = make([]int64, STACK_DEPTH)
	StackPointer = uint32(0)

	ReturnStack        = make([]uint32, RETURN_STACK_DEPTH)
	ReturnStackPointer = uint32(0)

	PC = uint32(0)
)

func Execute(instructions []*Instruction) {

	// fmt.Println("-------- LABELS --------")
	// litter.Dump(Labels)

	// Allocate the root context
	Memory[MemoryPointer] = make(map[uint16]int64)

	var currentPC uint32
	var currentInstruction *Instruction
	for PC < uint32(len(instructions)) {
		// fmt.Println("\n-------- STACK ---------")
		// fmt.Println(StackPointer)
		// litter.Dump(Stack[:3])
		// fmt.Println("\n------------------------")
		currentPC = PC
		currentInstruction = instructions[PC]
		//fmt.Println("EXECUTING: ", PC, currentInstruction.Opcode)
		executeFunctions[currentInstruction.Opcode](currentInstruction)
		// No PC manipulation happened
		if currentPC == PC {
			PC++
		}
	}
}

var (
	executeFunctions = map[OpcodeNumber]func(*Instruction){
		HALT:       executeHalt,
		ADD:        executeAdd,
		SUB:        executeSub,
		MUL:        executeMul,
		DIV:        executeDiv,
		MOD:        executeMod,
		NEG:        executeNeg,
		AND:        executeAnd,
		OR:         executeOr,
		NOT:        executeNot,
		XOR:        executeXor,
		CALL:       executeCall,
		JUMP:       executeJump,
		RETURN:     executeReturn,
		GREATER:    executeGreater,
		GREATER_EQ: executeGreaterEqual,
		LESS:       executeLess,
		LESS_EQ:    executeLessEqual,
		EQUAL:      executeEqual,
		NOT_EQUAL:  executeNotEqual,
		JUMPIF:     executeJumpIf,
		PUSH:       executePush,
		POP:        executePop,
		NOP:        executeNop,
	}
)

func executeAdd(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] + Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] + int64(inst.ImmediateValue)
	}
}

func executeSub(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] - Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] - int64(inst.ImmediateValue)
	}
}

func executeMul(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] * Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] * int64(inst.ImmediateValue)
	}
}

func executeDiv(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] / Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] / int64(inst.ImmediateValue)
	}
}

func executeMod(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] % Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] % int64(inst.ImmediateValue)
	}
}

func executeNeg(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			-1 * Memory[MemoryPointer][inst.SourceRegister1]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			-1 * int64(inst.ImmediateValue)
	}
}

func executeAnd(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] & Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] & int64(inst.ImmediateValue)
	}
}

func executeOr(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] | Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] | int64(inst.ImmediateValue)
	}
}

func executeNot(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			negateInt64(int64ToBool(int64(Memory[MemoryPointer][inst.SourceRegister1])))
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			negateInt64(int64ToBool(int64(inst.ImmediateValue)))
	}
}

func executeXor(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] ^ Memory[MemoryPointer][inst.SourceRegister2]
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			Memory[MemoryPointer][inst.SourceRegister1] ^ int64(inst.ImmediateValue)
	}
}

func executeCall(inst *Instruction) {
	// Store the PC at the top of the stack
	ReturnStack[ReturnStackPointer] = PC + 1
	ReturnStackPointer++
	// Set the PC to point at the label
	PC = Labels[inst.LabelImmediate]
	// Start a new context
	MemoryPointer++
	Memory[MemoryPointer] = make(map[uint16]int64)
}

func executeJump(inst *Instruction) {
	// Simply set PC to the label
	PC = Labels[inst.LabelImmediate]
}

func executeReturn(inst *Instruction) {
	// Set the PC to the last return address
	ReturnStackPointer--
	PC = ReturnStack[ReturnStackPointer]
	// Close the current context
	Memory[MemoryPointer] = nil
	MemoryPointer--
}

func executeGreater(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] > Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] > int64(inst.ImmediateValue))
	}
}

func executeGreaterEqual(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] >= Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] >= int64(inst.ImmediateValue))
	}
}

func executeLess(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] < Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] < int64(inst.ImmediateValue))
	}
}

func executeLessEqual(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] <= Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] <= int64(inst.ImmediateValue))
	}
}

func executeEqual(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] == Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] == int64(inst.ImmediateValue))
	}
}

func executeNotEqual(inst *Instruction) {
	if !inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] != Memory[MemoryPointer][inst.SourceRegister2])
	}
	if inst.IsImmediate {
		Memory[MemoryPointer][inst.DestinationRegister] =
			boolToInt64(Memory[MemoryPointer][inst.SourceRegister1] != int64(inst.ImmediateValue))
	}
}

func executeJumpIf(inst *Instruction) {
	if !inst.IsImmediate {
		if Memory[MemoryPointer][inst.SourceRegister1] == 0 {
			PC = Labels[inst.LabelImmediate]
		}
	}
	if inst.IsImmediate {
		if inst.ImmediateValue == 0 {
			PC = Labels[inst.LabelImmediate]
		}
	}
}

func executePush(inst *Instruction) {
	if inst.IsImmediate {
		Stack[StackPointer] = int64(inst.ImmediateValue)
	} else {
		Stack[StackPointer] = Memory[MemoryPointer][inst.SourceRegister1]
	}
	StackPointer++
}

func executePop(inst *Instruction) {
	StackPointer--
	Memory[MemoryPointer][inst.DestinationRegister] = Stack[StackPointer]
}

func executeHalt(inst *Instruction) {

	// fmt.Println("========= MEMORY")
	// litter.Dump(MemoryPointer)
	// litter.Dump(Memory[MemoryPointer])

	if !inst.IsImmediate {
		fmt.Printf("VMAGI exited with %d\n", Memory[MemoryPointer][inst.SourceRegister1])
		os.Exit(0)
	}
	if inst.IsImmediate {
		fmt.Printf("VMAGI exited with %d\n", inst.ImmediateValue)
		os.Exit(0)
	}
}

// executeNop does nothing
func executeNop(inst *Instruction) {}

func executeTODO(inst *Instruction) {
	panic("executed TODO")
}

func boolToInt64(what bool) int64 {
	if what {
		return 1
	}
	return 0
}

func int64ToBool(what int64) int64 {
	if what != 0 {
		return 1
	}
	return 0
}

func negateInt64(what int64) int64 {
	if what == 0 {
		return 1
	}
	return 1
}
