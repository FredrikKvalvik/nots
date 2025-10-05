package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Config struct {
	EditorCommand string `toml:"editor"`
	Pager         string `toml:"viewer"`
	RootDir       string `toml:"notes-dir"`

	DailyNameTemplate string `toml:"daily-name-template"`
	DailyDirName      string `toml:"daily-dir-name"`

	// could be nil, if so, default template defined in code
	SelectedTemplate *string `toml:"default-template"`

	// default open mode set how the default command resolves opening notes.
	// currently supports:
	//
	// "today" | "previous"
	DefaultOpenMode string `toml:"default-open-mode"`
}

func Load() (*Config, error) {
	filePath := resolveConfigPath()

	var conf = newDefaultConfig()
	_, err := toml.DecodeFile(filePath, &conf)
	if err != nil {
		return nil, err
	}

	conf.RootDir = resolveNotsDirPath(conf.RootDir)

	err = validateConfig(&conf)
	if err != nil {
		return nil, err
	}

	// create directories to make sure they are there
	err = createNotsDirectory(conf.RootDir)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func TemplateName(name string) string {
	return filepath.Join(resolveTemplatePath(), name)
}

// currenlty only checks ~/.config/nots/nots.toml
func resolveConfigPath() string {
	homedir := must(os.UserHomeDir())
	configPath := fmt.Sprintf("%s/.config/nots/nots.toml", homedir)
	return configPath
}

func resolveTemplatePath() string {
	cfgPath := resolveConfigPath()
	rootPath := filepath.Dir(cfgPath)
	return filepath.Join(rootPath, "templates")
}

func newDefaultConfig() Config {
	homedir := must(os.UserHomeDir())
	dirpath := must(filepath.Abs(filepath.Join(homedir, "nots")))

	return Config{
		EditorCommand: "$EDITOR",
		Pager:         "$PAGER",
		RootDir:       dirpath,

		DailyNameTemplate: "yyyy-mm-dd",
		DailyDirName:      "", // default to root
		SelectedTemplate:  nil,
		DefaultOpenMode:   "today",
	}
}

func resolveNotsDirPath(path string) string {
	// first we remove any leading/trailing whitespace
	path = strings.TrimSpace(path)

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
