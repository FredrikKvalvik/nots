package lister

import "strings"

type Path []string

func newPath(s string) Path {
	return strings.Split(s, "/")
}

func (p Path) String() string {
	return strings.Join(p, "/")
}
