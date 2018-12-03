package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type piece struct {
	row    int
	col    int
	height int
	width  int
}

func parsePuzzleLine(line string) (piece, error) {
	pattern, err := regexp.Compile(`(\d+),(\d+): (\d+)x(\d+)`)
	if err != nil {
		return piece{}, err
	}

	matches := pattern.FindStringSubmatch(line)
	parsedPiece := piece{}

	parsedPiece.col, err = strconv.Atoi(matches[1])
	if err != nil {
		return piece{}, err
	}

	parsedPiece.row, err = strconv.Atoi(matches[2])
	if err != nil {
		return piece{}, err
	}

	parsedPiece.width, err = strconv.Atoi(matches[3])
	if err != nil {
		return piece{}, err
	}

	parsedPiece.height, err = strconv.Atoi(matches[4])
	if err != nil {
		return piece{}, err
	}

	return parsedPiece, nil
}

func getMaxDimensions(pieces []piece) (int, int) {
	maxWidth := 0
	maxHeight := 0
	for _, checkPiece := range pieces {
		heightCandidate := checkPiece.row + checkPiece.height
		if heightCandidate > maxHeight {
			maxHeight = heightCandidate
		}

		widthCandidate := checkPiece.col + checkPiece.width
		if widthCandidate > maxWidth {
			maxWidth = widthCandidate
		}
	}

	return maxWidth, maxHeight
}

func insertPiece(insertingPiece piece, cloth [][]int) {
	maxCol := insertingPiece.col + insertingPiece.width
	maxRow := insertingPiece.row + insertingPiece.height
	for row := insertingPiece.row; row < maxRow; row++ {
		for col := insertingPiece.col; col < maxCol; col++ {
			cloth[row][col]++
		}
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
	rawPieces := strings.Split(string(inFileContents), "\n")
	// trim tailing newline
	rawPieces = rawPieces[:len(rawPieces)-1]
	pieces := make([]piece, 0, len(rawPieces))
	for _, rawPiece := range rawPieces {
		parsedPiece, err := parsePuzzleLine(rawPiece)
		if err != nil {
			panic(err)
		}

		pieces = append(pieces, parsedPiece)
	}

	width, height := getMaxDimensions(pieces)
	cloth := make([][]int, height)
	for i := range cloth {
		cloth[i] = make([]int, width)
	}

	for _, insertingPiece := range pieces {
		insertPiece(insertingPiece, cloth)
	}

	intersectCount := 0
	for i := range cloth {
		for j := range cloth[i] {
			if cloth[i][j] > 1 {
				intersectCount++
			}
		}
	}

	fmt.Println(intersectCount)
}
