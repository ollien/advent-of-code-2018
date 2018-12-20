package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	startingHealth      = 200
	attackPower         = 3
	noTargetFoundError  = "no target found"
	malformedInputError = "malformed input"
	targetChar          = 'x'
	wallChar            = '#'
	openChar            = '.'
	elfChar             = 'E'
	goblinChar          = 'G'
)

const (
	noWinner winner = iota
	elfWinner
	goblinWinner
)

type board [][]node
type nodeQueue []node
type winner int

// a list of nodes, sortable in reading order
type nodeList []node

type coordinate struct {
	row int
	col int
}

type node interface {
	setPos(coordinate)
	getPos() coordinate
	canTravelThrough() bool
}

type tile struct {
	position coordinate
	isWall   bool
}

type entity struct {
	position coordinate
	health   int
	isGoblin bool
}

func (q *nodeQueue) enqueue(newNode node) {
	*q = append(*q, newNode)
}

func (q *nodeQueue) dequeue() node {
	head := (*q)[0]
	*q = (*q)[1:]

	return head
}

func (list nodeList) Len() int {
	return len(list)
}

func (list nodeList) Less(i int, j int) bool {
	iPos := list[i].getPos()
	jPos := list[j].getPos()
	if iPos.row == jPos.row {
		return iPos.col < jPos.col
	}

	return iPos.row < jPos.row
}

func (list nodeList) Swap(i int, j int) {
	list[i], list[j] = list[j], list[i]
}

func (t *tile) setPos(newPos coordinate) {
	t.position = newPos
}

func (t *tile) getPos() coordinate {
	return t.position
}

func (t *tile) canTravelThrough() bool {
	return !t.isWall
}

func (e *entity) setPos(newPos coordinate) {
	e.position = newPos
}

func (e *entity) getPos() coordinate {
	return e.position
}

func (e *entity) canTravelThrough() bool {
	return false
}

func (e *entity) move(containingBoard board) {
	toVisit := nodeQueue{e}
	targetCandidates := nodeList{}
	distanceTable := make(map[node]int)
	distance := 0
	searchUntilDistance := -1
	searchedTiles := nodeList{}
	distanceTable[e] = 0
	for len(toVisit) > 0 && distance != searchUntilDistance {
		visitingNode := toVisit.dequeue()
		searchedTiles = append(searchedTiles, visitingNode)
		neighbors := containingBoard.getNeighbors(visitingNode.getPos())
		distance = distanceTable[visitingNode] + 1
		for _, neighbor := range neighbors {
			if _, wasVisited := distanceTable[neighbor]; !wasVisited {
				if entityNode, isEntity := neighbor.(*entity); isEntity && entityNode.isGoblin != e.isGoblin {
					toVisit.enqueue(neighbor)
					distanceTable[neighbor] = distance
					if visitingNode.canTravelThrough() {
						targetCandidates = append(targetCandidates, visitingNode)
						searchUntilDistance = distance + 1
					}
				} else if neighbor.canTravelThrough() {
					toVisit.enqueue(neighbor)
					distanceTable[neighbor] = distance
				}
			}
		}
	}

	if len(targetCandidates) == 0 {
		return
	}

	target := getBestTarget(targetCandidates, distanceTable)
	moveNode := getNextMove(containingBoard, target, distanceTable)
	newPos := moveNode.getPos()
	oldPos := e.position
	e.setPos(newPos)
	moveNode.setPos(oldPos)
	containingBoard[oldPos.row][oldPos.col] = moveNode
	containingBoard[newPos.row][newPos.col] = e
}

func (e *entity) attack(containingBoard board) bool {
	lowestHealthTarget := &entity{health: math.MaxInt32}
	neighbors := containingBoard.getNeighbors(e.getPos())
	sort.Sort(neighbors)
	for _, neighbor := range neighbors {
		if entityNode, isEntity := neighbor.(*entity); isEntity && entityNode.isGoblin != e.isGoblin {
			if entityNode.health < lowestHealthTarget.health {
				lowestHealthTarget = entityNode
			}
		}
	}
	if lowestHealthTarget.health == math.MaxInt32 {
		return false
	}

	lowestHealthTarget.health -= attackPower
	if lowestHealthTarget.health <= 0 {
		entityPos := lowestHealthTarget.getPos()
		containingBoard[entityPos.row][entityPos.col] = &tile{
			position: entityPos,
			isWall:   false,
		}
	}

	return true
}

func getBestTarget(targetCandidates nodeList, distanceTable map[node]int) node {
	sort.Sort(targetCandidates)
	lowestDistance := math.MaxInt32
	var bestNode node
	for _, target := range targetCandidates {
		if distance := distanceTable[target]; distance < lowestDistance {
			lowestDistance = distance
			bestNode = target
		}
	}

	return bestNode
}

func getNextMove(containingBoard board, target node, distanceTable map[node]int) node {
	possiblePaths := getPathsToTarget(containingBoard, target, distanceTable)

	pathStarts := nodeList{}
	for _, path := range possiblePaths {
		if len(path) > 0 {
			pathStarts = append(pathStarts, path[0])
		}
	}
	if len(pathStarts) > 0 {
		sort.Sort(pathStarts)
		// Get the first path in reading order
		return pathStarts[0]
	}

	return target
}

