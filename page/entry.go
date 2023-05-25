package page

import (
	"errors"
	"strings"
)

// Entry in page containing single post.
type Entry struct {
	Title             string   `yaml:"title"`
	Name              string   `yaml:"name"`
	URL               string   `yaml:"url"`
	Description       string   `yaml:"description"`
	Author            string   `yaml:"author"`
	ExampleImageURL   string   `yaml:"example_image_url"`
	ExampleContent    string   `yaml:"example_content"`
	ExampleContentExt string   `yaml:"example_content_ext"`
	ExampleOutput     string   `yaml:"example_output"`
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
	if len(v.Description) > 0 && !strings.HasSuffix(v.Description, ".") {
		return ErrInvalidEntryDescriptionEndsWithoutPeriod
	}
	return nil
}
