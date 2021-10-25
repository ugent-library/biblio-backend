package views

import "github.com/ugent-library/go-web/jsonapi"

type textData struct {
	Text     string
	Label    string
	Required bool
	Tooltip  string
}

type listData struct {
	List     []string
	Label    string
	Required bool
	Tooltip  string
}

type textFormData struct {
	Key      string
	Text     string
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type listFormValues struct {
	Key      string
	Value    string
	Selected bool
}

type listFormData struct {
	Key      string
	Values   []*listFormValues
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}

type boolData struct {
	Value    bool
	Label    string
	Required bool
	Tooltip  string
}

type CheckboxInput struct {
	Name     string
	Value    string
	Checked  bool
	Label    string
	Required bool
	Tooltip  string
	Cols     int
	HasError bool
	Error    jsonapi.Error
}
