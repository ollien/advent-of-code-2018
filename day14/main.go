package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const (
	usageString = "Usage: ./main number_of_scores partNum[,partNum]"
	score1      = 3
	score2      = 7
)

func calculateNewScores(scores []int, elf1Cursor int, elf2Cursor int) []int {
	total := scores[elf1Cursor] + scores[elf2Cursor]
	newScores := make([]int, 0, 2)
	if total == 0 {
		newScores = append(newScores, 0)
	} else {
		for total != 0 {
			newScores = append([]int{total % 10}, newScores...)
			total /= 10
		}
	}

	return newScores
}

func makeStringOfIntSlice(slice []int) string {
	buffer := bytes.NewBufferString("")
	for _, item := range slice {
		buffer.WriteString(strconv.Itoa(item))
	}

	return buffer.String()
}

func makeSolutionSlice(solutionString string) []int {
	solution := make([]int, len(solutionString))
	for i, char := range solutionString {
		solution[i] = int(char - '0')
	}

	return solution
}

func compareIntSlices(a []int, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func part1(numScores int) string {
	scores := make([]int, 2, numScores+10)
	scores[0] = score1
	scores[1] = score2
	elf1Cursor := 0
	elf2Cursor := 1
	for len(scores) < numScores+10 {
		scores = append(scores, calculateNewScores(scores, elf1Cursor, elf2Cursor)...)
		elf1Score := scores[elf1Cursor]
		elf2Score := scores[elf2Cursor]
		elf1Cursor = (elf1Cursor + elf1Score + 1) % len(scores)
		elf2Cursor = (elf2Cursor + elf2Score + 1) % len(scores)
	}

	resultString := ""
	for i := numScores; i < numScores+10; i++ {
		resultString += strconv.Itoa(scores[i])
	}

	return makeStringOfIntSlice(scores[numScores : numScores+10])
}

func part2(scoreString string) int {
	scores := make([]int, 2)
	scores[0] = score1
	scores[1] = score2
	elf1Cursor := 0
	elf2Cursor := 1
	itemIndex := -1
	solutionSlice := makeSolutionSlice(scoreString)
	for itemIndex == -1 {
		newScores := calculateNewScores(scores, elf1Cursor, elf2Cursor)
		// We must loop through each new score in order to make sure that we handle the cases of double digits
		for _, newScore := range newScores {
			scores = append(scores, newScore)
			if len(scores) >= len(solutionSlice) {
				// If the last len(solutionSlice) digits are the same sa our solution slice, we're done
				if compareIntSlices(scores[len(scores)-len(solutionSlice):], solutionSlice) {
					itemIndex = len(scores) - len(solutionSlice)
					break
				}
			}
		}
		elf1Score := scores[elf1Cursor]
		elf2Score := scores[elf2Cursor]
		elf1Cursor = (elf1Cursor + elf1Score + 1) % len(scores)
		elf2Cursor = (elf2Cursor + elf2Score + 1) % len(scores)
	}

	return itemIndex
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(usageString)
		return
	}

	rawNumScores := os.Args[1]
	partNum := os.Args[2]
	numScores, err := strconv.Atoi(rawNumScores)
	if err != nil {
		panic(err)
	}
	// Kind of ugly but it works for the simple argparsing of this...
	if partNum == "1" || partNum == "1,2" {
		fmt.Println(part1(numScores))
	}
	if partNum == "2" || partNum == "1,2" {
		fmt.Println(part2(rawNumScores))
	} else if partNum != "1" {
		fmt.Println(usageString)
	}
}
