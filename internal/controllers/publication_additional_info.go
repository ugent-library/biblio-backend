package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAdditionalInfo struct {
	Context
}

func NewPublicationAdditionalInfo(c Context) *PublicationAdditionalInfo {
	return &PublicationAdditionalInfo{c}
}

func (c *PublicationAdditionalInfo) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_show",
		views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAdditionalInfo) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_edit",
		views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			Form         *views.FormBuilder
			Vocabularies map[string][]string
		}{
			pub,
			views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationAdditionalInfo) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.Engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := forms.Decode(pub, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			"publication/additional_info/_edit",
			views.NewData(c.Render, r, struct {
				Publication  *models.Publication
				Form         *views.FormBuilder
				Vocabularies map[string][]string
			}{
				pub,
				views.NewFormBuilder(c.Render, locale.Get(r.Context()), formErrors),
				c.Engine.Vocabularies(),
			}),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK,
		"publication/additional_info/_update",
		views.NewData(c.Render, r, struct {
			Publication  *models.Publication
			Show         *views.ShowBuilder
			Vocabularies map[string][]string
		}{
			savedPub,
			views.NewShowBuilder(c.Render, locale.Get(r.Context())),
			c.Engine.Vocabularies(),
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}