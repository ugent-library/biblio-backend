package publicationediting

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/localize"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/snapstore"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

type BindLink struct {
	Position    int    `path:"position"`
	URL         string `form:"url"`
	Relation    string `form:"relation"`
	Description string `form:"description"`
}

type BindDeleteLink struct {
	Position int `path:"position"`
}

type YieldLinks struct {
	Context
}
type YieldAddLink struct {
	Context
	Form *form.Form
}
type YieldEditLink struct {
	Context
	Position int
	Form     *form.Form
}
type YieldDeleteLink struct {
	Context
	Position int
}

func (h *Handler) AddLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	form := linkForm(ctx, BindLink{
		Position: len(ctx.Publication.Link),
	}, nil)

	render.Layout(w, "show_modal", "publication/add_link", YieldAddLink{
		Context: ctx,
		Form:    form,
	})
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{Position: len(ctx.Publication.Link)}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	ctx.Publication.Link = append(
		ctx.Publication.Link,
		models.PublicationLink{
			URL:         b.URL,
			Relation:    b.Relation,
			Description: b.Description,
		},
	)

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		render.Layout(w, "refresh_modal", "publication/add_link", YieldAddLink{
			Context: ctx,
			Form:    linkForm(ctx, b, validationErrs.(validation.Errors)),
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) EditLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	link, err := ctx.Publication.GetLink(b.Position)
	if err != nil {
		render.BadRequest(w, r, err)
		return
	}

	b.URL = link.URL
	b.Description = link.Description
	b.Relation = link.Relation

	render.Layout(w, "show_modal", "publication/edit_link", YieldEditLink{
		Context:  ctx,
		Position: b.Position,
		Form:     linkForm(ctx, b, nil),
	})
}

func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	b := BindLink{}
	if err := bind.Request(r, &b, bind.Vacuum); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	link := models.PublicationLink{
		URL:         b.URL,
		Description: b.Description,
		Relation:    b.Relation,
	}

	if err := ctx.Publication.SetLink(b.Position, link); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if validationErrs := ctx.Publication.Validate(); validationErrs != nil {
		form := linkForm(ctx, b, validationErrs.(validation.Errors))

		render.Layout(w, "refresh_modal", "publication/edit_link", YieldEditLink{
			Context:  ctx,
			Position: b.Position,
			Form:     form,
		})
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func (h *Handler) ConfirmDeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	render.Layout(w, "show_modal", "publication/confirm_delete_link", YieldDeleteLink{
		Context:  ctx,
		Position: b.Position,
	})
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request, ctx Context) {
	var b BindDeleteLink
	if err := bind.Request(r, &b); err != nil {
		render.BadRequest(w, r, err)
		return
	}

	if err := ctx.Publication.RemoveLink(b.Position); err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	err := h.Repository.UpdatePublication(r.Header.Get("If-Match"), ctx.Publication)

	var conflict *snapstore.Conflict
	if errors.As(err, &conflict) {
		render.Layout(w, "refresh_modal", "error_dialog", ctx.Locale.T("publication.conflict_error"))
		return
	}

	if err != nil {
		render.InternalServerError(w, r, err)
		return
	}

	render.View(w, "publication/refresh_links", YieldLinks{
		Context: ctx,
	})
}

func linkForm(ctx Context, b BindLink, errors validation.Errors) *form.Form {
	l := ctx.Locale
	return form.New().
		WithTheme("default").
		WithErrors(localize.ValidationErrors(ctx.Locale, errors)).
		AddSection(
			&form.Text{
				Name:  "url",
				Value: b.URL,
				Label: l.T("builder.link.url"),
				Cols:  12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/url", b.Position),
				),
			},
			&form.Select{
				Name:    "relation",
				Value:   b.Relation,
				Label:   l.T("builder.link.relation"),
				Options: localize.VocabularySelectOptions(l, "publication_link_relations"),
				Cols:    12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/relation", b.Position),
				),
			},
			&form.Text{
				Name:  "description",
				Value: b.Description,
				Label: l.T("builder.link.description"),
				Cols:  12,
				Error: localize.ValidationErrorAt(
					l,
					errors,
					fmt.Sprintf("/link/%d/description", b.Position),
				),
			},
		)
}
