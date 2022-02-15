package controllers

import (
	"net/http"
	"time"

	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type DatasetDetails struct {
	Context
}

func NewDatasetDetails(c Context) *DatasetDetails {
	return &DatasetDetails{c}
}

func (c *DatasetDetails) Show(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/details/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		dataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) Edit(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	c.Render.HTML(w, http.StatusOK, "dataset/details/_edit", c.ViewData(r, struct {
		Dataset      *models.Dataset
		Show         *views.ShowBuilder
		Form         *views.FormBuilder
		Vocabularies map[string][]string
	}{
		dataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
		views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDetails) Update(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(dataset, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedDataset, err := c.Engine.UpdateDataset(dataset)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "dataset/details/_edit", c.ViewData(r, struct {
			Dataset      *models.Dataset
			Show         *views.ShowBuilder
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			dataset,
			views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
			views.NewFormBuilder(c.RenderPartial, locale.Get(r.Context()), formErrors),
			c.Engine.Vocabularies(),
		},
			views.Flash{Type: "error", Message: "There are some problems with your input"},
		),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/details/_update", c.ViewData(r, struct {
		Dataset *models.Dataset
		Show    *views.ShowBuilder
	}{
		savedDataset,
		views.NewShowBuilder(c.RenderPartial, locale.Get(r.Context())),
	},
		views.Flash{Type: "success", Message: "Details updated succesfully", DismissAfter: 5 * time.Second},
	),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}