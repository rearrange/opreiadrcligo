// Package template provides the embedded Markdown templates used when
// generating ADR files and the workspace index.
//
// Templates are stored as .tmpl files alongside this file and embedded into
// the binary at compile time using go:embed — no runtime file-system access
// is required to locate them.
package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
	"time"
)

//go:embed adr.md.tmpl
var adrTemplate string

//go:embed index.md.tmpl
var indexTemplate string

// ADRData contains the values interpolated into adr.md.tmpl.
type ADRData struct {
	Number int
	Title  string
	Date   string
	Author string
	Status string
}

// IndexData contains the values interpolated into index.md.tmpl.
type IndexData struct {
	Date string
}

// RenderADR executes the ADR template with the supplied data and returns
// the rendered Markdown string.
func RenderADR(data ADRData) (string, error) {
	return render("adr.md.tmpl", adrTemplate, data)
}

// RenderIndex executes the index template with today's date and returns
// the rendered Markdown string.
func RenderIndex() (string, error) {
	return render("index.md.tmpl", indexTemplate, IndexData{
		Date: time.Now().Format("2006-01-02"),
	})
}

// render is a shared helper that parses and executes a named template.
func render(name, tmplStr string, data any) (string, error) {
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %q: %w", name, err)
	}

	return buf.String(), nil
}
