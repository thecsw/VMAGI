package main

// LabelType denotes the label type, which is string.
type LabelType string

// LabelDepth shows how many labels we can allocate.
type LabelDepth uint16

// InstructionDepth shows how many instructions we can store.
type InstructionDepth uint32

// RegisterDepth shows how many variables we can store.
type RegisterDepth uint16

// ValueWidth shows how big our values are.
type ValueWidth int64

// Instruction is VMAGI's representation of an insntruction.
type Instruction struct {
	// Opcode is the iota enum opcode defined in `isa.go`.
	Opcode OpcodeNumber
	// NumberOperands shows how many operandns are expected.
	NumberOperands uint8
	// SourceRegister1 is the variable index in the context.
	SourceRegister1 RegisterDepth
	// SourceRegister2 is the variable index in the context.
	SourceRegister2 RegisterDepth
	// DestinationRegister is the destination variable's index in the context.
	DestinationRegister RegisterDepth
	// IsImmediate shows whether one of the operands is an immediate value.
	IsImmediate bool
	// ImmediateValue stores the immediate value passed in the instruction.
	ImmediateValue ValueWidth

	// LabelIndex is a temporary label index before label population happens.
	LabelIndex LabelDepth
	// LabelImmediate holds the insruction index if this instructions would
	// need to jump somewhere, like `call`, `jump`, etc.
	LabelImmediate InstructionDepth
}
