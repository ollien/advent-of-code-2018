package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	inputDelim            = ": "
	yBaseCaseMultiplicand = 16807
	xBaseCaseMultiplicand = 48271
	erosionMod            = 20183
	riskMod               = 3
)

type coordinate struct {
	x int
	y int
}

type caveSpec struct {
	depth  int
	target coordinate
}

func getInputLineValue(line string) string {
	components := strings.Split(line, inputDelim)

	return components[1]
}

func parseCoordinate(spec string) (coordinate, error) {
	components := strings.Split(spec, ",")
	if len(components) != 2 {
		return coordinate{}, errors.New("coordinate must have two components")
	}

	parsedComponents := make([]int, 0, 2)
	for _, component := range components {
		numComponent, err := strconv.Atoi(component)
		if err != nil {
			return coordinate{}, fmt.Errorf("could not parse coordinate component: %w", err)
		}

		parsedComponents = append(parsedComponents, numComponent)
	}

	return coordinate{
		x: parsedComponents[0],
		y: parsedComponents[1],
	}, nil
}

func parseInput(inputLines []string) (caveSpec, error) {
	if len(inputLines) > 2 {
		return caveSpec{}, errors.New("input must have two lines")
	}

	rawDepth := getInputLineValue(inputLines[0])
	depth, err := strconv.Atoi(rawDepth)
	if err != nil {
		return caveSpec{}, fmt.Errorf("could not parse depth: %w", err)
	}

	rawTarget := getInputLineValue(inputLines[1])
	target, err := parseCoordinate(rawTarget)
	if err != nil {
		return caveSpec{}, fmt.Errorf("could not parse target: %w", err)
	}

	return caveSpec{depth: depth, target: target}, nil
}

func calculateRisk(erosionLevel int) int {
	return erosionLevel % riskMod
}

func calculateErosionLevel(depth int, geologicalIndex int) int {
	return (geologicalIndex + depth) % erosionMod
}

func calculateTotalErosionLevelMemo(spec caveSpec, cursor coordinate, memo map[coordinate]int) int {
	if value, ok := memo[cursor]; ok {
		return value
	}

	var index int
	if cursor.x == 0 && cursor.y == 0 {
		index = 0
	} else if cursor.x == 0 {
		index = xBaseCaseMultiplicand * cursor.y
	} else if cursor.y == 0 {
		index = yBaseCaseMultiplicand * cursor.x
	} else {
		index = calculateTotalErosionLevelMemo(spec, coordinate{x: cursor.x - 1, y: cursor.y}, memo) *
			calculateTotalErosionLevelMemo(spec, coordinate{x: cursor.x, y: cursor.y - 1}, memo)
		if cursor == spec.target {
			index = 0
		}
	}

	erosionLevel := calculateErosionLevel(spec.depth, index)
	memo[cursor] = erosionLevel

	return erosionLevel
}

func calculateTotalErosionLevel(spec caveSpec, cursor coordinate) int {
	memo := map[coordinate]int{}
	return calculateTotalErosionLevelMemo(spec, cursor, memo)
}

func part1(spec caveSpec) int {
	totalRisk := 0
	// Since we're going over coordinates that will likely have been gone over, we need to keep a handle on the memo ourselves
	// (need is a strong term - in truth, it takes about 5s without, but why wait? :))
	memo := map[coordinate]int{}
	for x := 0; x <= spec.target.x; x++ {
		for y := 0; y <= spec.target.y; y++ {
			cursor := coordinate{x, y}
			erosionLevel := calculateTotalErosionLevelMemo(spec, cursor, memo)
			totalRisk += calculateRisk(erosionLevel)
		}
	}

	return totalRisk
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

	splitContents := strings.Split(string(inFileContents), "\n")
	splitContents = splitContents[:len(splitContents)-1]

	spec, err := parseInput(splitContents)
	if err != nil {
		panic(err)
	}

	fmt.Println(part1(spec))
}
