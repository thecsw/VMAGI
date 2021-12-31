package main

type OpcodeNumber int8

const (
	// HALT halts the VM with the value given: `halt (#VAR|VAL)`.
	HALT OpcodeNumber = iota
	// ADD adds two numbers: `add #VAR, #VAR, #VAR` or `addi #VAR, VAL, #VAR`.
	ADD
	// SUB subtracts two numbers: `sub #VAR, #VAR, #VAR` or `subi #VAR, VAL, #VAR`.
	SUB
	// MUL multiplies two numbers: `mul #VAR, #VAR, #VAR` or `muli #VAR, VAL, #VAR`.
	MUL
	// DIV divides two numbers: `div #VAR, #VAR, #VAR` or `divi #VAR, VAL, #VAR`.
	DIV
	// MOD mods two numbers: `mod #VAR, #VAR, #VAR` or `modi #VAR, VAL, #VAR`.
	MOD
	// NEG negates a number: `neg #VAR, #VAR` or `negi VAL, #VAR`.
	NEG
	// AND ands two values: `and #VAR, #VAR, #VAR` or `andi #VAR, VAL, #VAR`.
	AND
	// OR ors two values: `or #VAR, #VAR, #VAR` or `ori #VAR, VAL, #VAR`.
	OR
	// NOT nots a value: `not #VAR, #VAR` or `noti VAL, #VAR`.
	NOT
	// XOR xors two values: `xor #VAR, #VAR, #VAR` or `xori #VAR, VAL, #VAR`.
	XOR
	// CALL calls a function: `call @label`.
	CALL
	// JUMP jumps to a label: `jmp @label`.
	JUMP
	// RETURN returns from a function: `ret`.
	RETURN
	// GREATER compares: `gt #VAR, #VAR, #VAR` or `gti #VAR, VAL, #VAR`.
	GREATER
	// GREATER_EQ compares: `gte #VAR, #VAR, #VAR` or `gtei #VAR, VAL, #VAR`.
	GREATER_EQ
	// LESS compares: `lt #VAR, #VAR, #VAR` or `lti #VAR, VAL, #VAR`.
	LESS
	// LESS_EQ compares: `lte #VAR, #VAR, #VAR` or `ltei #VAR, VAL, #VAR`.
	LESS_EQ
	// EQUAL checks if equal: `eq #VAR, #VAR, #VAR` or `eqi #VAR, VAL, #VAR`.
	EQUAL
	// EQUAL checks if not: `eq #VAR, #VAR, #VAR` or `eqi #VAR, VAL, #VAR`.
	NOT_EQUAL
	// JUMPT jumps if val is 1: `jmpt #VAR, @label` or `jmpti VAL, @label`.
	JUMPT
	// JUMPF jumps if val is 0: `jmpz #VAR, @label` or `jmpzi VAL, @label`.
	JUMPF
	// PUSH pushes to the stack: `push #VAR` or `pushi VAL`.
	PUSH
	// POP pops from the stack: `pop #VAR`.
	POP
	// NOP does nothing: `nop`.
	NOP
)

var (
	// OpcodeStrings maps opcodes to their string representations in assembly.
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
		JUMPT:      "jmpt",
		JUMPF:      "jmpz",
		PUSH:       "push",
		POP:        "pop",
		NOP:        "nop",
	}

	// OpcodeNumOperands maps the opcodes to number of operands expected by them.
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
		JUMPT:      2, // jmmp s1, LABEL (jump to LABEL if 1)
		JUMPF:      2, // jmpz s1, LABEL (jump to LABEL if 0)
		PUSH:       1, // push VAL (s1)
		POP:        1, // pop MEM (d)
		NOP:        0, // nop
	}
)
