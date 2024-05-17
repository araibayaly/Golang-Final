package handlers

import (
	"net/http"
)

func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}