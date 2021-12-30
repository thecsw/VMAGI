package main

type OpcodeNumber int8

const (
	// Halt (return the value)
	HALT OpcodeNumber = iota
	// Math
	ADD
	SUB
	MUL
	DIV
	MOD
	NEG
	// Logic
	AND
	OR
	NOT
	XOR
	// Functions
	CALL
	JUMP
	RETURN
	// Comparisons
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ
	EQUAL
	NOT_EQUAL
	// Conditional jumps
	JUMPIF
	// Stack
	PUSH
	POP
	// Nothing
	NOP
)

var (
	OpcodeStrings = map[OpcodeNumber]string{
		HALT:       "halt",
		ADD:        "add",
		SUB:        "sub",
		MUL:        "mul",
		DIV:        "div",
		MOD:        "mod",
		NEG:        "neg",
		AND:        "and",
		OR:         "or",
		NOT:        "not",
		XOR:        "xor",
		CALL:       "call",
		JUMP:       "jmp",
		RETURN:     "ret",
		GREATER:    "gt",
		GREATER_EQ: "gte",
		LESS:       "lt",
		LESS_EQ:    "lte",
		EQUAL:      "eq",
		NOT_EQUAL:  "neq",
		JUMPIF:     "jmpz",
		PUSH:       "push",
		POP:        "pop",
		NOP:        "nop",
	}

	OpcodeNumOperands = map[OpcodeNumber]uint8{
		HALT:       1, // halt MEM or VAL
		ADD:        3, // add s1, s2, d
		SUB:        3, // sub s1, s2, d
		MUL:        3, // mul s1, s2, d
		DIV:        3, // div s1, s2, d
		MOD:        3, // mod s1, s2, d
		NEG:        2, // neg s1, d
		AND:        3, // and s1, s2, d
		OR:         3, // or  s1, s2, d
		NOT:        2, // not s1, s2, d
		XOR:        3, // xor s1, s2, d
		CALL:       1, // call LABEL
		JUMP:       1, // jmp LABEL
		RETURN:     0, // ret
		GREATER:    3, // gt  s1, s2, d
		GREATER_EQ: 3, // gte s1, s2, d
		LESS:       3, // lt  s1, s2, d
		LESS_EQ:    3, // lte s1, s2, d
		EQUAL:      3, // eq  s1, s2, d
		NOT_EQUAL:  3, // neq s1, s2, d
		JUMPIF:     2, // jmpz s1, LABEL (jump to LABEL if 0)
		PUSH:       1, // push VAL (s1)
		POP:        1, // pop MEM (d)
		NOP:        0, // nop
	}
)
