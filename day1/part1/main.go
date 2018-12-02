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

	total := 0
	for _, rawNum := range rawNums {
		num, err := strconv.Atoi(rawNum)
		if err != nil {
			panic(err)
		}

		total += num
	}
	fmt.Println(total)
}
