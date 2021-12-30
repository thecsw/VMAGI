package main

type Instruction struct {
	Opcode              OpcodeNumber
	NumRegs             uint8
	SourceRegister1     uint16
	SourceRegister2     uint16
	DestinationRegister uint16

	IsImmediate    bool
	ImmediateValue int32

	LabelImmediate string
}
