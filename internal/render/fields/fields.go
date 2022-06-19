package fields

import (
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/render"
)

type Fields struct {
	Theme  string
	Fields []Field
}

type Field interface {
	RenderHTML(*Fields, io.Writer) error
}

func (f *Fields) RenderHTML() (template.HTML, error) {
	var buf strings.Builder

	for _, field := range f.Fields {
		if err := field.RenderHTML(f, &buf); err != nil {
			return "", err
		}
	}

	return template.HTML(buf.String()), nil
}

type Text struct {
	Label         string
	Required      bool
	Tooltip       string
	Values        []string
	ValueTemplate string
}

type yieldHTML struct {
	Label    string
	Required bool
	Tooltip  string
	Values   []template.HTML
}

func (f *Text) RenderHTML(fields *Fields, w io.Writer) error {
	if f.ValueTemplate != "" {
		return f.renderWithValueTemplate(fields, w)
	}

	return render.Templates().ExecuteTemplate(w, path.Join("fields", fields.Theme, "text"), f)
}

func (f *Text) renderWithValueTemplate(fields *Fields, w io.Writer) error {
	y := yieldHTML{
		Label:    f.Label,
		Required: f.Required,
		Tooltip:  f.Tooltip,
	}

	for _, v := range f.Values {
		var buf strings.Builder
		if err := render.Templates().ExecuteTemplate(&buf, f.ValueTemplate, v); err != nil {
			return err
		}
		y.Values = append(y.Values, template.HTML(buf.String()))
	}

	return render.Templates().ExecuteTemplate(w, path.Join("fields", fields.Theme, "text"), y)
}