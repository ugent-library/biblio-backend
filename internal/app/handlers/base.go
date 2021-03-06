package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/render"
	"github.com/ugent-library/biblio-backend/internal/render/flash"
)

func init() {
	// register flash.Flash as a gob Type to make SecureCookieStore happy
	// see https://github.com/gin-contrib/sessions/issues/39
	gob.Register(flash.Flash{})
}

// TODO handlers should only have access to a url builder,
// the session and maybe the localizer
type BaseHandler struct {
	Router       *mux.Router
	SessionName  string
	SessionStore sessions.Store
	UserService  backends.UserService
	Localizer    *locale.Localizer
}

func (h BaseHandler) Wrap(fn func(http.ResponseWriter, *http.Request, BaseContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := h.NewContext(r, w)
		if err != nil {
			render.InternalServerError(w, r, err)
			return
		}
		fn(w, r, ctx)
	}
}

func (h BaseHandler) NewContext(r *http.Request, w http.ResponseWriter) (BaseContext, error) {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		return BaseContext{}, err
	}

	user, err := h.getUserFromSession(session, r, UserSessionKey)
	if err != nil {
		return BaseContext{}, err
	}

	originalUser, err := h.getUserFromSession(session, r, OriginalUserSessionKey)
	if err != nil {
		return BaseContext{}, err
	}

	flash, err := h.getFlashFromSession(session, r, w)
	if err != nil {
		return BaseContext{}, err
	}

	return BaseContext{
		CurrentURL:   r.URL,
		Flash:        flash,
		Locale:       h.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         user,
		OriginalUser: originalUser,
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}, nil
}

func (h BaseHandler) AddSessionFlash(r *http.Request, w http.ResponseWriter, f flash.Flash) error {
	session, err := h.SessionStore.Get(r, h.SessionName)
	if err != nil {
		return err
	}

	session.AddFlash(f, FlashSessionKey)

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (h BaseHandler) getFlashFromSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) ([]flash.Flash, error) {
	sessionFlashes := session.Flashes(FlashSessionKey)

	if err := session.Save(r, w); err != nil {
		return []flash.Flash{}, err
	}

	flashes := []flash.Flash{}
	for _, f := range sessionFlashes {
		flashes = append(flashes, f.(flash.Flash))
	}

	return flashes, nil
}

func (h BaseHandler) getUserFromSession(session *sessions.Session, r *http.Request, sessionKey string) (*models.User, error) {
	userID := session.Values[sessionKey]
	if userID == nil {
		return nil, nil
	}

	user, err := h.UserService.GetUser(userID.(string))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h BaseHandler) PathFor(name string, vars ...string) *url.URL {
	if route := h.Router.Get(name); route != nil {
		u, err := route.URLPath(vars...)
		if err != nil {
			log.Panic(fmt.Errorf("can't reverse route %s: %w", name, err))
		}
		return u
	}
	log.Panic(fmt.Errorf("route %s not found", name))
	return nil

}

func (h BaseHandler) URLFor(name string, vars ...string) *url.URL {
	if route := h.Router.Get(name); route != nil {
		u, err := route.URL(vars...)
		if err != nil {
			log.Panic(fmt.Errorf("can't reverse route %s: %w", name, err))
		}
		return u
	}
	log.Panic(fmt.Errorf("route %s not found", name))
	return nil
}

type BaseContext struct {
	CurrentURL   *url.URL
	Flash        []flash.Flash
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
}
