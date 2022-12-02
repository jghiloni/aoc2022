package utils

import (
	"bufio"
	"io"
)

func ReaderToLines(in io.Reader) ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
