package main

type LabelType string

type LabelDepth uint16
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
	IsImmediate         bool
	ImmediateValue      ImmediateWidth
	LabelIndex          LabelDepth
	LabelImmediate      InstructionDepth

	//Input string
}
