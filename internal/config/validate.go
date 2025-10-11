package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type validateFunc func(cfg *Config) error

var valdidationList = []validateFunc{
	// validate defaultOpenMode
	func(cfg *Config) error {
		valid := []string{
			string(OpenPrevious),
			string(OpenSeries),
		}

		if !slices.Contains(valid, string(cfg.DefaultOpenMode)) {
			return fmt.Errorf("invalid value for default-open-mode: %s, expects [%s]", cfg.DefaultOpenMode, strings.Join(valid, "|"))
		}
		return nil
	},
}

func validateConfig(cfg *Config) error {
	errs := []error{}

	for _, vfn := range valdidationList {
		err := vfn(cfg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
