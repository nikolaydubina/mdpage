package page

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/multierr"
)

// Entry in page containing single post.
type Entry struct {
	Title             string   `yaml:"title"`
	Description       string   `yaml:"description"`
	Author            string   `yaml:"author"`
	Source            string   `yaml:"source"`
	ExampleImageURL   string   `yaml:"example_image_url"`
	ExampleContentURL string   `yaml:"example_content_url"`
	Requirements      []string `yaml:"requirements"`
	Commands          []string `yaml:"commands"`
}

var (
	ErrInvalidEntryTitleEndsWithPeriod          = errors.New("title should not end with period")
	ErrInvalidEntryDescriptionEndsWithoutPeriod = errors.New("description should end with period")
)

func (v Entry) Validate() error {
	if len(v.Title) == 0 {
		return ErrTitleIsEmpty
	}
	if strings.HasSuffix(v.Title, ".") {
		return ErrInvalidEntryTitleEndsWithPeriod
	}
	if !strings.HasSuffix(v.Description, ".") {
		return ErrInvalidEntryDescriptionEndsWithoutPeriod
	}
	return nil
}

// Group is named ordered collection of Entries.
type Group struct {
	Title   string  `yaml:"title"`
	Entries []Entry `yaml:"entries"`
}

var (
	ErrTitleIsEmpty                    = errors.New("title is empty")
	ErrInvalidGroupTitleEndsWithPeriod = errors.New("title should not end with period")
)

func (v Group) Validate() error {
	if len(v.Title) == 0 {
		return ErrTitleIsEmpty
	}

	if strings.HasSuffix(v.Title, ".") {
		return ErrInvalidGroupTitleEndsWithPeriod
	}

	var errs []error
	for i, q := range v.Entries {
		if err := q.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid entry(%d): %w", i, err))
		}
	}
	return multierr.Combine(errs...)
}

// Page is full page contents.
type Page struct {
	Groups []Group `yaml:"groups"`
}

func (v Page) Validate() error {
	var errs []error
	for i, q := range v.Groups {
		if err := q.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid group(%d): %w", i, err))
		}
	}
	return multierr.Combine(errs...)
}
