package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	inputFormat         = "position=<%d, %d> velocity=<%d, %d>"
	malformedInputError = "malformed input"
)

type point struct {
	row int
	col int
}

type velocity struct {
	rowVelocity int
	colVelocity int
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}

	return n
}

func normalizeInputLines(points map[point][]velocity, rowOffset int, colOffset int) map[point][]velocity {
	normalizedPoints := make(map[point][]velocity, len(points))
	for storedPoint, storedVelocities := range points {
		storedPoint.row += rowOffset
		storedPoint.col += colOffset
		normalizedPoints[storedPoint] = storedVelocities
	}

	return normalizedPoints
}

func parseInput(rawPoints []string) (map[point][]velocity, error) {
	minRow := math.MaxInt32
	minCol := math.MaxInt32
	maxRow := 0
	maxCol := 0
	points := make(map[point][]velocity, len(rawPoints))
	for _, rawPoint := range rawPoints {
		parsedPoint := point{}
		parsedVelocity := velocity{}
		numMatched, err := fmt.Sscanf(rawPoint, inputFormat, &parsedPoint.col, &parsedPoint.row, &parsedVelocity.colVelocity, &parsedVelocity.rowVelocity)
		if err != nil {
			return nil, err
		} else if numMatched != 4 {
			return nil, fmt.Errorf(malformedInputError)
		}
		points[parsedPoint] = []velocity{parsedVelocity}
		if parsedPoint.row < minRow {
			minRow = parsedPoint.row
		}
		if parsedPoint.col < minCol {
			minCol = parsedPoint.col
		}
		if parsedPoint.row > maxRow {
			maxRow = parsedPoint.row
		}
		if parsedPoint.col > maxCol {
			maxCol = parsedPoint.col
		}
	}

	return normalizeInputLines(points, abs(minRow), abs(minCol)), nil
}
func findMinPos(points map[point][]velocity) (minRow int, minCol int, maxRow int, maxCol int) {
	minCol = math.MaxInt32
	maxCol = 0
	minRow = math.MaxInt32
	maxRow = 0
	for storedPoint := range points {
		if storedPoint.row < minRow {
			minRow = storedPoint.row
		}
		if storedPoint.row > maxRow {
			maxRow = storedPoint.row
		}
		if storedPoint.col < minCol {
			minCol = storedPoint.col
		}
		if storedPoint.col > maxCol {
			maxCol = storedPoint.col
		}
	}

	return
}

func movePoints(points map[point][]velocity) map[point][]velocity {
	updatedPoints := make(map[point][]velocity)
	for storedPoint, storedVelocities := range points {
		for _, storedVelocity := range storedVelocities {
			updatedPoint := storedPoint
			updatedPoint.row += storedVelocity.rowVelocity
			updatedPoint.col += storedVelocity.colVelocity
			_, ok := updatedPoints[updatedPoint]
			if ok {
				updatedPoints[updatedPoint] = append(updatedPoints[updatedPoint], storedVelocity)
			} else {
				updatedPoints[updatedPoint] = []velocity{storedVelocity}
			}
		}
	}

	return updatedPoints
}

func printBoard(points map[point][]velocity) {
	minRow, minCol, maxRow, maxCol := findMinPos(points)
	for row := minRow; row <= maxRow; row++ {
		lineBuffer := bytes.NewBufferString("")
		for col := minCol; col <= maxCol; col++ {
			if _, ok := points[point{row, col}]; ok {
				lineBuffer.WriteRune('#')
			} else {
				lineBuffer.WriteRune('.')
			}
		}
		fmt.Println(lineBuffer.String())
	}
}

func atLeastNInACol(points map[point][]velocity, n int) bool {
	cols := make(map[int][]int)
	for storedPoint := range points {
		if _, ok := cols[storedPoint.col]; ok {
			cols[storedPoint.col] = append(cols[storedPoint.col], storedPoint.row)
		} else {
			cols[storedPoint.col] = []int{storedPoint.row}
		}
	}
	for col := range cols {
		sort.Ints(cols[col])
		continuousCount := 1
		for i := range cols[col] {
			if i != 0 && cols[col][i-1] == cols[col][i]-1 {
				continuousCount++
			} else {
				continuousCount = 0
			}
			if continuousCount == n-1 {
				return true
			}
		}
	}

	return false
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
	rawPointInfo := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawPointInfo = rawPointInfo[:len(rawPointInfo)-1]
	points, err := parseInput(rawPointInfo)
	if err != nil {
		panic(err)
	}
	hourCount := 0
	// All letters are 10 chars high. The smallest one I found is just missing two chars
	for !atLeastNInACol(points, 8) {
		hourCount++
		points = movePoints(points)
	}
	fmt.Printf("Hour %d\n", hourCount)
	printBoard(points)
}
