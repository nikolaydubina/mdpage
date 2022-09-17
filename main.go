package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"github.com/nikolaydubina/mdpage/page"
	"github.com/nikolaydubina/mdpage/render"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	var (
		pageFilePath string
	)
	flag.StringVar(&pageFilePath, "page", "", "path to page file")
	flag.Parse()

	if pageFilePath == "" {
		log.Fatalf("filepath is missing")
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

	render := render.NewSimplePageRender("## Content")

	os.Stdout.WriteString(render.RenderPage(page))
}
