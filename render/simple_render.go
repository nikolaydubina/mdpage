package render

import (
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/nikolaydubina/mdpage/page"
)

type SimplePageRenderConfig struct {
	ContentTitle      string `json:"content_title"`
	RequirementsTitle string `json:"requirements_title"`
	ExampleTitle      string `json:"example_title"`
	EntryLinkPrefix   string `json:"entry_link_prefix"` // e.g. "➡"️
}

// TODO: grow capacity based on estimates
type SimplePageRender struct {
	config SimplePageRenderConfig
}

func NewSimplePageRender(config SimplePageRenderConfig) SimplePageRender {
	return SimplePageRender{config: config}
}

func (s SimplePageRender) RenderPage(page page.Page) string {
	var b strings.Builder
	b.WriteString(s.RenderHeader(page.Header))
	b.WriteRune('\n')
	b.WriteString(s.RenderPageSummary(page))
	b.WriteRune('\n')
	b.WriteString(s.RenderPageContent(page))
	b.WriteRune('\n')
	return b.String()
}

func (s SimplePageRender) RenderHeader(header page.Header) string {
	var b strings.Builder
	b.WriteString(getFileContent(header.ContentURL))
	return b.String()
}

func (s SimplePageRender) RenderPageSummary(page page.Page) string {
	var b strings.Builder

	b.WriteString("## ")
	b.WriteString(s.config.ContentTitle)
	b.WriteRune('\n')
	b.WriteRune('\n')

	entryPrefix := "   "
	for _, group := range page.Groups {
		b.WriteString(" - " + group.Title)
		b.WriteRune('\n')
		for _, entry := range group.Entries {
			b.WriteString(entryPrefix + "+ " + s.RenderSummaryEntry(entry))
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func (s SimplePageRender) RenderSummaryEntry(entry page.Entry) string {
	return "[" + s.config.EntryLinkPrefix + " " + entry.Title + "](" + makeMarkdownTitleLink(entry.Title) + ")"
}

func (s SimplePageRender) RenderPageContent(page page.Page) string {
	var b strings.Builder
	for _, group := range page.Groups {
		b.WriteString(s.RenderGroupContent(group))
	}
	return b.String()
}

func (s SimplePageRender) RenderGroupContent(group page.Group) string {
	var b strings.Builder

	b.WriteString("## ")
	b.WriteString(group.Title)
	b.WriteRune('\n')
	b.WriteRune('\n')

	for _, q := range group.Entries {
		b.WriteString(s.RenderEntryContent(q))
	}

	return b.String()
}

/*
ExampleContentURL string   `yaml:"example_content_url"`
*/

func (s SimplePageRender) RenderEntryContent(entry page.Entry) string {
	var b strings.Builder

	// title
	b.WriteString("### ")
	b.WriteString(s.config.EntryLinkPrefix)
	b.WriteRune(' ')
	b.WriteString(entry.Title)
	b.WriteRune('\n')
	b.WriteRune('\n')

	// description
	b.WriteString(entry.Description)

	// description: author
	if len(entry.Author) > 0 {
		b.WriteString(" — ")
		b.WriteString(entry.Author)
	}

	// description: source
	if len(entry.Source) > 0 {
		b.WriteString(" / ")
		b.WriteString(entry.Source)
	}

	// description: end
	b.WriteRune('\n')
	b.WriteRune('\n')

	// commands
	if len(entry.Commands) > 0 {
		b.WriteRune('\n')
		b.WriteString("```\n")

		for _, q := range entry.Commands {
			b.WriteString(q)
			b.WriteRune('\n')
		}

		b.WriteString("```")
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	// TODO: try fetch from HTTP too
	if entry.ExampleContentURL != "" {
		b.WriteString(renderContentFromFilePath(entry.ExampleContentURL))
		b.WriteRune('\n')
	}

	if entry.ExampleOutputURL != "" {
		b.WriteString(s.config.ExampleTitle)
		b.WriteRune('\n')
		b.WriteString(renderContentFromFilePath(entry.ExampleOutputURL))
		b.WriteRune('\n')
	}

	if entry.ExampleImageURL != "" {
		b.WriteString(renderExampleImage(entry.ExampleImageURL))
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	// requirements
	if len(entry.Requirements) > 0 {
		b.WriteString(s.config.RequirementsTitle)
		b.WriteRune('\n')
		b.WriteString("```\n")

		for _, q := range entry.Requirements {
			b.WriteString(q)
			b.WriteRune('\n')
		}

		b.WriteString("```")
		b.WriteRune('\n')
	}

	b.WriteRune('\n')

	return b.String()
}

func renderExampleImage(imageURL string) string {
	var b strings.Builder
	b.WriteString(`<div align="center">`)
	b.WriteString(`<img src="`)
	b.WriteString(imageURL)
	b.WriteString(`" style="margin: 8px; max-height: 640px;">`)
	b.WriteString(`</div>`)
	b.WriteRune('\n')
	return b.String()
}

// makeMarkdownTitleLink produces name that can be used to reference Markdown section for entry
func makeMarkdownTitleLink(s string) string {
	t := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(s)), " ", "-")

	// keep only alpha numerics and space
	var f strings.Builder
	f.Grow(len(t))
	for _, q := range t {
		if unicode.IsLetter(q) || unicode.IsDigit(q) || q == '-' {
			f.WriteRune(q)
		}
	}

	return "#-" + f.String()
}

// getFileContent unmodified bytes of file as is
func getFileContent(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("can not open file: %s", err)
	}
	s, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("can not read: %s", err)
	}
	return string(s)
}

// renderContentFromFilePath for local files
func renderContentFromFilePath(filePath string) string {
	var b strings.Builder

	b.WriteString("```")
	// markdown syntax
	if strings.HasSuffix(filePath, ".go") {
		b.WriteString("go")
	}
	b.WriteRune('\n')

	b.WriteString(getFileContent(filePath))

	b.WriteRune('\n')
	b.WriteString("```")
	b.WriteRune('\n')

	return b.String()
}
