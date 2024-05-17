package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/araibayaly/learnify/pkg/learnify/models"
	"github.com/araibayaly/learnify/pkg/learnify/repository"
	"github.com/gorilla/mux"
)

type LessonHandler struct {
    lessonRepo *repository.LessonRepository
}

func NewLessonHandler(db *sql.DB) *LessonHandler {
    return &LessonHandler{
        lessonRepo: repository.NewLessonRepository(db),
    }
}

func (h *LessonHandler) validateLessonInput(lesson *models.Lesson) error {
    if lesson.Title == "" {
        return fmt.Errorf("title is required")
    }
    if lesson.Content == "" {
        return fmt.Errorf("content is required")
    }
    if lesson.CourseID == 0 {
        return fmt.Errorf("course_id is required")
    }
    return nil
}

func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
    var lesson models.Lesson
    err := json.NewDecoder(r.Body).Decode(&lesson)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.validateLessonInput(&lesson); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    lesson.CreatedAt = time.Now()
    lesson.UpdatedAt = time.Now()

    err = h.lessonRepo.CreateLesson(&lesson)
    if err != nil {
        http.Error(w, "Error creating lesson: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(lesson)
}

func (h *LessonHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
        return
    }

    lesson, err := h.lessonRepo.GetLessonByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Lesson not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error retrieving lesson: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(lesson)
}

func (h *LessonHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
        return
    }

    var lesson models.Lesson
    err = json.NewDecoder(r.Body).Decode(&lesson)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.validateLessonInput(&lesson); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    lesson.ID = id
    lesson.UpdatedAt = time.Now()

    err = h.lessonRepo.UpdateLesson(&lesson)
    if err != nil {
        http.Error(w, "Error updating lesson: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(lesson)
}

func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
        return
    }

    err = h.lessonRepo.DeleteLesson(id)
    if err != nil {
        http.Error(w, "Error deleting lesson: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Lesson deleted successfully"})
}

func (h *LessonHandler) GetAllLessons(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    filter := make(map[string]string)
    for key := range query {
        if key != "sort" && key != "page" && key != "limit" {
            filter[key] = query.Get(key)
        }
    }

    sort := query.Get("sort")
    page, err := strconv.Atoi(query.Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(query.Get("limit"))
    if err != nil || limit < 1 {
        limit = 10
    }

    lessons, err := h.lessonRepo.GetAllLessons(filter, sort, page, limit)
    if err != nil {
        http.Error(w, "Error fetching lessons: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(lessons)
}
