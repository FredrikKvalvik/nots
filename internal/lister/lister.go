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

func listAllRecursive(path string, includeDirs bool) ([]Path, error) {
	paths := []Path{}
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if p == path {
			return nil
		}

		if d.IsDir() && !includeDirs {
			return nil
		}

		p = strings.TrimPrefix(p, path+"/")

		lpath := NewPath(p)
		// if the file is hidden, we ignore it
		if strings.HasPrefix(lpath[len(lpath)-1], ".") {
			slog.Debug("ignore hidden file", "path", p)
			return nil
		}
		paths = append(paths, lpath)
		return nil
	})

	return paths, err
}

func ListPaths(path string, includeDirs bool) ([]Path, error) {
	entries, err := listDir(path)
	if err != nil {
		return nil, err
	}

	names := []Path{}
	for _, ent := range entries {
		if ent.IsDir() && !includeDirs {
			continue
		}

		lpath := NewPath(ent.Name())
		if strings.HasPrefix(lpath[len(lpath)-1], ".") {
			slog.Debug("ignore hidden file", "path", lpath.String())
			continue
		}

		names = append(names, lpath)
	}

	slog.Info("list-dir", "list", names)
	return names, nil
}

// returns a list of paths where each item is a list of
// names to path to files in multiple directories
func ListPathsRecursive(path string, includeDirs bool) ([]Path, error) {
	paths, err := listAllRecursive(path, includeDirs)
	if err != nil {
		return nil, err
	}

	slog.Info("list-reccursive", "list", paths)
	return paths, nil
}
