/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

func openNote(path string) {
	if !checkFileExists(path) || checkFileEmpty(path) {

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		cobra.CheckErr(err)

		f := must(os.Create(path))
		tpl := templateNote()

		cobra.CheckErr(tpl.Execute(f, &TplArgs{Date: todayDate()}))
		cobra.CheckErr(f.Close())
	}

	slog.Debug("opening note", "path", path)
	currentState.PreviousNote = &path
	cobra.CheckErr(currentState.Save())

	editorName := os.ExpandEnv(cfg.EditorCommand)
	command := must(exec.LookPath(editorName))

	env := os.Environ()

	if err := syscall.Exec(command, []string{command, path}, env); err != nil {
		cobra.CheckErr(err)
	}

}

func getNoteContent(path string) string {
	b, err := os.ReadFile(path)
	cobra.CheckErr(err)
	return string(b)
}

func todayFilePath() string {
	return filePath(todayFileName())
}

func filePath(name string) string {
	return fmt.Sprintf("%s/%s", cfg.RootDir, name)

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

type TplArgs struct {
	Date string
}

func templateNote() *template.Template {
	tpl := template.Must(template.New("note").Parse(`# {{ .Date }}

## Notater

- ...
`))

	return tpl
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
