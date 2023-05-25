package render

import "github.com/nikolaydubina/mdpage/page"

func EnrichedTitleWithProjectNameLinkTextSummary(entry page.Entry) string {
	if entry.Name == "" {
		return entry.Title
	}
	return entry.Title + " with `" + entry.Name + "`"
}

func EnrichedTitleWithProjectNameLink(entry page.Entry) string {
	if entry.Name == "" {
		return entry.Title
	}
	if entry.URL == "" {
		return entry.Title + " with `" + entry.Name + "`"
	}
	return entry.Title + " with " + Link{Name: entry.Name, URL: entry.URL}.Render()
}

func EnrichedTitleWithProjectNameLinkText(entry page.Entry) string {
	if entry.Name == "" {
		return entry.Title
	}
	return entry.Title + " with " + entry.Name
}
