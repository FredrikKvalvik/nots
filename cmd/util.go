package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func openNote(path string) {
	if !checkFileExists(path) || checkFileEmpty(path) {
		slog.Debug("creating new file", "path", path)

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		cobra.CheckErr(err)

		f := must(os.Create(path))
		defer func() {
			_ = f.Close()
		}()

		if cfg.SelectedTemplate != "" {
			slog.Debug("found template", "template_name", cfg.SelectedTemplate)

			err := loadTemplate(cfg.SelectedTemplate, f)
			cobra.CheckErr(err)
		}
	}

	spawnEditor(path)
}

// helper for spawning editor process
func spawnEditor(path string) {
	slog.Debug("opening note", "path", path)
	currentState.PreviousNote = &path
	cobra.CheckErr(currentState.Save())

	editorName := os.ExpandEnv(cfg.EditorCommand)
	command := must(exec.LookPath(editorName))

	env := os.Environ()

	err := os.Chdir(cfg.RootDir)
	cobra.CheckErr(err)

	if err := syscall.Exec(command, []string{command, path}, env); err != nil {
		cobra.CheckErr(err)
	}
}

func getFileContent(path string) string {
	b, err := os.ReadFile(path)
	cobra.CheckErr(err)
	return string(b)
}

func todayFilePath() string {
	return absolutePath(todayFileName())
}

func absolutePath(name string) string {
	return filepath.Join(cfg.RootDir, name)
}

func todayDate() string {
	return time.Now().Format(time.DateOnly)
}

func todayFileName() string {
	date := todayDate()

	// TODO: implement today-file template parsing

	return fmt.Sprintf("%s.md", date)
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func checkFileEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return true
	}

	if info.Size() == 0 {
		return true
	}
	return false
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
