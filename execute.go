package main

const (
	STACK_DEPTH        = 10000
	RETURN_STACK_DEPTH = 10000

	CONTEXT_SIZE = 10
)

var (
	Memory = map[RegisterDepth]ValueWidth{}
	//ContextNumber = ContextDepth(0)
	MemoryOffset = RegisterDepth(0)

	Labels = map[LabelType]InstructionDepth{}

	Stack = &Stack64{}

	ReturnStack = &Stack32{}

	PC = InstructionDepth(0)

	HaltValue     ValueWidth
	Halted        bool
	ContextNumber = 0
)

func Execute(instructions []*Instruction) {

	Stack.Init(STACK_DEPTH)
	ReturnStack.Init(RETURN_STACK_DEPTH)

	var currentPC InstructionDepth
	var currentInstruction *Instruction
	for PC < InstructionDepth((len(instructions))) {
		currentPC = PC
		currentInstruction = instructions[PC]
		executeFunctions[currentInstruction.Opcode](currentInstruction)
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
	ReturnStack.Push(PC + 1)
	PC = Labels[inst.LabelImmediate]
	ContextNumber++
}

func executeJump(inst *Instruction) {
	PC = Labels[inst.LabelImmediate]
}

func executeReturn(inst *Instruction) {
	PC = ReturnStack.Pop()
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
		Stack.Push(getMemory(inst.SourceRegister1))
	} else {
		Stack.Push(ValueWidth(inst.ImmediateValue))
	}
}

func executePop(inst *Instruction) {
	setMemory(inst.DestinationRegister, Stack.Pop())
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
	Memory[RegisterDepth(ContextNumber)*CONTEXT_SIZE+reg] = val
}

func getMemory(reg RegisterDepth) ValueWidth {
	return Memory[RegisterDepth(ContextNumber)*CONTEXT_SIZE+reg]
}
