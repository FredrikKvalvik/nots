package util

import (
	"errors"
	"os"
	"strings"
)

// if data is piped to stdin, return true
func HasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// check s to see if it is a valid file name
func IsFileName(s string) bool {
	if strings.Contains(s, "/") {
		return false
	}

	if !strings.HasSuffix(s, ".md") {
		return false
	}

	return true
}

func IsFilePath(s string) bool {
	return !strings.HasSuffix(s, ".md")
}

func FileExists(s string) (bool, error) {
	_, err := os.Stat(s)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
