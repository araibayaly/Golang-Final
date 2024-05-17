package main

import (
	"net/http"

	"github.com/araibayaly/learnify/pkg/learnify/handlers"
	"github.com/araibayaly/learnify/pkg/learnify/middleware"
)

func (app *application) routes() {
    r := app.router

    h := &handlers.Application{DB: app.db}

    courseHandler := handlers.NewCourseHandler(app.db)
    lessonHandler := handlers.NewLessonHandler(app.db)
	authHandler := handlers.NewAuthHandler(app.db)

    r.NotFoundHandler = http.HandlerFunc(h.NotFoundResponse)
    r.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedResponse)

	r.HandleFunc("/api/v1/signup", authHandler.Signup).Methods("POST")
	r.HandleFunc("/api/v1/login", authHandler.Login).Methods("POST")

    r.HandleFunc("/api/v1/healthcheck", h.InfoHandler).Methods("GET")
    r.HandleFunc("/api/v1/dbtest", h.DbHandler).Methods("GET")

    // Secure routes with authentication middleware
    api := r.PathPrefix("/api/v1").Subrouter()
    api.Use(middleware.Authenticate)

    // Course routes
    api.HandleFunc("/courses", courseHandler.CreateCourse).Methods("POST")
	api.HandleFunc("/courses/{id}", courseHandler.GetCourse).Methods("GET")
	api.HandleFunc("/courses/{id}", courseHandler.UpdateCourse).Methods("PUT")
	api.HandleFunc("/courses/{id}", courseHandler.DeleteCourse).Methods("DELETE")
	api.HandleFunc("/courses", courseHandler.GetAllCourses).Methods("GET")

	// Lesson routes
	api.HandleFunc("/lessons", lessonHandler.CreateLesson).Methods("POST")
	api.HandleFunc("/lessons/{id}", lessonHandler.GetLesson).Methods("GET")
	api.HandleFunc("/lessons/{id}", lessonHandler.UpdateLesson).Methods("PUT")
	api.HandleFunc("/lessons/{id}", lessonHandler.DeleteLesson).Methods("DELETE")
	api.HandleFunc("/lessons", lessonHandler.GetAllLessons).Methods("GET")

}
