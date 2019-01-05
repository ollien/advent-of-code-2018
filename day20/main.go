package main

import (
	"container/heap"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

const (
	malformedInputError = "malformed input"
	northChar           = 'N'
	southChar           = 'S'
	eastChar            = 'E'
	westChar            = 'W'
	branchChar          = '|'
	branchStartChar     = '('
	branchEndChar       = ')'
	startChar           = '^'
	endChar             = '$'
	roomChar            = '.'
	noRoomChar          = ' '
	wallChar            = '#'
	verticalDoorChar    = '|'
	horizontalDoorChar  = '-'
	startPosChar        = 'X'
)

const (
	noDirection direction = iota
	northDirection
	southDirection
	eastDirection
	westDirection
)

type grid map[int]map[int]*node
type direction int
type nodeHeap []*node

type node struct {
	distance int
	north,
	south,
	east,
	west *node
}

type cursor struct {
	n        *node
	row, col int
}

type pathTree struct {
	path     []direction
	branches []*pathTree
	next     *pathTree
}

func newNode() *node {
	return &node{
		distance: math.MaxInt32,
	}
}

func newPathTree() *pathTree {
	return &pathTree{
		path:     []direction{},
		branches: []*pathTree{},
	}
}

func (c *cursor) updateCoord(dir direction) {
	switch dir {
	case northDirection:
		c.row--
	case eastDirection:
		c.col++
	case southDirection:
		c.row++
	case westDirection:
		c.col--
	}
}

// Attaches a ndoe to the graph, returning
func (n *node) attach(dir direction, newNode *node) {
	switch dir {
	case northDirection:
		n.north = newNode
		newNode.south = n
	case eastDirection:
		n.east = newNode
		newNode.west = n
	case southDirection:
		n.south = newNode
		newNode.north = n
	case westDirection:
		n.west = newNode
		newNode.east = n
	}
}

func (n *node) makeNeighborMap() map[direction]*node {
	neighbors := make(map[direction]*node, 4)
	if n.north != nil {
		neighbors[northDirection] = n.north
	}
	if n.east != nil {
		neighbors[eastDirection] = n.east
	}
	if n.south != nil {
		neighbors[southDirection] = n.south
	}
	if n.west != nil {
		neighbors[westDirection] = n.west
	}

	return neighbors
}

func (h *nodeHeap) Len() int {
	return len(*h)
}

func (h *nodeHeap) Less(i, j int) bool {
	return (*h)[i].distance < (*h)[j].distance
}

func (h *nodeHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *nodeHeap) Push(item interface{}) {
	*h = append(*h, item.(*node))
}

func (h *nodeHeap) Pop() interface{} {
	var tailNode *node
	tailNode, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]
	return tailNode
}

func (h *nodeHeap) updateDistance(n *node, newDistance int) {
	for i, storedNode := range *h {
		if storedNode == n {
			n.distance = newDistance
			heap.Fix(h, i)
			break
		}
	}
}

func (g grid) print() {
	g.printCursor(0, 0)
}

func (g grid) printCursor(rowCursor, colCursor int) {
	// Get all of the keys in sorted order
	rows := make([]int, 0, len(g))
	minRow := math.MaxInt32
	maxRow := 0
	minCol := math.MaxInt32
	maxCol := 0
	for row := range g {
		rows = append(rows, row)
		if row < minRow {
			minRow = row
		}
		if row > maxRow {
			maxRow = row
		}
		for col := range g[row] {
			if col < minCol {
				minCol = col
			}
			if col > maxCol {
				maxCol = col
			}
		}
	}
	sort.Ints(rows)

	for i, row := 0, minRow; row <= maxRow+1; i++ {
		if i%2 == 0 || row == maxRow+1 {
			for j, col := 0, minCol; col <= maxCol; j++ {
				if j%2 == 0 {
					fmt.Printf("%c", wallChar)
					continue
				}
				if node, haveNode := g[row-1][col]; !haveNode || node.south == nil {
					fmt.Printf("%c", wallChar)
				} else {
					fmt.Printf("%c", horizontalDoorChar)
				}
				col++
			}
			fmt.Printf("%c", wallChar)
		} else {
			for col := minCol; col <= maxCol; col++ {
				if col == minCol {
					fmt.Printf("%c", wallChar)
				}
				node, haveNode := g[row][col]
				if haveNode {
					if row == 0 && col == 0 || row == rowCursor && col == colCursor {
						fmt.Printf("%c", startPosChar)
					} else {
						fmt.Printf("%c", roomChar)
					}
					if node.east == nil {
						fmt.Printf("%c", wallChar)
					} else {
						fmt.Printf("%c", verticalDoorChar)
					}
				} else {
					fmt.Printf("%c ", noRoomChar)
				}
			}
		}
		if i%2 != 0 || row == maxRow+1 {
			row++
		}
		fmt.Print("\n")
	}
}

func (g grid) flatten() []*node {
	nodes := []*node{}
	for row := range g {
		for col := range g[row] {
			nodes = append(nodes, g[row][col])
		}
	}
	return nodes
}

func getClosingBranchIndex(s string) (int, error) {
	skipCount := 0
	for i, char := range s[1:] {
		if char == branchStartChar {
			skipCount++
		} else if char == branchEndChar && skipCount == 0 {
			return i + 1, nil
		} else if char == branchEndChar {
			skipCount--
		}
	}

	return -1, nil
}

