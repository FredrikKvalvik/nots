package cmd

import (
	"os"

	"github.com/fredrikkvalvik/nots/internal/config"
	"github.com/fredrikkvalvik/nots/pkg/template"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(TemplateCmd())
}

func TemplateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "template",
		Short:             "evaluate templates to be used as text or fragments for notes",
		Example:           "nots template <template_name> [--option]",
		ValidArgsFunction: fileCompleter(config.TemplateDir()),

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			err := os.Chdir(config.TemplateDir())
			cobra.CheckErr(err)

			name := args[0]
			template, err := template.FromPath(name)
			cobra.CheckErr(err)

			// print the output the os.Stdout, or report an error
			cobra.CheckErr(template.ExecuteWriter(os.Stdout))
		},
	}

	return cmd
}

func ptr[T any](v T) *T {
	return &v
}
