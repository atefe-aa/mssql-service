package templates

import (
	"embed"
	"html/template"
	"io"
	"log"
)

//go:embed *.html
var templateFS embed.FS

type Templates struct {
	templates *template.Template
}

func New() (*Templates, error) {
	tpl, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		return nil, err
	}
	return &Templates{templates: tpl}, nil
}

// NewFromDir loads templates from a directory (for development)
func NewFromDir(dir string) (*Templates, error) {
	tpl, err := template.ParseGlob(dir + "/*.html")
	if err != nil {
		return nil, err
	}
	return &Templates{templates: tpl}, nil
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("Template error: %v", err)
	}
}