package utils

func FromXY(x, y int) uint64 {
	return uint64((x&0xffffffff)<<32) + uint64(uint32(y))
}

func ToXY(u uint64) (int, int) {
	x := int32(u >> 32)
	y := int32(u & 0xffffffff)
	return int(x), int(y)
}
