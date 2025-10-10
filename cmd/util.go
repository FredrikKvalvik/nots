package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// open note with selected template if it does not exist, else
// we open the existing note
func openNoteWithSelectedTemplate(absolutePath string) {
	openNote(absolutePath, cfg.SelectedTemplate)
}

func openNote(absolutePath string, templateName string) {
	if !checkFileExists(absolutePath) || checkFileEmpty(absolutePath) {
		slog.Debug("creating new file", "path", absolutePath)

		ensureFilepathExists(absolutePath)
		f := must(os.Create(absolutePath))
		defer func() {
			_ = f.Close()
		}()

		if cfg.SelectedTemplate != "" {
			slog.Debug("found template", "template_name", templateName)

			err := loadTemplate(templateName, f)
			cobra.CheckErr(err)
		}
	}

	spawnEditor(absolutePath)
}

// ensures that the filepath exists by creating the parent directories
func ensureFilepathExists(absolutePath string) {
	err := os.MkdirAll(filepath.Dir(absolutePath), os.ModePerm)
	cobra.CheckErr(err)
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
	if !strings.HasPrefix(name, cfg.RootDir) {
		return filepath.Join(cfg.RootDir, name)
	}
	return name
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
