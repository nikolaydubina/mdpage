package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/nikolaydubina/mdpage/page"
	"github.com/nikolaydubina/mdpage/render"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	var (
		pageFilePath         string
		renderConfigFilePath string
	)
	flag.StringVar(&pageFilePath, "page", "", "path to page file")
	flag.StringVar(&renderConfigFilePath, "config", "", "path to render config file")
	flag.Parse()

	if pageFilePath == "" {
		log.Fatalf("page filepath is missing")
	}
	if renderConfigFilePath == "" {
		log.Fatalf("render config filepath is missing")
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

	renderConfigFile, err := os.Open(renderConfigFilePath)
	if err != nil {
		log.Fatalf("can not open file: %s", err)
	}

	var renderConfig render.SimplePageRenderConfig
	if err := json.NewDecoder(renderConfigFile).Decode(&renderConfig); err != nil {
		log.Fatalf("can not parse render config: %s", err)
	}

	render := render.NewSimplePageRender(renderConfig)
	os.Stdout.WriteString(render.RenderPage(page))
}
