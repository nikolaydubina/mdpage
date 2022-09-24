package render

import "strings"

const GitHubBaseURL = "https://github.com/"

// GitHubAuthorRenderer detects if this is GitHub user URL and creates link with at tag for it.
// Useful to shorten text len.
type GitHubAuthorRenderer struct {
	Prefix string // e.g. "@"
}

func (r GitHubAuthorRenderer) Render(s string) string {
	if !strings.HasPrefix(s, GitHubBaseURL) {
		return s
	}
	if len(s) <= len(GitHubBaseURL) {
		return s
	}
	h := strings.TrimSpace(s[len(GitHubBaseURL):])

	// not simple user link
	if strings.Contains(h, "/") {
		return s
	}

	return "[" + r.Prefix + h + "](" + s + ")"
}
