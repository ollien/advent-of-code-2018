package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	score1 = 3
	score2 = 7
)

func calculateNewScores(scores []int, elf1Cursor int, elf2Cursor int) []int {
	total := scores[elf1Cursor] + scores[elf2Cursor]
	newScores := make([]int, 0, 2)
	if total == 0 {
		newScores = append(newScores, 0)
	} else {
		for total != 0 {
			newScores = append([]int{total % 10}, newScores...)
			total = total / 10
		}
	}

	return append(scores, newScores...)
}

func part1(numScores int) string {
	scores := make([]int, 2, numScores+10)
	scores[0] = score1
	scores[1] = score2
	elf1Cursor := 0
	elf2Cursor := 1
	for len(scores) < numScores+10 {
		scores = calculateNewScores(scores, elf1Cursor, elf2Cursor)
		elf1Score := scores[elf1Cursor]
		elf2Score := scores[elf2Cursor]
		elf1Cursor = (elf1Cursor + elf1Score + 1) % len(scores)
		elf2Cursor = (elf2Cursor + elf2Score + 1) % len(scores)
	}

	resultString := ""
	for i := numScores; i < numScores+10; i++ {
		resultString += strconv.Itoa(scores[i])
	}

	return resultString
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main number_of_scores")
		return
	}

	rawNumScores := os.Args[1]
	numScores, err := strconv.Atoi(rawNumScores)
	if err != nil {
		panic(err)
	}
	fmt.Println(part1(numScores))
}
