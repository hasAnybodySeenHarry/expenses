package main

import (
	"errors"
	"net/http"
	"time"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Password != "", "password", "must be provided")

	if !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.invalidCredentials(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	ok, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		app.invalidCredentials(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, data.ScopeAuthentication, 24*time.Hour)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"token": token, "user": user}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
