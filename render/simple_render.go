package render

import (
	"strings"
	"unicode"

	"github.com/nikolaydubina/mdpage/page"
)

type AuthorRenderer interface{ Render(s string) string }

type SimplePageRender struct {
	b              *strings.Builder
	authorRenderer AuthorRenderer
}

func NewSimplePageRender(authorRenderer AuthorRenderer) SimplePageRender {
	return SimplePageRender{
		b:              &strings.Builder{},
		authorRenderer: authorRenderer,
	}
}

func (s SimplePageRender) estimatePageSize(page page.Page) int {
	c := 0
	c += 1000 // header
	for _, q := range page.Groups {
		c += len(q.Entries) * 1000
	}
	return c
}

func (s SimplePageRender) RenderPage(page page.Page) string {
	s.b.Reset()
	s.b.Grow(s.estimatePageSize(page))

	s.RenderHeader(page.Header)
	s.b.WriteRune('\n')
	s.RenderPageSummary(page)
	s.b.WriteRune('\n')
	s.RenderPageContent(page)
	s.b.WriteRune('\n')
	return s.b.String()
}

func (s SimplePageRender) RenderHeader(header string) { s.b.WriteString(header) }

func (s SimplePageRender) RenderPageSummary(page page.Page) {
	s.b.WriteString("## ")
	s.b.WriteString(page.Contents.Title)
	s.b.WriteRune('\n')
	s.b.WriteRune('\n')

	entryPrefix := "   "
	for _, group := range page.Groups {
		s.b.WriteString(" - " + group.Title)
		s.b.WriteRune('\n')
		for _, entry := range group.Entries {
			s.b.WriteString(entryPrefix + "+ " + s.RenderSummaryEntry(entry, page.EntryConfig))
			s.b.WriteRune('\n')
		}
	}
}

func (s SimplePageRender) RenderSummaryEntry(entry page.Entry, config page.EntryConfig) string {
	return "[" + config.TitlePrefix + " " + entry.Title + "](" + makeMarkdownTitleLink(entry.Title) + ")"
}

func (s SimplePageRender) RenderPageContent(page page.Page) {
	for _, group := range page.Groups {
		s.RenderGroupContent(group, page.EntryConfig, page.Contents)
	}
}

func (s SimplePageRender) RenderGroupContent(group page.Group, config page.EntryConfig, configContent page.ContentsConfig) {
	s.b.WriteString("## ")
	s.b.WriteString(group.Title)
	s.b.WriteRune('\n')
	s.b.WriteRune('\n')

	for _, q := range group.Entries {
		s.RenderEntryContent(q, config, configContent)
	}
}

type PageLink struct {
	Name string
	URL  string
}

func (v PageLink) Render(b *strings.Builder) {
	b.WriteString("[" + v.Name + "]")
	b.WriteString("(" + v.URL + ")")
}

func (s SimplePageRender) RenderEntryContent(entry page.Entry, config page.EntryConfig, configContent page.ContentsConfig) {
	// title
	s.b.WriteString("### ")
	if config.Back != "" {
		url := "#" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(configContent.Title)), " ", "-")
		v := PageLink{Name: config.Back, URL: url}
		v.Render(s.b)
	}
	s.b.WriteString(config.TitlePrefix)
	s.b.WriteRune(' ')
	s.b.WriteString(entry.Title)
	s.b.WriteRune('\n')
	s.b.WriteRune('\n')

	// description
	s.b.WriteString(entry.Description)

	// description: author
	if len(entry.Author) > 0 {
		s.b.WriteString(" — ")
		s.b.WriteString(s.authorRenderer.Render(entry.Author))
	}

	// description: source
	if len(entry.Source) > 0 {
		s.b.WriteString(" / ")
		s.b.WriteString(entry.Source)
	}

	// description: end
	s.b.WriteRune('\n')
	s.b.WriteRune('\n')

	// commands
	if len(entry.Commands) > 0 {
		s.b.WriteRune('\n')
		s.b.WriteString("```\n")

		for _, q := range entry.Commands {
			s.b.WriteString(q)
			s.b.WriteRune('\n')
		}

		s.b.WriteString("```")
		s.b.WriteRune('\n')
		s.b.WriteRune('\n')
	}

	if entry.ExampleContent != "" {
		s.b.WriteString("```")
		if entry.ExampleContentExt != "" {
			s.b.WriteString(entry.ExampleContentExt)
		}
		s.b.WriteRune('\n')
		s.b.WriteString(entry.ExampleContent)
		s.b.WriteString("```")
		s.b.WriteRune('\n')
		s.b.WriteRune('\n')
	}

	if entry.ExampleOutput != "" {
		s.b.WriteString(config.Example.Title)
		s.b.WriteRune('\n')
		s.b.WriteString("```\n")
		s.b.WriteString(entry.ExampleOutput)
		s.b.WriteString("```")
		s.b.WriteRune('\n')
		s.b.WriteRune('\n')
	}

	if entry.ExampleImageURL != "" {
		s.b.WriteString(renderExampleImage(entry.ExampleImageURL))
		s.b.WriteRune('\n')
		s.b.WriteRune('\n')
	}

	// requirements
	if len(entry.Requirements) > 0 {
		s.b.WriteString(config.Requirements.Title)
		s.b.WriteRune('\n')
		s.b.WriteString("```\n")

		for _, q := range entry.Requirements {
			s.b.WriteString(q)
			s.b.WriteRune('\n')
		}

		s.b.WriteString("```")
		s.b.WriteRune('\n')
	}

	s.b.WriteRune('\n')
}

func renderExampleImage(imageURL string) string {
	var b strings.Builder
	b.Grow(1000)
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
