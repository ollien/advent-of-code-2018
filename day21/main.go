package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	ipFormat          = "#ip %d"
	instructionFormat = "%s %d %d %d"
)

var errMalformedInput = errors.New("malformed input")
var errNoResult = errors.New("no result")

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
		return 0, nil, errMalformedInput
	}

	for i, rawInstruction := range rawInstructions[1:] {
		var operationName string
		var arg1, arg2, arg3 int
		numMatched, err := fmt.Sscanf(rawInstruction, instructionFormat, &operationName, &arg1, &arg2, &arg3)
		if err != nil {
			return 0, nil, err
		} else if numMatched != 4 {
			return 0, nil, errMalformedInput
		}

		operation, ok := operations[operationName]
		if !ok {
			return 0, nil, errMalformedInput
		}

		instructions[i] = makeInstruction(operation, arg1, arg2, arg3)
	}

	return instructionPointerIndex, instructions, nil
}

// runAsDebug naively runs the elfcode program and prints a trace - is not used for solution but was used for debugging.
// It is left in to help future readers determine how the puzzle was solved
func runAsDebug(registers registerSet, instructionPointerIndex int, instructions []instruction, rawInstructions []string) int {
	for registers[instructionPointerIndex] < len(instructions) {
		fmt.Print(registers, " => ")
		instructionIndex := registers[instructionPointerIndex]
		registers = instructions[instructionIndex](registers)
		fmt.Println(instructionIndex, "-", rawInstructions[instructionIndex], "=>", registers)
		registers[instructionPointerIndex]++
	}

	return registers[0]
}

// Run the program, taking the registers, the instruction pointer index, the instructions to execute, their unparsed equivalent, and a callback that will return a solution to the puzzle
func run(registers registerSet, instructionPointerIndex int, instructions []instruction, rawInstructions []string, findAnswer func(currentRegisters registerSet, rawInstrunction string, numInstructionsRun int) (int, error)) (int, error) {
	n := 0
	for registers[instructionPointerIndex] < len(instructions) {
		instructionIndex := registers[instructionPointerIndex]
		registers = instructions[instructionIndex](registers)
		registers[instructionPointerIndex]++
		n++
		result, err := findAnswer(registers, rawInstructions[instructionIndex], n)
		if err == nil {
			return result, nil
		} else if err != nil && !errors.Is(err, errNoResult) {
			return 0, fmt.Errorf("Could not find solution: %w", err)
		}
	}

	return 0, errNoResult
}

// Run the program, taking the registers, the instruction pointer index, the instructions to execute, their unparsed equivalent, and a callback that will return a solution to the puzzle
func getNumInstructionsRun(registers registerSet, instructionPointerIndex int, instructions []instruction) int {
	n := 0
	for registers[instructionPointerIndex] < len(instructions) {
		instructionIndex := registers[instructionPointerIndex]
		registers = instructions[instructionIndex](registers)
		registers[instructionPointerIndex]++
		n++
	}

	return n
}

func part1(registers registerSet, instructionPointerIndex int, instructions []instruction, rawInstructions []string) int {
	res, err := run(registers, instructionPointerIndex, instructions, rawInstructions, func(currentRegisters registerSet, currentRawInstruction string, _ int) (int, error) {
		components := strings.Split(currentRawInstruction, " ")
		// The input has an "eqrr" instruction that compares to register 0. Find it.
		if components[0] != "eqrr" {
			return 0, errNoResult
		}

		for _, component := range components[1:3] {
			if component == "0" {
				return currentRegisters[5], nil
			}
		}

		return 0, errNoResult
	})

	if err != nil {
		panic("Could not find solution to part 1 " + err.Error())
	}

	return res
}

func part2(registers registerSet, instructionPointerIndex int, instructions []instruction, rawInstructions []string) int {
	// This is awful and takes about five minutes
	// Eliminating the split would probably speed it up, but I'd either need to refactor this entire program (since I can't compare functions) or go based on instruction IDs, which is brittle.
	seen := map[int]int{}
	_, err := run(registers, instructionPointerIndex, instructions, rawInstructions, func(currentRegisters registerSet, currentRawInstruction string, numInstructionsRun int) (int, error) {
		components := strings.Split(currentRawInstruction, " ")
		if components[0] != "eqrr" {
			return 0, errNoResult
		}

		valid := false
		for _, component := range components[1:3] {
			if component == "0" {
				valid = true
			}
		}

		if !valid {
			return 0, errNoResult
		}

		val := currentRegisters[5]
		if _, ok := seen[val]; ok {
			fmt.Println("Found duplicate")
			return -1, nil
		}

		seen[val] = numInstructionsRun

		return 0, errNoResult
	})

	takenMin := false
	min := 0
	bestValue := 0
	for value, count := range seen {
		if !takenMin || (count > min || (count == min && value < bestValue)) {
			takenMin = true
			min = count
			bestValue = value
		}
	}

	if err != nil {
		panic("Could not find solution to part 2 " + err.Error())
	}

	return bestValue
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main in_file")
		return
	}

	inFile := os.Args[1]
	inFileContents, err := ioutil.ReadFile(inFile)
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

	var registers registerSet
	fmt.Println(part1(registers, instructionPointerIndex, instructions, rawInstructions[1:]))
	fmt.Println(part2(registers, instructionPointerIndex, instructions, rawInstructions[1:]))
}
