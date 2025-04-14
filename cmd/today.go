/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// todayCmd represents the today command
var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "A brief description of your command",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		openTodaysNote()
	},
}

func openTodaysNote() {
	path := todayFilePath()

	if !checkFileExists(path) || checkFileEmpty(path) {
		f := must(os.Create(path))
		tpl := templateNote()

		cobra.CheckErr(tpl.Execute(f, &TplArgs{Date: todayDate()}))
		cobra.CheckErr(f.Close())

	}

	editorName := os.ExpandEnv(cfg.editorCommand)
	command := must(exec.LookPath(editorName))

	env := os.Environ()

	if err := syscall.Exec(command, []string{command, path}, env); err != nil {
		cobra.CheckErr(err)
	}
}

func printTodaysNote() {
	path := todayFilePath()
	b, err := os.ReadFile(path)
	cobra.CheckErr(err)

	fmt.Fprintln(os.Stdout, string(b))
}

func todayFilePath() string {
	return fmt.Sprintf("%s/%s", cfg.dir, todayFileName())
}

func todayDate() string {
	return time.Now().Format(time.DateOnly)
}

func todayFileName() string {
	date := todayDate()

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

func init() {
	rootCmd.AddCommand(todayCmd)
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
	// if err := tpl.Execute(&buf, struct {
	// 	Date string
	// }{date}); err != nil {
	// 	panic(err)
	// }

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
