package main

import (
	"fmt"
	"os"
	"strconv"
)

const gridSize = 300

type summedAreaTable [][]int

// Make a table where each location represents the scores up and to the left
func makeSummedAreaTable(serialNumber int) summedAreaTable {
	table := make([][]int, gridSize)
	for row := range table {
		table[row] = make([]int, gridSize)
		for col := range table[row] {
			score := getTileScore(row+1, col+1, serialNumber)
			if row-1 >= 0 && col-1 >= 0 {
				score += table[row-1][col] + table[row][col-1] - table[row-1][col-1]
			} else if row-1 >= 0 {
				score += table[row-1][col]
			} else if col-1 >= 0 {
				score += table[row][col-1]
			}
			table[row][col] = score
		}
	}

	return table
}

func getTileScore(row int, col int, serialNumber int) int {
	rackID := col + 10
	powerLevel := rackID * (rackID*row + serialNumber)
	powerLevel = powerLevel / 100 % 10
	return powerLevel - 5
}

func getSquareScore(areaTable summedAreaTable, startRow int, startCol int, size int, serialNumber int) (score int) {
	bottomRightRow := startRow + (size - 1)
	bottomRightCol := startCol + (size - 1)
	if bottomRightCol >= gridSize || bottomRightRow >= gridSize {
		return 0
	}

	bottomRightScore := areaTable[bottomRightRow][bottomRightCol]

	leftBoundRow := bottomRightRow
	leftBoundCol := startCol - 1
	leftBoundScore := 0
	if leftBoundRow >= 0 && leftBoundCol >= 0 {
		leftBoundScore = areaTable[leftBoundRow][leftBoundCol]
	}

	topBoundRow := startRow - 1
	topBoundCol := bottomRightCol
	topBoundScore := 0
	if topBoundRow >= 0 && topBoundCol >= 0 {
		topBoundScore = areaTable[topBoundRow][topBoundCol]
	}

	topLeftBoundScore := 0
	if startRow > 0 && startCol > 0 {
		topLeftBoundScore = areaTable[startRow-1][startCol-1]
	}

	/*
	 * cxxb
	 * xxxx
	 * xxxx
	 * dxxa
	 * Perform a - b - c + d; the +d eliminates double counting form the other two subtractions
	 */
	return bottomRightScore - leftBoundScore - topBoundScore + topLeftBoundScore
}

func part1(areaTable summedAreaTable, serialNumber int) (bestRow int, bestCol int) {
	maxScore := 0
	for i := 1; i < gridSize; i++ {
		for j := 1; j < gridSize; j++ {
			score := getSquareScore(areaTable, i-1, j-1, 3, serialNumber)
			if score > maxScore {
				maxScore = score
				bestRow = i
				bestCol = j
			}
		}
	}

	return
}

func part2(areaTable summedAreaTable, serialNumber int) (bestRow int, bestCol int, bestSize int) {
	maxScore := 0
	for squareSize := 1; squareSize <= gridSize; squareSize++ {
		for i := 1; i < gridSize; i++ {
			for j := 1; j < gridSize; j++ {
				score := getSquareScore(areaTable, i-1, j-1, squareSize, serialNumber)
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

	areaTable := makeSummedAreaTable(serialNumber)
	bestRow, bestCol := part1(areaTable, serialNumber)
	fmt.Printf("%d,%d\n", bestCol, bestRow)
	bestRow, bestCol, bestSize := part2(areaTable, serialNumber)
	fmt.Printf("%d,%d,%d\n", bestCol, bestRow, bestSize)
}
