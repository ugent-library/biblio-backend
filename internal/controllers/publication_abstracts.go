package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/context"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-locale/locale"
	"github.com/ugent-library/go-web/forms"
	"github.com/ugent-library/go-web/jsonapi"
	"github.com/unrolled/render"
)

type PublicationAbstracts struct {
	Context
}

func NewPublicationAbstracts(c Context) *PublicationAbstracts {
	return &PublicationAbstracts{c}
}

// Show the "Add abstract" modal
func (c *PublicationAbstracts) AddAbstract(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	abstract := &models.Text{}

	c.Render.HTML(w, http.StatusOK, "publication/abstracts/_form", views.NewData(c.Render, r, struct {
		PublicationID string
		Abstract      *models.Text
		Form          *views.FormBuilder
		Vocabularies  map[string][]string
	}{
		id,
		abstract,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Save an abstract to Librecat
func (c *PublicationAbstracts) CreateAbstract(w http.ResponseWriter, r *http.Request) {
	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts, *abstract)
	pub.Abstract = abstracts

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK, "publication/abstracts/_form", views.NewData(c.Render, r, struct {
			PublicationID string
			Abstract      *models.Text
			Form          *views.FormBuilder
			Vocabularies  map[string][]string
		}{
			savedPub.ID,
			abstract,
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

	w.Header().Set("HX-Trigger", "PublicationCreateAbstract")

	c.Render.HTML(w, http.StatusOK,
		"publication/abstracts/_table_body",
		views.NewData(c.Render, r, struct {
			Publication *models.Publication
		}{
			savedPub,
		}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// Show the "Edit abstract" modal
func (c *PublicationAbstracts) EditAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	abstract := &pub.Abstract[rowDelta]

	c.Render.HTML(w, http.StatusOK, "publication/abstracts/_form_edit", views.NewData(c.Render, r, struct {
		PublicationID string
		Delta         string
		Abstract      *models.Text
		Form          *views.FormBuilder
		Vocabularies  map[string][]string
	}{
		pub.ID,
		muxRowDelta,
		abstract,
		views.NewFormBuilder(c.Render, locale.Get(r.Context()), nil),
		c.Engine.Vocabularies(),
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Save the updated abstract to Librecat
func (c *PublicationAbstracts) UpdateAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstract := &models.Text{}

	if err := forms.Decode(abstract, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts[rowDelta] = *abstract
	pub.Abstract = abstracts

	savedPub, err := c.Engine.UpdatePublication(pub)

	if formErrors, ok := err.(jsonapi.Errors); ok {
		c.Render.HTML(w, http.StatusOK,
			"publication/abstracts/_form_edit",
			views.NewData(c.Render, r, struct {
				PublicationID string
				Delta         string
				Abstract      *models.Text
				Form          *views.FormBuilder
				Vocabularies  map[string][]string
			}{
				savedPub.ID,
				strconv.Itoa(rowDelta),
				abstract,
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

	w.Header().Set("HX-Trigger", "PublicationUpdateAbstract")

	c.Render.HTML(w, http.StatusOK, "publication/abstracts/_table_body", views.NewData(c.Render, r, struct {
		Publication *models.Publication
	}{
		savedPub,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Show the "Confirm remove" modal
func (c *PublicationAbstracts) ConfirmRemoveFromPublication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	muxRowDelta := mux.Vars(r)["delta"]

	c.Render.HTML(w, http.StatusOK, "publication/abstracts/_modal_confirm_removal", views.NewData(c.Render, r, struct {
		ID  string
		Key string
	}{
		id,
		muxRowDelta,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

// // Remove an abstract from Librecat
func (c *PublicationAbstracts) RemoveAbstract(w http.ResponseWriter, r *http.Request) {
	muxRowDelta := mux.Vars(r)["delta"]
	rowDelta, _ := strconv.Atoi(muxRowDelta)

	pub := context.GetPublication(r.Context())

	abstracts := make([]models.Text, len(pub.Abstract))
	copy(abstracts, pub.Abstract)

	abstracts = append(abstracts[:rowDelta], abstracts[rowDelta+1:]...)
	pub.Abstract = abstracts

	// TODO: error handling
	c.Engine.UpdatePublication(pub)

	w.Header().Set("HX-Trigger", "PublicationRemoveAbstract")

	// Empty content, denotes we deleted the record
	fmt.Fprintf(w, "")
}
