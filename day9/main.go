package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	lineFormat          = "%d players; last marble is worth %d points"
	malformedInputError = "malformed input"
)

type node struct {
	value int
	next  *node
	prev  *node
}

func newCircularLinkedList() *node {
	n := new(node)
	n.next = n
	n.prev = n

	return n
}

func parseInput(input string) (numPlayers, numMarbles int, err error) {
	numMatched, err := fmt.Sscanf(input, lineFormat, &numPlayers, &numMarbles)
	if err != nil {
		return 0, 0, err
	} else if numMatched != 2 {
		return 0, 0, fmt.Errorf(malformedInputError)
	}

	return
}

func getMax(arr []int) (max int) {
	for _, item := range arr {
		if item > max {
			max = item
		}
	}

	return
}

func getElementByOffset(n *node, offset int) *node {
	cursor := n
	negative := false
	numJumps := offset
	if offset < 0 {
		numJumps *= -1
		negative = true
	}

	for i := 0; i < numJumps; i++ {
		if negative {
			cursor = cursor.prev
		} else {
			cursor = cursor.next
		}
	}

	return cursor
}

func runGame(numPlayers int, numMarbles int) int {
	scores := make([]int, numPlayers)
	currentMarble := newCircularLinkedList()
	currentPlayer := 0
	for nextMarble := 1; nextMarble <= numMarbles; nextMarble++ {
		if nextMarble%23 == 0 {
			removeMarble := getElementByOffset(currentMarble, -7)
			removeMarble.prev.next = removeMarble.next
			scores[currentPlayer] += nextMarble + removeMarble.value
			currentMarble = removeMarble.next
		} else {
			newMarble := &node{
				value: nextMarble,
				next:  currentMarble.next.next,
				prev:  currentMarble.next,
			}
			newMarble.next.prev = newMarble
			currentMarble.next.next = newMarble
			currentMarble = newMarble
		}
		currentPlayer = (currentPlayer + 1) % numPlayers
	}

	return getMax(scores)
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

	input := strings.TrimSuffix(string(inFileContents), "\n")
	numPlayers, numMarbles, err := parseInput(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(runGame(numPlayers, numMarbles))
	fmt.Println(runGame(numPlayers, numMarbles*100))
}
