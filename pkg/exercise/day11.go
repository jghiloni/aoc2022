package exercise

import (
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

const self int64 = 1<<63 - 1

var monkeyColors = []string{
	"bg-white;black",
	"bold;red",
	"bold;bright-green",
	"bold;bright-yellow",
	"bold;bright-blue",
	"bold;bright-magenta",
	"bold;bright-cyan",
	"bold;bright-white",
}

type day11 struct{}

type monkeyOp = func(int64) int64
type monkeyTest = func(int64) bool

type monkey struct {
	inspections int64
	worries     *utils.Queue[int64]
	operation   monkeyOp
	test        monkeyTest
	passTarget  int
	failTarget  int
}

type troop []*monkey

func init() {
	Register("day11", day11{})
}

func (d day11) Part1(input io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	inputQueue := utils.NewQueue(lines...)
	monkeys, err := investigateMonkeys(inputQueue, nil)
	if err != nil {
		return nil, err
	}

	return d.runCommands(monkeys, output, 20, func(i int64) int64 {
		return i / 3
	}), nil
}

func (d day11) Part2(input io.Reader, output *log.Logger) (any, error) {
	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	inputQueue := utils.NewQueue(lines...)

	divisors := []int64{}
	monkeys, err := investigateMonkeys(inputQueue, func(line string) {
		if strings.HasPrefix(line, "Test: divisible by ") {
			divStr := strings.TrimPrefix(line, "Test: divisible by ")
			divisor, err := strconv.ParseInt(divStr, 10, 64)
			if err != nil {
				return
			}
			divisors = append(divisors, divisor)
		}
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(divisors, func(i, j int) bool {
		return divisors[i] < divisors[j]
	})

	lcd := divisors[len(divisors)-1]
	for i := 0; i < len(divisors)-1; i++ {
		for j := i; j < len(divisors)-1; j++ {
			if lcd%divisors[j] == 0 {
				continue
			}

			lcd *= divisors[i]
			break
		}
	}

	output.Printf("We have test divisors of %v, will use a stress management factor of %d", divisors, lcd)
	return d.runCommands(monkeys, log.New(io.Discard, "", 0), 10000, func(i int64) int64 {
		return i % lcd
	}), nil
}

func (d day11) runCommands(monkeys troop, output *log.Logger, rounds int, destressor monkeyOp) int64 {
	monkeys.log(output)
	for i := 0; i < rounds; i++ {
		monkeys.monkeyAround(destressor)
		output.Println()
		output.Printf(`After {{ colorize "bold;bright-cyan" "%d" }} rounds, this is the monkeys' business`, i)
		monkeys.log(output)
	}

	return monkeys.monkeyBusiness(output)
}

func investigateMonkeys(inputQueue *utils.Queue[string], observer func(string)) (troop, error) {
	monkeys := make([]*monkey, 0, 8)
	var m *monkey
	for line, size := inputQueue.Pop(); size >= 0; line, size = inputQueue.Pop() {
		line = strings.TrimSpace(line)
		if observer != nil {
			observer(line)
		}
		if line == "" {
			monkeys = append(monkeys, m)
			continue
		}

		if strings.HasPrefix(line, "Monkey") {
			m = &monkey{
				worries: utils.NewQueue[int64](),
			}
			continue
		}

		if strings.HasPrefix(line, "Starting items: ") {
			worries := strings.Split(strings.TrimPrefix(line, "Starting items: "), ", ")
			for _, worry := range worries {
				w, err := strconv.ParseInt(worry, 10, 64)
				if err != nil {
					return nil, err
				}

				m.worries.Push(w)
			}
			continue
		}

		if strings.HasPrefix(line, "Operation: new = old ") {
			operaTorAnd := strings.TrimPrefix(line, "Operation: new = old ")
			operator := operaTorAnd[0]
			operandStr := operaTorAnd[2:]
			operand := self
			if operandStr != "old" {
				var err error
				if operand, err = strconv.ParseInt(operandStr, 10, 64); err != nil {
					return nil, err
				}
			}

			m.operation = func(item int64) int64 {
				o := operand
				if o == self {
					o = item
				}

				switch operator {
				case '*':
					return item * o
				case '+':
					return item + o
				case '-':
					return item - o
				case '/':
					return item / o
				default:
					return self
				}
			}

			continue
		}

		if strings.HasPrefix(line, "Test: divisible by ") {
			divStr := strings.TrimPrefix(line, "Test: divisible by ")
			divisor, err := strconv.ParseInt(divStr, 10, 64)
			if err != nil {
				return nil, err
			}

			m.test = func(i int64) bool {
				return i%divisor == 0
			}

			continue
		}

		if strings.HasPrefix(line, "If ") {
			words := strings.Split(line, " ")
			target, err := strconv.Atoi(words[len(words)-1])
			if err != nil {
				return nil, err
			}

			condition := strings.TrimSuffix(words[1], ":")
			if condition == "true" {
				m.passTarget = target
			} else {
				m.failTarget = target
			}

			continue
		}
	}

	monkeys = append(monkeys, m)
	return troop(monkeys), nil
}

func (t troop) monkeyAround(destressor monkeyOp) {
	for _, m := range t {
		for worry, i := m.worries.Pop(); i >= 0; worry, i = m.worries.Pop() {
			m.inspections++
			worry = m.operation(worry)
			worry = destressor(worry)

			if m.test(worry) {
				t[m.passTarget].worries.Push(worry)
			} else {
				t[m.failTarget].worries.Push(worry)
			}
		}
	}
}

func (t troop) monkeyBusiness(log *log.Logger) int64 {
	sort.Sort(sort.Reverse(t))
	t.log(log)
	return t[0].inspections * t[1].inspections
}

func (t troop) Len() int {
	return len(t)
}

func (t troop) Less(i, j int) bool {
	return t[i].inspections < t[j].inspections
}

func (t troop) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t troop) log(out *log.Logger) {
	for i, m := range t {
		out.Printf(`Monkey {{ colorize %[1]q "%[2]d" }} has inspected %[4]d items and has items with worries {{ colorize %[1]q "[%[3]s]" }}`, monkeyColors[i%len(monkeyColors)], i, m.worries.Join(", "), m.inspections)
	}
}
