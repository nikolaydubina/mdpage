package page

import (
	"errors"
	"fmt"
)

// Page is full page contents.
type Page struct {
	Header      string         `yaml:"header"`
	Groups      []Group        `yaml:"groups"`
	Contents    ContentsConfig `yaml:"contents"`
	EntryConfig EntryConfig    `yaml:"entry"`
}

func (v Page) Validate() error {
	var errs []error
	for i, q := range v.Groups {
		if err := q.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid group(%d): %w", i, err))
		}
	}
	return errors.Join(errs...)
}

type ContentsConfig struct {
	Title string `yaml:"title"`
}

type EntryRequirementsConfig struct {
	Title string `yaml:"title"`
}

type EntryExampleConfig struct {
	Title string `yaml:"title"`
}

type EntryConfig struct {
	TitlePrefix  string                  `yaml:"title_prefix"`
	Back         string                  `yaml:"back"`
	Requirements EntryRequirementsConfig `yaml:"requirements"`
	Example      EntryExampleConfig      `yaml:"example"`
}
