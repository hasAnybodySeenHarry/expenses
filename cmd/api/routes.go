package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)
	router.NotFound = http.HandlerFunc(app.notFound)

	// meta
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerHandler)

	return router
}
