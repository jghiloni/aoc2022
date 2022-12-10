package utils

import "fmt"

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d✕%d", c.X, c.Y)
}

type RopeBridge struct {
	knots  []Coordinate
	visits []map[Coordinate]int
}

func NewRopeBridge(numKnots int) *RopeBridge {
	b := &RopeBridge{
		knots:  make([]Coordinate, numKnots),
		visits: make([]map[Coordinate]int, numKnots),
	}

	b.visit()
	return b
}

func (r *RopeBridge) MoveHead(direction string, distance int) {
	xIncrementor := 0
	yIncrementor := 0

	switch direction {
	case "L":
		xIncrementor = -1
	case "R":
		xIncrementor = 1
	case "U":
		yIncrementor = -1
	case "D":
		yIncrementor = 1
	}

	for step := 0; step < distance; step++ {
		r.knots[0].X += xIncrementor
		r.knots[0].Y += yIncrementor
		r.reposition()
		r.visit()
	}
}

func (r *RopeBridge) String() string {
	h := fmt.Sprintf("H(%s)", r.knots[0])
	t := fmt.Sprintf("T(%s)", r.knots[len(r.knots)-1])

	return fmt.Sprintf(`{{ colorize "bold;bright-magenta" %q }},{{ colorize "bold;bright-green" %q }}`, h, t)
}

func (r *RopeBridge) UniqueVisits() []int {
	visits := make([]int, len(r.knots))
	for i, knotMap := range r.visits {
		visits[i] = len(knotMap)
	}

	return visits
}

func (r *RopeBridge) isKnotInPosition(i int) bool {
	if i == 0 {
		return true
	}
	xDistance := r.knots[i].X - r.knots[i-1].X
	yDistance := r.knots[i].Y - r.knots[i-1].Y

	return Abs(xDistance) <= 1 && Abs(yDistance) <= 1
}

func (r *RopeBridge) reposition() {
	for i := 1; i < len(r.knots); i++ {
		if r.isKnotInPosition(i) {
			continue
		}

		r.repositionKnot(i)
	}
}

func (r *RopeBridge) repositionKnot(i int) {
	head := r.knots[i-1]
	tail := r.knots[i]

	vecDeltaX := head.X - tail.X
	magDeltaX := Abs(vecDeltaX)
	vecDeltaY := head.Y - tail.Y
	magDeltaY := Abs(vecDeltaY)
	moveX := 0
	moveY := 0
	if magDeltaX > 1 {
		moveX = vecDeltaX / magDeltaX
		if vecDeltaY != 0 {
			moveY = vecDeltaY / magDeltaY
		}
	} else if magDeltaY != 0 {
		moveY = vecDeltaY / magDeltaY
		if vecDeltaX != 0 {
			moveX = vecDeltaX / magDeltaX
		}
	}

	tail.X += moveX
	tail.Y += moveY

	r.knots[i] = tail
}

func (r *RopeBridge) visit() {
	for i, k := range r.knots {
		if r.visits[i] == nil {
			r.visits[i] = map[Coordinate]int{k: 0}
		}

		r.visits[i][k]++
	}
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}
