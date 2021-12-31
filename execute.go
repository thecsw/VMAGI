package main

const (
	// MEMORY_DEPTH regulates on how many variables should we
	// support in the whole VM, all contexts go here.
	MEMORY_DEPTH = 100000

	// STACK_DEPTH regulates the initial stack size for all
	// the stack operations we want to do.
	STACK_DEPTH = 100000

	// RETURN_STACK_DEPTH regulates on how many function calls
	// can we make and store their return addresses.
	RETURN_STACK_DEPTH = STACK_DEPTH

	// CONTEXT_SIZE regulates on how many variables we can have
	// per context, can be easily changed, this also maves the offset
	// accordingly
	CONTEXT_SIZE = 10
)

var (
	// Memory is just a slab array on all our variables and contexts,
	// notice that it's continuous.
	Memory = make([]ValueWidth, MEMORY_DEPTH)
	// Stack is the stack used by functions to pass values across, such
	// as arguments and returns.
	Stack = &Stack64{}
	// ReturnStack stores the instruction addresses of where the program
	// counter should return to after a function returns.
	ReturnStack = &Stack32{}
	// PC is the Program Counter that points at the current instruction
	PC = InstructionDepth(0)
	// currInst is the current instruction executed
	currInst *Instruction
	// HaltValue is the value VMAGI should exit with
	HaltValue ValueWidth
	// Halted is a flag, if set true, VMAGI shuts down with HaltValue printed
	Halted bool

	// ContextNumber shows how deep are we in the function call stack
	ContextNumber = 0
	// MemoryOffset is the offset for contexts to refer to memory
	MemoryOffset = RegisterDepth(0)

	// dst is the destination register of the current instruction
	dst RegisterDepth
	// src is the first source register of the current instruction
	src1 RegisterDepth
	// src is the second source register of the current instruction
	src2 RegisterDepth
	// isImm indicates whether the current instruction has an
	// immediate field as one of the operands
	isImm bool
	// imm is the immediate value that has been passed as an operand
	imm ValueWidth
)

// Execute takes a slice of instruction pointers and runs them
// one-by-one with PC and enabling all the optimizations
func Execute(instructions []*Instruction) {
	// Init the execution function array generated from the map
	for i := OpcodeNumber(0); i <= NOP; i++ {
		executeFunctions[i] = executeFunctionsMap[i]
	}
	// Enable the optimizer
	optimizerInitialize()

	// Init the stacks
	Stack.Init(STACK_DEPTH)
	ReturnStack.Init(RETURN_STACK_DEPTH)

	// Start going through instructions
	for {
		// Retrieve and unmarshal the current instruction
		currInst = instructions[PC]
		dst = currInst.DestinationRegister
		src1 = currInst.SourceRegister1
		src2 = currInst.SourceRegister2
		isImm = currInst.IsImmediate

		// Bump the program counter
		PC++
		// Run the actual instruction
		executeFunctions[currInst.Opcode]()
		// Shut down if `halt VAL` is run
		if Halted {
			return
		}
	}
}

var (
	// executeFunctionsMap is a nice wap to map opcodes
	// to the functions that should run on the opcode.
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

	// executeFunctions is the slice version of executeFunctionsMap for
	// performance improvements, as indexing is faster than hashing lookup.
	executeFunctions = make([]func(), len(executeFunctionsMap))
)

// executeAdd adds the first two operands and stores in the third one.
func executeAdd() {
	if !isImm {
		setMemory(dst, getMemory(src1)+getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)+currInst.ImmediateValue)
	}
}

// executeSub subtracts the second operand from the first one and stores in the third.
func executeSub() {
	if !isImm {
		setMemory(dst, getMemory(src1)-getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)-currInst.ImmediateValue)
	}
}

// executeMul multiplies the first two operands and stores the result in the third.
func executeMul() {
	if !isImm {
		setMemory(dst, getMemory(src1)*getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)*currInst.ImmediateValue)
	}
}

// executeDiv divides the first operand by the second one and stores in the third,
// the emulator will crash if the second operand is a zero.
func executeDiv() {
	if !isImm {
		setMemory(dst, getMemory(src1)/getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)/currInst.ImmediateValue)
	}
}

// executeMod find the remainder when dividing the first operand by the second one,
// by storing the result in the third.
func executeMod() {
	if !isImm {
		setMemory(dst, getMemory(src1)%getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)%currInst.ImmediateValue)
	}
}

// executeNeg negates the first operand and stores the result in the second one.
func executeNeg() {
	if !isImm {
		setMemory(dst, -1*getMemory(src1))
	} else {
		setMemory(dst, -1*currInst.ImmediateValue)
	}
}

// executeAnd finds a bitwise AND of the first and second operand, stores in the third.
func executeAnd() {
	if !isImm {
		setMemory(dst, getMemory(src1)&getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)&currInst.ImmediateValue)
	}
}

// executeOr finds a bitwise Or of the first and second operand, stores in the third.
func executeOr() {
	if !isImm {
		setMemory(dst, getMemory(src1)|getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)|currInst.ImmediateValue)
	}
}

