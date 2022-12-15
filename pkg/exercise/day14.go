package exercise

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day14 struct{}

const (
	air   = '.'
	sand  = 'o'
	stone = '#'
	path  = '~'
)

type caveSlice struct {
	p                [][]byte
	minX             int
	maxX             int
	maxY             int
	currentPath      []uint64
	hasInfiniteFloor bool
}

func init() {
	Register("day14", day14{})
}

func (d day14) Part1(input io.Reader, output *log.Logger) (any, error) {
	caveMap, err := d.drawMap(input, 0)
	if err != nil {
		return nil, err
	}

	output.Println("starting map")
	fmt.Fprintf(output.Writer(), "%s\n", caveMap)

	total := 0
	for caveMap.dropGrain(500, 0) {
		total++
	}

	output.Println("ending map")
	fmt.Fprintf(output.Writer(), "%s\n", caveMap)

	return total, nil
}

func (d day14) Part2(input io.Reader, output *log.Logger) (any, error) {
	caveMap, err := d.drawMap(input, 2)
	if err != nil {
		return nil, err
	}

	output.Println("starting map")
	fmt.Fprintf(output.Writer(), "%s\n", caveMap)

	total := 1
	for caveMap.dropGrain(500, 0) {
		total++
	}

	output.Println("ending map")
	fmt.Fprintf(output.Writer(), "%s\n", caveMap)

	return total, nil
}

func (d day14) drawMap(in io.Reader, floorHeight int) (*caveSlice, error) {
	lines, err := utils.ReaderToLines(in)
	if err != nil {
		return nil, err
	}

	c := new(caveSlice)
	xStones := map[int]bool{}
	yStones := map[int]bool{}

	for _, line := range lines {
		coords := strings.Split(line, " -> ")
		for _, coord := range coords {
			x, y, err := parseCoordinates(coord)
			if err != nil {
				return nil, err
			}

			xStones[x] = true
			yStones[y] = true
		}
	}

	c.minX = math.MaxInt
	for x := range xStones {
		if x < c.minX {
			c.minX = x
			continue
		}

		if x > c.maxX {
			c.maxX = x
			continue
		}
	}

	for y := range yStones {
		if y > c.maxY {
			c.maxY = y
			continue
		}
	}

	c.hasInfiniteFloor = floorHeight > 0

	c.maxY += floorHeight

	width := c.maxX - c.minX + 1
	height := c.maxY + 1

	c.p = make([][]byte, height)
	for y := 0; y < height; y++ {
		c.p[y] = bytes.Repeat([]byte{air}, width)
	}

	for _, line := range lines {
		coords := strings.Split(line, " -> ")
		var (
			start string
			end   string
		)

		start = coords[0]
		coords = coords[1:]
		for len(coords) > 0 {
			end = coords[0]
			coords = coords[1:]

			fillers, err := getPath(start, end)
			if err != nil {
				return nil, err
			}

			for _, chash := range fillers {
				x, y := utils.ToXY(chash)
				if y == height-1 {
					log.Print()
				}
				c.p[y][x-c.minX] = stone
			}

			start = end
		}
	}

	if floorHeight > 0 {
		c.p[height-1] = bytes.Repeat([]byte{stone}, width)
	}
	return c, nil
}

func (c *caveSlice) dropGrain(startX, startY int) bool {
	if c.hasInfiniteFloor {
		return c.dropGrainWithBottom(startX, startY)
	}

	return c.dropGrainWithoutBottom(startX, startY)
}

func (c *caveSlice) dropGrainWithoutBottom(startX, startY int) bool {
	x, y := startX, startY
	c.currentPath = []uint64{utils.FromXY(startX, startY)}
	for {
		nx, ny := c.nextMove(x, y)
		c.currentPath = append(c.currentPath, utils.FromXY(nx, ny))

		if nx < c.minX || nx > c.maxX {
			return c.finalizePath()
		}

		if ny == c.maxY {
			c.p[ny][nx-c.minX] = sand
			return c.finalizePath()
		}

		if x == nx && y == ny {
			c.p[ny][nx-c.minX] = sand
			return true
		}

		x, y = nx, ny
	}
}