// getPathsToTarget gets all of the shortest path to a target given the distances to each node
func getPathsToTarget(containingBoard board, target node, distanceTable map[node]int) []nodeList {
	if distanceTable[target] == 0 {
		return []nodeList{}
	}

	pathStarts := nodeList{}
	targetNeighbors := containingBoard.getNeighbors(target.getPos())
	for _, neighbor := range targetNeighbors {
		if distance, ok := distanceTable[neighbor]; ok && distance == distanceTable[target]-1 {
			pathStarts = append(pathStarts, neighbor)
		}
	}
	sort.Sort(pathStarts)

	paths := []nodeList{}
	for _, pathStart := range pathStarts {
		possiblePaths := getPathsToTarget(containingBoard, pathStart, distanceTable)
		if len(possiblePaths) == 0 {
			// Create a list for this path, and then continue on.
			paths = append(paths, nodeList{})
			continue
		}
		for _, possiblePath := range possiblePaths {
			path := append(possiblePath, pathStart)
			paths = append(paths, path)
		}
	}

	return paths
}

// print outputs the board to stdout, with any targets marked with targetChar
func (b board) print(targets nodeList) {
	for row, boardRow := range b {
		for col, node := range boardRow {
			printedTarget := false
			for _, target := range targets {
				pos := target.getPos()
				if pos.row == row && pos.col == col {
					fmt.Printf("%c", targetChar)
					printedTarget = true
				}
			}
			if printedTarget {
				continue
			}
			switch n := node.(type) {
			case *tile:
				if n.isWall {
					fmt.Printf("%c", wallChar)
				} else {
					fmt.Printf("%c", openChar)
				}
			case *entity:
				if n.isGoblin {
					fmt.Printf("%c", goblinChar)
				} else {
					fmt.Printf("%c", elfChar)
				}
			}
		}
		fmt.Print("\n")
	}
}

func (b board) getNeighbors(pos coordinate) nodeList {
	neighbors := make(nodeList, 0, 4)
	for dRow := -1; dRow <= 1; dRow++ {
		for dCol := -1; dCol <= 1; dCol++ {
			if dRow == dCol || -dRow == dCol {
				continue
			}

			row := pos.row + dRow
			col := pos.col + dCol
			if row >= 0 && row < len(b) && col >= 0 && col < len(b[row]) {
				neighbors = append(neighbors, b[row][col])
			}
		}
	}

	return neighbors
}

func (b board) getWinner() winner {
	currentWinner := noWinner
	for row := range b {
		for _, memberNode := range b[row] {
			if entityNode, isEntity := memberNode.(*entity); isEntity {
				if currentWinner == goblinWinner && !entityNode.isGoblin {
					return noWinner
				} else if currentWinner == elfWinner && entityNode.isGoblin {
					return noWinner
				} else if currentWinner == noWinner && entityNode.isGoblin {
					currentWinner = goblinWinner
				} else if currentWinner == noWinner && !entityNode.isGoblin {
					currentWinner = elfWinner
				}
			}
		}
	}

	return currentWinner
}

func parseInput(rawBoard []string) (board, nodeList, error) {
	parsedBoard := make(board, len(rawBoard))
	entities := nodeList{}
	for row := range parsedBoard {
		parsedBoard[row] = make([]node, len(rawBoard[row]))
		for col, char := range rawBoard[row] {
			pos := coordinate{row, col}
			if char == wallChar || char == openChar {
				parsedBoard[row][col] = &tile{
					position: pos,
					isWall:   char == wallChar,
				}
			} else if char == elfChar || char == goblinChar {
				parsedBoard[row][col] = &entity{
					position: pos,
					health:   startingHealth,
					isGoblin: char == goblinChar,
				}
				entities = append(entities, parsedBoard[row][col])
			} else {
				return nil, nil, errors.New(malformedInputError)
			}
		}
	}

	return parsedBoard, entities, nil
}

func runSimulation(b board, entities nodeList) (winner, int) {
	roundCount := 0
	roundWinner := noWinner
	for roundWinner == noWinner {
		sort.Sort(entities)
		finishedRoundEarly := false
		for _, e := range entities {
			entity := e.(*entity)
			if entity.health <= 0 {
				continue
			}
			if !entity.attack(b) {
				entity.move(b)
				entity.attack(b)
				roundWinner = b.getWinner()
				if roundWinner != noWinner {
					finishedRoundEarly = true
					break
				}
			}
		}
		if !finishedRoundEarly {
			roundCount++
		}
	}
	healthTotal := 0
	for _, e := range entities {
		entity := e.(*entity)
		if entity.health > 0 {
			healthTotal += entity.health
		}
	}

	return roundWinner, roundCount * healthTotal
}

func part1(b board, entities nodeList) (outcome int) {
	_, outcome = runSimulation(b, entities)
	return
}

func main() {
	if len(os.Args) != 2 {
		return
	}

	inFile := os.Args[1]
	inFileContents, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}
	rawBoard := strings.Split(string(inFileContents), "\n")
	// trim tailing newline
	rawBoard = rawBoard[:len(rawBoard)-1]
	parsedBoard, entities, err := parseInput(rawBoard)
	if err != nil {
		panic(err)
	}

	fmt.Println(part1(parsedBoard, entities))
}
