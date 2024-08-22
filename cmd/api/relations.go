package main

import "net/http"

func (app *application) getRelations(w http.ResponseWriter, r *http.Request) {
	u := app.getUser(r)

	users, err := app.models.Users.GetAll(u.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
