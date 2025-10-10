package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
			s, err := state.Load(cfg)
			cobra.CheckErr(err)

			if s.PreviousNote == nil {
				fmt.Println("Error: no previous file is detected")
				os.Exit(1)
			}

			absoulutePath := filepath.Clean(*s.PreviousNote)

			switch true {
			case view:
				viewNote(absoulutePath)

			default:
				openNoteWithSelectedTemplate(absoulutePath)
			}

		},
	}

	cmd.Flags().BoolVar(&view, "view", view, "view previous opened note")

	return cmd
}
