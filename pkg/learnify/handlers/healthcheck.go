package handlers

import (
	"database/sql"
	"net/http"

	"github.com/araibayaly/learnify/pkg/learnify"
)

type Application struct {
	DB *sql.DB
}

func (app *Application) InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := learnify.Info()

	w.Header().Set("Content-Type", "application/json")

	_, err := w.Write([]byte(info))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
