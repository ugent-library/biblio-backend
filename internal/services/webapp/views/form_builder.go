package views

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type formData struct {
	Name            string
	values          []string
	Label           string
	ShowLabel       bool
	Tooltip         string
	Placeholder     string
	Required        bool
	Checked         bool
	Max             string
	Min             string
	Choices         []string
	ChoicesLabels   []string
	EmptyChoice     bool
	Cols            int
	Rows            int
	AutocompleteURL string
	Error           string
	errorPointer    string
	ID              string
	Disabled        bool
}

func (f *formData) Value() string {
	if len(f.values) > 0 {
		return f.values[0]
	}
	return ""
}

func (f *formData) Values() []string {
	return f.values
}

type formOption func(*formData)

type formLocaleOption func(string) string

type FormBuilder struct {
	renderPartial func(string, interface{}) (template.HTML, error)
	locale        *locale.Locale
	localeScope   string
	Errors        validation.Errors
}

func NewFormBuilder(r func(string, interface{}) (template.HTML, error), l *locale.Locale, e validation.Errors) *FormBuilder {
	return &FormBuilder{
		renderPartial: r,
		locale:        l,
		localeScope:   "builder",
		Errors:        e,
	}
}

func (b *FormBuilder) newFormData(opts []formOption) *formData {
	d := &formData{}

	for _, opt := range opts {
		opt(d)
	}

	if d.errorPointer == "" {
		d.errorPointer = "/" + strings.ReplaceAll(d.Name, ".", "/")
	}

	if formErr := b.ErrorFor(d.errorPointer); formErr != nil {
		d.Error = d.Label + " " + b.locale.TranslateScope("validation", formErr.Code)
	}

	return d
}

func (b *FormBuilder) ErrorFor(pointer string) *validation.Error {
	for _, err := range b.Errors {
		if err.Pointer == pointer {
			return err
		}
	}
	return nil
}

func (b *FormBuilder) Locale(scope string) formLocaleOption {
	return func(str string) string {
		return b.locale.TranslateScope(scope, str)
	}
}

func (b *FormBuilder) LanguageName() formLocaleOption {
	return func(str string) string {
		return b.locale.LanguageName(str)
	}
}

func (b *FormBuilder) Name(name string, localeOpts ...formLocaleOption) formOption {
	return func(d *formData) {
		d.Name = name

		if len(localeOpts) > 0 {
			opt := localeOpts[0]
			if lbl := opt(name); lbl != "" {
				d.Label = lbl
			}
		}
		if d.Label == "" {
			d.Label = b.locale.TranslateScope(b.localeScope, d.Name)
		}
	}
}

func (b *FormBuilder) Value(v string) formOption {
	return func(d *formData) {
		d.values = []string{v}
	}
}

func (b *FormBuilder) Values(v []string) formOption {
	return func(d *formData) {
		d.values = v
	}
}

func (b *FormBuilder) ShowLabel(v bool) formOption {
	return func(d *formData) {
		d.ShowLabel = v
	}
}

func (b *FormBuilder) Tooltip(v string) formOption {
	return func(d *formData) {
		d.Tooltip = v
	}
}

func (b *FormBuilder) Placeholder(v string) formOption {
	return func(d *formData) {
		d.Placeholder = v
	}
}

func (b *FormBuilder) Required() formOption {
	return func(d *formData) {
		d.Required = true
	}
}

func (b *FormBuilder) Checked(v bool) formOption {
	return func(d *formData) {
		d.Checked = v
	}
}

func (b *FormBuilder) Max(v string) formOption {
	return func(d *formData) {
		d.Max = v
	}
}

func (b *FormBuilder) Min(v string) formOption {
	return func(d *formData) {
		d.Min = v
	}
}

func (b *FormBuilder) ID(v string) formOption {
	return func(d *formData) {
		d.ID = v
	}
}

func (b *FormBuilder) Disabled(v bool) formOption {
	return func(d *formData) {
		d.Disabled = v
	}
}

func (b *FormBuilder) Choices(choices []string, localeOpts ...formLocaleOption) formOption {
	return func(d *formData) {
		d.Choices = make([]string, len(choices))
		d.ChoicesLabels = make([]string, len(choices))
		copy(d.Choices, choices)
		copy(d.ChoicesLabels, choices)

		if len(localeOpts) > 0 {
			opt := localeOpts[0]
			for i, c := range choices {
				if lbl := opt(c); lbl != "" {
					d.ChoicesLabels[i] = lbl
				}
			}
		}
	}
}

func (b *FormBuilder) EmptyChoice() formOption {
	return func(d *formData) {
		d.EmptyChoice = true
	}
}

func (b *FormBuilder) Cols(num int) formOption {
	return func(d *formData) {
		d.Cols = num
	}
}

func (b *FormBuilder) Rows(num int) formOption {
	return func(d *formData) {
		d.Rows = num
	}
}

func (b *FormBuilder) ErrorPointer(ptr string) formOption {
	return func(d *formData) {
		d.errorPointer = ptr
	}
}

func (b *FormBuilder) AutocompleteURL(v string) formOption {
	return func(d *formData) {
		d.AutocompleteURL = v
	}
}

func (b *FormBuilder) Render(tmpl string, opts ...formOption) (template.HTML, error) {
	return b.renderPartial(fmt.Sprintf("form_builder/%s", tmpl), b.newFormData(opts))
}

func (b *FormBuilder) Text(opts ...formOption) (template.HTML, error) {
	return b.Render("_text", opts...)
}

func (b *FormBuilder) TextRepeat(opts ...formOption) (template.HTML, error) {
	return b.Render("_text_repeat", opts...)
}

func (b *FormBuilder) TextArea(opts ...formOption) (template.HTML, error) {
	return b.Render("_text_area", opts...)
}

func (b *FormBuilder) Checkbox(opts ...formOption) (template.HTML, error) {
	return b.Render("_checkbox", opts...)
}

func (b *FormBuilder) List(opts ...formOption) (template.HTML, error) {
	return b.Render("_list", opts...)
}

func (b *FormBuilder) ListRepeat(opts ...formOption) (template.HTML, error) {
	return b.Render("_list_repeat", opts...)
}

func (b *FormBuilder) RadioButtonGroup(opts ...formOption) (template.HTML, error) {
	return b.Render("_radio_button_group", opts...)
}

func (b *FormBuilder) Date(opts ...formOption) (template.HTML, error) {
	return b.Render("_date", opts...)
}