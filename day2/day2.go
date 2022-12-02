package day2

import (
	"fmt"
	"log"
	"os"

	"github.com/jghiloni/aoc2022/utils"
)

type play int

const (
	rock play = iota + 1
	paper
	scissors
)

type result byte

const (
	lose result = 'X'
	draw result = 'Y'
	win  result = 'Z'
)

func (r result) String() string {
	switch r {
	case win:
		return "win"
	case lose:
		return "lose"
	case draw:
		return "draw"
	default:
		panic(fmt.Sprintf("invalid strategy %s", string(byte(r))))
	}
}

func (p play) CompareTo(q play) int {
	switch p {
	case rock:
		switch q {
		case rock:
			return 3
		case paper:
			return 0
		case scissors:
			return 6
		}
	case paper:
		switch q {
		case rock:
			return 6
		case paper:
			return 3
		case scissors:
			return 0
		}
	case scissors:
		switch q {
		case rock:
			return 0
		case paper:
			return 6
		case scissors:
			return 3
		}
	}

	return -1
}

func (p play) String() string {
	switch p {
	case rock:
		return "rock"
	case paper:
		return "paper"
	case scissors:
		return "scissors"
	}

	panic("invalid play")
}

func letterToPlay(letter byte) play {
	switch letter {
	case 'A', 'X':
		return rock
	case 'B', 'Y':
		return paper
	case 'C', 'Z':
		return scissors
	}

	return -1
}

func runPlay(them, us play) int {
	log.Printf("Play: Them (%s) vs Us (%s)", them, us)
	return int(us) + us.CompareTo(them)
}

func satisfyStrategy(theirPlay play, strategy result) play {
	switch strategy {
	case win:
		switch theirPlay {
		case rock:
			return paper
		case paper:
			return scissors
		case scissors:
			return rock
		default:
			panic(fmt.Sprintf("invalid opponent play %s", theirPlay))
		}
	case lose:
		switch theirPlay {
		case rock:
			return scissors
		case paper:
			return rock
		case scissors:
			return paper
		default:
			panic(fmt.Sprintf("invalid opponent play %s", theirPlay))
		}
	case draw:
		return theirPlay
	default:
		panic(fmt.Sprintf("invalid strategy %s", string(strategy)))
	}
}

func Part1() (string, error) {
	file, err := os.Open("day2/input.txt")
	if err != nil {
		return "", err
	}

	lines, err := utils.ReaderToLines(file)
	if err != nil {
		return "", err
	}

	totalScore := 0
	for _, line := range lines {
		var (
			them play
			us   play
		)

		bytes := []byte(line)
		them, us = letterToPlay(bytes[0]), letterToPlay(bytes[2])

		totalScore += runPlay(them, us)
	}

	return fmt.Sprintf("%d", totalScore), nil
}

func Part2() (string, error) {
	file, err := os.Open("day2/input.txt")
	if err != nil {
		return "", err
	}

	lines, err := utils.ReaderToLines(file)
	if err != nil {
		return "", err
	}

	totalScore := 0
	for _, line := range lines {
		var (
			them play
			us   play
		)

		bytes := []byte(line)
		them = letterToPlay(bytes[0])
		strategy := result(bytes[2])
		us = satisfyStrategy(them, strategy)

		log.Printf("their play: %s, our strategy: %s, our play: %s", them, strategy, us)

		totalScore += runPlay(them, us)
	}

	return fmt.Sprintf("%d", totalScore), nil
}
