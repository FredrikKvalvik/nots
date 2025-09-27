package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fredrikkvalvik/nots/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolVarP(&printFileDir, "dir", "d", false, "print the notes directory path")
	rootCmd.Flags().BoolVarP(&printFilePath, "file", "f", false, "print the notes file path")
	rootCmd.Flags().BoolVarP(&printContent, "print", "p", false, "print the content of the file")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notes",
	Short: "utility for managing daily notes",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		if util.HasStdinData() {
			rootHandleStdin(cmd, args)
			return
		}
		rootHandleCmds(cmd, args)
	},
}

func rootHandleStdin(_ *cobra.Command, _ []string) {
	b, err := io.ReadAll(os.Stdin)
	cobra.CheckErr(err)

	fileName := string(b)
	fileName = strings.TrimSpace(fileName)

	if util.IsFileName(fileName) {
		openNote(filePath(fileName))
		return
	}

	if util.IsFilePath(fileName) {
		openNote(fileName)
		return
	}

	cobra.CheckErr(fmt.Errorf("could not resolve input: %s", fileName))
}

func rootHandleCmds(_ *cobra.Command, _ []string) {
	switch true {
	case printContent:
		printTodaysNote()
		return

	case printFilePath:
		_, _ = fmt.Fprintln(os.Stdout, todayFilePath())
		return

	case printFileDir:
		_, _ = fmt.Fprintln(os.Stdout, cfg.Dir)
		return

	default:
		openTodaysNote()
	}

}

var (
	printFilePath = false
	printFileDir  = false
	printContent  = false
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
