package main

const (
	MEMORY_DEPTH       = 100000
	STACK_DEPTH        = 100000
	RETURN_STACK_DEPTH = STACK_DEPTH

	CONTEXT_SIZE = 10
)

var (
	Memory      = make([]ValueWidth, MEMORY_DEPTH)
	Stack       = &Stack64{}
	ReturnStack = &Stack32{}
	PC          = InstructionDepth(0)
	currInst    *Instruction
	HaltValue   ValueWidth
	Halted      bool

	ContextNumber = 0
	MemoryOffset  = RegisterDepth(0)

	dst   RegisterDepth
	src1  RegisterDepth
	src2  RegisterDepth
	isImm bool
	imm   ValueWidth
)

func Execute(instructions []*Instruction) {
	// Init the execution function array generated from the map
	for i := OpcodeNumber(0); i <= NOP; i++ {
		executeFunctions[i] = executeFunctionsMap[i]
	}
	optimizerInitialize()

	// Init the stacks
	Stack.Init(STACK_DEPTH)
	ReturnStack.Init(RETURN_STACK_DEPTH)

	// Start going through instructions
	for {
		currInst = instructions[PC]
		dst = currInst.DestinationRegister
		src1 = currInst.SourceRegister1
		src2 = currInst.SourceRegister2
		isImm = currInst.IsImmediate

		PC++
		executeFunctions[currInst.Opcode]()
		if Halted {
			return
		}
	}
}

var (
	executeFunctionsMap = map[OpcodeNumber]func(){
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

	executeFunctions = make([]func(), len(executeFunctionsMap))
)

func executeAdd() {
	if !isImm {
		setMemory(dst, getMemory(src1)+getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)+currInst.ImmediateValue)
	}
}

func executeSub() {
	if !isImm {
		setMemory(dst, getMemory(src1)-getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)-currInst.ImmediateValue)
	}
}

func executeMul() {
	if !isImm {
		setMemory(dst, getMemory(src1)*getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)*currInst.ImmediateValue)
	}
}

func executeDiv() {
	if !isImm {
		setMemory(dst, getMemory(src1)/getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)/currInst.ImmediateValue)
	}
}

func executeMod() {
	if !isImm {
		setMemory(dst, getMemory(src1)%getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)%currInst.ImmediateValue)
	}
}

func executeNeg() {
	if !isImm {
		setMemory(dst, -1*getMemory(src1))
	} else {
		setMemory(dst, -1*currInst.ImmediateValue)
	}
}

func executeAnd() {
	if !isImm {
		setMemory(dst, getMemory(src1)&getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)&currInst.ImmediateValue)
	}
}

func executeOr() {
	if !isImm {
		setMemory(dst, getMemory(src1)|getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)|currInst.ImmediateValue)
	}
}

func executeNot() {
	if !isImm {
		setMemory(dst, negateValue(valueToBool(getMemory(src1))))
	} else {
		setMemory(dst, negateValue(valueToBool(currInst.ImmediateValue)))
	}
}

func executeXor() {
	if !isImm {
		setMemory(dst, getMemory(src1)^getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)^currInst.ImmediateValue)
	}
}

func executePush() {
	if !isImm {
		Stack.Push(getMemory(src1))
	} else {
		Stack.Push(currInst.ImmediateValue)
	}
}

func executePop() {
	setMemory(dst, Stack.Pop())
}

func executeCall() {
	// If optimizer found a cached call for this function,
	// then it already simulated the function. We leave
	if cacheFunctionCall(currInst.LabelImmediate) {
		return
	}

	ReturnStack.Push(PC)
	PC = currInst.LabelImmediate
	LastFunctionCalled = PC
	ContextNumber++
	MemoryOffset = RegisterDepth(ContextNumber)*CONTEXT_SIZE - 1
	// See if we need to bump up the memory
	if int(MemoryOffset) >= len(Memory) {
		newMemory := make([]ValueWidth, len(Memory)*2)
		copy(newMemory, Memory)
		Memory = newMemory
	}

	optimizerAnalyzeCallStack()
}

func executeReturn() {
	PC = ReturnStack.Pop()
	ContextNumber--
	MemoryOffset = RegisterDepth(ContextNumber)*CONTEXT_SIZE - 1

	optimizerAnalyzeReturnStack()
}

func executeJump() {
	PC = InstructionDepth(currInst.LabelImmediate)
}

func executeGreater() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) > getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) > currInst.ImmediateValue))
	}
}

func executeGreaterEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) >= getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) >= currInst.ImmediateValue))
	}
}

func executeLess() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) < getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) < currInst.ImmediateValue))
	}
}

func executeLessEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) <= getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) <= currInst.ImmediateValue))
	}
}

func executeEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) == getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) == currInst.ImmediateValue))
	}
}

func executeNotEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) != getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) != currInst.ImmediateValue))
	}
}

func executeJumpIfTrue() {
	if !isImm {
		if getMemory(src1) != 0 {
			PC = currInst.LabelImmediate
		}
	} else {
		if imm != 0 {
			PC = currInst.LabelImmediate
		}
	}
}

func executeJumpIfFalse() {
	if !isImm {
		if getMemory(src1) == 0 {
			PC = currInst.LabelImmediate
		}
	} else {
		if imm == 0 {
			PC = currInst.LabelImmediate
		}
	}
}

func executeHalt() {
	if !isImm {
		HaltValue = getMemory(src1)
	} else {
		HaltValue = currInst.ImmediateValue
	}
	Halted = true
}

// executeNop does nothing
func executeNop() {}

func executeTODO() {
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
	Memory[MemoryOffset+reg] = val
}

func getMemory(reg RegisterDepth) ValueWidth {
	return Memory[MemoryOffset+reg]
}
