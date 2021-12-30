package main

type LabelType string

type ContextDepth uint16
type InstructionDepth uint32
type RegisterDepth uint16

type ImmediateWidth uint32
type ValueWidth int64

type Instruction struct {
	Opcode              OpcodeNumber
	NumberOperands      uint8
	SourceRegister1     RegisterDepth
	SourceRegister2     RegisterDepth
	DestinationRegister RegisterDepth

	IsImmediate    bool
	ImmediateValue ImmediateWidth

	LabelImmediate LabelType

	//	Input string
}
