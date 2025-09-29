package completers

import (
	"slices"

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
	toCompletePath := lister.NewPath(toComplete)
	toCompleteLength := len(toCompletePath)

	completions := []string{}

	for _, path := range list {
		// remove the file name, if one is there
		if len(path) < toCompleteLength {
			continue
		}

		// remove the file at the end of the path, if one is there
		dir := path.PopFile()
		// this means we removed a file and are left with "nothing", meaning the file was at root. move on.
		if len(dir) == 0 {
			continue
		}

		// if toComplete is empty, we simple add all paths first item
		if toCompleteLength == 0 {
			completions = append(completions, dir[0])
			continue
		}

		// if we are further down, we need to find out if the toComplete is a subPath of
		// the path we are comparing. then we return the item of the same index as
		// the last element in toCompletePath
		if dir.ContainsSubset(toCompletePath) {
			completions = append(completions, dir[toCompleteLength-1])
			continue
		}
	}

	return slices.Compact(completions)

}
