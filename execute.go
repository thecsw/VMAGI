package main

const (
	STACK_DEPTH        = 10000
	RETURN_STACK_DEPTH = 10000
)

var (
	Memory        = map[ContextDepth](map[RegisterDepth]ValueWidth){}
	ContextNumber = ContextDepth(0)

	Labels = map[LabelType]InstructionDepth{}

	Stack        = make([]ValueWidth, STACK_DEPTH)
	StackPointer = uint32(0)

	ReturnStack        = make([]InstructionDepth, RETURN_STACK_DEPTH)
	ReturnStackPointer = uint32(0)

	PC = InstructionDepth(0)

	HaltValue ValueWidth
	Halted    bool
)

func Execute(instructions []*Instruction) {

	// fmt.Println("-------- LABELS --------")
	// litter.Dump(Labels)

	// Allocate the root context
	Memory[ContextNumber] = make(map[RegisterDepth]ValueWidth)

	var currentPC InstructionDepth
	var currentInstruction *Instruction
	for PC < InstructionDepth((len(instructions))) {
		currentPC = PC
		currentInstruction = instructions[PC]

		// fmt.Println("\n-------- STACK ---------")
		// fmt.Println(StackPointer)
		// litter.Dump(Stack[:3])
		// fmt.Println("\n------ REGISTERS -------")
		// fmt.Println(ContextNumber)
		// litter.Dump(Memory[ContextNumber])
		// fmt.Println("TO EXECUTE: ", currentInstruction.Input)

		executeFunctions[currentInstruction.Opcode](currentInstruction)
		// No PC manipulation happened
		if currentPC == PC {
			PC++
		}
		if Halted {
			return
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
		JUMPT:      executeJumpIfTrue,
		JUMPF:      executeJumpIfFalse,
		PUSH:       executePush,
		POP:        executePop,
		NOP:        executeNop,
	}
)

func executeAdd(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)+getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)+ValueWidth(inst.ImmediateValue))
	}
}

func executeSub(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)-getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)-ValueWidth(inst.ImmediateValue))
	}
}

func executeMul(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)*getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)*ValueWidth(inst.ImmediateValue))
	}
}

func executeDiv(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)/getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)/ValueWidth(inst.ImmediateValue))
	}
}

func executeMod(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)%getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)%ValueWidth(inst.ImmediateValue))
	}
}

func executeNeg(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, -1*getMemory(inst.SourceRegister1))
	} else {
		setMemory(inst.DestinationRegister, -1*ValueWidth(inst.ImmediateValue))
	}
}

func executeAnd(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)&getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)&ValueWidth(inst.ImmediateValue))
	}
}

func executeOr(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)|getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)|ValueWidth(inst.ImmediateValue))
	}
}

func executeNot(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, negateValue(valueToBool(getMemory(inst.SourceRegister1))))
	} else {
		setMemory(inst.DestinationRegister, negateValue(valueToBool(ValueWidth((inst.ImmediateValue)))))
	}
}

func executeXor(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)^getMemory(inst.SourceRegister2))
	} else {
		setMemory(inst.DestinationRegister, getMemory(inst.SourceRegister1)^ValueWidth(inst.ImmediateValue))
	}
}

func executeCall(inst *Instruction) {
	// Store the PC at the top of the stack
	ReturnStack[ReturnStackPointer] = PC + 1
	ReturnStackPointer++
	// Set the PC to point at the label
	PC = Labels[inst.LabelImmediate]
	// Start a new context
	ContextNumber++
	Memory[ContextNumber] = make(map[RegisterDepth]ValueWidth)
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
	Memory[ContextNumber] = nil
	ContextNumber--
}

func executeGreater(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) > ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) > ValueWidth(inst.ImmediateValue)))
	}
}

func executeGreaterEqual(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) >= ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) >= ValueWidth(inst.ImmediateValue)))
	}
}

func executeLess(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) < ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) < ValueWidth(inst.ImmediateValue)))
	}
}

func executeLessEqual(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) <= ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) <= ValueWidth(inst.ImmediateValue)))
	}
}

func executeEqual(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) == ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) == ValueWidth(inst.ImmediateValue)))
	}
}

func executeNotEqual(inst *Instruction) {
	if !inst.IsImmediate {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) != ValueWidth(inst.ImmediateValue)))
	} else {
		setMemory(inst.DestinationRegister, boolToValue(getMemory(inst.SourceRegister1) != ValueWidth(inst.ImmediateValue)))
	}
}

func executeJumpIfTrue(inst *Instruction) {
	if !inst.IsImmediate {
		if getMemory(inst.SourceRegister1) != 0 {
			PC = Labels[inst.LabelImmediate]
		}
	} else {
		if inst.ImmediateValue != 0 {
			PC = Labels[inst.LabelImmediate]
		}
	}
}

func executeJumpIfFalse(inst *Instruction) {
	if !inst.IsImmediate {
		if getMemory(inst.SourceRegister1) == 0 {
			PC = Labels[inst.LabelImmediate]
		}
	} else {
		if inst.ImmediateValue == 0 {
			PC = Labels[inst.LabelImmediate]
		}
	}
}

func executePush(inst *Instruction) {
	if !inst.IsImmediate {
		Stack[StackPointer] = getMemory(inst.SourceRegister1)
	} else {
		Stack[StackPointer] = ValueWidth(inst.ImmediateValue)
	}
	StackPointer++
}

func executePop(inst *Instruction) {
	StackPointer--
	Memory[ContextNumber][inst.DestinationRegister] = Stack[StackPointer]
}

func executeHalt(inst *Instruction) {
	if !inst.IsImmediate {
		HaltValue = getMemory(inst.SourceRegister1)
	} else {
		HaltValue = ValueWidth(inst.ImmediateValue)
	}
	Halted = true
}

// executeNop does nothing
func executeNop(inst *Instruction) {}

func executeTODO(inst *Instruction) {
	panic("executed TODO")
}

func boolToValue(what bool) ValueWidth {
	if what {
		return 1
	}
	return 0
}

func valueToBool(what ValueWidth) ValueWidth {
	if what != 0 {
		return 1
	}
	return 0
}

func negateValue(what ValueWidth) ValueWidth {
	if what == 0 {
		return 1
	}
	return 1
}

func setMemory(reg RegisterDepth, val ValueWidth) {
	Memory[ContextNumber][reg] = val
}

func getMemory(reg RegisterDepth) ValueWidth {
	return Memory[ContextNumber][reg]
}
