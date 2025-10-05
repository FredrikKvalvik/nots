package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fredrikkvalvik/nots/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ViewCmd())
}

func ViewCmd() *cobra.Command {
	var toStdOut bool

	cmd := &cobra.Command{
		Use:               "view [file]",
		Example:           "echo 'filename.md' | nots view" + "\n" + "nots view filename.md",
		ValidArgsFunction: fileCompleter(cfg.RootDir),

		Short:   "view a specified note",
		Aliases: []string{"ls"},

		Run: func(cmd *cobra.Command, args []string) {
			if util.HasStdinData() {
				fmt.Println("not implemented")
				return
			} else if len(args) == 1 {
				path := strings.TrimSpace(args[0])

				if !util.IsFilePath(path) {
					cobra.CheckErr(fmt.Errorf("expects a valid file name, got=%s", path))
				}

				if toStdOut {
					content := getFileContent(path)
					fmt.Print(content)
					return
				} else {
					viewNote(absolutePath(path))
				}
			}

			fmt.Println("invalid use")
			_ = cmd.Help()
		},
	}
	cmd.Flags().BoolVar(&toStdOut, "stdout", false, "prints the file to stdout instead of pager")

	return cmd
}

// open viewer and view file
func viewNote(absoulutePath string) {

	ok, err := util.FileExists(absoulutePath)
	cobra.CheckErr(err)
	if !ok {
		cobra.CheckErr(fmt.Errorf("no note with name=%s", absoulutePath))
	}

	spawnViewer(absoulutePath)

}

func spawnViewer(filepath string) {
	pagerName := os.ExpandEnv(cfg.Pager)
	cmdWithArgs := strings.Split(pagerName, " ")
	cmdName, cmdArgs := cmdWithArgs[0], cmdWithArgs[1:]

	command, err := exec.LookPath(cmdName)
	cobra.CheckErr(err)

	cmd := append([]string{cmdName}, cmdArgs...)
	cmd = append(cmd, filepath)
	slog.Debug(fmt.Sprintf("cmd: %v\n", cmd))

	if err := syscall.Exec(command, cmd, os.Environ()); err != nil {
		cobra.CheckErr(err)
	}
}
