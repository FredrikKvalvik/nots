package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/lister"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	var fullPath bool
	var listAll bool
	var listRecursive bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all notes",
		// ValidArgsFunction: completers.DirListCompleter(cfg),
		ValidArgsFunction: dirCompleter(cfg.RootDir),
		Aliases:           []string{"ls"},

		Run: func(cmd *cobra.Command, args []string) {
			dirPath := cfg.RootDir
			if len(args) > 0 {
				dirPath = filepath.Join(dirPath, args[0])
			}

			var err error
			var entries []lister.Path
			if listRecursive {
				entries, err = lister.ListPathsRecursive(dirPath, listAll)
				cobra.CheckErr(err)
			} else {
				entries, err = lister.ListPaths(dirPath, listAll)
				cobra.CheckErr(err)
			}

			var str strings.Builder
			if fullPath {
				for _, path := range entries {
					fmt.Fprintf(&str, "%s/%s\n", cfg.RootDir, path.String())
				}
			} else {
				for _, path := range entries {
					fmt.Fprintln(&str, path.String())
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
