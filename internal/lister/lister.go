package lister

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func listDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func listAllRecursive(path string, includeDirs bool) ([][]string, error) {
	paths := [][]string{}
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if p == path {
			return nil
		}

		if d.IsDir() && !includeDirs {
			return nil
		}

		p = strings.TrimPrefix(p, path+"/")

		paths = append(paths, strings.Split(p, "/"))
		return nil
	})

	return paths, nil
}

func ListFilesInDir(path string) ([]string, error) {
	entries, err := listDir(path)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, ent := range entries {
		if !ent.IsDir() {
			names = append(names, ent.Name())
		}
	}

	return names, nil
}

func ListDir(path string, includeDirs bool) ([][]string, error) {
	entries, err := listDir(path)
	if err != nil {
		return nil, err
	}

	names := [][]string{}
	for _, ent := range entries {
		if ent.IsDir() && !includeDirs {
			continue
		}
		names = append(names, []string{ent.Name()})
	}

	return names, nil
}

// returns a list of paths where each item is a list of
// names to path to files in multiple directories
func ListPathsRecursive(path string, includeDirs bool) ([][]string, error) {
	paths, err := listAllRecursive(path, includeDirs)
	if err != nil {
		return nil, err
	}

	slog.Info("list-reccursive", "list", paths)
	return paths, nil
}
