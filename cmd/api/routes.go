package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/films", app.requirePermission("films:read", app.listFilmsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/films", app.requirePermission("films:write", app.createFilmHandler))
	router.HandlerFunc(http.MethodGet, "/v1/films/:id", app.requirePermission("films:read", app.showFilmHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/films/:id", app.requirePermission("films:write", app.updateFilmHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/films/:id", app.requirePermission("films:write", app.deleteFilmHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}