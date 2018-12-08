package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	instructionStringFormat = "Step %s must be finished before step %s can begin."
	malformedError          = "malformed input"
	noeEntrypointError      = "no entrypoint"
	numWorkers              = 5
)

type instructionList map[string][]string
type workerList []workerJob

type workerJob struct {
	name           string
	stepsRemaining int
}

func parseInstructions(rawInstructions []string) (instructionList, error) {
	instructions := make(instructionList)
	for _, rawInstruction := range rawInstructions {
		var dependencyName, instructionName string
		numMatched, err := fmt.Sscanf(rawInstruction, instructionStringFormat, &dependencyName, &instructionName)
		if err != nil {
			return nil, err
		} else if numMatched != 2 {
			return nil, fmt.Errorf(malformedError)
		}
		_, hasDependency := instructions[dependencyName]
		if !hasDependency {
			instructions[dependencyName] = make([]string, 0)
		}
		_, hasInstruction := instructions[instructionName]
		if hasInstruction {
			instructions[instructionName] = append(instructions[instructionName], dependencyName)
		} else {
			instructions[instructionName] = []string{dependencyName}
		}
	}

	for instructionName := range instructions {
		sort.Strings(instructions[instructionName])
	}

	return instructions, nil
}

func (instructions instructionList) findReadySteps() []string {
	entrypointCandidates := make([]string, 0)
	for instructionName := range instructions {
		if len(instructions[instructionName]) == 0 {
			entrypointCandidates = append(entrypointCandidates, instructionName)
		}
	}

	sort.Strings(entrypointCandidates)

	return entrypointCandidates
}

func (instructions instructionList) markAsDone(doneInstructionName string) {
	for instructionName, instruction := range instructions {
		doneIndex := -1
		for i, dependencyName := range instruction {
			if dependencyName == doneInstructionName {
				doneIndex = i
				break
			}
		}
		if doneIndex != -1 {
			instructions[instructionName] = append(instruction[:doneIndex], instruction[doneIndex+1:]...)
		}
	}

	delete(instructions, doneInstructionName)
}

func (instructions instructionList) resolveDependencies() string {
	instructionSet := ""
	for len(instructions) > 0 {
		for _, readyStep := range instructions.findReadySteps() {
			instructions.markAsDone(readyStep)
			instructionSet += readyStep
		}
	}

	return instructionSet
}

func (workers workerList) findReadyWorker() int {
	for worker, job := range workers {
		if job.stepsRemaining == 0 {
			return worker
		}
	}

	return -1
}

func (workers workerList) allIdle() bool {
	for _, job := range workers {
		if job.stepsRemaining > 0 {
			return false
		}
	}

	return true
}

func (workers workerList) getDoneWorkers() []int {
	jobs := make([]int, 0)
	for i, job := range workers {
		if job.stepsRemaining == 0 && job.name != "" {
			jobs = append(jobs, i)
		}
	}

	return jobs
}

func (workers workerList) makingDependency(desiredJob string, instructions instructionList) bool {
	for _, dependency := range instructions[desiredJob] {
		for _, job := range workers {
			if job.name == dependency {
				return true
			}
		}
	}

	return false
}

func (workers workerList) work() {
	stepSize := workers.getMinStepsRemaining()
	for i := range workers {
		if workers[i].stepsRemaining > 0 {
			workers[i].stepsRemaining -= stepSize
		}
	}
}

func (workers workerList) getMinStepsRemaining() int {
	minTime := math.MaxInt32
	allZero := true
	for _, worker := range workers {
		if worker.stepsRemaining != 0 {
			allZero = false
		}
		if worker.stepsRemaining < minTime && worker.stepsRemaining > 0 {
			minTime = worker.stepsRemaining
		}
	}

	if allZero {
		return 0
	} else {
		return minTime
	}
}

// Relieve any workers who have finished their work
func (workers workerList) relieve() {
	doneJobs := workers.getDoneWorkers()
	for _, doneJob := range doneJobs {
		workers[doneJob].name = ""
	}
}

func part1(instructions instructionList) string {
	return instructions.resolveDependencies()
}

func part2(instructions instructionList) int {
	allInstructions := make(instructionList)
	for instruction, dependencies := range instructions {
		allInstructions[instruction] = make([]string, len(dependencies))
		copy(allInstructions[instruction], dependencies)
	}
	time := 0
	workers := make(workerList, numWorkers)
	workQueue := instructions.findReadySteps()
	// While workers are working or there is new work to be done
	for len(workQueue) > 0 || !workers.allIdle() {
		workers.work()
		workers.relieve()
		if len(workQueue) == 0 {
			workQueue = instructions.findReadySteps()
		}
		availableWorker := workers.findReadyWorker()
		// Store all items that cannot be worked on yet
		invalidItems := make([]string, 0)
		for len(workQueue) > 0 && availableWorker != -1 {
			var job string
			job, workQueue = workQueue[0], workQueue[1:]
			if workers.makingDependency(job, allInstructions) {
				invalidItems = append(invalidItems, job)
				continue
			}
			workers[availableWorker].name = job
			workers[availableWorker].stepsRemaining = int(rune(job[0])-'A') + 61
			instructions.markAsDone(job)
			availableWorker = workers.findReadyWorker()
		}
		// Add the items that cannot be worked on yet to the queue
		workQueue = append(invalidItems, workQueue...)
		time += workers.getMinStepsRemaining()
	}

	return time
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
	rawInstructions := strings.Split(string(inFileContents), "\n")
	// trim trailing newline
	rawInstructions = rawInstructions[:len(rawInstructions)-1]

	instructions, err := parseInstructions(rawInstructions)
	if err != nil {
		panic(err)
	}
	fmt.Println(part1(instructions))

	// Reparse after depleting instruction list
	instructions, err = parseInstructions(rawInstructions)
	if err != nil {
		panic(err)
	}
	fmt.Println(part2(instructions))

}
