package cmd

import (
	"os"
	"path/filepath"
)

type config struct {
	editorCommand string
	dir           string
}

var cfg = newConfig()

func newConfig() *config {
	homedir := must(os.UserHomeDir())
	dirpath := must(filepath.Abs(filepath.Join(homedir, "nots")))

	return &config{
		editorCommand: "$EDITOR",
		dir:           dirpath,
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
