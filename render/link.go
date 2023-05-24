package render

import (
	"io"
	"strings"
	"unicode"
)

type Link struct {
	Name string
	URL  string
}

func (v Link) RenderTo(out io.StringWriter) {
	out.WriteString("[" + v.Name + "]")
	out.WriteString("(" + v.URL + ")")
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
