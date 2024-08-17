package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if input.Password == "" {
		app.badRequest(w, r, errors.New("password cannot be blank"))
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Create(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequest(w, r, errors.New("user with this email already exists in the system"))
		default:
			app.serverError(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, data.ScopeActivation, 15*time.Minute)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.background(func() {
		err := app.mailer.Send("user-welcome", user.Name, user.Email, token.Plaintext)
		if err != nil {
			app.logger.Println(fmt.Errorf("%s", err))
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"user": user,
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
