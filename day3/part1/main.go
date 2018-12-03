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
	x      int
	y      int
	width  int
	height int
}

func parsePuzzleLine(line string) (piece, error) {
	pattern, err := regexp.Compile(`(\d+),(\d+): (\d+)x(\d+)`)
	if err != nil {
		return piece{}, err
	}

	matches := pattern.FindStringSubmatch(line)
	parsedPiece := piece{}

	parsedPiece.x, err = strconv.Atoi(matches[1])
	if err != nil {
		return piece{}, err
	}

	parsedPiece.y, err = strconv.Atoi(matches[2])
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
		heightCandidate := checkPiece.y + checkPiece.height
		if heightCandidate > maxHeight {
			maxHeight = heightCandidate
		}

		widthCandidate := checkPiece.x + checkPiece.width
		if widthCandidate > maxWidth {
			maxWidth = widthCandidate
		}
	}

	return maxWidth, maxHeight
}

func insertPiece(insertingPiece piece, cloth [][]int) {
	maxX := insertingPiece.x + insertingPiece.width
	maxY := insertingPiece.y + insertingPiece.height
	for y := insertingPiece.y; y < maxY; y++ {
		for x := insertingPiece.x; x < maxX; x++ {
			cloth[y][x]++
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
