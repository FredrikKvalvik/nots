package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/config"
	"github.com/fredrikkvalvik/nots/pkg/template"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
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
			name := args[0]
			err := loadTemplate(name, os.Stdout)
			cobra.CheckErr(err)
		},
	}

	return cmd
}

// load the template into w. the write might be successful
// and still return and error. all other errors are os errors and can be checked for
func loadTemplate(name string, w io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(config.TemplateDir())
	if err != nil {
		return err
	}

	template, err := template.FromPath(name)
	if err != nil {
		return err
	}

	ts := &templateSymbolsFuncs{
		filename: strings.TrimSuffix(name, ".md"),
	}

	ts.registerSymbols(template)

	// print the output the os.Stdout, or report an error
	err = template.ExecuteWriter(w)
	if err != nil {
		slog.Error("falied to execute template", "error", err)
		return err
	}

	// change back to where we were before to not mess with execution
	err = os.Chdir(cwd)
	if err != nil {
		return err
	}
	return nil
}

// use type to keep namespace clean
type templateSymbolsFuncs struct {
	filename string
}

func (ts *templateSymbolsFuncs) registerSymbols(t *template.Template) {
	t.RegisterFnValue("joke", "calls an API and returns a random joke", ts.joke)
	t.RegisterStringValue("filename", ts.filename)
}

// uses a jokes api to get a random joke
func (ts *templateSymbolsFuncs) joke() (object.Object, error) {
	baseUrl := "https://v2.jokeapi.dev"
	url := fmt.Sprintf("%s/joke/Programming?format=txt&blacklistFlags=nsfw,racist,sexist", baseUrl)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &object.ObjectString{Val: string(bytes)}, nil
}
