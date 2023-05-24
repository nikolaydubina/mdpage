package render

import (
	"io"
	"strings"
	"unicode"

	"github.com/nikolaydubina/mdpage/page"
)

type AuthorRenderer interface{ Render(s string) string }

type SimplePageRender struct {
	authorRenderer AuthorRenderer
}

func NewSimplePageRender(authorRenderer AuthorRenderer) SimplePageRender {
	return SimplePageRender{
		authorRenderer: authorRenderer,
	}
}

func (s SimplePageRender) RenderTo(out io.StringWriter, page page.Page) {
	s.RenderHeader(out, page.Header)
	out.WriteString("\n")
	s.RenderPageSummary(out, page)
	out.WriteString("\n")
	s.RenderPageContent(out, page)
	out.WriteString("\n")
}

func (s SimplePageRender) RenderHeader(out io.StringWriter, header string) { out.WriteString(header) }

func (s SimplePageRender) RenderPageSummary(out io.StringWriter, page page.Page) {
	out.WriteString("## ")
	out.WriteString(page.Contents.Title)
	out.WriteString("\n")
	out.WriteString("\n")

	entryPrefix := "   "
	for _, group := range page.Groups {
		out.WriteString(" - " + group.Title)
		out.WriteString("\n")
		for _, entry := range group.Entries {
			out.WriteString(entryPrefix + "+ " + RenderSummaryEntry(out, entry, page.EntryConfig))
			out.WriteString("\n")
		}
	}
}

func RenderSummaryEntry(out io.StringWriter, entry page.Entry, config page.EntryConfig) string {
	return "[" + config.TitlePrefix + " " + entry.Title + "](" + makeMarkdownTitleLink(entry.Title) + ")"
}

func (s SimplePageRender) RenderPageContent(out io.StringWriter, page page.Page) {
	for _, group := range page.Groups {
		s.RenderGroupContent(out, group, page.EntryConfig, page.Contents)
	}
}

func (s SimplePageRender) RenderGroupContent(out io.StringWriter, group page.Group, config page.EntryConfig, configContent page.ContentsConfig) {
	out.WriteString("## ")
	out.WriteString(group.Title)
	out.WriteString("\n")
	out.WriteString("\n")

	for _, q := range group.Entries {
		if group.Type == page.MarkdownListGroupType {
			s.RenderMarkdownListGroupEntryContent(out, q, config, configContent)
			continue
		}

		s.RenderEntryContent(out, q, config, configContent)
	}
}

func (s SimplePageRender) RenderMarkdownListGroupEntryContent(out io.StringWriter, entry page.Entry, config page.EntryConfig, configContent page.ContentsConfig) {
	out.WriteString("- ")
	out.WriteString("[" + entry.Title + "](" + entry.URL + ")")
	out.WriteString("\n")
	out.WriteString("\n")
}

type PageLink struct {
	Name string
	URL  string
}

func (v PageLink) RenderTo(out io.StringWriter) {
	out.WriteString("[" + v.Name + "]")
	out.WriteString("(" + v.URL + ")")
}

func (s SimplePageRender) RenderEntryContent(out io.StringWriter, entry page.Entry, config page.EntryConfig, configContent page.ContentsConfig) {
	// title
	out.WriteString("### ")
	if config.Back != "" {
		url := "#" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(configContent.Title)), " ", "-")
		v := PageLink{Name: config.Back, URL: url}
		v.RenderTo(out)
	}
	out.WriteString(config.TitlePrefix)
	out.WriteString(" ")
	out.WriteString(entry.Title)
	out.WriteString("\n")
	out.WriteString("\n")

	// description
	out.WriteString(entry.Description)

	// description: author
	if len(entry.Author) > 0 {
		out.WriteString(" — ")
		out.WriteString(s.authorRenderer.Render(entry.Author))
	}

	// description: source
	if len(entry.Source) > 0 {
		out.WriteString(" / ")
		out.WriteString(entry.Source)
	}

	// description: end
	out.WriteString("\n")
	out.WriteString("\n")

	// commands
	if len(entry.Commands) > 0 {
		out.WriteString("\n")
		out.WriteString("```\n")

		for _, q := range entry.Commands {
			out.WriteString(q)
			out.WriteString("\n")
		}

		out.WriteString("```")
		out.WriteString("\n")
		out.WriteString("\n")
	}

	if entry.ExampleContent != "" {
		out.WriteString("```")
		if entry.ExampleContentExt != "" {
			out.WriteString(entry.ExampleContentExt)
		}
		out.WriteString("\n")
		out.WriteString(entry.ExampleContent)
		out.WriteString("```")
		out.WriteString("\n")
		out.WriteString("\n")
	}

	if entry.ExampleOutput != "" {
		out.WriteString(config.Example.Title)
		out.WriteString("\n")
		out.WriteString("```\n")
		out.WriteString(entry.ExampleOutput)
		out.WriteString("```")
		out.WriteString("\n")
		out.WriteString("\n")
	}

	if entry.ExampleImageURL != "" {
		Image{ImageURL: entry.ExampleImageURL}.RenderTo(out)
		out.WriteString("\n")
		out.WriteString("\n")
	}

	// requirements
	if len(entry.Requirements) > 0 {
		out.WriteString(config.Requirements.Title)
		out.WriteString("\n")
		out.WriteString("```\n")

		for _, q := range entry.Requirements {
			out.WriteString(q)
			out.WriteString("\n")
		}

		out.WriteString("```")
		out.WriteString("\n")
	}

	out.WriteString("\n")
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
