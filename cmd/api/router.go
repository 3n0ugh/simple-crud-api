package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) SetRouter() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/book", app.handleGetBook)
	router.HandlerFunc(http.MethodPost, "/v1/book", app.handleAddBook)
	router.HandlerFunc(http.MethodDelete, "/v1/book/:id", app.handleDeleteBook)
	router.HandlerFunc(http.MethodPut, "/v1/book/:id", app.handleUpdateBook)

	return router
}