func (c *caveSlice) dropGrainWithBottom(startX, startY int) bool {
	x, y := startX, startY
	for {
		nx, ny := c.nextMove(x, y)

		if nx == startX && ny == startY {
			c.p[startY][startX-c.minX] = sand
			return false
		}

		if nx < c.minX && ny > y {
			c.minX--
			for y := range c.p {
				c.p[y] = append([]byte{air}, c.p[y]...)
			}
			c.p[c.maxY][0] = stone
		}

		if nx > c.maxX && ny > y {
			c.maxX++
			for y := range c.p {
				c.p[y] = append(c.p[y], air)
			}
			c.p[c.maxY][len(c.p[c.maxY])-1] = stone
		}

		if x == nx && y == ny {
			c.p[ny][nx-c.minX] = sand
			return true
		}

		x, y = nx, ny
	}
}

func (c *caveSlice) finalizePath() bool {
	for _, seg := range c.currentPath {
		x, y := utils.ToXY(seg)
		if x >= c.minX {
			c.p[y][x-c.minX] = path
		}
	}

	return false
}

func (c *caveSlice) nextMove(x, y int) (int, int) {
	if c.hasInfiniteFloor && y == c.maxY-1 {
		return x, y
	}

	options := []uint64{utils.FromXY(x, y+1), utils.FromXY(x-1, y+1), utils.FromXY(x+1, y+1)}
	for _, option := range options {
		nx, ny := utils.ToXY(option)
		if nx < c.minX || nx > c.maxX || ny > c.maxY {
			return nx, ny
		}

		if c.p[ny][nx-c.minX] == air {
			return nx, ny
		}
	}

	return x, y
}

// func (c *caveSlice) countSand() int {
// 	total := 0
// 	for y := range c.p {
// 		for x := range c.p[y] {
// 			if c.p[y][x] == sand {
// 				total++
// 			}
// 		}
// 	}
// 	return total
// }

func (c *caveSlice) String() string {
	space := string(bytes.Repeat([]byte{' '}, c.maxX-c.minX-1))
	lines := make([]string, c.maxY+5)
	lines[0] = fmt.Sprintf("    %d%s%d", c.minX/100, space, c.maxX/100)
	lines[1] = fmt.Sprintf("    %d%s%d", (c.minX%100)/10, space, (c.maxX%100)/10)
	lines[2] = fmt.Sprintf("    %d%s%d", c.minX%10, space, c.maxX%10)
	for y := 0; y <= c.maxY; y++ {
		lines[3+y] = fmt.Sprintf("%3d %s", y, c.p[y])
	}

	return strings.Join(lines, "\n")
}

func getPath(start, end string) ([]uint64, error) {
	var (
		sx  int
		sy  int
		ex  int
		ey  int
		err error
	)

	if sx, sy, err = parseCoordinates(start); err != nil {
		return nil, err
	}

	if ex, ey, err = parseCoordinates(end); err != nil {
		return nil, err
	}

	p := []uint64{}

	x, y := sx, sy
	if sx == ex && sy != ey {
		ystep := (ey - sy) / utils.Abs(ey-sy)
		for y != ey+ystep {
			p = append(p, utils.FromXY(x, y))
			y += ystep
		}
	} else if sx != ex && sy == ey {
		xstep := (ex - sx) / utils.Abs(ex-sx)
		for x != ex+xstep {
			p = append(p, utils.FromXY(x, y))
			x += xstep
		}
	}

	return p, nil
}

func parseCoordinates(xy string) (int, int, error) {
	var x, y int
	if parsed, err := fmt.Sscanf(xy, "%d,%d", &x, &y); parsed != 2 || err != nil {
		if err != nil {
			return -1, -1, err
		}

		return -1, -1, errors.New("could not parse coordinate")
	}

	return x, y, nil
}
