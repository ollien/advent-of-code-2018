package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type coordinate struct {
	row int
	col int
}

func parseCoords(rawCoords []string) ([]coordinate, error) {
	coords := make([]coordinate, 0, len(rawCoords))
	for _, rawCoordPair := range rawCoords {
		var coordPair coordinate
		numMatched, err := fmt.Sscanf(rawCoordPair, "%d, %d", &coordPair.col, &coordPair.row)
		if err != nil {
			return nil, err
		} else if numMatched != 2 {
			return nil, fmt.Errorf("malformed input")
		}

		coords = append(coords, coordPair)
	}

	return coords, nil
}

func getMaxRowAndCol(coords []coordinate) (maxRow int, maxCol int) {
	for _, coord := range coords {
		if coord.row > maxRow {
			maxRow = coord.row
		}
		if coord.col > maxCol {
			maxCol = coord.col
		}
	}

	return
}

// findClosest the index to the closest coordinate from the input file, or -1 if two are equally distant
func findClosest(loc coordinate, coords []coordinate) int {
	minIndex := -1
	minDistance := math.MaxInt32
	multipleDistances := false
	for i, coord := range coords {
		rowDistance := float64(coord.row - loc.row)
		colDistance := float64(coord.col - loc.col)
		distance := int(math.Abs(rowDistance)) + int(math.Abs(colDistance))
		if distance < minDistance {
			minDistance = distance
			minIndex = i
			multipleDistances = false
		} else if distance == minDistance {
			multipleDistances = true
		}
	}

	if multipleDistances {
		return -1
	} else {
		return minIndex
	}
}

func populateBoardWithNearest(board [][]int, coords []coordinate) {
	for i := range board {
		for j := range board[i] {
			board[i][j] = findClosest(coordinate{i, j}, coords)
		}
	}
}

// isBounded returns true if a board has a finite area
func isBounded(board [][]int, loc coordinate) bool {
	if board[loc.row][0] == board[loc.row][loc.col] || board[loc.row][len(board[loc.row])-1] == board[loc.row][loc.col] {
		return false
	}
	if board[0][loc.col] == board[loc.row][loc.col] || board[len(board)-1][loc.col] == board[loc.row][loc.col] {
		return false
	}

	return true
}

// getAreas Returns the areas of all regions, with each region at the index of the returned slice
func getAreas(board [][]int, numCoords int) []int {
	areas := make([]int, numCoords)
	for i := range board {
		for j := range board[i] {
			closest := board[i][j]
			if closest != -1 {
				areas[closest]++
			}
		}
	}

	return areas
}

func part1(board [][]int, coords []coordinate) int {
	populateBoardWithNearest(board, coords)
	areas := getAreas(board, len(coords))
	largestArea := 0
	for i, area := range areas {
		if isBounded(board, coords[i]) && area > largestArea {
			largestArea = area
		}
	}

	return largestArea
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
	rawCoords := strings.Split(string(inFileContents), "\n")
	// trim tailing newline
	rawCoords = rawCoords[:len(rawCoords)-1]
	coords, err := parseCoords(rawCoords)
	if err != nil {
		panic(err)
	}

	maxRow, maxCol := getMaxRowAndCol(coords)
	board := make([][]int, maxRow+1)
	for i := range board {
		board[i] = make([]int, maxCol+1)
	}

	fmt.Println(part1(board, coords))
}
