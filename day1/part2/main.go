package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
	rawNums := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawNums = rawNums[:len(rawNums)-1]

	totals := []int{0}
	lastTotal := 0
	for i := 0; ; i = (i + 1) % len(rawNums) {
		num, err := strconv.Atoi(rawNums[i])
		if err != nil {
			panic(err)
		}
		newTotal := lastTotal + num
		for _, total := range totals {
			if newTotal == total {
				fmt.Println(newTotal)
				return
			}
		}

		totals = append(totals, newTotal)
		lastTotal = newTotal
	}
}
