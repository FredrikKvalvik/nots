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
	vh := viewHandler{}

	cmd := &cobra.Command{
		Use:     "view",
		Short:   "view a specified note",
		Aliases: []string{"ls"},

		Run: func(cmd *cobra.Command, args []string) {
			if util.HasStdinData() {
				vh.viewHandleStdin(cmd, args)
				return
			} else if len(args) == 1 {
				vh.viewHandleCmd(cmd, args)
				return
			}

			cobra.CheckErr("invalid use")
		},
	}
	cmd.Flags().BoolVar(&vh.toStdOut, "stdout", false, "prints the file to stdout instead of pager")

	return cmd
}

type viewHandler struct {
	toStdOut bool
}

func (vh *viewHandler) viewHandleStdin(cmd *cobra.Command, args []string) {}

func (vh *viewHandler) viewHandleCmd(_ *cobra.Command, args []string) {
	filename := strings.TrimSpace(args[0])

	if !util.IsFileName(filename) {
		cobra.CheckErr(fmt.Errorf("expects a valid file name, got=%s", filename))
	}
	filePath := filePath(filename)

	ok, err := util.FileExists(filePath)
	cobra.CheckErr(err)
	if !ok {
		cobra.CheckErr(fmt.Errorf("no note with name=%s", filePath))
	}

	content := getNoteContent(filePath)

	if vh.toStdOut {
		fmt.Printf("hello!")
		fmt.Print(content)
		return
	}

	pagerName := os.ExpandEnv(cfg.Pager)
	cmdWithArgs := strings.Split(pagerName, " ")
	cmdName, cmdArgs := cmdWithArgs[0], cmdWithArgs[1:]

	command, err := exec.LookPath(cmdName)
	cobra.CheckErr(err)

	cmd := append([]string{cmdName}, cmdArgs...)
	cmd = append(cmd, filePath)
	slog.Debug(fmt.Sprintf("cmd: %v\n", cmd))

	if err := syscall.Exec(command, cmd, os.Environ()); err != nil {
		cobra.CheckErr(err)
	}
}

// func viewHandleStdin(cmd *cobra.Command, args []string) {}
// func viewHandleCmd(cmd *cobra.Command, args []string) {
// 	filename := strings.TrimSpace(args[0])

// 	if !util.IsFileName(filename) {
// 		cobra.CheckErr(fmt.Errorf("expects a valid file name, got=%s", filename))
// 	}

// 	ok, err := util.FileExists(filename)
// 	cobra.CheckErr(err)
// 	if !ok {
// 		cobra.CheckErr(fmt.Errorf("no note with name=%s", filename))
// 	}

// 	if
// }
