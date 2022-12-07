package communicator

import (
	"errors"
	"path"
	"sort"
	"strings"
)

type CommunicatorDirEntry interface {
	Name() string
	IsDirectory() bool
	Entries() []CommunicatorDirEntry
	Size() int
	Add(...CommunicatorDirEntry) bool
	Walk(WalkFunc) error
	Parent() CommunicatorDirEntry
	AbsPath() string
	setParent(p CommunicatorDirEntry) error
}

type communicatorFile struct {
	name   string
	size   int
	parent CommunicatorDirEntry
}

type communicatorDir struct {
	name       string
	entries    map[string]CommunicatorDirEntry
	cachedSize int
	parent     CommunicatorDirEntry
}

type WalkFunc func(CommunicatorDirEntry) error

var ErrSkipDir = errors.New("skip")
var ErrCantWalkFile = errors.New("files are unwalkable")
var ErrFileCantHaveChildren = errors.New("files cannot have child entries")

func NewDirectory(name string, entries ...CommunicatorDirEntry) CommunicatorDirEntry {
	d := &communicatorDir{
		name:       name,
		entries:    make(map[string]CommunicatorDirEntry),
		cachedSize: -1,
	}

	d.Add(entries...)
	return d
}

func (d *communicatorDir) Name() string {
	return d.name
}

func (d *communicatorDir) IsDirectory() bool {
	return true
}

func (d *communicatorDir) Entries() []CommunicatorDirEntry {
	entries := []CommunicatorDirEntry{}
	for _, entry := range d.entries {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	return entries
}

func (d *communicatorDir) Size() int {
	if d.cachedSize < 0 {
		// recursively calculate the size of all files in this directory. DO NOT include directories in the count because they will be accounted for
		// in Walk's recursive calling
		total := 0
		d.Walk(func(entry CommunicatorDirEntry) error {
			if !entry.IsDirectory() {
				total += entry.Size()
			}

			return nil
		})

		d.cachedSize = total
	}

	return d.cachedSize
}

func (d *communicatorDir) Add(items ...CommunicatorDirEntry) bool {
	for _, entry := range items {
		if _, ok := d.entries[entry.Name()]; ok {
			return false
		}
		d.cachedSize = -1
		entry.setParent(d)
		d.entries[entry.Name()] = entry
	}

	return true
}

func (d *communicatorDir) Walk(visitor WalkFunc) error {
	if err := visitor(d); err != nil {
		if errors.Is(err, ErrSkipDir) {
			return nil
		}

		return err
	}

	for _, entry := range d.Entries() {
		if entry.IsDirectory() {
			if err := entry.Walk(visitor); err != nil && !errors.Is(err, ErrSkipDir) {
				return err
			}

			continue
		}

		if err := visitor(entry); err != nil {
			return err
		}
	}

	return nil
}

func (d *communicatorDir) Parent() CommunicatorDirEntry {
	return d.parent
}

func (d *communicatorDir) AbsPath() string {
	if d.parent == nil {
		return d.name
	}

	return path.Join(d.parent.AbsPath(), d.name)
}

func (d *communicatorDir) setParent(p CommunicatorDirEntry) error {
	if p.IsDirectory() {
		d.parent = p
		return nil
	}

	return ErrFileCantHaveChildren
}

func NewFile(name string, size int) CommunicatorDirEntry {
	return &communicatorFile{
		name: name,
		size: size,
	}
}

func (c *communicatorFile) Name() string {
	return c.name
}

func (c *communicatorFile) IsDirectory() bool {
	return false
}

func (c *communicatorFile) Entries() []CommunicatorDirEntry {
	return nil
}

func (c *communicatorFile) Size() int {
	return c.size
}

func (c *communicatorFile) Add(e ...CommunicatorDirEntry) bool {
	return false
}

func (c *communicatorFile) Walk(_ WalkFunc) error {
	return ErrCantWalkFile
}

func (c *communicatorFile) Parent() CommunicatorDirEntry {
	return c.parent
}

func (c *communicatorFile) AbsPath() string {
	if c.parent == nil {
		return c.name
	}

	return path.Join(c.parent.AbsPath(), c.name)
}

func (c *communicatorFile) setParent(p CommunicatorDirEntry) error {
	if p.IsDirectory() {
		c.parent = p
		return nil
	}

	return ErrFileCantHaveChildren
}
