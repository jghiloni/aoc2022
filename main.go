package main

import (
	"fmt"
	"log"

	"github.com/jghiloni/aoc2022/day2"
)

func main() {
	answer, err := day2.Part1()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)

	answer, err = day2.Part2()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)
}
