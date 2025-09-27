package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/lister"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ListCmd())
}

func ListCmd() *cobra.Command {
	var fullPath bool
	var listAll bool
	var listRecursive bool

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list all notes",
		Aliases: []string{"ls"},

		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cfg.Dir
			if len(args) > 0 {
				dirPath = filepath.Join(dirPath, args[0])
			}

			var err error
			var entries [][]string
			if listRecursive {
				entries, err = lister.ListPathsRecursive(dirPath)
				cobra.CheckErr(err)
			} else {
				ent, err := lister.ListDir(dirPath, listAll)
				cobra.CheckErr(err)
				for _, e := range ent {
					entries = append(entries, []string{e})
				}
			}

			var str strings.Builder
			if fullPath {
				for _, name := range entries {

					fmt.Fprintf(&str, "%s/%s\n", cfg.Dir, strings.Join(name, "/"))
				}
			} else {
				for _, name := range entries {
					fmt.Fprintln(&str, strings.Join(name, "/"))
				}
			}

			_, err = fmt.Fprint(os.Stdout, str.String())
			cobra.CheckErr(err)
		},
	}

	cmd.Flags().BoolVarP(&fullPath, "full-path", "v", false, "prints the files with absoule paths")
	cmd.Flags().BoolVarP(&listAll, "all", "a", false, "include directories in list")
	cmd.Flags().BoolVarP(&listRecursive, "recursive", "r", false, "list recursivly")

	return cmd
}
