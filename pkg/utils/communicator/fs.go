package communicator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jghiloni/aoc2022/pkg/utils"
)

var ErrNotExist = errors.New("ENOTEXIST")
var ErrNotDirectory = errors.New("ENOTDIR")

type CommunicatorFileSystem struct {
	rootDir CommunicatorDirEntry
	wdStack *utils.Stack[CommunicatorDirEntry]
}

func NewFS(root CommunicatorDirEntry) *CommunicatorFileSystem {
	wdStack := utils.NewStack[CommunicatorDirEntry]()
	wdStack.Push(root)

	return &CommunicatorFileSystem{
		rootDir: root,
		wdStack: wdStack,
	}
}

func BuildFromReplay(input io.Reader) (*CommunicatorFileSystem, error) {
	replayBuffer := bufio.NewReader(input)
	fs := NewFS(NewDirectory("/"))

	for {
		lineBytes, _, err := replayBuffer.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		line := strings.TrimSpace(string(lineBytes))
		if line == "" {
			continue
		}

		if line[0] == '$' {
			if err := fs.processCommand(line, replayBuffer); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unexpected replay output %q", line)
		}
	}

	return fs, nil
}

func (fs *CommunicatorFileSystem) Size() int {
	return fs.rootDir.Size()
}

func (fs *CommunicatorFileSystem) WalkFS(visitor WalkFunc) error {
	fs.clearWdStack()
	return fs.rootDir.Walk(visitor)
}

func (fs *CommunicatorFileSystem) clearWdStack() {
	for fs.wdStack.Size() > 1 {
		fs.wdStack.Pop()
	}
}

func (fs *CommunicatorFileSystem) processCommand(line string, replayBuffer *bufio.Reader) error {
	commandLine := line[2:] // remove the '$ '
	commandName := commandLine[0:2]

	switch commandName {
	case "cd":
		dirname := commandLine[3:]
		if dirname == ".." {
			fs.wdStack.Pop()
			return nil
		}

		if dirname == "/" {
			fs.clearWdStack()
			return nil
		}

		cwd := fs.wdStack.Peek()
		for _, entry := range cwd.Entries() {
			if entry.Name() == dirname {
				if !entry.IsDirectory() {
					return ErrNotDirectory
				}

				fs.wdStack.Push(entry)
				return nil
			}
		}

		return ErrNotExist
	case "ls":
		cwd := fs.wdStack.Peek()
		for {
			peeked, err := replayBuffer.Peek(1)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return err
			}

			if string(peeked) == "$" {
				return nil
			}

			listingLine, _, err := replayBuffer.ReadLine()
			if err != nil {
				return err
			}

			lsParts := strings.SplitN(string(listingLine), " ", 2)
			switch lsParts[0] {
			case "dir":
				cwd.Add(NewDirectory(lsParts[1]))
			default:
				size, err := strconv.Atoi(lsParts[0])
				if err != nil {
					return err
				}

				cwd.Add(NewFile(lsParts[1], size))
			}
		}
	}

	return nil
}
