package render

import (
	"io"
	"log"
	"os"
	"strings"
)

// getFileContent unmodified bytes of file as is
func getFileContent(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("can not open file: %s", err)
	}
	s, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("can not read: %s", err)
	}
	return string(s)
}

// renderContentFromFilePath for local files
func renderContentFromFilePath(filePath string) string {
	var b strings.Builder

	b.WriteString("```")
	// markdown syntax
	if strings.HasSuffix(filePath, ".go") {
		b.WriteString("go")
	}
	b.WriteRune('\n')

	b.WriteString(getFileContent(filePath))

	b.WriteRune('\n')
	b.WriteString("```")
	b.WriteRune('\n')

	return b.String()
}
