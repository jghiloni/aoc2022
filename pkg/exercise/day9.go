package exercise

import (
	"fmt"
	"io"
	"log"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day9 struct{}

var startingPoint = utils.Coordinate{X: 0, Y: 0}

func init() {
	Register("day9", day9{})
}

func (d day9) Part1(input io.Reader, output *log.Logger) (any, error) {
	return runCommands(input, output, 2)
}

func (d day9) Part2(input io.Reader, output *log.Logger) (any, error) {
	return runCommands(input, output, 10)
}

func runCommands(input io.Reader, output *log.Logger, numKnots int) (any, error) {
	rope := utils.NewRopeBridge(numKnots)

	commands, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	for _, command := range commands {
		var (
			direction string
			distance  int
		)
		fmt.Sscanf(command, "%s %d", &direction, &distance)
		p1 := rope.String()
		rope.MoveHead(direction, distance)
		p2 := rope.String()

		output.Printf(`Move %s {{ colorize "bold;bright-cyan" %q }} to %s`, p1, command, p2)
	}

	visits := rope.UniqueVisits()
	answer := visits[len(visits)-1]
	output.Printf(`The tail of the rope visited {{ colorize "bold;yellow" "%d" }} unique locations`, answer)

	return answer, nil
}
