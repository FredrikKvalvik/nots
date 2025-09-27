package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type config struct {
	EditorCommand string `toml:"editor"`
	Dir           string `toml:"notes-dir"`
}

var cfg *config = loadConfig()

func loadConfig() *config {
	filePath := resolveConfigPath()

	var conf = newDefaultConfig()
	_, err := toml.DecodeFile(filePath, &conf)
	cobra.CheckErr(err)

	conf.Dir = resolveNotsDirPath(conf.Dir)

	// create directories to make sure they are there
	err = createNotsDirectory(conf.Dir)
	cobra.CheckErr(err)

	return &conf
}

// currenlty only checks ~/.config/nots/nots.toml
func resolveConfigPath() string {
	homedir := must(os.UserHomeDir())
	configPath := fmt.Sprintf("%s/.config/nots/nots.toml", homedir)
	return configPath
}

func newDefaultConfig() config {
	homedir := must(os.UserHomeDir())
	dirpath := must(filepath.Abs(filepath.Join(homedir, "nots")))

	return config{
		EditorCommand: "$EDITOR",
		Dir:           dirpath,
	}
}

func resolveNotsDirPath(path string) string {
	// first we remove any leading/trailing whitespace
	path = strings.TrimSpace(path)

	fmt.Printf("path: %v\n", path)
	// user refers the their home dir
	if strings.HasPrefix(path, "~/") {

		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		path = strings.TrimPrefix(path, "~/")
		path, err = filepath.Abs(filepath.Join(homeDir, path))
		cobra.CheckErr(err)
	}

	// we remove any trailing slashes
	path = strings.TrimSuffix(path, "/")
	fmt.Printf("path: %v\n", path)

	return path
}
func createNotsDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
