package exercise

import (
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

type day1 struct{}

func init() {
	Register("day1", day1{})
}

func (d day1) Part1(stdin io.Reader, stdout io.Writer, stderr io.Writer) (any, error) {
	calorieCounts, err := getCalorieCounts(stdin)
	if err != nil {
		return "", err
	}

	var maxCalCount uint64
	var richestElf int
	for i, count := range calorieCounts {
		if count > maxCalCount {
			log.Printf("Elf #%d has %d calories in their bag, which is the most we've seen so far ...", i, count)
			maxCalCount = count
			richestElf = i
		}
	}

	log.Printf("Elf %d has %d calories in their bag!", richestElf, maxCalCount)
	return maxCalCount, nil
}
func (d day1) Part2(stdin io.Reader, stdout io.Writer, stderr io.Writer) (any, error) {
	calorieCounts, err := getCalorieCounts(stdin)
	if err != nil {
		return "", err
	}

	// because we only care about totals, just sort the slice
	sort.Slice(calorieCounts, func(i, j int) bool {
		// sort in reverse
		return calorieCounts[i] > calorieCounts[j]
	})

	log.Printf("The top 3 elves' bags have %d, %d, and %d calories", calorieCounts[0], calorieCounts[1], calorieCounts[2])
	return calorieCounts[0] + calorieCounts[1] + calorieCounts[2], nil
}

func getCalorieCounts(input io.Reader) ([]uint64, error) {
	lines, err := utils.ReaderToLines(input)
	if err != nil {
		return nil, err
	}

	calorieCounts := []uint64{}
	var calTotal uint64
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			calorieCounts = append(calorieCounts, calTotal)
			calTotal = 0
			continue
		}

		lineVal, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}

		calTotal += lineVal
	}

	if calTotal > 0 {
		calorieCounts = append(calorieCounts, calTotal)
	}

	return calorieCounts, nil
}
