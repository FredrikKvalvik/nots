package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fredrikkvalvik/nots/internal/state"
	"github.com/spf13/cobra"
)

type openPrevious struct {
	view bool
}

func PreviousCmd() *cobra.Command {
	prev := openPrevious{}
	cmd := &cobra.Command{
		Use:   "previous",
		Short: "open the previous note.",

		Args:    cobra.NoArgs,
		Aliases: []string{"p", "prev"},

		Run: prev.openPreviousNoteCmd,
	}

	cmd.Flags().BoolVar(&prev.view, "view", prev.view, "view previous opened note")

	return cmd
}

func (op *openPrevious) openPreviousNoteCmd(cmd *cobra.Command, args []string) {
	s, err := state.Load(cfg)
	cobra.CheckErr(err)

	if s.PreviousNote == nil {
		fmt.Println("Error: no previous file is detected")
		os.Exit(1)
	}

	absoulutePath := filepath.Clean(*s.PreviousNote)

	switch true {
	case op.view:
		viewNote(absoulutePath)

	default:
		openNoteWithSelectedTemplate(absoulutePath)
	}

}
