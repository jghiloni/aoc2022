package exercise

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
)

type day6 struct{}

func init() {
	Register("day6", day6{})
}

func (d day6) Part1(input io.Reader, output *log.Logger) (any, error) {
	buf := bufio.NewReader(input)
	return findProtocolMarker("transmission", 4, buf, output)
}

func (d day6) Part2(input io.Reader, output *log.Logger) (any, error) {
	buf := bufio.NewReader(input)
	return findProtocolMarker("message", 14, buf, output)
}

func findProtocolMarker(markerType string, markerLen int, protocolBuffer *bufio.Reader, output *log.Logger) (int, error) {
	cur := make([]rune, markerLen)
	var err error

	var idx int
	for idx = 0; idx < markerLen; idx++ {
		cur[idx], _, err = protocolBuffer.ReadRune()
		if err != nil {
			return -1, err
		}
	}

	var latest rune
	for {
		if isUnique(string(cur)) {
			output.Printf(`Current protocol buffer {{ colorize "bold;yellow" %q }} is unique and indicates the start of %s`, string(cur), markerType)
			return idx, nil
		}

		latest, _, err = protocolBuffer.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				break
			}
			output.Printf(`{{ colorize "bold;red" "An error occurred reading the protocol buffer: %v" }}`, err)
			return -1, err
		}

		cur = append(cur[1:], latest)
		idx++
	}

	err = fmt.Errorf("no start of protocol %s found", markerType)
	output.Printf(`{{ colorize "red" %q }}`, err.Error())
	return -1, err
}

func isUnique(s string) bool {
	found := map[rune]bool{}
	for _, r := range s {
		if found[r] {
			return false
		}
		found[r] = true
	}

	return true
}
