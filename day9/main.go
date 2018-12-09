package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	lineFormat          = "%d players; last marble is worth %d points"
	malformedInputError = "malformed input"
)

func parseInput(input string) (numPlayers, numMarbles int, err error) {
	numMatched, err := fmt.Sscanf(input, lineFormat, &numPlayers, &numMarbles)
	if err != nil {
		return 0, 0, err
	} else if numMatched != 2 {
		return 0, 0, fmt.Errorf(malformedInputError)
	}

	return
}

func getMax(arr []int) (max int) {
	for _, item := range arr {
		if item > max {
			max = item
		}
	}

	return
}

func runGame(numPlayers int, numMarbles int) int {
	scores := make([]int, numPlayers)
	board := make([]int, 1, numMarbles)
	currentMarbleIndex := 0
	currentPlayer := 0
	for nextMarble := 1; nextMarble <= numMarbles; nextMarble++ {
		if nextMarble%23 == 0 {
			removeIndex := (currentMarbleIndex - 7 + len(board)) % len(board)
			scores[currentPlayer] += board[removeIndex] + nextMarble
			board = append(board[:removeIndex], board[removeIndex+1:]...)
			currentMarbleIndex = removeIndex % len(board)
		} else {
			insertIndex := (currentMarbleIndex + 2) % len(board)
			if len(board) == 1 {
				insertIndex = 1
			} else if currentMarbleIndex+2 == len(board) {
				insertIndex = len(board)
			}
			tailElements := append([]int{nextMarble}, board[insertIndex:]...)
			board = append(board[:insertIndex], tailElements...)
			currentMarbleIndex = insertIndex
		}
		currentPlayer = (currentPlayer + 1) % numPlayers
	}

	return getMax(scores)
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

	input := strings.TrimSuffix(string(inFileContents), "\n")
	numPlayers, numMarbles, err := parseInput(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(runGame(numPlayers, numMarbles))
}
