package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// only used for part 2
type node struct {
	value       int
	childValues []int
}

func sumUntil(arr []int, max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += arr[i]
	}

	return total
}

func findTreeTotal(tree []int, numNodes int) (int, int) {
	total := 0
	cursor := 0
	for i := 0; i < numNodes; i++ {
		children, metadataCount := tree[0], tree[1]
		// Remove the tree header
		tree = tree[2:]
		n, subtotal := findTreeTotal(tree, children)
		// Remove the part that the subtree parsed
		tree = tree[n:]
		// Sum and remove the metadata part from the tree
		total += sumUntil(tree, metadataCount) + subtotal
		tree = tree[metadataCount:]
		// Advance the cursor that we will return past the metadata and the subtrees
		cursor += metadataCount + 2 + n
	}

	return cursor, total
}

func getChildValuesFromMetadata(tree []int, numValues int, children []node) int {
	total := 0
	for i := 0; i < numValues; i++ {
		childIndex := tree[i] - 1
		if childIndex < len(children) {
			total += children[childIndex].value
		}
	}

	return total
}

func calculateRootNode(tree []int, numNodes int) (int, []node) {
	nodes := make([]node, numNodes)
	cursor := 0
	for i := 0; i < numNodes; i++ {
		numChildren, metadataCount := tree[0], tree[1]
		// Remove the tree header
		tree = tree[2:]
		n, children := calculateRootNode(tree, numChildren)
		// Remove the part that the subtree parsed
		tree = tree[n:]
		if numChildren == 0 {
			// Extract the value of the metadata
			nodes[i].value = sumUntil(tree, metadataCount)
		} else {
			nodes[i].value = getChildValuesFromMetadata(tree, metadataCount, children)
			// Copy the values of the children into our node
			nodes[i].childValues = make([]int, numChildren)
			for j := 0; j < numChildren; j++ {
				nodes[i].childValues[j] = children[j].value
			}
		}
		tree = tree[metadataCount:]
		// Advance the cursor that we will return past the metadata and the subtrees
		cursor += metadataCount + 2 + n
	}

	return cursor, nodes
}

func part1(tree []int) int {
	_, total := findTreeTotal(tree, 1)
	return total
}

func part2(tree []int) int {
	_, rootedTree := calculateRootNode(tree, 1)
	return rootedTree[0].value
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

	rawTree := strings.TrimSuffix(string(inFileContents), "\n")
	tree := make([]int, 0)
	for _, item := range strings.Split(rawTree, " ") {
		result, err := strconv.Atoi(item)
		if err != nil {
			panic(err)
		}
		tree = append(tree, result)
	}
	fmt.Println(part1(tree))
	fmt.Println(part2(tree))
}
