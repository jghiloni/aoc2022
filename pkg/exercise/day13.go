package exercise

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"reflect"
	"sort"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day13 struct{}

func init() {
	Register("day13", day13{})
}

func (d day13) Part1(input io.Reader, output *log.Logger) (any, error) {
	q, err := d.loadPairs(input)
	output.Printf(`there are {{ colorize "bold;yellow" "%d" }} packet pairs`, q.Size())
	if err != nil {
		return nil, err
	}

	indices := []int{}
	index := 0
	for q.Size() > 0 {
		index++

		left, _ := q.Pop()
		right, _ := q.Pop()

		output.Printf(`LEFT : {{ colorize "bold;bright-cyan" "%v" }}`, left)
		output.Printf(`RIGHT: {{ colorize "bold;bright-green" "%v" }}`, right)

		if left.Compare(right, output) < 0 {
			indices = append(indices, index)
		}
	}

	output.Printf(`{{ colorize "bold;yellow" "%d" }} packet pairs in order: {{ colorize "bold;bright-magenta" "%v" }}`, len(indices), indices)
	return sumInts(indices...), nil
}

func (d day13) Part2(input io.Reader, output *log.Logger) (any, error) {
	q, err := d.loadPairs(input)
	if err != nil {
		return nil, err
	}

	div2 := packet([]any{[]any{float64(2)}})
	div6 := packet([]any{[]any{float64(6)}})

	q.PushAll(div2, div6)

	allPackets := q.Slice()
	sort.Slice(allPackets, func(i, j int) bool {
		pi, pj := allPackets[i], allPackets[j]
		return pi.Compare(pj, nil) < 0
	})

	var (
		i2 int
		i6 int
	)

	for i, p := range allPackets {
		output.Printf(`{{ colorize "cyan" "%03d" }}: {{ colorize "magenta" "%v" }}`, i+1, p)
		if i == 113 || i == 202 {
			output.Println("debugger")
		}
		if i2 == 0 && reflect.DeepEqual(p, div2) {
			i2 = i + 1
			continue
		}

		if i6 == 0 && reflect.DeepEqual(p, div6) {
			i6 = i + 1
			continue
		}
	}

	return i2 * i6, nil
}

func (day13) loadPairs(in io.Reader) (*utils.Queue[packet], error) {
	lines, err := utils.ReaderToLines(in)
	if err != nil {
		return nil, err
	}

	q := utils.NewQueue[packet]()
	for i := 0; i < len(lines); i += 3 {
		l, r := lines[i], lines[i+1]

		var (
			left  packet
			right packet
		)

		if err = json.Unmarshal([]byte(l), &left); err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(r), &right); err != nil {
			return nil, err
		}

		q.PushAll(left, right)
	}

	return q, nil
}

type packet []any
type packetIterator struct {
	current int
	frames  []any
}

func (x packet) Compare(y packet, out *log.Logger) int {
	xi := x.Frames()
	yi := y.Frames()

	if out == nil {
		out = log.New(io.Discard, "", 0)
	}

	var (
		xn bool
		yn bool
	)
	for {
		xn, yn = xi.Next(), yi.Next()
		if !xn || !yn {
			break
		}
		xc, xcraw := xi.Frame()
		yc, ycraw := yi.Frame()
		out.Printf(`Comparing {{ colorize "bold;bright-blue" "%v" }} and {{ colorize "bold;bright-red" "%v" }}`, xcraw, ycraw)
		if len(xc) == 1 && len(yc) == 1 {
			if fx, xok := xcraw.(float64); xok {
				if fy, yok := ycraw.(float64); yok {
					if fx == fy {
						continue
					}

					// will return either -1 or 1
					return int(fx-fy) / int(math.Abs(fx-fy))
				}
			}
		}

		if v := xc.Compare(yc, out); v != 0 {
			return v
		}
	}

	if !xn && yn {
		return -1
	}

	if xn && !yn {
		return 1
	}

	return 0
}

func (p packet) Frames() *packetIterator {
	return &packetIterator{
		current: 0,
		frames:  p,
	}
}

func (i *packetIterator) Next() bool {
	return i.current < len(i.frames)
}

func (i *packetIterator) Frame() (packet, any) {
	f := i.frames[i.current]
	i.current++
	if s, ok := f.([]any); ok {
		return packet(s), f
	}

	return packet([]any{f}), f
}

func sumInts(ints ...int) int {
	total := 0
	for _, i := range ints {
		total += i
	}

	return total
}
