package handlers

import (
	"encoding/json"
	"net/http"
)

func (app *Application) DbHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	err := app.DB.QueryRow("SELECT 'Hello, world!'").Scan(&result)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": result}
	json.NewEncoder(w).Encode(response)
}