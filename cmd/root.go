package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "notes",
	Short: "A brief description of your application",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case printContent:
			printTodaysNote()
			return

		case printFilePath:
			fmt.Fprintln(os.Stdout, todayFilePath())
			return

		case printFileDir:
			fmt.Fprintln(os.Stdout, cfg.dir)
			return

		default:
			openTodaysNote()
		}
	},
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

func init() {
	rootCmd.PersistentFlags().BoolVarP(&printFileDir, "dir", "d", false, "print the notes directory path")
	rootCmd.PersistentFlags().BoolVarP(&printFilePath, "file", "f", false, "print the notes file path")
	rootCmd.PersistentFlags().BoolVarP(&printContent, "print", "p", false, "print the content of the file")
}
