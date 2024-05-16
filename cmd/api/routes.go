package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/films", app.createFilmHandler)
	router.HandlerFunc(http.MethodGet, "/v1/films/:id", app.showFilmHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/films/:id", app.updateFilmHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/films/:id", app.deleteFilmHandler)

	return router
}