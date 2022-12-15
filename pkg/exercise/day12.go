package exercise

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

var hash = utils.FromXY

var neighborDeltas = [4][2]int{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

var arrows = []rune{'\uffe9', '\uffea', '\uffeb', '\uffec'}

type day12 struct{}

func init() {
	Register("day12", day12{})
}

func (d day12) Part1(input io.Reader, output *log.Logger) (any, error) {
	return d.shortestPath(input, output, 'S')
}

func (d day12) Part2(input io.Reader, output *log.Logger) (any, error) {
	return d.shortestPath(input, output, 'a')
}

func (d day12) shortestPath(in io.Reader, out *log.Logger, start byte) (int, error) {
	mapBytes, err := io.ReadAll(in)
	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(mapBytes), "\n")
	sources := []uint64{}
	target := uint64(0)

	vertices := map[uint64]int{}

	for y, line := range lines {
		for x, r := range []byte(line) {
			h := hash(x, y)
			if r == start {
				// start on an edge
				if y == 0 || x == 0 || y == len(lines)-1 || x == len(line)-1 {
					sources = append(sources, h)
					vertices[h] = 0
					continue
				}
			}

			if r == 'E' {
				target = h
				vertices[h] = 25
				continue
			}

			vertices[h] = int(r - 'a')
		}
	}

	var shortest *utils.Stack[uint64]
	shortestLen := math.MaxInt
	for _, source := range sources {
		x, y := utils.ToXY(source)
		out.Printf(`Starting path at {{ colorize "bold;bright-magenta" "(%d,%d)" }}`, x, y)
		solver := newPathSolver(vertices, source, target)
		solver.breadthFirst()

		stack, err := solver.walkBack(shortestLen)
		if err != nil {
			out.Printf(`{{ colorize "bold;red" %q }}`, err.Error())
			continue
		}

		out.Printf(`Found path with length {{ colorize "bold;bright-cyan" "%d" }}`, stack.Size())
		if stack != nil && (shortest == nil || shortest.Size() > stack.Size()) {
			shortest = stack
			shortestLen = stack.Size()
			out.Printf("This is now the shortest path")
		}
	}

	if shortest == nil {
		return -1, errors.New("no path found")
	}

	len := shortest.Size()
	for shortest.Size() > 0 {
		v, _ := shortest.Pop()
		x, y := utils.ToXY(v)
		next := shortest.Peek()
		nx, ny := utils.ToXY(next)
		line := []rune(lines[y])

		r := line[x]
		if nx < x {
			r = arrows[0]
		} else if nx > x {
			r = arrows[2]
		} else if ny < y {
			r = arrows[1]
		} else if ny > y {
			r = arrows[3]
		}

		line[x] = rune(r)
		lines[y] = string(line)
	}

	out.Println(`{{colorize "bold;bright-cyan" "FINAL MAP"}}`)
	printFrame(lines, out.Writer())

	return len, nil
}

type pathSolver struct {
	vertices   map[uint64]int
	vertexList []uint64
	source     uint64
	target     uint64
	dists      map[uint64]int
	prev       map[uint64]uint64
}

func newPathSolver(vertices map[uint64]int, source uint64, target uint64) *pathSolver {
	p := &pathSolver{
		vertices:   map[uint64]int{},
		vertexList: make([]uint64, 0, len(vertices)),
		source:     source,
		target:     target,
		dists:      map[uint64]int{},
		prev:       map[uint64]uint64{},
	}

	// copy the vertices map for safety
	for h, k := range vertices {
		p.vertices[h] = k
		p.vertexList = append(p.vertexList, h)
		p.dists[h] = math.MaxInt
	}

	p.dists[source] = 0
	return p
}

func (p pathSolver) Len() int {
	return len(p.vertexList)
}

func (p pathSolver) Less(i, j int) bool {
	hi, hj := p.vertexList[i], p.vertexList[j]
	return p.dists[hi] < p.dists[hj]
}

func (p *pathSolver) Swap(i, j int) {
	p.vertexList[i], p.vertexList[j] = p.vertexList[j], p.vertexList[i]
}

func (p *pathSolver) Pop() (uint64, int) {
	sort.Sort(p)

	h := p.vertexList[0]
	k := p.vertices[h]

	p.vertexList = p.vertexList[1:]

	return h, k
}

func (p *pathSolver) breadthFirst() {
	for p.Len() > 0 {
		u, uh := p.Pop()

		n := p.neighbors(u)
		if len(n) == 0 {
			return
		}

		for _, v := range n {
			vh := p.vertices[v]
			alt := p.dists[u] + (vh - uh) + 1
			if alt < p.dists[v] {
				p.dists[v] = alt
				p.prev[v] = u
				if v == p.target {
					return
				}
			}
		}
	}
}

func (p *pathSolver) walkBack(maxlen int) (*utils.Stack[uint64], error) {
	s := utils.NewStack[uint64]()

	u := p.target
	var ok bool
	if u, ok = p.prev[u]; ok || u == p.source {
		for ok {
			s.Push(u)
			if s.Size() > maxlen {
				return nil, fmt.Errorf("this path is longer than the shortest path, don't bother continuing")
			}
			if v, ok := p.prev[u]; ok && u == v {
				x, y := utils.ToXY(u)
				return nil, fmt.Errorf("cycle found at (%d,%d)", x, y)
			}
			u, ok = p.prev[u]
		}
	}

	return s, nil
}

func (p *pathSolver) neighbors(u uint64) []uint64 {
	n := []uint64{}
	for _, delta := range neighborDeltas {
		if v, ok := p.isValidNeighbor(u, delta[0], delta[1]); ok {
			n = append(n, v)
		}
	}

	return n
}

func (p *pathSolver) isValidNeighbor(u uint64, dx, dy int) (uint64, bool) {
	ux, uy := utils.ToXY(u)
	v := hash(ux+dx, uy+dy)

	uh := p.vertices[u]
	if vh, ok := p.vertices[v]; ok && vh <= (uh+1) {
		return v, true
	}

	return 0, false
}

func printFrame(lines []string, out io.Writer) {
	for _, line := range lines {
		for _, r := range line {
			if r < arrows[0] {
				fmt.Fprint(out, string(r))
				continue
			}

			fmt.Fprintf(out, `{{ colorize "bg-cyan;black" %q }}`, string(r))
		}
		out.Write([]byte("\n"))
	}
}
