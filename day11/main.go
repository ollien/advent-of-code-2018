package main

import (
	"fmt"
	"os"
	"strconv"
)

const gridSize = 300

func getSquareScore(startRow int, startCol int, size int, serialNumber int) (score int) {
	for i := startRow; i < startRow+size; i++ {
		for j := startCol; j < startCol+size; j++ {
			rackID := j + 10
			powerLevel := rackID * (rackID*i + serialNumber)
			powerLevel = powerLevel / 100 % 10
			score += powerLevel - 5
		}
	}

	return
}

func part1(serialNumber int) (bestRow int, bestCol int) {
	maxScore := 0
	for i := 1; i <= gridSize; i++ {
		for j := 1; j <= gridSize; j++ {
			score := getSquareScore(i, j, 3, serialNumber)
			if score > maxScore {
				maxScore = score
				bestRow = i
				bestCol = j
			}
		}
	}

	return
}

func part2(serialNumber int) (bestRow int, bestCol int, bestSize int) {
	maxScore := 0
	for squareSize := 1; squareSize <= gridSize; squareSize++ {
		for i := 1; i <= gridSize; i++ {
			for j := 1; j <= gridSize; j++ {
				score := getSquareScore(i, j, squareSize, serialNumber)
				if score > maxScore {
					maxScore = score
					bestRow = i
					bestCol = j
					bestSize = squareSize
				}
			}
		}
	}

	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main serial_number")
		return
	}

	rawSerialNumber := os.Args[1]
	serialNumber, err := strconv.Atoi(rawSerialNumber)
	if err != nil {
		panic(err)
	}

	bestRow, bestCol := part1(serialNumber)
	fmt.Printf("%d,%d\n", bestCol, bestRow)
	bestRow, bestCol, bestSize := part2(serialNumber)
	fmt.Printf("%d,%d,%d\n", bestCol, bestRow, bestSize)
}
