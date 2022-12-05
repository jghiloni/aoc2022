package exercise

import (
	"io"
	"log"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day3 struct{}

// add a : at the beginning to 1-index the letters
const priorities = ":abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	Register("day3", day3{})
}

func (d day3) Part1(stdin io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(stdin)
	if err != nil {
		return nil, err
	}

	totalPriorities := 0
	for _, line := range lines {
		intersection := getCompartmentIntersection(line)
		output.Printf("Compartments share the following items: %q\n", intersection)
		totalPriorities += calculatePriorities(intersection)
	}

	return totalPriorities, nil
}

func (d day3) Part2(stdin io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(stdin)
	if err != nil {
		return nil, err
	}

	totalPriorities := 0
	for i := 0; i < len(lines); i += 3 {
		badge := getBadge(lines[i], lines[i+1], lines[i+2])
		if badge <= 0 {
			continue
		}

		output.Printf("Bags %d, %d, and %d have badge %s\n", i, i+1, i+2, string(badge))
		totalPriorities += strings.IndexRune(priorities, badge)
	}

	return totalPriorities, nil
}

func getCompartmentIntersection(line string) string {
	halfway := len(line) / 2
	part1, part2 := line[0:halfway], line[halfway:]

	seen := ""
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(part2, r) && !strings.ContainsRune(seen, r) {
			seen += string(r)
			return r
		}

		return -1
	}, part1)
}

func calculatePriorities(items string) int {
	total := 0
	for _, r := range items {
		total += strings.IndexRune(priorities, r)
	}

	return total
}

func getBadge(bag1, bag2, bag3 string) rune {
	for _, r := range bag1 {
		if strings.ContainsRune(bag2, r) && strings.ContainsRune(bag3, r) {
			return r
		}
	}

	return -1
}
