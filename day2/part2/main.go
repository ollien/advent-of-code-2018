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
	for _, box1 := range boxes {
		for _, box2 := range boxes {
			diffCount := getNumDifferentLetters(box1, box2)
			if diffCount == 1 {
				fmt.Println(getLettersInCommon(box1, box2))
				return
			}
		}
	}
}

func getLettersInCommon(box1, box2 string) string {
	if len(box1) != len(box2) {
		return ""
	}

	commonLetters := ""
	for i, letter := range box1 {
		if box1[i] == box2[i] {
			commonLetters += string(letter)
		}
	}

	return commonLetters
}

func getNumDifferentLetters(box1, box2 string) int {
	if len(box1) != len(box2) {
		return -1
	}
	diffCount := 0
	commonString := getLettersInCommon(box1, box2)
	if len(commonString) == len(box1)-1 {
		diffCount++
	}

	return diffCount
}
