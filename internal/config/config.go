package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ pflag.Value = (*NotsOpenMode)(nil)

type NotsOpenMode string

// Set implements pflag.Value.
func (n *NotsOpenMode) Set(s string) error {
	*n = NotsOpenMode(s)
	return nil
}

// String implements pflag.Value.
func (n NotsOpenMode) String() string {
	return string(n)
}

// Type implements pflag.Value.
func (n NotsOpenMode) Type() string {
	return "NotsOpenMode"
}

const (
	OpenPrevious NotsOpenMode = "previous"
	OpenSeries   NotsOpenMode = "series"
)

type NoteSeries struct {
	// required. the name of the series.
	SeriesName string `toml:"name"`

	// optional. if set, controls what template to use for new notes in the series.
	TemplateName string `toml:"template"`

	// required. takes a single template expression that we evaluate with the template package,
	// and check to see of a file with the evaluated name exists in the series directory
	SeriesFilenameExpression string `toml:"filename-expression"`

	// optional. where to put the note series. defaults to the series-name.
	DirName string `toml:"directory"`
}

type Config struct {
	EditorCommand string `toml:"editor"`
	Pager         string `toml:"viewer"`
	RootDir       string `toml:"notes-dir"`

	// could be nil, if so, default template defined in code
	SelectedTemplate string `toml:"default-template"`

	// default open mode set how the default command resolves opening notes.
	// currently supports:
	//
	// "series" | "previous"
	// defaults to previous
	DefaultOpenMode NotsOpenMode `toml:"default-open-mode"`

	// must be set when using default-open-mode set to "series" to control what note-series to use
	OpenModeSeries string `toml:"open-mode-series"`

	NoteSeries []NoteSeries `toml:"note-series"`
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

func TemplateDir() string {
	return resolveTemplatePath()
}

// currenlty only checks ~/.config/nots/nots.toml
func resolveConfigPath() string {
	homedir := must(os.UserHomeDir())
	configPath := fmt.Sprintf("%s/.config/nots/nots.toml", homedir)
	return configPath
}

// creates the template dir if it does not aleady exist
func resolveTemplatePath() string {
	cfgPath := resolveConfigPath()
	rootPath := filepath.Dir(cfgPath)
	templatePath := filepath.Join(rootPath, "templates")
	_ = os.MkdirAll(templatePath, 0646)
	return templatePath
}

func newDefaultConfig() Config {
	homedir := must(os.UserHomeDir())
	dirpath := must(filepath.Abs(filepath.Join(homedir, "nots")))

	return Config{
		EditorCommand: "$EDITOR",
		Pager:         "$PAGER",
		RootDir:       dirpath,

		SelectedTemplate: "",
		DefaultOpenMode:  "previous",
		OpenModeSeries:   "",
		NoteSeries:       []NoteSeries{},
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
func ptr[T any](v T) *T {
	return &v
}
