package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/config"
	"github.com/fredrikkvalvik/nots/pkg/template"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(SeriesCmd())
}

func SeriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "series",
		Aliases: []string{"serie", "s"},
		Short:   "open a note in a series",
		Example: "nots series <series_name> [--option]",
		// returns a list of the defined series
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			validArgs := []string{}
			for _, series := range cfg.NoteSeries {
				validArgs = append(validArgs, series.SeriesName)
			}
			return validArgs, cobra.ShellCompDirectiveDefault
		},

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			var series *config.NoteSeries
			for _, s := range cfg.NoteSeries {
				if s.SeriesName == name {
					series = &s
					break
				}
			}
			if series == nil {
				cobra.CheckErr(fmt.Errorf("there are no note-series name '%s'", name))
				return
			}

			// having multiple expressions in the filenameExpression is an error. make sure the user doesn't escape
			if strings.Contains(series.SeriesFilenameExpression, "{{") || strings.Contains(series.SeriesFilenameExpression, "}}") {
				cobra.CheckErr(fmt.Errorf("invalid format on series filename expression. expression should not contain opening/closing expression signifiers"))
			}

			filename, err := template.Eval(fmt.Sprintf("{{ %s }}", series.SeriesFilenameExpression))
			cobra.CheckErr(err)

			if !strings.HasSuffix(filename, ".md") {
				filename += ".md"
			}

			seriesDir := series.SeriesName
			if series.DirName != "" {
				seriesDir = series.DirName
			}

			absolutePath := filepath.Join(cfg.RootDir, seriesDir, filename)

			openNote(absolutePath, series.TemplateName)
		},
	}

	return cmd
}
