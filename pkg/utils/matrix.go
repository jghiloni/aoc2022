package utils

import (
	"errors"
	"io"
	"sort"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "top"
	case Down:
		return "bottom"
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return "INVALID"
	}
}

var Directions = []Direction{Up, Down, Left, Right}
var ErrOutOfBounds = errors.New("out of bounds")

type TreeCanopy [][]byte

func ReadMatrix(input io.Reader) (TreeCanopy, error) {
	m := make([][]byte, 0, 50)

	lines, err := ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		m = append(m, []byte(line))
	}

	return TreeCanopy(m), nil
}

func (t TreeCanopy) IsDirectionalMaximum(row, col int, direction Direction) (bool, error) {
	if row < 1 || row >= len(t)-1 {
		return false, ErrOutOfBounds
	}

	if col < 1 || col >= len(t[row])-1 {
		return false, ErrOutOfBounds
	}

	sightLine := t.getSightLine(row, col, direction)
	sort.Sort(sort.Reverse(sort.IntSlice(sightLine)))

	digit := int(t[row][col] - '0')
	return digit > sightLine[0], nil
}

func (t TreeCanopy) getSightLine(row, col int, direction Direction) []int {
	var (
		incrementor int
		begin       int
		end         int
		line        []byte
	)

	// make the end one PAST the last allowable choice to account for the edge
	switch direction {
	case Up:
		incrementor = -1
		begin = row - 1
		end = -1
		line = t.getColumn(col)
	case Down:
		incrementor = 1
		begin = row + 1
		end = len(t)
		line = t.getColumn(col)
	case Left:
		incrementor = -1
		begin = col - 1
		end = -1
		line = t.getRow(row)
	case Right:
		incrementor = 1
		begin = col + 1
		end = len(t[col])
		line = t.getRow(row)
	}

	sightLine := make([]int, 0, len(line)-begin)
	for i := begin; i != end; i += incrementor {
		sightLine = append(sightLine, int(line[i]-'0'))
	}

	return sightLine
}

func (t TreeCanopy) GetVisibilityScore(row, col int, direction Direction) int {
	curHeight := int(t[row][col] - '0')
	sightLine := t.getSightLine(row, col, direction)

	var distance int
	var height int
	for distance, height = range sightLine {
		if curHeight <= height {
			break
		}
	}

	return distance + 1
}

func (t TreeCanopy) getRow(row int) []byte {
	return t[row]
}

func (t TreeCanopy) getColumn(col int) []byte {
	column := make([]byte, len(t))
	for i := range t {
		column[i] = t[i][col]
	}

	return column
}
