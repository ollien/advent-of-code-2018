package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func sumUntil(arr []int, max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += arr[i]
	}

	return total
}

func parseTree(tree []int, numNodes int) (int, int) {
	total := 0
	cursor := 0
	for i := 0; i < numNodes; i++ {
		children, metadataCount := tree[0], tree[1]
		// Remove the tree header
		tree = tree[2:]
		n, subtotal := parseTree(tree, children)
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

func part1(tree []int) int {
	_, total := parseTree(tree, 1)
	return total
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
}
