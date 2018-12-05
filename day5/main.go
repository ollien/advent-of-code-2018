package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func shouldAnihalate(chain string) bool {
	firstElement := rune(chain[0])
	secondElement := rune(chain[1])
	if firstElement == secondElement {
		return false
	}
	if firstElement == unicode.ToLower(secondElement) || firstElement == unicode.ToUpper(secondElement) {
		return true
	}

	return false
}

func findReactingIndex(polymer string) int {
	for i := 0; i < len(polymer)-1; i++ {
		if shouldAnihalate(polymer[i : i+2]) {
			return i
		}
	}

	return -1
}

func performReaction(polymer string, reactIndex int) string {
	result := polymer[:reactIndex]
	if reactIndex == len(polymer)-2 {
		return result
	} else {
		return result + polymer[reactIndex+2:]
	}

}

func part1(polymer string) int {
	reactedPolymer := polymer
	reactIndex := findReactingIndex(reactedPolymer)
	for reactIndex != -1 {
		reactedPolymer = performReaction(reactedPolymer, reactIndex)
		reactIndex = findReactingIndex(reactedPolymer)
	}

	return len(reactedPolymer)
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

	polymer := strings.TrimSuffix(string(inFileContents), "\n")
	fmt.Println(part1(polymer))
}
