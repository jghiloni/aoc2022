package exercise

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day5 struct{}

func init() {
	Register("day5", day5{})
}

func (d day5) Part1(input io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	stacks, actions := prepareStacks(lines)
	for _, action := range actions {
		if stacks, err = performActionWithCrateMaster9000(action, stacks); err != nil {
			output.Printf(`{{ colorize "red" "An error occurred performing action [%s]: %v" }}`, action, err)
			return nil, err
		}

		output.Printf(`Action %q on the CrateMaster 9000 resulted in Top Line {{ colorize "bold;bright-green" %q }}`, action, topLine(stacks))
	}

	return topLine(stacks), nil
}

func (d day5) Part2(input io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	stacks, actions := prepareStacks(lines)
	for _, action := range actions {
		if stacks, err = performActionWithCrateMaster9001(action, stacks); err != nil {
			output.Printf(`{{ colorize "red" "An error occurred performing action %q: %v" }}`, action, err)
			return nil, err
		}

		output.Printf(`Action %q on the CrateMaster 9001 resulted in Top Line {{ colorize "bold;bright-green" %q }}`, action, topLine(stacks))
	}

	return topLine(stacks), nil
}

func prepareStacks(lines []string) ([]*utils.Stack[byte], []string) {
	initialInput := []string{}
	stacks := []*utils.Stack[byte]{}

	for len(lines) > 0 {
		line := lines[0]
		lines = lines[1:]
		if strings.TrimSpace(line) == "" {
			break
		}

		initialInput = append([]string{line}, initialInput...)
	}

	for i := 1; i < len(initialInput[0]); i += 4 {
		stack := utils.NewStack[byte]()
		for j := 1; j < len(initialInput); j++ {
			if !unicode.IsUpper(rune(initialInput[j][i])) {
				continue
			}

			stack.Push(initialInput[j][i])
		}
		stacks = append(stacks, stack)
	}

	return stacks, lines
}

func performActionWithCrateMaster9000(line string, stacks []*utils.Stack[byte]) ([]*utils.Stack[byte], error) {
	var (
		iterations int
		source     int
		target     int
	)

	if _, err := fmt.Sscanf(line, "move %d from %d to %d", &iterations, &source, &target); err != nil {
		return stacks, err
	}

	for i := 0; i < iterations; i++ {
		b, _ := stacks[source-1].Pop()
		stacks[target-1].Push(b)
	}

	return stacks, nil
}

func performActionWithCrateMaster9001(line string, stacks []*utils.Stack[byte]) ([]*utils.Stack[byte], error) {
	var (
		iterations int
		source     int
		target     int
	)

	if _, err := fmt.Sscanf(line, "move %d from %d to %d", &iterations, &source, &target); err != nil {
		return stacks, err
	}

	toMove := make([]byte, iterations)
	for i := 0; i < iterations; i++ {
		toMove[i], _ = stacks[source-1].Pop()
	}

	for i := iterations - 1; i >= 0; i-- {
		stacks[target-1].Push(toMove[i])
	}

	return stacks, nil
}

func topLine(stacks []*utils.Stack[byte]) string {
	b := &bytes.Buffer{}
	for _, s := range stacks {
		b.WriteByte(s.Peek())
	}

	return b.String()
}
