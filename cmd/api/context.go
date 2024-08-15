package main

import (
	"context"
	"net/http"

	"harry2an.com/expenses/internal/data"
)

type contextKey string

const userCtx = contextKey("user")

func (app *application) setUser(r *http.Request, u *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtx, u)
	return r.WithContext(ctx)
}

func (app *application) getUser(r *http.Request) *data.User {
	u, ok := r.Context().Value(userCtx).(*data.User)
	if !ok {
		panic("invalid user pointer inside the context")
	}

	return u
}
