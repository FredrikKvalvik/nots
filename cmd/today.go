package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(TodayCmd())
}

func TodayCmd() *cobra.Command {
	var view bool

	cmd := &cobra.Command{
		Use:   "today",
		Short: "open the todays note.",

		Args: cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			slog.Debug("running today command")

			path := todayFilePath()

			slog.Debug("today file path", "path", path)
			switch true {
			case view:
				viewNote(path)

			default:
				openNote(path)
			}

		},
	}

	cmd.Flags().BoolVar(&view, "view", view, "view previous opened note")

	return cmd
}
