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
		if binarySearch(totals, newTotal) != -1 {
			fmt.Println(newTotal)
			return
		}

		totals = insertSorted(totals, newTotal)
		lastTotal = newTotal
	}
}

func insertSorted(nums []int, newNum int) []int {
	for i, num := range nums {
		if newNum < num {
			return append(nums[:i], append([]int{newNum}, nums[i:]...)...)
		}
	}

	return append(nums, newNum)
}

func binarySearch(nums []int, target int) int {
	hitList := make([]int, len(nums))
	copy(hitList, nums)
	offset := 0
	for len(hitList) != 0 {
		midIndex := len(hitList) / 2
		if target > hitList[midIndex] {
			offset += midIndex + 1
			hitList = hitList[midIndex+1:]
		} else if target < hitList[midIndex] {
			hitList = hitList[:midIndex]
		} else {
			return midIndex + offset
		}
	}

	return -1
}
