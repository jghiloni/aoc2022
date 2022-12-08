package exercise

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day8 struct{}

func init() {
	Register("day8", day8{})
}

var colorMap = map[utils.Direction]string{
	utils.Up:    "green",
	utils.Down:  "magenta",
	utils.Left:  "cyan",
	utils.Right: "red",
}

func (d day8) Part1(input io.Reader, output *log.Logger) (any, error) {
	canopy, err := utils.ReadMatrix(input)
	if err != nil {
		return nil, err
	}

	// This calculates the perimeter while eliminating duplicate counts
	totalVisible := (len(canopy) * 2) + ((len(canopy[0]) - 2) * 2)
	for row := 1; row < len(canopy)-1; row++ {
		for col := 1; col < len(canopy[row])-1; col++ {
			visible := false
			for _, direction := range utils.Directions {
				v, err := canopy.IsDirectionalMaximum(row, col, direction)
				if err != nil {
					output.Printf(`{{ colorize "bold;red" "An error occurred: %v" }}`, err)
					return nil, err
				}

				if v {
					output.Printf(`Tree at {{ colorize "bold;yellow" "%02d✕%02d" }} is visible from the {{ colorize "bold;bright-cyan" "%s" }}`, row+1, col+1, direction)
				}

				visible = visible || v
			}

			if visible {
				totalVisible++
			}
		}
	}

	output.Printf(`There are {{ colorize "bold;yellow" "%d" }} visible trees`, totalVisible)
	return totalVisible, nil
}

func (d day8) Part2(input io.Reader, output *log.Logger) (any, error) {
	canopy, err := utils.ReadMatrix(input)
	if err != nil {
		return nil, err
	}

	highestVisibilityScore := 0
	for row := 1; row < len(canopy)-1; row++ {
		for col := 1; col < len(canopy[row])-1; col++ {
			visibilityScore := 1
			msg := fmt.Sprintf(`The tree at {{ colorize "bold;yellow" "%02d✕%02d" }} has a visibility score of (`, row, col)
			for _, direction := range utils.Directions {
				dirScore := canopy.GetVisibilityScore(row, col, direction)
				msg = fmt.Sprintf(`%s{{ colorize "bold;bright-%s" "%02d" }} * `, msg, colorMap[direction], dirScore)

				visibilityScore *= dirScore
			}

			output.Printf(`%s) = {{ colorize "bold;bright-yellow" "%d" }}`, strings.TrimSuffix(msg, " * "), visibilityScore)
			if visibilityScore > highestVisibilityScore {
				highestVisibilityScore = visibilityScore
				output.Printf(`The current highest visibility score of {{ colorize "bold;bright-cyan" "%d" }} is at tree {{ colorize "bold;yellow" "%02d✕%02d" }}`, highestVisibilityScore, row, col)
			}
		}
	}

	return highestVisibilityScore, nil
}
