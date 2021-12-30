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

	LabelMap    = map[string]LabelDepth{}
	Labels      = make([]InstructionDepth, 100)
	FoundLabels = 1
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
			addNewLabel(matches[0][1])
			setLabel(matches[0][1], InstructionDepth(len(instructions)-1))
		}
	}
	// Fill out labelImmediates

	for i, v := range instructions {
		if v.LabelIndex > 0 {
			instructions[i].LabelImmediate = Labels[v.LabelIndex]
		}
	}

	//	litter.Dump(LabelMap)
	// litter.Dump(Labels)
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
		Opcode:         opcode,
		NumberOperands: numRegs,
		IsImmediate:    strings.Contains(line, "i "),
		Input:          line,
	}

	//fmt.Println(line, " -- ", numRegs, " -- ", opcode)

	if what.NumberOperands == 3 {
		if !what.IsImmediate {
			matches := ThreeRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("three register regex match failed")
			}
			registers := matches[0]
			what.SourceRegister1 = stringToRegister(registers[1])
			what.SourceRegister2 = stringToRegister(registers[2])
			what.DestinationRegister = stringToRegister(registers[3])
		}
		if what.IsImmediate {
			matches := ThreeImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("three immediate regex match failed")
			}
			values := matches[0]
			what.SourceRegister1 = stringToRegister(values[1])
			what.ImmediateValue = stringToImmediate(values[2])
			what.DestinationRegister = stringToRegister(values[3])
		}
		return
	}

	if what.NumberOperands == 2 {
		// See if we're doing the JUMPIF
		if what.Opcode == JUMPF || what.Opcode == JUMPT {
			if !what.IsImmediate {
				matches := ConditionalJumpRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("conditional jump register regex match failed")
				}
				values := matches[0]
				what.SourceRegister1 = stringToRegister(values[1])
				what.LabelIndex = addNewLabel(values[2])
				//what.LabelImmediate = labelToInstruction(values[2])
			}
			if what.IsImmediate {
				matches := ConditionalJumpImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("conditional jump immediate regex match failed")
				}
				values := matches[0]
				what.ImmediateValue = stringToImmediate(values[1])
				what.LabelIndex = addNewLabel(values[2])
				//what.LabelImmediate = labelToInstruction(values[2])
			}
			return
		}

		if !what.IsImmediate {
			matches := TwoRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("two register regex match failed")
			}
			registers := matches[0]
			what.SourceRegister1 = stringToRegister(registers[1])
			what.DestinationRegister = stringToRegister(registers[2])
		}
		if what.IsImmediate {
			matches := TwoImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("two immediate regex match failed")
			}
			values := matches[0]
			what.ImmediateValue = stringToImmediate(values[1])
			what.DestinationRegister = stringToRegister(values[2])
		}
		return
	}

	if what.NumberOperands == 1 {
		// Call or jump
		if what.Opcode == CALL || what.Opcode == JUMP {
			matches := LabelImmediateRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("call regex match failed")
			}
			what.LabelIndex = addNewLabel(matches[0][1])
			//what.LabelImmediate = addNewLabel(matches[0][1])
			return
		}
		// Push
		if what.Opcode == PUSH {
			if what.IsImmediate {
				matches := OneImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("push immediate regex match failed")
				}
				what.ImmediateValue = stringToImmediate(matches[0][1])
			}
			if !what.IsImmediate {
				matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("push register regex match failed")
				}
				what.SourceRegister1 = stringToRegister(matches[0][1])
			}
			return
		}
		// Pop
		if what.Opcode == POP {
			matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
			if len(matches) < 1 {
				panic("pop register regex match failed")
			}
			what.DestinationRegister = stringToRegister(matches[0][1])
			return
		}
		// Halt
		if what.Opcode == HALT {
			if what.IsImmediate {
				matches := OneImmediateRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("halt immediate regex match failed")
				}
				what.ImmediateValue = stringToImmediate(matches[0][1])
			}
			if !what.IsImmediate {
				matches := OneRegisterRegex.FindAllStringSubmatch(line, -1)
				if len(matches) < 1 {
					panic("halt register regex match failed")
				}
				what.SourceRegister1 = stringToRegister(matches[0][1])
			}
			return
		}
	}
	return
}

var dummyInt = 0

func stringToRegister(register string) RegisterDepth {
	dummyInt, _ = strconv.Atoi(register)
	return RegisterDepth(dummyInt)
}

func stringToImmediate(what string) ImmediateWidth {
	dummyInt, _ = strconv.Atoi(what)
	return ImmediateWidth(dummyInt)
}

func addNewLabel(what string) LabelDepth {
	if v, ok := LabelMap[what]; !ok {
		LabelMap[what] = LabelDepth(FoundLabels)
		FoundLabels++
		return LabelDepth(FoundLabels) - 1
	} else {
		return v
	}
}

func setLabel(what string, instruction InstructionDepth) {
	Labels[LabelMap[what]] = instruction
}

func labelToInstruction(what string) InstructionDepth {
	return Labels[LabelMap[what]]
}
