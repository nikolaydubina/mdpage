package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/nikolaydubina/mdpage/page"
	"github.com/nikolaydubina/mdpage/render"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	var (
		pageFilePath   string
		outputFilePath string
	)

	flag.StringVar(&pageFilePath, "page", "", "path to page file")
	flag.StringVar(&outputFilePath, "o", "", "path to output")
	flag.Parse()

	if pageFilePath == "" {
		log.Fatalf("page filepath is missing")
	}

	pageFile, err := os.Open(pageFilePath)
	if err != nil {
		log.Fatalf("can not open file: %s", err)
	}

	var page page.Page
	if err := yaml.NewDecoder(pageFile).Decode(&page); err != nil {
		log.Fatalf("can not parse yaml: %s", err)
	}

	if err := page.Validate(); err != nil {
		log.Fatalf("invalid page: %s", err)
	}

	render := render.NewSimplePageRender(
		render.GitHubAuthorRenderer{Prefix: "@"},
	)

	var out io.StringWriter = os.Stdout

	if outputFilePath != "" {
		out, err = os.Create(outputFilePath)
		if err != nil {
			log.Fatalf("can not open file(%s): %s", outputFilePath, err)
		}
	}

	out.WriteString(render.RenderPage(page))
}
