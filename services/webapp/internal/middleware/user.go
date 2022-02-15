package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

func SetUser(e *engine.Engine, sessionName string, sessionStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionStore.Get(r, sessionName)
			if err != nil {
				// TODO
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userID := session.Values["user_id"]
			if userID == nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := e.GetUser(userID.(string))
			if err != nil {
				// TODO
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			originalUserID := session.Values["original_user_id"]
			if originalUserID != nil {
				originalUser, err := e.GetUser(originalUserID.(string))
				if err != nil {
					log.Printf("get user error: %s", err)
					// TODO
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				c := context.WithOriginalUser(r.Context(), originalUser)
				r = r.WithContext(c)
			}

			c := context.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireUser(redirectURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if context.GetUser(r.Context()) == nil {
				http.Redirect(w, r, redirectURL, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}