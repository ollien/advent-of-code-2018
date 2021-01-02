package main

import (
	"container/heap"
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
	toolSwitchTime        = 7
	movementTime          = 1
)

type tool int

const (
	toolTorch tool = iota
	toolClimbingGear
	toolNothing
)

var toolByRisk = map[int][]tool{
	// rocky
	0: {toolClimbingGear, toolTorch},
	// wet
	1: {toolClimbingGear, toolNothing},
	// narrow
	2: {toolTorch, toolNothing},
}

type coordinate struct {
	x int
	y int
}

type caveSpec struct {
	depth  int
	target coordinate
}

// represents a single movement from room to room
type movement struct {
	endingState playerState
	time        int
}

// represents the playerState of the player after moving
// This does _NOT_ include time, and should not, as equality must represent being at the same place, regardless fo time.
type playerState struct {
	destination coordinate
	currentTool tool
}

// Implements container/heap.Interface
// The highest priority item is the movement with the lowest distance
type movementPriorityQueue struct {
	backingSlice []movement
	times        map[playerState]int
}

func (queue *movementPriorityQueue) Len() int {
	return len(queue.backingSlice)
}

func (queue *movementPriorityQueue) Less(i, j int) bool {
	movement1 := queue.backingSlice[i]
	movement2 := queue.backingSlice[j]
	movement1Time := queue.times[movement1.endingState]
	movement2Time := queue.times[movement2.endingState]
	return movement1Time < movement2Time
}

func (queue *movementPriorityQueue) Swap(i, j int) {
	queue.backingSlice[i], queue.backingSlice[j] = queue.backingSlice[j], queue.backingSlice[i]
}

func (queue *movementPriorityQueue) Push(x interface{}) {
	queue.backingSlice = append(queue.backingSlice, x.(movement))
}

func (queue *movementPriorityQueue) Pop() interface{} {
	el := queue.backingSlice[queue.Len()-1]
	queue.backingSlice = queue.backingSlice[:queue.Len()-1]

	return el
}

func (queue *movementPriorityQueue) contains(m movement) bool {
	for _, item := range queue.backingSlice {
		if item == m {
			return true
		}
	}

	return false
}

func (queue *movementPriorityQueue) containsMovementWithDestination(dest coordinate) bool {
	for _, item := range queue.backingSlice {
		if item.endingState.destination == dest {
			return true
		}
	}

	return false
}

func (c coordinate) calculateNeighbors() []coordinate {
	candidates := []coordinate{
		{
			x: c.x + 1,
			y: c.y,
		},
		{
			x: c.x,
			y: c.y + 1,
		},
		{
			x: c.x - 1,
			y: c.y,
		},
		{
			x: c.x,
			y: c.y - 1,
		},
	}

	res := make([]coordinate, 0, 4)
	for _, candidate := range candidates {
		if candidate.x < 0 || candidate.y < 0 {
			continue
		}

		res = append(res, candidate)
	}

	return res
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

// Not used: kept for illustrative purposes
func calculateTotalErosionLevel(spec caveSpec, cursor coordinate) int {
	memo := map[coordinate]int{}
	return calculateTotalErosionLevelMemo(spec, cursor, memo)
}

// Find the possible tools that can be used given a current set of possible tools at a location
// Takes the erosion memo into account. Could technically be omitted but it would take too long without an existing memo.
func getPossibleToolsForLocation(spec caveSpec, location coordinate, erosionMemo map[coordinate]int) []tool {
	erosionLevel := calculateTotalErosionLevelMemo(spec, location, erosionMemo)
	risk := calculateRisk(erosionLevel)

	return toolByRisk[risk]
}

// Find the possible tools that can be used given a current set of possible tools and moving from point A to point B.
// Takes the erosion memo into account. Could technically be omitted but it would take too long without an existing memo.
func findPossibleToolsForMovement(spec caveSpec, origin coordinate, destination coordinate, erosionMemo map[coordinate]int) []tool {
	// There will only ever be 2 tools at maximum (which is when the other location is of the same type)
	// This could potentially be micro-optimized but eh
	possibleTools := make([]tool, 0, 2)
	for _, neighborTool := range getPossibleToolsForLocation(spec, origin, erosionMemo) {
		for _, currentLocationTool := range getPossibleToolsForLocation(spec, destination, erosionMemo) {
			if neighborTool == currentLocationTool {
				possibleTools = append(possibleTools, currentLocationTool)
			}
		}
	}

	return possibleTools
}

// Find all possible movements from the visiting location
// Takes the erosion memo into account. Could technically be omitted but it would take too long without an existing memo.
func calculateMovementCandidates(spec caveSpec, visiting movement, erosionMemo map[coordinate]int) []movement {
	movementCandidates := []movement{}
	for _, neighbor := range visiting.endingState.destination.calculateNeighbors() {
		for _, toolCandidate := range findPossibleToolsForMovement(spec, visiting.endingState.destination, neighbor, erosionMemo) {
			time := toolSwitchTime + movementTime
			if toolCandidate == visiting.endingState.currentTool {
				time -= toolSwitchTime
			}

			candidate := movement{
				endingState: playerState{
					destination: neighbor,
					currentTool: toolCandidate,
				},
				time: time,
			}

			// We need to force a switch to the torch if the destination is the torch
			if candidate.endingState.destination == spec.target && candidate.endingState.currentTool != toolTorch {
				candidate.time += toolSwitchTime
				candidate.endingState.currentTool = toolTorch
			}

			movementCandidates = append(movementCandidates, candidate)
		}
	}

	return movementCandidates
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

func part2(spec caveSpec) int {
	erosionMemo := map[coordinate]int{}
	times := map[playerState]int{
		{destination: coordinate{x: 0, y: 0}, currentTool: toolTorch}: 0,
	}
	toVisit := &movementPriorityQueue{
		backingSlice: []movement{},
		times:        times,
	}

	heap.Init(toVisit)
	heap.Push(toVisit, movement{
		endingState: playerState{
			destination: coordinate{x: 0, y: 0},
			currentTool: toolTorch,
		},
		time: 0,
	})
	// Uses Djikstra, specifically a variant that terminates once we find the specific item, rather than searching the full graph (since that's infinite)
	for toVisit.Len() > 0 && !toVisit.containsMovementWithDestination(spec.target) {
		visiting := heap.Pop(toVisit).(movement)
		visitingState := playerState{destination: visiting.endingState.destination, currentTool: visiting.endingState.currentTool}
		movementCandidates := calculateMovementCandidates(spec, visiting, erosionMemo)
		currentTime, ok := times[visitingState]
		if !ok {
			panic("this can literally never happen, but somehow, we attempted to visit a location without a ne entry in the time table")
		}

		for _, candidate := range movementCandidates {
			candidateTime := currentTime + candidate.time
			candidateState := playerState{destination: candidate.endingState.destination, currentTool: candidate.endingState.currentTool}

			currentTimeForCandidate, haveTimeForCandidate := times[candidateState]
			if haveTimeForCandidate && candidateTime >= currentTimeForCandidate {
				continue
			}

			times[candidateState] = candidateTime
			if !toVisit.contains(candidate) {
				heap.Push(toVisit, candidate)
			}
		}

	}

	return times[playerState{destination: spec.target, currentTool: toolTorch}]
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
	fmt.Println(part2(spec))
}
