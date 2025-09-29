// state handles bookkeeping when using nots. Things like "previous note", and other stuff
package state

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fredrikkvalvik/nots/internal/config"
)

const STATE_FILE_NAME = ".nstate"

type State struct {
	PreviousNote *string `toml:"previous-note"`

	// IDEA: add private field to check if state is stale (some other task has saved data)
	// and could be used to avoid possible overwrites
	fileHash string
	cfg      *config.Config
}

func Load(cfg *config.Config) (*State, error) {
	err := ensureFile(cfg.RootDir)
	if err != nil {
		return nil, err
	}

	s := State{
		cfg: cfg,
	}
	_, err = toml.DecodeFile(filepath.Join(cfg.RootDir, STATE_FILE_NAME), &s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state file: %w", err)
	}

	return &s, nil
}

// Save the current state.
func (s *State) Save() error {
	slog.Debug("saving state", "prev_file", s.PreviousNote)

	f, err := os.OpenFile(stateFileName(s.cfg.RootDir), os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(f).Encode(s)
	if err != nil {
		return err
	}

	return nil
}

// called on load to ensure that the a file exists for us to parse to simplify
// the loading process
func ensureFile(root string) error {
	if !fileExists(stateFileName(root)) {
		// we create the file and hope no error occurs
		return createFile(root)
	}
	return nil
}

// helper for creating the file for ensureFile
func createFile(root string) error {
	slog.Debug("crating new state file")

	file, err := os.Create(filepath.Join(root, STATE_FILE_NAME))
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer file.Close()
	return nil
}

func stateFileName(root string) string {
	return filepath.Join(root, STATE_FILE_NAME)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
