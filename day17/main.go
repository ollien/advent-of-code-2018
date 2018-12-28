// This solution is slow and not ironclad, but it works. Probably could improve it using a 2D array, but it took a long time to figure it out as is.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

const (
	noSuchRangeError    = "no such range"
	malformedInputError = "malformed input"
	xYRangeFormat       = "x=%d, y=%d..%d"
	yXRangeFormat       = "y=%d, x=%d..%d"
	initialCol          = 500
)

const (
	cantOverflow overflowCondition = iota
	wontOverflow
	willOverflow
)

type overflowCondition int
type rangeList []positionRange

type positionRange struct {
	high   int
	low    int
	static bool
}

type board struct {
	positions      map[int]rangeList
	minCol, maxCol int
	minRow, maxRow int
}

type positionCursor struct {
	row int
	col int
}

func appendToRangeList(ranges rangeList, newRange positionRange) []positionRange {
	for i := range ranges {
		if ranges[i].low <= newRange.low && ranges[i].high >= newRange.high {
			// If our new set is entirely contained within the existing set, we're done.
			ranges[i].static = (newRange.static || ranges[i].static)
			return ranges
		} else if ranges[i].low >= newRange.low && ranges[i].high <= newRange.high {
			// If our existing set is entirely contained within the new set, take the new set.
			ranges[i].low = newRange.low
			ranges[i].high = newRange.high
			ranges[i].static = (newRange.static || ranges[i].static)
			return ranges
		} else if ranges[i].low == newRange.high || ranges[i].low-1 == newRange.high {
			ranges[i].low = newRange.low
			ranges[i].static = (newRange.static || ranges[i].static)
			return ranges
		} else if ranges[i].high == newRange.low || ranges[i].high+1 == newRange.low {
			ranges[i].high = newRange.high
			ranges[i].static = (newRange.static || ranges[i].static)
			return ranges
		}
	}

	return append(ranges, newRange)
}

// parseInput takes the input for the problem and produces a map of clay rows to ranges of columns (i.e. y values to x value ranges)
func parseInput(input []string) (board, error) {
	parsedBoard := board{
		positions: map[int]rangeList{},
		minCol:    math.MaxInt32,
		maxCol:    0,
		minRow:    math.MaxInt32,
		maxRow:    0,
	}
	for _, line := range input {
		var coord, rangeMin, rangeMax int
		numMatched, err := fmt.Sscanf(line, yXRangeFormat, &coord, &rangeMin, &rangeMax)
		if err == nil && numMatched != 3 {
			return board{}, errors.New(malformedInputError)
		} else if err == nil {
			parsedBoard.positions[coord] = appendToRangeList(parsedBoard.positions[coord], positionRange{low: rangeMin, high: rangeMax, static: true})
			if rangeMax > parsedBoard.maxCol {
				parsedBoard.maxCol = rangeMax
			}
			if rangeMin < parsedBoard.minCol {
				parsedBoard.minCol = coord
			}
			if coord > parsedBoard.maxRow {
				parsedBoard.maxRow = coord
			}
			if coord < parsedBoard.minRow {
				parsedBoard.minRow = coord
			}
			continue
		}

		numMatched, err = fmt.Sscanf(line, xYRangeFormat, &coord, &rangeMin, &rangeMax)
		if err == nil && numMatched != 3 {
			return board{}, errors.New(malformedInputError)
		} else if err == nil {
			if coord > parsedBoard.maxCol {
				parsedBoard.maxCol = coord
			}
			if coord < parsedBoard.minCol {
				parsedBoard.minCol = coord
			}
			for yCoord := rangeMin; yCoord <= rangeMax; yCoord++ {
				parsedBoard.positions[yCoord] = appendToRangeList(parsedBoard.positions[yCoord], positionRange{low: coord, high: coord})
				if yCoord > parsedBoard.maxRow {
					parsedBoard.maxRow = yCoord
				}
				if yCoord < parsedBoard.minRow {
					parsedBoard.minRow = yCoord
				}
			}
		} else {
			return board{}, err
		}
	}

	// Account for possibility that water flows off edges
	parsedBoard.maxCol++
	parsedBoard.minCol--

	return parsedBoard, nil
}

func (b board) getStatusAt(row int, col int) (occupied, static bool) {
	colRanges, haveRow := b.positions[row]
	if !haveRow {
		return false, false
	}

	for _, colRange := range colRanges {
		if col >= colRange.low && col <= colRange.high {
			return true, colRange.static
		}
	}

	return false, false
}

