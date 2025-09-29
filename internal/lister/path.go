package lister

import (
	"strings"

	"github.com/fredrikkvalvik/nots/internal/util"
)

type Path []string

func NewPath(s string) Path {
	if len(s) == 0 {
		return Path{}
	}

	return strings.Split(s, "/")
}

func (p Path) String() string {
	return strings.Join(p, "/")
}

// return true if the path is equal to parts of p
func (p Path) ContainsSubset(path Path) bool {
	if len(p) < len(path) {
		return false
	}
	return comparePaths(path, p[:len(path)])
}

// returns true if path is equal to p, up to n
func (p Path) IsEqualUpTo(path Path, n int) bool {
	if len(path) < n || len(p) < n {
		return false
	}

	return comparePaths(p[:n], path[:n])
}

func (p Path) IsEqual(path Path) bool {
	return comparePaths(p, path)
}

// return a new path with the last element removed
func (p Path) Pop() Path {
	return p[:len(p)-1]
}

// appends name to Path. Only legal when path is not a filePath (last element is a file name). Will panic if IsFilePath == true
func (p Path) Append(name string) Path {
	if p.IsFilePath() {
		panic("trying to append name to filepath is not valid")
	}
	return append(p, name)
}

// return true if the last element in Path is a valid filename (ending in .md)
func (p Path) IsFilePath() bool {
	return util.IsFileName(p[len(p)-1])
}

// returns a new path where the file name is removed. if no file is there, do nothing
func (p Path) PopFile() Path {
	if p.IsFilePath() {
		return p.Pop()
	}
	return p
}

// return true if p1 and p2 are equal
func comparePaths(p1, p2 Path) bool {
	if len(p1) != len(p2) {
		return false
	}

	for idx, part := range p1 {
		if part != p2[idx] {
			// if any part is not equal to the same part of p, return false
			return false
		}
	}
	return true
}
