package render

import (
	"strings"
	"unicode"

	"github.com/nikolaydubina/mdpage/page"
)

type SimplePageRender struct {
	contentTitle string
}

func NewSimplePageRender(
	contentTitle string,
) SimplePageRender {
	return SimplePageRender{contentTitle: contentTitle}
}

func (s SimplePageRender) RenderPage(page page.Page) string {
	return s.RenderPageSummary(page) + "\n" + s.RenderPageContent(page)
}

func (s SimplePageRender) RenderPageSummary(page page.Page) string {
	var c string

	c += s.contentTitle + "\n"
	c += "\n"

	entryPrefix := "   "
	for _, group := range page.Groups {
		c += " - " + group.Title + "\n"
		for _, entry := range group.Entries {
			c += entryPrefix + "+ " + s.RenderSummaryEntry(entry) + "\n"
		}
	}

	return c
}

func (s SimplePageRender) RenderSummaryEntry(entry page.Entry) string {
	return "[➡ " + entry.Title + "](" + makeMarkdownTitleLink(entry.Title) + ")"
}

func (s SimplePageRender) RenderPageContent(page page.Page) string {
	return ""
}

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
