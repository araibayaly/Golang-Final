package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/araibayaly/learnify/pkg/learnify/models"
)

type LessonRepository struct {
    DB *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
    return &LessonRepository{DB: db}
}

func (r *LessonRepository) GetLessonByID(id int) (*models.Lesson, error) {
    lesson := &models.Lesson{}
    err := r.DB.QueryRow("SELECT id, title, content, course_id, created_at, updated_at FROM lessons WHERE id = $1", id).Scan(
        &lesson.ID, &lesson.Title, &lesson.Content, &lesson.CourseID, &lesson.CreatedAt, &lesson.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("lesson not found")
        }
        return nil, err
    }
    return lesson, nil
}

func (r *LessonRepository) CreateLesson(lesson *models.Lesson) error {
    err := r.DB.QueryRow(
        "INSERT INTO lessons (title, content, course_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        lesson.Title, lesson.Content, lesson.CourseID, lesson.CreatedAt, lesson.UpdatedAt).Scan(&lesson.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *LessonRepository) UpdateLesson(lesson *models.Lesson) error {
    _, err := r.DB.Exec(
        "UPDATE lessons SET title = $1, content = $2, course_id = $3, updated_at = $4 WHERE id = $5",
        lesson.Title, lesson.Content, lesson.CourseID, lesson.UpdatedAt, lesson.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *LessonRepository) DeleteLesson(id int) error {
    _, err := r.DB.Exec("DELETE FROM lessons WHERE id = $1", id)
    if err != nil {
        return err
    }
    return nil
}

func (r *LessonRepository) GetAllLessons(filter map[string]string, sort string, page, limit int) ([]*models.Lesson, error) {
    query := "SELECT id, title, content, course_id, created_at, updated_at FROM lessons"
    var args []interface{}
    var conditions []string
    var placeholderCount int

    for key, value := range filter {
        placeholderCount++
        conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", key, placeholderCount))
        args = append(args, "%"+value+"%")
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    if sort != "" {
        query += " ORDER BY " + sort
    }

    placeholderCount++
    query += fmt.Sprintf(" LIMIT $%d", placeholderCount)
    args = append(args, limit)

    placeholderCount++
    query += fmt.Sprintf(" OFFSET $%d", placeholderCount)
    args = append(args, (page-1)*limit)

    rows, err := r.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    lessons := []*models.Lesson{}
    for rows.Next() {
        lesson := &models.Lesson{}
        err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.Content, &lesson.CourseID, &lesson.CreatedAt, &lesson.UpdatedAt)
        if err != nil {
            return nil, err
        }
        lessons = append(lessons, lesson)
    }

    return lessons, nil
}
