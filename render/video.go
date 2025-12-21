package render

import (
	"io"
	"regexp"
)

var youtubeRegex = regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtube\.com/embed/|youtu\.be/)([a-zA-Z0-9_-]{11})`)

type Video struct {
	VideoURL string
}

func YouTubeVideoID(url string) string {
	matches := youtubeRegex.FindStringSubmatch(url)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (v Video) RenderTo(out io.StringWriter) {
	videoID := YouTubeVideoID(v.VideoURL)
	if videoID == "" {
		return
	}

	out.WriteString(`<div align="center">`)
	out.WriteString(`<a href="`)
	out.WriteString(v.VideoURL)
	out.WriteString(`"><img src="https://img.youtube.com/vi/`)
	out.WriteString(videoID)
	out.WriteString(`/0.jpg" style="margin: 8px; max-height: 640px;"></a>`)
	out.WriteString(`</div>`)
	out.WriteString("\n")
}