func getDirectionFromChar(char byte) (direction, error) {
	switch char {
	case northChar:
		return northDirection, nil
	case eastChar:
		return eastDirection, nil
	case southChar:
		return southDirection, nil
	case westChar:
		return westDirection, nil
	}

	return noDirection, errors.New(malformedInputError)
}

func parseInput(rawRegex string) ([]*pathTree, error) {
	head := newPathTree()
	paths := []*pathTree{}
	treeCursor := head
	var prevCursor *pathTree

	for i := 0; i < len(rawRegex); i++ {
		char := rawRegex[i]
		if char == branchStartChar {
			branchEndOffset, err := getClosingBranchIndex(rawRegex[i:])
			if err != nil {
				return nil, err
			}
			newBranches, err := parseInput(rawRegex[i+1:])
			if err != nil {
				return nil, err
			}
			treeCursor.branches = newBranches
			paths = append(paths, treeCursor)
			i += branchEndOffset
			prevCursor = treeCursor
			treeCursor = newPathTree()
		} else if char == branchChar {
			branch := newPathTree()
			branch.path = make([]direction, len(treeCursor.path))
			copy(branch.path, treeCursor.path)
			treeCursor.path = treeCursor.path[:0]
			paths = append(paths, branch)
		} else if char == branchEndChar {
			if len(treeCursor.path) > 0 || rawRegex[i-1] == branchChar {
				branch := newPathTree()
				branch.path = treeCursor.path
				paths = append(paths, branch)
			}
			return paths, nil
		} else {
			if prevCursor != nil && prevCursor.next == nil {
				prevCursor.next = treeCursor
			}
			dir, err := getDirectionFromChar(char)
			if err != nil {
				return nil, err
			}
			treeCursor.path = append(treeCursor.path, dir)
		}
	}

	return []*pathTree{head}, nil
}

func makeGraph(rawRegex string, headCursor cursor, roomGrid grid) (*node, grid, int, error) {
	graphCursor := headCursor
	// Keep track of the spaces we've already allocated so we can circle back to existing rooms
	// While this typing might look excessive, it makes assignment a bit more optimal
	if roomGrid == nil {
		roomGrid = make(grid)
		roomGrid[headCursor.row] = make(map[int]*node)
		roomGrid[headCursor.row][headCursor.col] = headCursor.n
	}
	for i := 0; i < len(rawRegex); i++ {
		char := rawRegex[i]
		if char == branchStartChar {
			_, _, skipCount, err := makeGraph(rawRegex[i+1:], graphCursor, roomGrid)
			if err != nil {
				return nil, nil, 0, err
			}
			i += skipCount
		} else if char == branchChar {
			// By resetting to the head no matter what when we hit a branch char, we are abusing a property of the detours that they always seem to circle back upon themselves.
			graphCursor = headCursor
		} else if char == branchEndChar {
			return headCursor.n, roomGrid, i + 1, nil
		} else {
			dir, err := getDirectionFromChar(char)
			if err != nil {
				return nil, nil, 0, err
			}

			graphCursor.updateCoord(dir)
			if roomGrid[graphCursor.row] == nil {
				roomGrid[graphCursor.row] = make(map[int]*node)
			}
			if roomGrid[graphCursor.row][graphCursor.col] == nil {
				roomGrid[graphCursor.row][graphCursor.col] = newNode()
			}
			gridNode := roomGrid[graphCursor.row][graphCursor.col]
			graphCursor.n.attach(dir, gridNode)
			graphCursor.n = gridNode
		}
	}

	return headCursor.n, roomGrid, 0, nil
}

// getShortestDistances gets the shortest distance to every node from a given head.
// Currently implemented as Dijkstra's algorithm
func getShortestDistances(head *node, nodeList []*node) map[*node]int {
	distances := map[*node]int{head: 0}
	distanceHeap := nodeHeap{}
	heap.Init(&distanceHeap)
	for _, storedNode := range nodeList {
		heap.Push(&distanceHeap, storedNode)
	}
	for len(distanceHeap) > 0 {
		nextNode := heap.Pop(&distanceHeap).(*node)
		neighbors := nextNode.makeNeighborMap()
		for _, neighbor := range neighbors {
			distanceCandidate := nextNode.distance + 1
			if distanceCandidate < neighbor.distance {
				distanceHeap.updateDistance(neighbor, distanceCandidate)
				distances[neighbor] = distanceCandidate
				neighbor.distance = distanceCandidate
			}
		}
	}

	return distances
}

func part1(distances map[*node]int) int {
	maxDistance := 0
	for _, distance := range distances {
		if distance > maxDistance {
			maxDistance = distance
		}
	}

	return maxDistance
}

func part2(distances map[*node]int) (count int) {
	for _, distance := range distances {
		if distance >= 1000 {
			count++
		}
	}

	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main in_file")
		return
	}

	inFile := os.Args[1]
	inputFileContents, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}
	rawRegex := string(inputFileContents)
	if rawRegex[0] != startChar || rawRegex[len(rawRegex)-2] != endChar {
		panic(malformedInputError)
	}

	rawRegex = rawRegex[1 : len(rawRegex)-2]
	tree, err := parseInput(string(rawRegex))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *(tree[0]))
}
