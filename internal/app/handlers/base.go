package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
)

type Base struct {
	SessionName  string
	SessionStore sessions.Store
	Localizer    *locale.Localizer
}

type BaseContext struct {
	Locale       *locale.Locale
	User         *models.User
	OriginalUser *models.User
	CSRFToken    string
	CSRFTag      template.HTML
}

func NewBaseContext(base Base, r *http.Request) BaseContext {
	return BaseContext{
		Locale:       base.Localizer.GetLocale(r.Header.Get("Accept-Language")),
		User:         context.GetUser(r.Context()),
		OriginalUser: context.GetOriginalUser(r.Context()),
		CSRFToken:    csrf.Token(r),
		CSRFTag:      csrf.TemplateField(r),
	}
}

func (c BaseContext) T(key string, args ...any) string {
	return c.Locale.T(key, args...)
}

func (c BaseContext) TS(scope, key string, args ...any) string {
	return c.Locale.TS(scope, key, args...)
}
