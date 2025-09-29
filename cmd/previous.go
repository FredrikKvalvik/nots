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

			openNote(*s.PreviousNote)
		},
	}

	return cmd
}
