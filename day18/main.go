package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	malformedInputError   = "malformed input"
	openChar              = '.'
	treeChar              = '|'
	lumberChar            = '#'
	breedCount            = 3
	treeFillCount         = 3
	stationaryRequirement = 1
	part1Ticks            = 10
	part2Ticks            = 1000000000
)

const (
	openState boardState = iota
	treeState
	lumberState
)

type boardState int
type board [][]boardState

func (b board) print() {
	for row := range b {
		for _, state := range b[row] {
			printChar := openChar
			switch state {
			case treeState:
				printChar = treeChar
			case lumberState:
				printChar = lumberChar
			}
			fmt.Printf("%c", printChar)
		}
		fmt.Print("\n")
	}
}

func (b board) getValue() int {
	numTrees := 0
	numLumber := 0
	for row := range b {
		for _, state := range b[row] {
			switch state {
			case treeState:
				numTrees++
			case lumberState:
				numLumber++
			}
		}
	}

	return numTrees * numLumber
}

func (b board) clone() board {
	clonedBoard := make(board, len(b))
	for row := range b {
		clonedBoard[row] = make([]boardState, len(b[row]))
		copy(clonedBoard[row], b[row])
	}

	return clonedBoard
}

func (b board) getAdjacentCounts(row, col int) (trees, lumberyards int) {
	for dRow := -1; dRow <= 1; dRow++ {
		for dCol := -1; dCol <= 1; dCol++ {
			checkRow := row + dRow
			checkCol := col + dCol
			if (dRow == 0 && dCol == 0) || checkRow < 0 || checkCol < 0 || checkRow >= len(b) || checkCol >= len(b[row]) {
				continue
			}
			switch b[checkRow][checkCol] {
			case treeState:
				trees++
			case lumberState:
				lumberyards++
			}
		}
	}
	return
}

func (b board) tick() board {
	readBoard := b.clone()
	for row := range b {
		for col := range b[row] {
			adjacentTrees, adjacentLumberyards := readBoard.getAdjacentCounts(row, col)
			if readBoard[row][col] == openState && adjacentTrees >= breedCount {
				b[row][col] = treeState
			} else if readBoard[row][col] == treeState && adjacentLumberyards >= treeFillCount {
				b[row][col] = lumberState
			} else if readBoard[row][col] == lumberState && !(adjacentLumberyards >= stationaryRequirement && adjacentTrees >= stationaryRequirement) {
				b[row][col] = openState
			}
		}
	}

	return readBoard
}

func (b board) isIdentical(b2 board) bool {
	if b == nil && b2 == nil {
		return true
	} else if b == nil || b2 == nil {
		return false
	}

	for row := range b {
		for col := range b[row] {
			if b[row][col] != b2[row][col] {
				return false
			}
		}
	}

	return true
}

func parseBoard(rawBoard []string) (board, error) {
	parsedBoard := make(board, len(rawBoard))
	for row, boardLine := range rawBoard {
		parsedBoard[row] = make([]boardState, len(boardLine))
		for col, boardChar := range boardLine {
			switch boardChar {
			case treeChar:
				parsedBoard[row][col] = treeState
			case lumberChar:
				parsedBoard[row][col] = lumberState
			case openChar:
				parsedBoard[row][col] = openState
			default:
				return nil, errors.New(malformedInputError)
			}
		}
	}

	return parsedBoard, nil
}

func runSimulation(parsedBoard board, numTicks int) int {
	pastBoards := []board{}
	cycleLowerBound := -1
	cycleUpperBound := -1
	for tick := 0; tick < numTicks; tick++ {
		lastBoard := parsedBoard.tick()
		pastBoards = append(pastBoards, lastBoard)
		// Check all past borads. If we've already encountered this board, we are bound to repeat ourselves. We can calculate the final score from here.
		for i, pastBoard := range pastBoards {
			if parsedBoard.isIdentical(pastBoard) {
				cycleUpperBound = len(pastBoards)
				cycleLowerBound = i
				break
			}
		}
		if cycleLowerBound != -1 {
			break
		}
	}

	// If we have a cycle, we need to calculate the score at the end of it
	if cycleLowerBound == -1 {
		return parsedBoard.getValue()
	} else {
		cycleRange := cycleUpperBound - cycleLowerBound
		// find the cycle that will put us just over the number of ticks
		// rearrange cycleLowerBound + cycleRange * n >= numTicks for n, n >= (numTicks - cycleLoewrBound)/cycleRange
		// Need to take ceiling of division, hence the funky add by ones.
		numCyclesToOvershoot := ((numTicks-cycleLowerBound-1)/cycleRange + 1)
		overshootTicks := numCyclesToOvershoot*cycleRange + cycleLowerBound
		numOver := overshootTicks - numTicks
		// Work backwards from the top number of cycles based on th enumber we went over
		boardIndex := ((cycleUpperBound - numOver) + cycleUpperBound) % cycleUpperBound

		return pastBoards[boardIndex].getValue()
	}
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
	rawBoard := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawBoard = rawBoard[:len(rawBoard)-1]
	parsedBoard, err := parseBoard(rawBoard)
	if err != nil {
		panic(err)
	}
	fmt.Println(runSimulation(parsedBoard, part1Ticks))

	parsedBoard, err = parseBoard(rawBoard)
	if err != nil {
		panic(err)
	}
	fmt.Println(runSimulation(parsedBoard, part2Ticks))
}
