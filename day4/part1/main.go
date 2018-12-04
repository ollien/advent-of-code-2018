package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type guardAction int

const (
	startShift guardAction = iota
	fallAsleep
	wakeUp
)

const (
	timeFormat         = "2006-01-02 15:04"
	logLineFormat      = "[%s] %s"
	asleepTrigger      = "falls asleep"
	wakeupTrigger      = "wakes up"
	shiftTriggerFormat = "Guard #%d begins shift"
	malformedLineError = "malformed line"
)

type logLine struct {
	actionTime time.Time
	guardID    int
	action     guardAction
}

// A list that implements sort's Interface interface
type logData []logLine

func (data logData) Len() int {
	return len(data)
}

func (data logData) Swap(i int, j int) {
	data[i], data[j] = data[j], data[i]
}

func (data logData) Less(i int, j int) bool {
	return data[i].actionTime.Before(data[j].actionTime)
}

func parseLogLine(line string) (logLine, error) {
	lineComponents := strings.Split(line, "] ")
	rawTime, action := lineComponents[0][1:], lineComponents[1]
	parsedTime, err := time.Parse(timeFormat, rawTime)
	if err != nil {
		return logLine{}, err
	}

	lineInfo := logLine{actionTime: parsedTime, guardID: -1}
	if action == asleepTrigger {
		lineInfo.action = fallAsleep
		return lineInfo, nil
	} else if action == wakeupTrigger {
		lineInfo.action = wakeUp
		return lineInfo, nil
	}

	numMatched, err := fmt.Sscanf(action, shiftTriggerFormat, &lineInfo.guardID)
	if err != nil {
		return logLine{}, err
	} else if numMatched != 1 {
		return logLine{}, fmt.Errorf(malformedLineError)
	}

	return lineInfo, nil
}

func parseLog(logLines []string) ([]logLine, error) {
	result := make(logData, 0, len(logLines))
	for _, line := range logLines {
		lineInfo, err := parseLogLine(line)
		if err != nil {
			return nil, err
		}
		result = append(result, lineInfo)
	}

	sort.Sort(result)

	return result, nil
}

// Gets the amount of time the guard is asleep and their sleepiest time
func getSleepInfo(sleepTimes []int) (totalSleepCount int, sleepiestTime int) {
	sleepiestCount := 0
	for sleepTime, sleepCount := range sleepTimes {
		totalSleepCount += sleepCount
		if sleepCount > sleepiestCount {
			sleepiestCount = sleepCount
			sleepiestTime = sleepTime
		}
	}
	return
}

func constructSleepLog(log logData, sleepLog map[int][]int) {
	currentGuard := -1
	sleepTime := time.Time{}
	for _, lineInfo := range log {
		if lineInfo.guardID != -1 {
			currentGuard = lineInfo.guardID
		}
		if sleepLog[currentGuard] == nil {
			sleepLog[currentGuard] = make([]int, 60)
		}
		if lineInfo.action == fallAsleep {
			sleepTime = lineInfo.actionTime
		} else if lineInfo.action == wakeUp {
			for i := sleepTime.Minute(); i < lineInfo.actionTime.Minute(); i++ {
				sleepLog[currentGuard][i]++
			}
		}
	}
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
	logLines := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	logLines = logLines[:len(logLines)-1]

	sleepLog := make(map[int][]int)
	parsedLog, err := parseLog(logLines)
	if err != nil {
		panic(err)
	}

	constructSleepLog(parsedLog, sleepLog)

	sleepiestGuard := -1
	sleepiestGuardTime := 0
	sleepiestTimeForGuard := -1
	for guardID, guardInfo := range sleepLog {
		totalSleepTime, sleepiestTime := getSleepInfo(guardInfo)
		if totalSleepTime > sleepiestGuardTime {
			sleepiestGuard = guardID
			sleepiestGuardTime = totalSleepTime
			sleepiestTimeForGuard = sleepiestTime
		}
	}

	fmt.Println(sleepiestGuard * sleepiestTimeForGuard)
}
