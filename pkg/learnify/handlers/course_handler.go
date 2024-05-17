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

type CourseHandler struct {
    courseRepo *repository.CourseRepository
}

func NewCourseHandler(db *sql.DB) *CourseHandler {
    return &CourseHandler{
        courseRepo: repository.NewCourseRepository(db),
    }
}

func (h *CourseHandler) validateCourseInput(course *models.Course) error {
    if course.Title == "" {
        return fmt.Errorf("title is required")
    }
    if course.Description == "" {
        return fmt.Errorf("description is required")
    }
    return nil
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    var course models.Course
    err := json.NewDecoder(r.Body).Decode(&course)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.validateCourseInput(&course); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userID, err := strconv.Atoi(r.Header.Get("UserID"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusInternalServerError)
        return
    }

    course.TeacherID = userID
    course.CreatedAt = time.Now()
    course.UpdatedAt = time.Now()

    err = h.courseRepo.CreateCourse(&course)
    if err != nil {
        http.Error(w, "Error creating course: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid course ID", http.StatusBadRequest)
        return
    }

    course, err := h.courseRepo.GetCourseByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Course not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error retrieving course: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid course ID", http.StatusBadRequest)
        return
    }

    var course models.Course
    err = json.NewDecoder(r.Body).Decode(&course)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.validateCourseInput(&course); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    course.ID = id
    course.UpdatedAt = time.Now()

    err = h.courseRepo.UpdateCourse(&course)
    if err != nil {
        http.Error(w, "Error updating course: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid course ID", http.StatusBadRequest)
        return
    }

    err = h.courseRepo.DeleteCourse(id)
    if err != nil {
        http.Error(w, "Error deleting course: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
}

func (h *CourseHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
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

    courses, err := h.courseRepo.GetAllCourses(filter, sort, page, limit)
    if err != nil {
        http.Error(w, "Error fetching courses: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(courses)
}
