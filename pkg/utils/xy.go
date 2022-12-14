package utils

import "math"

func FromXY(x, y int) int64 {
	return int64((x << 32) + y)
}

func ToXY(i int64) (int, int) {
	return int(i >> 32), int(i & math.MaxInt32)
}
