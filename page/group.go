package page

import (
	"errors"
	"fmt"
	"strings"
)

// Group is named ordered collection of Entries.
type Group struct {
	Title   string    `yaml:"title"`
	Type    GroupType `yaml:"type"`
	Entries []Entry   `yaml:"entries"`
}

type GroupType string

const (
	MarkdownListGroupType GroupType = "md-list"
)

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
	return errors.Join(errs...)
}
