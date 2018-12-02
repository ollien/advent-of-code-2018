package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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
	boxes := strings.Split(string(inFileContents), "\n")
	// trim tailing newline
	boxes = boxes[:len(boxes)-1]
	twoCount, threeCount := 0, 0
	for _, box := range boxes {
		twoLetter, threeLetter := getLetters(box)
		if twoLetter != 0 {
			twoCount++
		}
		if threeLetter != 0 {
			threeCount++
		}
	}
	checksum := threeCount * twoCount
	fmt.Println(checksum)
}

// getLetters returns (letter that appears twice, letter that appears thrice)
func getLetters(boxString string) (rune, rune) {
	counts := make(map[rune]int)
	for _, letter := range boxString {
		if _, ok := counts[letter]; !ok {
			counts[letter] = 0
		}
		counts[letter]++
	}
	// Get a letter that occurs three times or two times
	var twoLetter, threeLetter rune
	for letter, count := range counts {
		if count == 2 {
			twoLetter = letter
		} else if count == 3 {
			threeLetter = letter
		}
	}

	return twoLetter, threeLetter
}
