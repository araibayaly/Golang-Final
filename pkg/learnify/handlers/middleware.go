package handlers

import (
	"net/http"
)

func (app *Application) Authenticate(handler http.Handler) http.Handler {
	return handler
}