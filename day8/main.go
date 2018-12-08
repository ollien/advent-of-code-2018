package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

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

func parseTree(tree []int, numNodes int) (int, int, []node) {
	nodes := make([]node, numNodes)
	cursor := 0
	total := 0
	for i := 0; i < numNodes; i++ {
		numChildren, metadataCount := tree[0], tree[1]
		// Remove the tree header
		tree = tree[2:]
		n, subtotal, children := parseTree(tree, numChildren)
		// Remove the part that the subtree parsed
		tree = tree[n:]
		// Get the data total of the metadata
		metadataValue := sumUntil(tree, metadataCount)
		total += subtotal + metadataValue
		if numChildren == 0 {
			// Extract the value of the metadata
			nodes[i].value = metadataValue
		} else {
			nodes[i].value = getChildValuesFromMetadata(tree, metadataCount, children)
			// Copy the values of the children into our node
			nodes[i].childValues = make([]int, numChildren)
			for j := 0; j < numChildren; j++ {
				nodes[i].childValues[j] = children[j].value
			}
		}
		// Remove the metadata from the tree and advance the cursor accordingly
		tree = tree[metadataCount:]
		cursor += metadataCount + 2 + n
	}

	return cursor, total, nodes
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
	_, total, rootedTree := parseTree(tree, 1)
	fmt.Println(total)
	fmt.Println(rootedTree[0].value)
}
