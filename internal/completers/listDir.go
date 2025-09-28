package completers

import (
	"path/filepath"

	"github.com/fredrikkvalvik/nots/internal/config"
	"github.com/fredrikkvalvik/nots/internal/lister"
	"github.com/spf13/cobra"
)

func DirListCompleter(cfg *config.Config) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := lister.ListPathsRecursive(cfg.Dir, true)
		cobra.CheckErr(err)
		return directoryList(list, toComplete), cobra.ShellCompDirectiveDefault
	}
}

func directoryList(list []lister.Path, toComplete string) []string {
	// toCompletePath := lister.NewPath(toComplete)

	completions := []string{}

	for _, path := range list {
		p := path.String()
		p = filepath.Dir(p)
		if p == "." {
			continue
		}

		// skip the current directory

		// path = path.Pop()
		completions = append(completions, p)

		// if path.ContainsSubset(toCompletePath) {
		// 	completions = append(completions, filepath.Dir(path.String()))
		// }
	}

	return completions

}