// executeNot finds a bitwise NOT of the first operand, stores in the second.
func executeNot() {
	if !isImm {
		setMemory(dst, negateValue(valueToBool(getMemory(src1))))
	} else {
		setMemory(dst, negateValue(valueToBool(currInst.ImmediateValue)))
	}
}

// executeXor finds a bitwise XOR of the first and second opcode, stores in the third.
func executeXor() {
	if !isImm {
		setMemory(dst, getMemory(src1)^getMemory(src2))
	} else {
		setMemory(dst, getMemory(src1)^currInst.ImmediateValue)
	}
}

// executePush pushes the first operand onto the stack.
func executePush() {
	if !isImm {
		Stack.Push(getMemory(src1))
	} else {
		Stack.Push(currInst.ImmediateValue)
	}
}

// executePop pops the top of the stack into the first operand (must be a variable)
func executePop() {
	setMemory(dst, Stack.Pop())
}

// executeCall starts a function call by saving the return stack, switching to the
// new context, and checking with the optimizer if the call is even necessary.
func executeCall() {
	// If optimizer found a cached call for this function,
	// then it already simulated the function. We leave
	if cacheFunctionCall(currInst.LabelImmediate) {
		return
	}
	// Save the return address
	ReturnStack.Push(PC)
	// Jump to the function start
	PC = currInst.LabelImmediate
	// Optimizer business, we need to keep track of this
	LastFunctionCalled = PC
	// Switch to a new memory context
	ContextNumber++
	// Calculate the new memory offset with the new context
	MemoryOffset = RegisterDepth(ContextNumber)*CONTEXT_SIZE - 1
	// See if we need to bump up the memory
	if int(MemoryOffset) >= len(Memory) {
		newMemory := make([]ValueWidth, len(Memory)*2)
		copy(newMemory, Memory)
		Memory = newMemory
	}
	// Call the optimizer to analyze the stack at function call
	optimizerAnalyzeCallStack()
}

// executeReturn returns from a function by jumping back into the old
// context and instruction that we left off at.
func executeReturn() {
	// Pop the last return point
	PC = ReturnStack.Pop()
	// Switch back the context
	ContextNumber--
	// Calculate the new offset
	MemoryOffset = RegisterDepth(ContextNumber)*CONTEXT_SIZE - 1
	// Let the optimizer analyze the stack with return values at return
	optimizerAnalyzeReturnStack()
}

// executeJump simply jumps to the given label's instruction
func executeJump() {
	PC = InstructionDepth(currInst.LabelImmediate)
}

// executeGreater computes if first operand is greater than the second one, stores 0/1 in third.
func executeGreater() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) > getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) > currInst.ImmediateValue))
	}
}

// executeGreaterEqual computes if first operand is greater or equal than the second one, stores 0/1 in third.
func executeGreaterEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) >= getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) >= currInst.ImmediateValue))
	}
}

// executeLess computes if first operand is less than the second one, stores 0/1 in third.
func executeLess() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) < getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) < currInst.ImmediateValue))
	}
}

// executeLessEqual computes if first operand is less or equal than the second one, stores 0/1 in third.
func executeLessEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) <= getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) <= currInst.ImmediateValue))
	}
}

// executeEqual computes if first operand is equal to the second one, stores 0/1 in third.
func executeEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) == getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) == currInst.ImmediateValue))
	}
}

// executeEqual computes if first operand is not equal to the second one, stores 0/1 in third.
func executeNotEqual() {
	if !isImm {
		setMemory(dst, boolToValue(getMemory(src1) != getMemory(src2)))
	} else {
		setMemory(dst, boolToValue(getMemory(src1) != currInst.ImmediateValue))
	}
}

// executeJumpIfTrue jumps to the label if the first operand is not zero.
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

// executeJumpIfFalse jumps to the label if the first operand is zero.
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

// executeHalt sets the halt value to the first operand and sends a signal
// to shut down the emulator, where it will show the value returned.
func executeHalt() {
	if !isImm {
		HaltValue = getMemory(src1)
	} else {
		HaltValue = currInst.ImmediateValue
	}
	Halted = true
}

// executeNop does nothing.
func executeNop() {}

// executeTODO is a development function.
func executeTODO() {
	panic("executed TODO")
}

// boolToValue promotes a bool to uint64.
func boolToValue(what bool) ValueWidth {
	if what {
		return 1
	}
	return 0
}

// valueToBool converts any non-zero value to 1.
func valueToBool(what ValueWidth) ValueWidth {
	if what != 0 {
		return 1
	}
	return 0
}

// negateValue converts any non-zero value to 0, and 0 to 1.
func negateValue(what ValueWidth) ValueWidth {
	if what == 0 {
		return 1
	}
	return 1
}

// setMemory saves the given value at a variable index (1-10) in the
// current context.
func setMemory(reg RegisterDepth, val ValueWidth) {
	Memory[MemoryOffset+reg] = val
}

// getMemory returns given variable index's (1-10) value from the current
// memory context.
func getMemory(reg RegisterDepth) ValueWidth {
	return Memory[MemoryOffset+reg]
}
