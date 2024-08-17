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

	// transactions
	router.HandlerFunc(http.MethodPost, "/v1/transactions", app.mustAuth(app.createTransactionHandler))

	// debts
	router.HandlerFunc(http.MethodGet, "/v1/debts", app.mustAuth(app.showDebtsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/debts", app.mustAuth(app.createDebtHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/debts/:id", app.mustAuth(app.deleteDebtHandler))

	// auth
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.loginHandler)

	return router
}
