package completers

import (
	"log/slog"

	"github.com/fredrikkvalvik/nots/internal/config"
	"github.com/fredrikkvalvik/nots/internal/lister"
	"github.com/spf13/cobra"
)

func FileListCompleter(cfg *config.Config) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fileList(cfg), cobra.ShellCompDirectiveDefault
	}
}

func fileList(cfg *config.Config) []string {
	list, err := lister.ListPathsRecursive(cfg.Dir, false)
	cobra.CheckErr(err)

	completions := []string{}

	for _, path := range list {
		p := path.String()
		completions = append(completions, p)
	}
	slog.Debug("completer running...", "result", completions)

	return completions

}
