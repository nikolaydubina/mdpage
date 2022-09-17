## 🍙 mdpage

This CLI tool generates one-page lists with summary based on YAML.
This tool provides:

* consistency of formatting for large Markdown pages
* automatic Markdown generation
* parametrization
* separation of data definition (YAML) from representation (Markdown)

## Tests

```go
go run ../main.go -page page.yaml -config render_config.json > README.md
```