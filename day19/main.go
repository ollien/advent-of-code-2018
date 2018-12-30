package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	malformedInputError = "malformed input"
	ipFormat            = "#ip %d"
	instructionFormat   = "%s %d %d %d"
)

// Can't use a constant for a map - this is our next best thing
var operations = map[string]deviceFunction{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}

type registerSet [6]int
type deviceFunction = func(registerSet, int, int, int) registerSet
type instruction = func(registerSet) registerSet

func addr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return addi(registers, register1, registers[register2], destinationRegister)
}

func addi(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	registers[destinationRegister] = registers[register1] + value
	return registers
}

func mulr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return muli(registers, register1, registers[register2], destinationRegister)
}

func muli(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	registers[destinationRegister] = registers[register1] * value
	return registers
}

func banr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return bani(registers, register1, registers[register2], destinationRegister)
}

func bani(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	registers[destinationRegister] = registers[register1] & value
	return registers
}

func borr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return bori(registers, register1, registers[register2], destinationRegister)
}

func bori(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	registers[destinationRegister] = registers[register1] | value
	return registers
}

func setr(registers registerSet, register1 int, na int, destinationRegister int) registerSet {
	return seti(registers, registers[register1], na, destinationRegister)
}

func seti(registers registerSet, value int, na int, destinationRegister int) registerSet {
	registers[destinationRegister] = value
	return registers
}

// Not a valid opcode, but is helpful for the purposes of writing gtir, gtri, and gtrr
func gtii(registers registerSet, value1 int, value2 int, destinationRegister int) registerSet {
	output := 0
	if value1 > value2 {
		output = 1
	}
	registers[destinationRegister] = output
	return registers
}

func gtir(registers registerSet, value int, register1 int, destinationRegister int) registerSet {
	return gtii(registers, value, registers[register1], destinationRegister)
}

func gtri(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	return gtii(registers, registers[register1], value, destinationRegister)
}

func gtrr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return gtii(registers, registers[register1], registers[register2], destinationRegister)
}

// Not a valid opcode, but is helpful for the purposes of writing equir, eqri, eqrr
func eqii(registers registerSet, value1 int, value2 int, destinationRegister int) registerSet {
	output := 0
	if value1 == value2 {
		output = 1
	}
	registers[destinationRegister] = output
	return registers
}

func eqir(registers registerSet, value int, register1 int, destinationRegister int) registerSet {
	return eqii(registers, value, registers[register1], destinationRegister)
}

func eqri(registers registerSet, register1 int, value int, destinationRegister int) registerSet {
	return eqii(registers, registers[register1], value, destinationRegister)
}

func eqrr(registers registerSet, register1 int, register2 int, destinationRegister int) registerSet {
	return eqii(registers, registers[register1], registers[register2], destinationRegister)
}

func makeInstruction(deviceFunc deviceFunction, arg1, arg2, arg3 int) instruction {
	return func(registers registerSet) registerSet {
		return deviceFunc(registers, arg1, arg2, arg3)
	}
}

func parseInput(rawInstructions []string) (int, []instruction, error) {
	instructions := make([]instruction, len(rawInstructions)-1)
	var instructionPointerIndex int
	numMatched, err := fmt.Sscanf(rawInstructions[0], ipFormat, &instructionPointerIndex)
	if err != nil {
		return 0, nil, err
	} else if numMatched != 1 {
		return 0, nil, errors.New(malformedInputError)
	}

	for i, rawInstruction := range rawInstructions[1:] {
		var operationName string
		var arg1, arg2, arg3 int
		numMatched, err := fmt.Sscanf(rawInstruction, instructionFormat, &operationName, &arg1, &arg2, &arg3)
		if err != nil {
			return 0, nil, err
		} else if numMatched != 4 {
			return 0, nil, errors.New(malformedInputError)
		}

		operation, ok := operations[operationName]
		if !ok {
			return 0, nil, errors.New(malformedInputError)
		}

		instructions[i] = makeInstruction(operation, arg1, arg2, arg3)
	}

	return instructionPointerIndex, instructions, nil
}

func runElfcode(registers registerSet, instructionPointerIndex int, instructions []instruction, rawInstructions []string) int {
	for registers[instructionPointerIndex] < len(instructions) {
		fmt.Print(registers, " => ")
		instructionIndex := registers[instructionPointerIndex]
		registers = instructions[instructionIndex](registers)
		fmt.Println(rawInstructions[instructionIndex], "=>", registers)
		registers[instructionPointerIndex]++
	}

	return registers[0]
}

func main() {
	// Part 2 is solved by observing that the value in register 0 is simply the sum of the factors of the value that get put into register 2. In part 1 this is 1030, in part 2 this is 10550400
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./main in_file register_0_value")
		fmt.Println("Will print trace of program")
		return
	}

	inFile := os.Args[1]
	inFileContents, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}

	var registers registerSet
	registers[0], err = strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	rawInstructions := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawInstructions = rawInstructions[:len(rawInstructions)-1]
	instructionPointerIndex, instructions, err := parseInput(rawInstructions)
	if err != nil {
		panic(err)
	}
	fmt.Println(runElfcode(registers, instructionPointerIndex, instructions, rawInstructions[1:]))
}
