package page

import (
	"errors"
	"strings"
)

// Entry in page containing single post.
type Entry struct {
	Title             string   `yaml:"title"`
	Description       string   `yaml:"description"`
	Author            string   `yaml:"author"`
	Source            string   `yaml:"source"`
	ExampleImageURL   string   `yaml:"example_image_url"`
	ExampleContentURL string   `yaml:"example_content_url"`
	ExampleOutputURL  string   `yaml:"example_output_url"`
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
	if !strings.HasSuffix(v.Description, ".") {
		return ErrInvalidEntryDescriptionEndsWithoutPeriod
	}
	return nil
}
