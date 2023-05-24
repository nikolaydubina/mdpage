package render

import (
	"io"
)

type Image struct {
	ImageURL string
}

func (v Image) RenderTo(out io.StringWriter) {
	out.WriteString(`<div align="center">`)
	out.WriteString(`<img src="`)
	out.WriteString(v.ImageURL)
	out.WriteString(`" style="margin: 8px; max-height: 640px;">`)
	out.WriteString(`</div>`)
	out.WriteString("\n")
}
