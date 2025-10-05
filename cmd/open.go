package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(OpenCmd())
}

func OpenCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "open note with spesified name. (must end with .md)",

		ValidArgsFunction: fileCompleter(cfg.RootDir),
		Args:              cobra.RangeArgs(0, 1),
		Aliases:           []string{"o"},

		Run: func(cmd *cobra.Command, args []string) {
			var input string
			if util.HasStdinData() {
				b, err := io.ReadAll(os.Stdin)
				cobra.CheckErr(err)

				input = strings.TrimSpace(string(b))
			} else {
				if len(args) < 1 {
					cobra.CheckErr("you need to specify a file name")
					cobra.CheckErr(cmd.Help())
					os.Exit(1)
				}
				input = args[0]
			}
			input = strings.TrimSpace(input)

			if !util.IsFilePath(input) {
				cobra.CheckErr(fmt.Errorf("invalid file path: %s", input))
			}
			fp := absolutePath(input)

			openNote(fp)

		},
	}

	return cmd
}
