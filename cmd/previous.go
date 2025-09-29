package cmd

import (
	"fmt"
	"os"

	"github.com/fredrikkvalvik/nots/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(PreviousCmd())
}

func PreviousCmd() *cobra.Command {
	var view bool

	cmd := &cobra.Command{
		Use:   "previous",
		Short: "open the previous note.",

		Args:    cobra.NoArgs,
		Aliases: []string{"p", "prev"},

		Run: func(cmd *cobra.Command, args []string) {
			// TODO: implement logic for opening the previous note
			// It should store its state in nots root dir.

			s, err := state.Load(cfg)
			cobra.CheckErr(err)

			if s.PreviousNote == nil {
				fmt.Println("Error: no previous file is detected")
				os.Exit(1)
			}

			switch true {
			case view:
				viewNote(*s.PreviousNote)

			default:
				openNote(*s.PreviousNote)
			}

		},
	}

	cmd.Flags().BoolVar(&view, "view", view, "view previous opened note")

	return cmd
}
