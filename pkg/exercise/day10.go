package exercise

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day10 struct{}

func init() {
	Register("day10", day10{})
}

type observer = func(*screen)

type screen struct {
	cycle   int
	x       int
	observe func(*screen)
}

func (s *screen) noop() {
	s.observe(s)
	s.cycle++
}

func (s *screen) addX(i int) {
	s.noop()
	s.noop()
	s.x += i
}

func (d day10) Part1(input io.Reader, output *log.Logger) (any, error) {
	observedSignalStrengths := []int{}

	observe := func(s *screen) {
		if s.cycle <= 220 && s.cycle%40 == 20 {
			output.Printf(`The strength at cycle {{colorize "bold;magenta" "%d"}} is {{colorize "bold;cyan" "%d"}}`, s.cycle, s.cycle*s.x)
			observedSignalStrengths = append(observedSignalStrengths, s.cycle*s.x)
		}

	}

	if err := d.runCommands(input, output, observe); err != nil {
		output.Printf(`{{ colorize "bold;red" "an error occurred: %v" }}`, err)
		return nil, err
	}

	answer := sum(observedSignalStrengths...)
	output.Printf(`The combined signal strengths total {{ colorize "bold;bright-yellow" "%d" }}`, answer)
	return answer, nil
}

const (
	lit   = '#'
	unlit = '.'
)

type crt struct {
	bitmap [][]byte
	sprite int
}

func newCrt() *crt {
	c := new(crt)
	c.reset()
	return c
}

func (c *crt) reset() {
	c.sprite = 1
	c.bitmap = make([][]byte, 6)

	for row := 0; row < 6; row++ {
		c.bitmap[row] = make([]byte, 40)
	}
}

func (c *crt) paint(reg int) {
	row := c.sprite / 40
	col := c.sprite % 40

	c.bitmap[row][col] = unlit
	if utils.Abs(reg-col) <= 1 {
		c.bitmap[row][col] = lit
	}
}

func (c *crt) String() string {
	lines := make([]string, len(c.bitmap))
	for i, row := range c.bitmap {
		lines[i] = string(row)
	}

	return strings.Join(lines, "\n")
}

const banner = `

{{ colorize "bold;underline;bright-white" "                DISPLAY                 " }}`

func (d day10) Part2(input io.Reader, output *log.Logger) (any, error) {
	c := newCrt()
	observe := func(s *screen) {
		c.sprite = s.cycle - 1
		c.paint(s.x)
	}

	if err := d.runCommands(input, output, observe); err != nil {
		return nil, err
	}

	fmt.Fprintln(output.Writer(), banner)
	fmt.Fprintf(output.Writer(), `{{ colorize "bold;bright-white" %q }}`, c)
	fmt.Fprint(output.Writer(), "\n\n")

	output.Println("Calculate the answer from the output")
	return "", nil
}

func (d day10) runCommands(input io.Reader, output *log.Logger, observe observer) error {
	scr := &screen{
		cycle:   1,
		x:       1,
		observe: observe,
	}

	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return err
	}

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		output.Printf(`command {{ colorize "bold;bright-green" %q }} at cycle {{ colorize "bold;bright-magenta" "%d" }}, register {{ colorize "bold;bright-cyan" "%d" }}`, line, scr.cycle, scr.x)
		switch line[0:4] {
		case "noop":
			scr.noop()
		case "addx":
			op, err := strconv.Atoi(line[5:])
			if err != nil {
				return err
			}
			scr.addX(op)
		}
	}
	return nil
}

func sum(nums ...int) int {
	x := 0
	for _, n := range nums {
		x += n
	}

	return x
}
