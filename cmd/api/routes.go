package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"harry2an.com/expenses/cmd/proto/users"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)
	router.NotFound = http.HandlerFunc(app.notFound)

	// meta
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheck)

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.register)

	// relations
	router.HandlerFunc(http.MethodGet, "/v1/users", app.mustAuth(app.getRelations))

	// transactions
	router.HandlerFunc(http.MethodPost, "/v1/transactions", app.mustAuth(app.createTransaction))
	router.HandlerFunc(http.MethodGet, "/v1/transactions/:id", app.mustAuth(app.showTransactions))

	// debts
	router.HandlerFunc(http.MethodGet, "/v1/debts", app.mustAuth(app.showDebts))
	router.HandlerFunc(http.MethodPost, "/v1/debts", app.mustAuth(app.createDebt))
	router.HandlerFunc(http.MethodDelete, "/v1/debts/:id", app.mustAuth(app.deleteDebt))

	// auth
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.login)

	return app.enableCORS(router)
}

func (app *application) registerGRPCservers(server *grpc.Server) {
	userService := &userServiceServer{models: &app.models}
	users.RegisterUserServiceServer(server, userService)
}