func (b board) findContainingRange(row int, col int) positionRange {
	resultRange := positionRange{}
	for colCursor := col; colCursor >= b.minCol; colCursor-- {
		if occupied, _ := b.getStatusAt(row, colCursor); occupied {
			resultRange.low = colCursor + 1
			break
		}
	}

	for colCursor := col; colCursor <= b.maxCol; colCursor++ {
		if occupied, _ := b.getStatusAt(row, colCursor); occupied {
			resultRange.high = colCursor - 1
			break
		}
	}

	return resultRange
}

func (b board) clone() board {
	clonedBoard := board{
		positions: make(map[int]rangeList, len(b.positions)),
		minCol:    b.minCol,
		maxCol:    b.maxCol,
		maxRow:    b.maxRow,
	}

	for row, positions := range b.positions {
		clonedBoard.positions[row] = make(rangeList, len(positions))
		copy(clonedBoard.positions[row], positions)
	}

	return clonedBoard
}

func (b board) areIdentical(b2 board) bool {
	if b.minCol != b2.minCol || b.maxCol != b2.minCol || b.maxRow != b2.maxRow || len(b.positions) != len(b2.positions) {
		return false
	}

	for row := range b.positions {
		if len(b.positions[row]) != len(b2.positions[row]) {
			return false
		}
		for i := range b.positions[row] {
			if b.positions[row][i] != b2.positions[row][i] {
				return false
			}
		}
	}
	return true
}

func (b board) getNumTilesOccupied() (total int, numStatic int) {
	for row := b.minRow; row <= b.maxRow; row++ {
		for col := b.minCol; col <= b.maxCol; col++ {
			occupied, static := b.getStatusAt(row, col)
			if occupied {
				total++
			}
			if static {
				numStatic++
			}
		}
	}

	return
}

