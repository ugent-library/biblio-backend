package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/engine"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
)

func SetDataset(e *engine.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dataset, err := e.StorageService.GetDataset(mux.Vars(r)["id"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			c := context.WithDataset(r.Context(), dataset)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireCanViewDataset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetDataset(c)
		user := context.GetUser(c)

		if user.CanViewDataset(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanEditDataset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetDataset(c)
		user := context.GetUser(c)

		if user.CanEditDataset(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanPublishDataset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetDataset(c)
		user := context.GetUser(c)

		if user.CanPublishDataset(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func RequireCanDeleteDataset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		pub := context.GetDataset(c)
		user := context.GetUser(c)

		if user.CanDeleteDataset(pub) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}
