package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	instructionStringFormat = "Step %s must be finished before step %s can begin."
	malformedError          = "malformed input"
	noeEntrypointError      = "no entrypoint"
)

type instructionList map[string][]string

func parseInstructions(rawInstructions []string) (instructionList, error) {
	instructions := make(instructionList)
	for _, rawInstruction := range rawInstructions {
		var dependencyName, instructionName string
		numMatched, err := fmt.Sscanf(rawInstruction, instructionStringFormat, &dependencyName, &instructionName)
		if err != nil {
			return nil, err
		} else if numMatched != 2 {
			return nil, fmt.Errorf(malformedError)
		}
		_, hasDependency := instructions[dependencyName]
		if !hasDependency {
			instructions[dependencyName] = make([]string, 0)
		}
		_, hasInstruction := instructions[instructionName]
		if hasInstruction {
			instructions[instructionName] = append(instructions[instructionName], dependencyName)
		} else {
			instructions[instructionName] = []string{dependencyName}
		}
	}

	for instructionName := range instructions {
		sort.Strings(instructions[instructionName])
	}

	return instructions, nil
}

func (instructions instructionList) findReadySteps() []string {
	entrypointCandidates := make([]string, 0)
	for instructionName := range instructions {
		if len(instructions[instructionName]) == 0 {
			entrypointCandidates = append(entrypointCandidates, instructionName)
		}
	}

	sort.Strings(entrypointCandidates)

	return entrypointCandidates
}

func (instructions instructionList) markAsDone(doneInstructionName string) {
	for instructionName, instruction := range instructions {
		doneIndex := -1
		for i, dependencyName := range instruction {
			if dependencyName == doneInstructionName {
				doneIndex = i
				break
			}
		}
		if doneIndex != -1 {
			instructions[instructionName] = append(instruction[:doneIndex], instruction[doneIndex+1:]...)
		}
	}

	delete(instructions, doneInstructionName)
}

func (instructions instructionList) resolveDependencies() string {
	instructionSet := ""
	for len(instructions) > 0 {
		for _, readyStep := range instructions.findReadySteps() {
			instructions.markAsDone(readyStep)
			instructionSet += readyStep
		}
	}

	return instructionSet
}

func part1(instructions instructionList) string {
	return instructions.resolveDependencies()
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

	instructions, err := parseInstructions(rawInstructions)
	if err != nil {
		panic(err)
	}

	fmt.Println(part1(instructions))

}
