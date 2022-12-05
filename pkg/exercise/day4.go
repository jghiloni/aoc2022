package exercise

import (
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day4 struct{}

type assignment struct {
	startSection int
	endSection   int
}

func (a assignment) size() int {
	return a.endSection - a.startSection
}

func (a assignment) contains(b assignment) bool {
	return b.startSection >= a.startSection &&
		b.endSection <= a.endSection
}

func assignmentFromString(str string) (assignment, error) {
	parts := strings.SplitN(str, "-", 2)
	a := assignment{}
	var err error
	if a.startSection, err = strconv.Atoi(parts[0]); err != nil {
		return a, err
	}

	if a.endSection, err = strconv.Atoi(parts[1]); err != nil {
		return a, err
	}

	return a, nil
}

func processLine(line string) (smaller assignment, larger assignment, err error) {
	parts := strings.Split(line, ",")
	smaller, err = assignmentFromString(parts[0])
	if err != nil {
		return
	}

	larger, err = assignmentFromString(parts[1])
	if err != nil {
		return
	}

	if larger.size() < smaller.size() {
		smaller, larger = larger, smaller
	}

	err = nil
	return
}

func hasTotalOverlap(line string) (bool, error) {
	smaller, larger, err := processLine(line)
	if err != nil {
		return false, err
	}

	return larger.contains(smaller), nil
}

func hasPartialOverlap(line string) (bool, error) {
	a, b, err := processLine(line)
	if err != nil {
		return false, err
	}

	switch {
	case a.startSection >= b.startSection && a.startSection <= b.endSection:
		return true, nil
	case a.endSection >= b.startSection && a.endSection <= b.endSection:
		return true, nil
	case b.startSection >= a.startSection && b.startSection <= a.endSection:
		return true, nil
	case b.endSection >= a.startSection && b.endSection <= a.endSection:
		return true, nil
	default:
		return false, nil
	}
}

func init() {
	Register("day4", day4{})
}

func (d day4) Part1(stdin io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(stdin)
	if err != nil {
		return nil, err
	}

	answer := 0
	for i, line := range lines {
		overlap, err := hasTotalOverlap(line)
		if err != nil {
			output.Printf(`{{ colorize "bold;red" "An error occurred: %v" }}`, err)
			return nil, err
		}

		if overlap {
			answer++
			output.Printf("There is total overlap in line %d: %s. Current total: %d", i+1, line, answer)
		}
	}

	return answer, nil
}

func (d day4) Part2(stdin io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(stdin)
	if err != nil {
		return nil, err
	}

	answer := 0
	for i, line := range lines {
		overlap, err := hasPartialOverlap(line)
		if err != nil {
			output.Printf(`{{ colorize "bold;red" "An error occurred: %v" }}`, err)
			return nil, err
		}

		if overlap {
			answer++
			output.Printf("There is partial overlap in line %d: %s. Current total: %d", i+1, line, answer)
		}
	}

	return answer, nil
}
