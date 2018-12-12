package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	malformedInputError = "malformed input"
	initialStateDelim   = ": "
	stateDelim          = " => "
	liveChar            = '#'
	deadChar            = '.'
	numSteps            = 20
)

func parseInput(inputLines []string) (string, map[string]bool, error) {
	states := make(map[string]bool, len(inputLines)-2)
	initialStateComponents := strings.Split(inputLines[0], initialStateDelim)
	if len(initialStateComponents) != 2 {
		return "", nil, fmt.Errorf(malformedInputError)
	}

	initialState := initialStateComponents[1]

	for _, line := range inputLines[2:] {
		lineComponents := strings.Split(line, stateDelim)
		if len(lineComponents) != 2 {
			return "", nil, fmt.Errorf(malformedInputError)
		}
		resultChar := lineComponents[1][0]
		if resultChar == liveChar || resultChar == deadChar {
			states[lineComponents[0]] = (resultChar == liveChar)
		} else {
			return "", nil, fmt.Errorf(malformedInputError)
		}
	}

	return initialState, states, nil
}

// runStep calculates the new step and returns the new string and the number of pots added to the left
func runStep(state string, states map[string]bool) (string, int) {
	// Add .... to either side - this ensures we cover the case of ....# and #....
	parseState := "...." + state + "...."
	finalState := bytes.NewBufferString("")
	firstPotPos := -1
	for i := range parseState[2 : len(parseState)-2] {
		parseComponent := parseState[i : i+5]
		if states[parseComponent] {
			finalState.WriteRune(liveChar)
			if firstPotPos == -1 {
				firstPotPos = i + 2
			}
		} else if firstPotPos != -1 {
			finalState.WriteRune(deadChar)
		}
	}

	trimmedFinalState := strings.TrimRight(finalState.String(), ".")
	return trimmedFinalState, 4 - firstPotPos
}

func part1(initialState string, states map[string]bool) int {
	currentState := initialState
	leftPotPos := 0
	for i := 0; i < numSteps; i++ {
		var leftPots int
		currentState, leftPots = runStep(currentState, states)
		leftPotPos += leftPots
	}

	total := 0
	for i, char := range currentState {
		if char == liveChar {
			total += i - leftPotPos
		}
	}

	return total
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
	inputLines := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	inputLines = inputLines[:len(inputLines)-1]
	initialState, states, err := parseInput(inputLines)
	if err != nil {
		panic(err)
	}
	fmt.Println(initialState)
	fmt.Println(part1(initialState, states))
}
