package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	malformedInputError = "malformed input"
	beforeFormat        = "Before: [%d, %d, %d, %d]"
	instructionFormat   = "%d %d %d %d"
	afterFormat         = "After: [%d, %d, %d, %d]"
)

// Can't use a constant for an array - this is our next best thing
var operations = [...]deviceFunction{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}

type registerSet [4]int
type instruction [4]int
type deviceFunction = func(registerSet, int, int, int) registerSet

type note struct {
	before registerSet
	input  instruction
	after  registerSet
}

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

// Gets the indices in the operations array of the operations that match the note
func (n note) getMatchingOperations() []int {
	matchingOperations := []int{}
	for i, operation := range operations {
		if operation(n.before, n.input[1], n.input[2], n.input[3]) == n.after {
			matchingOperations = append(matchingOperations, i)
		}
	}

	return matchingOperations
}

func parseInput(input []string) ([]note, []instruction, error) {
	notes, rawPart2Input, err := parsePart1Input(input)
	if err != nil {
		return nil, nil, err
	}
	instructions, err := parsePart2Input(rawPart2Input)
	if err != nil {
		return nil, nil, err
	}

	return notes, instructions, nil
}

func parsePart1Input(rawNotes []string) ([]note, []string, error) {
	lastLineLength := -1
	notes := []note{}
	currentNote := note{}
	lastIndex := -1
	for i, line := range rawNotes {
		lastIndex = i
		if lastLineLength == 0 && len(line) == 0 {
			break
		}
		lastLineLength = len(line)
		if lastLineLength == 0 {
			continue
		}

		var registers registerSet
		var ins instruction

		numMatched, err := fmt.Sscanf(line, beforeFormat, &registers[0], &registers[1], &registers[2], &registers[3])
		if err == nil && numMatched != 4 {
			return nil, nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.before = registers
			continue
		}

		numMatched, err = fmt.Sscanf(line, instructionFormat, &ins[0], &ins[1], &ins[2], &ins[3])
		if err == nil && numMatched != 4 {
			return nil, nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.input = ins
			continue
		}

		numMatched, err = fmt.Sscanf(line, afterFormat, &registers[0], &registers[1], &registers[2], &registers[3])
		if err == nil && numMatched != 4 {
			return nil, nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.after = registers
			notes = append(notes, currentNote)
		} else {
			// If we have an error at this point, something is actually wrong.
			return nil, nil, err
		}
	}

	return notes, rawNotes[lastIndex+2:], nil
}

func parsePart2Input(rawInstructions []string) ([]instruction, error) {
	instructions := make([]instruction, 0, len(rawInstructions))
	for _, line := range rawInstructions {
		var ins instruction
		numMatched, err := fmt.Sscanf(line, instructionFormat, &ins[0], &ins[1], &ins[2], &ins[3])
		if err != nil {
			return nil, err
		} else if numMatched != 4 {
			return nil, errors.New(malformedInputError)
		}
		instructions = append(instructions, ins)
	}

	return instructions, nil
}

func intersection(arr1 []int, arr2 []int) []int {
	resultSlice := []int{}
	for _, item1 := range arr1 {
		found := false
		for _, item2 := range arr2 {
			if item1 == item2 {
				found = true
				break
			}
		}
		if found {
			resultSlice = append(resultSlice, item1)
		}
	}

	return resultSlice
}

func allListsEmpty(itemMap map[int][]int) bool {
	for _, itemList := range itemMap {
		if len(itemList) != 0 {
			return false
		}
	}

	return true
}

func getOneItemList(itemMap map[int][]int) int {
	for k, itemList := range itemMap {
		if len(itemList) == 1 {
			return k
		}
	}

	return -1
}

func removeItemFromAllLists(itemMap map[int][]int, target int) {
	for k, itemList := range itemMap {
		index := -1
		for i, item := range itemList {
			if item == target {
				index = i
				break
			}
		}
		if index != -1 {
			itemMap[k] = append(itemMap[k][:index], itemMap[k][index+1:]...)
		}
	}
}

func getOpcodesFromPossibilities(possibleOpcodes map[int][]int) map[int]deviceFunction {
	opcodes := make(map[int]deviceFunction, len(possibleOpcodes))
	for !allListsEmpty(possibleOpcodes) {
		// Find an opcode with only one function as a possibility - that way we know for certain this opcode maps to this function.
		nextOpcode := getOneItemList(possibleOpcodes)
		functionID := possibleOpcodes[nextOpcode][0]
		opcodes[nextOpcode] = operations[functionID]
		// Remove all other instances of this item from other lists, thus creating another one function mapping.
		removeItemFromAllLists(possibleOpcodes, functionID)
	}

	return opcodes
}

// part1 runs the part1 simulation, and also gives us the possibilities for opcodes, a needed part for part2
func part1(notes []note) (int, map[int][]int) {
	opcodeMappings := map[int][]int{}
	total := 0
	for _, note := range notes {
		matchingOperations := note.getMatchingOperations()
		if len(matchingOperations) >= 3 {
			total++
		}
		opcode := note.input[0]
		// If an opcode is already stored, we only need to remove the opcodes that weren't already there
		if _, haveOpcode := opcodeMappings[opcode]; haveOpcode {
			opcodeMappings[opcode] = matchingOperations
			opcodeMappings[opcode] = intersection(matchingOperations, opcodeMappings[opcode])
		} else {
			opcodeMappings[opcode] = matchingOperations
		}
	}

	return total, opcodeMappings
}

func part2(instructions []instruction, opcodes map[int]deviceFunction) int {
	var registers registerSet
	for _, instruction := range instructions {
		opcode := instruction[0]
		registers = opcodes[opcode](registers, instruction[1], instruction[2], instruction[3])
	}

	return registers[0]
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
	rawNotes := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawNotes = rawNotes[:len(rawNotes)-1]
	part1Notes, part2Instructions, err := parseInput(rawNotes)
	if err != nil {
		panic(err)
	}

	part1Result, opcodePossibilities := part1(part1Notes)
	fmt.Println(part1Result)
	opcodes := getOpcodesFromPossibilities(opcodePossibilities)
	fmt.Println(part2(part2Instructions, opcodes))
}