func printSimulation(clayBoard board, streamBoard board) {
	for row := 0; row <= clayBoard.maxRow; row++ {
		fmt.Printf("%3d", row)
		for col := clayBoard.minCol; col <= clayBoard.maxCol; col++ {
			clayBoardOccupied, _ := clayBoard.getStatusAt(row, col)
			streamBoardOccupied, streamIsStatic := streamBoard.getStatusAt(row, col)
			if row == 0 && col == 500 {
				fmt.Print("+")
			} else if clayBoardOccupied {
				fmt.Print("#")
			} else if streamBoardOccupied && !streamIsStatic {
				fmt.Print("|")
			} else if streamBoardOccupied {
				fmt.Print("~")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func getOverflowCondition(row int, col int, clayBoard board, streamBoard board) overflowCondition {
	streamBoardOccupied, streamIsStatic := streamBoard.getStatusAt(row+1, col)
	clayBoardOccupied, _ := clayBoard.getStatusAt(row+1, col)
	clayBoardOccupiedAbove, _ := clayBoard.getStatusAt(row, col)
	streamBoardOccupiedAbove, streamIsStaticAbove := streamBoard.getStatusAt(row, col)
	if (!streamBoardOccupied && !clayBoardOccupied) || (streamBoardOccupiedAbove && !streamIsStaticAbove && !clayBoardOccupied && !streamIsStatic) {
		return willOverflow
	} else if clayBoardOccupiedAbove {
		return cantOverflow
	} else {
		return wontOverflow
	}
}

func findOverflowCandidates(streamCursor positionCursor, clayBoard board, streamBoard board) []positionCursor {
	overflowCursors := []positionCursor{}
	for colCursor := streamCursor.col; colCursor >= clayBoard.minCol; colCursor-- {
		streamCondition := getOverflowCondition(streamCursor.row, colCursor, clayBoard, streamBoard)
		if streamCondition == willOverflow {
			overflowCursors = append(overflowCursors, positionCursor{row: streamCursor.row, col: colCursor})
			break
		} else if streamCondition == cantOverflow {
			break
		}
	}

	for colCursor := streamCursor.col; colCursor <= clayBoard.maxCol; colCursor++ {
		streamCondition := getOverflowCondition(streamCursor.row, colCursor, clayBoard, streamBoard)
		if streamCondition == willOverflow {
			overflowCursors = append(overflowCursors, positionCursor{row: streamCursor.row, col: colCursor})
			break
		} else if streamCondition == cantOverflow {
			break
		}
	}

	return overflowCursors
}

func findCursorCandidates(streamBoard board) []positionCursor {
	cursors := []positionCursor{}
	for row, colRanges := range streamBoard.positions {
		for _, colRange := range colRanges {
			for col := colRange.low; col <= colRange.high; col++ {
				streamBoardOccupiedBelow, streamIsStaticBelow := streamBoard.getStatusAt(row+1, col)
				streamBoardOccupied, streamIsStatic := streamBoard.getStatusAt(row, col)
				_, streamIsStaticLeft := streamBoard.getStatusAt(row, col+1)
				_, streamIsStaticRight := streamBoard.getStatusAt(row, col-1)
				// We have a candidate for flow if a stream is occupied, not static, and it is not surrounded by flowing water.
				if streamBoardOccupied && !streamIsStatic && (!streamBoardOccupiedBelow || streamIsStaticBelow) && !(streamIsStaticLeft && streamIsStaticRight) {
					cursors = append(cursors, positionCursor{row: row, col: col})
				}
			}
		}
	}

	return cursors
}

func flow(clayBoard board) (total int, numStatic int) {
	streamBoard := board{
		maxRow:    0,
		minRow:    clayBoard.minRow,
		minCol:    clayBoard.minCol,
		maxCol:    clayBoard.maxCol,
		positions: map[int]rangeList{0: []positionRange{positionRange{low: initialCol, high: initialCol}}},
	}
	lastBoard := board{}

	var lastCursors []positionCursor
	for !lastBoard.areIdentical(streamBoard) {
		lastBoard = streamBoard.clone()
		streamCursors := findCursorCandidates(streamBoard)
		for _, streamCursor := range streamCursors {
			streamBoardOccupied, streamIsStatic := streamBoard.getStatusAt(streamCursor.row+1, streamCursor.col)
			clayBoardOccupied, _ := clayBoard.getStatusAt(streamCursor.row+1, streamCursor.col)
			containingRange := clayBoard.findContainingRange(streamCursor.row, streamCursor.col)
			if (streamBoardOccupied && streamIsStatic) || clayBoardOccupied {
				overflowCursors := findOverflowCandidates(streamCursor, clayBoard, streamBoard)
				if len(overflowCursors) == 0 {
					containingRange.static = true
					streamBoard.positions[streamCursor.row] = appendToRangeList(streamBoard.positions[streamCursor.row], containingRange)
					streamCursor.row--
				} else if len(overflowCursors) == 1 {
					overflowCol := overflowCursors[0].col
					// If of stream flows left, set the low to where we flowed
					if overflowCol <= streamCursor.col {
						containingRange.low = overflowCol
					} else {
						containingRange.high = overflowCol
					}
					overflowRow := overflowCursors[0].row
					streamBoard.positions[overflowRow] = appendToRangeList(streamBoard.positions[overflowRow], containingRange)
				} else {
					// The most overflow cursors we can have is two
					firstCol := overflowCursors[0].col
					secondCol := overflowCursors[1].col
					flowRange := positionRange{}
					if firstCol > secondCol {
						flowRange.high = firstCol
						flowRange.low = secondCol
					} else {
						flowRange.high = secondCol
						flowRange.low = firstCol
					}
					overflowRow := overflowCursors[0].row
					streamBoard.positions[overflowRow] = appendToRangeList(streamBoard.positions[overflowRow], flowRange)
				}
			} else if !streamBoardOccupied && !streamIsStatic {
				newRange := positionRange{low: streamCursor.col, high: streamCursor.col}
				// We can't flow outside the board
				if streamCursor.row <= clayBoard.maxRow {
					streamCursor.row++
				}
				streamBoard.positions[streamCursor.row] = appendToRangeList(streamBoard.positions[streamCursor.row], newRange)
			}
			if streamCursor.row > streamBoard.maxRow {
				streamBoard.maxRow = streamCursor.row
			}
		}

		// If, after flowing, we're in the same state as last time, we're done.
		if len(lastCursors) == len(streamCursors) {
			different := false
			for i := range lastCursors {
				if lastCursors[i] != streamCursors[i] {
					different = true
					break
				}
			}
			if !different {
				// Delete the extra flow
				delete(streamBoard.positions, streamBoard.maxRow)
				break
			}
		}
		lastCursors = streamCursors
	}

	return streamBoard.getNumTilesOccupied()
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
	rawTileInfo := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawTileInfo = rawTileInfo[:len(rawTileInfo)-1]
	parsedBoard, err := parseInput(rawTileInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println(flow(parsedBoard))
}
