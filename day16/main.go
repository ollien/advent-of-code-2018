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
	return bori(registers, register1, registers[register2], destinationRegister)
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

func (n note) getNumMatchingOperations() (count int) {
	operations := []deviceFunction{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
	for _, operation := range operations {
		if operation(n.before, n.input[1], n.input[2], n.input[3]) == n.after {
			count++
		}
	}
	return
}

func parsePart1Input(rawNotes []string) ([]note, error) {
	lastLineLength := -1
	notes := []note{}
	currentNote := note{}
	for _, line := range rawNotes {
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
			return nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.before = registers
			continue
		}

		numMatched, err = fmt.Sscanf(line, instructionFormat, &ins[0], &ins[1], &ins[2], &ins[3])
		if err == nil && numMatched != 4 {
			return nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.input = ins
			continue
		}

		numMatched, err = fmt.Sscanf(line, afterFormat, &registers[0], &registers[1], &registers[2], &registers[3])
		if err == nil && numMatched != 4 {
			return nil, errors.New(malformedInputError)
		} else if err == nil {
			currentNote.after = registers
			notes = append(notes, currentNote)
		} else {
			// If we have an error at this point, something is actually wrong.
			return nil, err
		}
	}

	return notes, nil
}

func part1(notes []note) (total int) {
	for _, note := range notes {
		numOps := note.getNumMatchingOperations()
		if numOps >= 3 {
			total++
		}
	}
	return
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
	part1Notes, err := parsePart1Input(rawNotes)
	if err != nil {
		panic(err)
	}
	fmt.Println(part1(part1Notes))
}
