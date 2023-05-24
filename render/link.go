package render

import "io"

type Link struct {
	Name string
	URL  string
}

func (v Link) RenderTo(out io.StringWriter) {
	out.WriteString("[" + v.Name + "]")
	out.WriteString("(" + v.URL + ")")
}
