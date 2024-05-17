package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	config *Config
	router *mux.Router
	db     *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := LoadConfig()
	db, err := openDB(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: config,
		router: mux.NewRouter(),
		db:     db,
	}
	
	
	app.routes()

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", app.router)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
