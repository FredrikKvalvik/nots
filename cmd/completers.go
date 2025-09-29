package cmd

import (
	"path/filepath"

	"github.com/fredrikkvalvik/nots/internal/lister"
	"github.com/spf13/cobra"
)

func fileListCompleter(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := lister.ListPathsRecursive(cfg.RootDir, false)
	cobra.CheckErr(err)

	completions := []string{}

	for _, path := range list {
		p := path.String()
		completions = append(completions, p)
	}

	return completions, cobra.ShellCompDirectiveDefault
}

func dirListCompleter(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	list, err := lister.ListPathsRecursive(cfg.RootDir, false)
	cobra.CheckErr(err)

	toCompletePath := lister.NewPath(toComplete)

	completions := []string{}

	for _, path := range list {
		p := path.String()
		p = filepath.Dir(p)
		if p == "." {
			continue
		}

		// skip the current directory

		completions = append(completions, p)

		if path.ContainsSubset(toCompletePath) {
			completions = append(completions, filepath.Dir(path.String()))
		}
	}

	return completions, cobra.ShellCompDirectiveDefault
}
