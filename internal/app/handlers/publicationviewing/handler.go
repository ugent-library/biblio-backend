package publicationviewing

import (
	"net/http"

	"github.com/ugent-library/biblio-backend/internal/app/handlers"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/bind"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
)

type Handler struct {
	handlers.BaseHandler
	Repository backends.Repository
}

type Context struct {
	handlers.BaseContext
	Publication *models.Publication
	RedirectURL string
}

func (h *Handler) Wrap(fn func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return h.BaseHandler.Wrap(func(w http.ResponseWriter, r *http.Request, ctx handlers.BaseContext) {
		if ctx.User == nil {
			render.Unauthorized(w, r)
			return
		}

		p, err := h.Repository.GetPublication(bind.PathValues(r).Get("id"))
		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}

		if !ctx.User.CanViewPublication(p) {
			render.Forbidden(w, r)
			return
		}

		redirectURL := r.URL.Query().Get("redirect-url")
		if redirectURL == "" {
			redirectURL = h.PathFor("publications").String()
		}

		fn(w, r, Context{
			BaseContext: ctx,
			Publication: p,
			RedirectURL: redirectURL,
		})
	})
}
