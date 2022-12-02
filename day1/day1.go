package day1

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/utils"
)

func getCalorieCounts() ([]uint64, error) {
	file, err := os.Open("day1/input.txt")
	if err != nil {
		return nil, err
	}

	lines, err := utils.ReaderToLines(file)
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

func Part1() (string, error) {
	calorieCounts, err := getCalorieCounts()
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
	return fmt.Sprintf("%d", maxCalCount), nil
}

func Part2() (string, error) {
	calorieCounts, err := getCalorieCounts()
	if err != nil {
		return "", err
	}

	// because we only care about totals, just sort the slice
	sort.Slice(calorieCounts, func(i, j int) bool {
		// sort in reverse
		return calorieCounts[i] > calorieCounts[j]
	})

	log.Printf("The top 3 elves' bags have %d, %d, and %d calories", calorieCounts[0], calorieCounts[1], calorieCounts[2])
	return fmt.Sprintf("%d", calorieCounts[0]+calorieCounts[1]+calorieCounts[2]), nil
}
