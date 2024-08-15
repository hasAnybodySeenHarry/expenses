package main

import (
	"errors"
	"net/http"
	"strings"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

const AuthorizationHeader = "Authorization"

func (app *application) mustAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", AuthorizationHeader)

		h := r.Header.Get(AuthorizationHeader)
		if h == "" {
			app.invalidAuthToken(w, r)
			return
		}

		segs := strings.Split(h, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			app.invalidAuthToken(w, r)
			return
		}

		token := segs[1]
		v := validator.New()
		if data.ValidateToken(v, token); !v.Validate() {
			app.failedValidation(w, r, v.Errors)
			return
		}

		user, err := app.models.Users.GetForToken(token, data.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecord):
				app.invalidAuthToken(w, r)
			default:
				app.serverError(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, app.setUser(r, user))
	}
}
