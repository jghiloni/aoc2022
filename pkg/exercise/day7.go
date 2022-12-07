package exercise

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/jghiloni/aoc2022/pkg/utils/communicator"
)

type day7 struct{}

func init() {
	Register("day7", day7{})
}

func (d day7) Part1(input io.Reader, output *log.Logger) (any, error) {
	fs, err := communicator.BuildFromReplay(input)
	if err != nil {
		output.Printf(`{{ colorize "bold:red" "An error occurred reading the command replay: %v"}}`, err)
		return nil, err
	}

	maxSize := 100000
	total := 0
	fs.WalkFS(func(entry communicator.CommunicatorDirEntry) error {
		if entry.IsDirectory() {
			name := fmt.Sprintf(`{{ colorize "green" %q }}`, entry.AbsPath())
			output.Printf("Calculating size of %s", name)
			size := entry.Size()
			output.Printf(`%s contains {{ colorize "bold;yellow" "%d" }} bytes`, name, size)

			if size <= maxSize {
				total += size
				output.Printf(`%s is below the max size of {{ colorize "bold;bright-blue" "%d" }} bytes. The current total bytes of directories below the max size is {{ colorize "bold;yellow" "%d" }}`, name, maxSize, total)
			}
		}

		return nil
	})

	output.Printf(`Total size of all directories that contain less than {{ colorize "bold;cyan" "%d" }} bytes is {{ colorize "bold;yellow" "%d" }}`, maxSize, total)
	return total, err
}

func (d day7) Part2(input io.Reader, output *log.Logger) (any, error) {
	fs, err := communicator.BuildFromReplay(input)
	if err != nil {
		output.Printf(`{{ colorize "bold:red" "An error occurred reading the command replay: %v"}}`, err)
		return nil, err
	}

	totalSize := 70000000
	minRequiredFreeSpace := 30000000
	fsSize := fs.Size()

	freeSpace := totalSize - fsSize
	requiredDeletionThreshold := minRequiredFreeSpace - freeSpace
	output.Printf(`the communicator requires {{colorize "bold;bright-cyan" "%d"}} free bytes but only has {{ colorize "bold;red" "%d"}} free. We must free up {{colorize "bold;yellow" "%d"}} bytes`, minRequiredFreeSpace, freeSpace, requiredDeletionThreshold)

	deletionCandidates := []communicator.CommunicatorDirEntry{}
	fs.WalkFS(func(entry communicator.CommunicatorDirEntry) error {
		if entry.IsDirectory() {
			size := entry.Size()
			if size == fsSize {
				return nil // skip the root dir
			}
			message := fmt.Sprintf(`Deleting {{colorize "bold;green" %q}} would free up`, entry.AbsPath())
			if size >= requiredDeletionThreshold {
				deletionCandidates = append(deletionCandidates, entry)
				message = fmt.Sprintf(`%s {{colorize "bold;yellow" "%d"}} bytes more than necessary`, message, size-requiredDeletionThreshold)
			} else {
				message = fmt.Sprintf(`%s {{colorize "bold;bright-red" "%d"}} too few bytes`, message, requiredDeletionThreshold-size)
			}

			output.Println(message)
		}

		return nil
	})

	if len(deletionCandidates) == 0 {
		err := errors.New("no directories were big enough to free up enough space")
		output.Printf(`{{colorize "bold;red" %q}}`, err.Error())
		return nil, err
	}

	sort.Slice(deletionCandidates, func(i, j int) bool {
		return deletionCandidates[i].Size() < deletionCandidates[j].Size()
	})

	output.Println(`{{colorize "bold;italic;bright-white" " ========================== RESULTS ========================== " }}`)
	for _, dc := range deletionCandidates {
		output.Printf(`Deleting {{colorize "bold;green" %q }} would leave us with {{colorize "bold;bright-cyan" "%d"}} bytes free`, dc.AbsPath(), totalSize-fsSize-dc.Size())
	}

	output.Printf(`The smallest directory that will free up enough space is {{colorize "green" %q}} with {{colorize "bold;yellow" "%d"}} bytes`, deletionCandidates[0].AbsPath(), deletionCandidates[0].Size())
	return deletionCandidates[0].Size(), nil
}
