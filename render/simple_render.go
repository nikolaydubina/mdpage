package render

import (
	"io"
	"strings"

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

func (s SimplePageRender) RenderPageSummary(out io.StringWriter, p page.Page) {
	out.WriteString("## ")
	out.WriteString(p.Contents.Title)
	out.WriteString("\n\n")

	entryPrefix := "   "
	for _, group := range p.Groups {
		out.WriteString(" - " + group.Title)
		out.WriteString("\n")
		for _, entry := range group.Entries {
			if group.Type == page.MarkdownListGroupType {
				out.WriteString(entryPrefix + "+ " + Link{Name: p.EntryConfig.TitlePrefix + " " + entry.Title, URL: RelativeLink(group.Title)}.Render())
				out.WriteString("\n")
				continue
			}

			out.WriteString(entryPrefix + "+ " + RenderSummaryEntry(out, entry, p.EntryConfig))
			out.WriteString("\n")
		}
	}
}

func RenderSummaryEntry(out io.StringWriter, entry page.Entry, config page.EntryConfig) string {
	title := EnrichedTitleWithProjectNameLinkTextSummary(entry)
	titleEntry := EnrichedTitleWithProjectNameLinkText(entry) // special version to avoid url
	return "[" + config.TitlePrefix + " " + title + "](" + makeMarkdownTitleLink(titleEntry) + ")"
}

func (s SimplePageRender) RenderPageContent(out io.StringWriter, page page.Page) {
	for _, group := range page.Groups {
		s.RenderGroupContent(out, group, page.EntryConfig, page.Contents)
	}
}

func (s SimplePageRender) RenderGroupContent(out io.StringWriter, group page.Group, config page.EntryConfig, configContent page.ContentsConfig) {
	out.WriteString("## ")
	out.WriteString(group.Title)
	out.WriteString("\n\n")

	for _, q := range group.Entries {
		if group.Type == page.MarkdownListGroupType {
			s.RenderMarkdownListGroupEntryContent(out, q, config, configContent)
			continue
		}
		s.RenderEntryContent(out, q, config, configContent)
		out.WriteString("\n")
	}
}

func (s SimplePageRender) RenderMarkdownListGroupEntryContent(out io.StringWriter, entry page.Entry, config page.EntryConfig, configContent page.ContentsConfig) {
	out.WriteString("- ")
	out.WriteString("[" + entry.Title + "](" + entry.URL + ")")
	out.WriteString("\n\n")
}

func RelativeLink(configTitle string) string {
	return "#" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(configTitle)), " ", "-")
}

func (s SimplePageRender) RenderEntryContent(out io.StringWriter, entry page.Entry, config page.EntryConfig, configContent page.ContentsConfig) {
	// title
	out.WriteString("### ")
	if config.Back != "" {
		v := Link{Name: config.Back, URL: RelativeLink(configContent.Title)}
		v.RenderTo(out)
	}
	out.WriteString(config.TitlePrefix)
	out.WriteString(" ")
	out.WriteString(EnrichedTitleWithProjectNameLink(entry))
	out.WriteString("\n\n")

	// description
	out.WriteString(entry.Description)
	if len(entry.Author) > 0 {
		out.WriteString(" — ")
		out.WriteString(s.authorRenderer.Render(entry.Author))
	}
	out.WriteString("\n\n")

	if len(entry.Commands) > 0 {
		out.WriteString("\n")
		out.WriteString("```\n")

		for _, q := range entry.Commands {
			out.WriteString(q)
			out.WriteString("\n")
		}

		out.WriteString("```")
		out.WriteString("\n\n")
	}

	if entry.ExampleContent != "" {
		out.WriteString("```")
		if entry.ExampleContentExt != "" {
			out.WriteString(entry.ExampleContentExt)
		}
		out.WriteString("\n")
		out.WriteString(entry.ExampleContent)
		out.WriteString("```")
		out.WriteString("\n\n")
	}

	if entry.ExampleOutput != "" {
		out.WriteString(config.Example.Title)
		out.WriteString("\n")
		out.WriteString("```\n")
		out.WriteString(entry.ExampleOutput)
		out.WriteString("```")
		out.WriteString("\n\n")
	}

	if entry.ExampleImageURL != "" {
		Image{ImageURL: entry.ExampleImageURL}.RenderTo(out)
		out.WriteString("\n\n")
	}

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
}
