package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	ThreeRegisterRegex  = regexp.MustCompile(`\s*#(\d+)\s*,\s*#(\d+)\s*,\s*#(\d+)\s*`)
	ThreeImmediateRegex = regexp.MustCompile(`i \s*#(\d+)\s*,\s*(\d+)\s*,\s*#(\d+)\s*`)

	TwoRegisterRegex  = regexp.MustCompile(`\s*#(\d+)\s*,\s*#(\d+)\s*`)
	TwoImmediateRegex = regexp.MustCompile(`i \s*#(\d+)\s*,\s*(\d+)\s*`)

	OneRegisterRegex  = regexp.MustCompile(`\s*#(\d+)\s*`)
	OneImmediateRegex = regexp.MustCompile(`i \s*(\d+)\s*`)

	LabelDeclarationRegex = regexp.MustCompile(`\s*([^:]+):\s*`)
	LabelImmediateRegex   = regexp.MustCompile(`\s*@([^:]+)\s*`)

	ConditionalJumpRegisterRegex  = regexp.MustCompile(`\s*#(\d+)\s*,\s*@([^:]+)\s*`)
	ConditionalJumpImmediateRegex = regexp.MustCompile(`i \s*(\d+)\s*,\s*@([^:]+)\s*`)

	CommentRegex = regexp.MustCompile(`--.+`)
)

// ParseInput takes lines of assembly and returns instructions
func ParseInput(lines []string) []*Instruction {
	instructions := make([]*Instruction, 0)

	for _, line := range lines {
		// Clean from comments
		line = CommentRegex.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		// Parse instruction and add it
		instructions = append(instructions, ParseInstruction(line))
		// Try to extract labels if there are any
		matches := LabelDeclarationRegex.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			Labels[matches[0][1]] = uint32(len(instructions) - 1)
		}
	}

	return instructions
}

// ParseInstruction gets a line and returns an instruction
func ParseInstruction(line string) *Instruction {
	for i := OpcodeNumber(0); i <= NOP; i++ {
		// Specific instructions with 0 operands
		if strings.HasSuffix(line, "nop") {
			return formInstruction(line, OpcodeNumOperands[NOP], NOP)
		}
		if strings.HasSuffix(line, "ret") {
			return formInstruction(line, OpcodeNumOperands[RETURN], RETURN)
		}
		if strings.Contains(line, OpcodeStrings[i]+" ") ||
			strings.Contains(line, OpcodeStrings[i]+"i ") {
			return formInstruction(line, OpcodeNumOperands[i], i)
		}
	}
	panic("illegal instruction while parsing")
}

func formInstruction(line string, numRegs uint8, opcode OpcodeNumber) (what *Instruction) {
	what = &Instruction{
		Opcode:      opcode,
		NumRegs:     numRegs,
		IsImmediate: strings.Contains(line, "i "),
	}

	//fmt.Println(line, " -- ", numRegs, " -- ", opcode)

	if what.NumRegs == 3 {
		if !what.IsImmediate {
			matches := ThreeRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("three register regex match failed")
			}
			registers := matches[0]
			what.SourceRegister1 = strToRegNum(registers[1])
			what.SourceRegister2 = strToRegNum(registers[2])
			what.DestinationRegister = strToRegNum(registers[3])
		}
		if what.IsImmediate {
			matches := ThreeImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("three immediate regex match failed")
			}
			values := matches[0]
			what.SourceRegister1 = strToRegNum(values[1])
			what.ImmediateValue = strToInt32(values[2])
			what.DestinationRegister = strToRegNum(values[3])
		}
		return
	}

	if what.NumRegs == 2 {
		// See if we're doing the JUMPIF
		if what.Opcode == JUMPIF {
			if !what.IsImmediate {
				matches := ConditionalJumpRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("conditional jump register regex match failed")
				}
				values := matches[0]
				what.SourceRegister1 = strToRegNum(values[1])
				what.LabelImmediate = values[2]
			}
			if what.IsImmediate {
				matches := ConditionalJumpImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("conditional jump immediate regex match failed")
				}
				values := matches[0]
				what.ImmediateValue = strToInt32(values[1])
				what.LabelImmediate = values[2]
			}
			return
		}

		if !what.IsImmediate {
			matches := TwoRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("two register regex match failed")
			}
			registers := matches[0]
			what.SourceRegister1 = strToRegNum(registers[1])
			what.DestinationRegister = strToRegNum(registers[2])
		}
		if what.IsImmediate {
			matches := TwoImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("two immediate regex match failed")
			}
			values := matches[0]
			what.ImmediateValue = strToInt32(values[1])
			what.DestinationRegister = strToRegNum(values[2])
		}
		return
	}

	if what.NumRegs == 1 {
		// Call or jump
		if what.Opcode == CALL || what.Opcode == JUMP {
			matches := LabelImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("call regex match failed")
			}
			what.LabelImmediate = matches[0][1]
			return
		}
		// Push
		if what.Opcode == PUSH {
			if what.IsImmediate {
				matches := OneImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("push immediate regex match failed")
				}
				what.ImmediateValue = strToInt32(matches[0][1])
			}
			if !what.IsImmediate {
				matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("push register regex match failed")
				}
				what.SourceRegister1 = strToRegNum(matches[0][1])
			}
			return
		}
		// Pop
		if what.Opcode == POP {
			matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("pop register regex match failed")
			}
			what.DestinationRegister = strToRegNum(matches[0][1])
			return
		}
		// Halt
		if what.Opcode == HALT {
			if what.IsImmediate {
				matches := OneImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("halt immediate regex match failed")
				}
				what.ImmediateValue = strToInt32(matches[0][1])
			}
			if !what.IsImmediate {
				matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("halt register regex match failed")
				}
				what.SourceRegister1 = strToRegNum(matches[0][1])
			}
			return
		}
	}
	return
}

var dummyInt = 0

func strToRegNum(register string) uint16 {
	dummyInt, _ = strconv.Atoi(register)
	return uint16(dummyInt)
}

func strToInt32(what string) int32 {
	dummyInt, _ = strconv.Atoi(what)
	return int32(dummyInt)
}
