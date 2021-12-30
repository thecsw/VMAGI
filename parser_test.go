package main

import (
	"reflect"
	"testing"
)

func TestParseInstruction(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want *Instruction
	}{
		{
			"three register add",
			args{"add #1, #2, #3"},
			&Instruction{
				Opcode:              ADD,
				NumberOperands:      3,
				SourceRegister1:     1,
				SourceRegister2:     2,
				DestinationRegister: 3,
				IsImmediate:         false,
				ImmediateValue:      0,
			},
		},
		{
			"immediate add",
			args{"addi #1, 2, #3"},
			&Instruction{
				Opcode:              ADD,
				NumberOperands:      2,
				SourceRegister1:     1,
				SourceRegister2:     0,
				DestinationRegister: 3,
				IsImmediate:         true,
				ImmediateValue:      2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInstruction(tt.args.line); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("ParseInstruction() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
