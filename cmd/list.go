package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ListCmd())
}

func ListCmd() *cobra.Command {
	var fullPath bool

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list all notes",
		Aliases: []string{"ls"},

		Run: func(cmd *cobra.Command, args []string) {
			files, err := os.ReadDir(cfg.Dir)
			cobra.CheckErr(err)

			var str strings.Builder
			if fullPath {
				for _, f := range files {
					fmt.Fprintf(&str, "%s/%s\n", cfg.Dir, f.Name())
				}
			} else {
				for _, f := range files {
					if f.IsDir() {
						continue
					}
					fmt.Fprintln(&str, f.Name())
				}
			}

			_, err = fmt.Fprint(os.Stdout, str.String())
			cobra.CheckErr(err)
		},
	}

	cmd.Flags().BoolVarP(&fullPath, "full-path", "v", false, "prints the files with absoule paths")

	return cmd
}
