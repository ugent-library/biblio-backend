package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/internal/views"
	"github.com/ugent-library/go-web/forms"
	"github.com/unrolled/render"
)

type PublicationsDetails struct {
	engine *engine.Engine
	render *render.Render
}

func NewPublicationsDetails(e *engine.Engine, r *render.Render) *PublicationsDetails {
	return &PublicationsDetails{
		engine: e,
		render: r,
	}
}

func (c *PublicationsDetails) Show(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, 200,
		fmt.Sprintf("publication/details/_%s", pub.Type),
		views.NewPublicationData(r, c.render, pub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationsDetails) OpenForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	c.render.HTML(w, 200,
		fmt.Sprintf("publication/details/_%s_edit_form", pub.Type),
		views.NewPublicationForm(r, c.render, pub, nil),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *PublicationsDetails) SaveForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pub, err := c.engine.GetPublication(id)
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

	savedPub, formErrors := c.engine.UpdatePublication(pub)

	if formErrors != nil {
		c.render.HTML(w, 200,
			fmt.Sprintf("publication/details/_%s_edit_form", pub.Type),
			views.NewPublicationForm(r, c.render, pub, formErrors),
			render.HTMLOptions{Layout: "layouts/htmx"},
		)

		return
	}

	c.render.HTML(w, 200,
		fmt.Sprintf("publication/details/_%s_edit_submit", savedPub.Type),
		views.NewPublicationData(r, c.render, savedPub),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}
