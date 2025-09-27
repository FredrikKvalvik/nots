package lister

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func listDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func listAllFiles(path string) ([]string, error) {
	paths := []string{}
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if p == path {
			return nil
		}

		paths = append(paths, p)
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

func ListDir(path string, includeDirs bool) ([]string, error) {
	entries, err := listDir(path)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, ent := range entries {
		if ent.IsDir() && !includeDirs {
			continue
		}
		names = append(names, ent.Name())
	}

	return names, nil
}

// returns a list of paths where each item is a list of
// names to path to files in multiple directories
func ListPathsRecursive(path string) ([][]string, error) {
	paths, err := listAllFiles(path)
	if err != nil {
		return nil, err
	}

	pathsSplit := make([][]string, len(paths))

	for idx, p := range paths {
		paths[idx] = strings.TrimPrefix(p, path+"/")
		pathsSplit[idx] = strings.Split(paths[idx], "/")
	}

	return pathsSplit, nil
}
