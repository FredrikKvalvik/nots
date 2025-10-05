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
		Use:     "template",
		Short:   "evaluate templates to be used as text or fragments for notes",
		Example: "nots template <template_name> [--option]",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			templatePath := config.TemplateName(name)

			template, err := template.FromPath(templatePath)
			cobra.CheckErr(err)

			// print the output the os.Stdout, or report an error
			cobra.CheckErr(template.ExecuteWriter(os.Stdout))
		},
	}

	return cmd
}
