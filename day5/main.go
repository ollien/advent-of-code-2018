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

func performReaction(polymer string) int {
	reaction_buffer := make([]rune, 0, len(polymer))
	for _, chr := range polymer {
		if len(reaction_buffer) > 0 && shouldAnihalate(string(chr)+string(reaction_buffer[len(reaction_buffer)-1])) {
			reaction_buffer = reaction_buffer[:len(reaction_buffer)-1]
		} else {
			reaction_buffer = append(reaction_buffer, chr)
		}
	}

	return len(reaction_buffer)
}

func part1(polymer string) int {
	return performReaction(polymer)
}

func part2(polymer string) int {
	smallestLen := len(polymer)
	for element := 'a'; element <= 'z'; element++ {
		strippedPolymer := polymer
		strippedPolymer = strings.Replace(strippedPolymer, string(element), "", -1)
		strippedPolymer = strings.Replace(strippedPolymer, strings.ToUpper(string(element)), "", -1)
		reactedLen := part1(strippedPolymer)
		if reactedLen < smallestLen {
			smallestLen = reactedLen
		}
	}

	return smallestLen
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
	fmt.Println(part2(polymer))
}
