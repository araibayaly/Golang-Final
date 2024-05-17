package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/araibayaly/learnify/pkg/learnify/models"
)

type CourseRepository struct {
    DB *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
    return &CourseRepository{DB: db}
}

func (r *CourseRepository) GetCourseByID(id int) (*models.Course, error) {
    course := &models.Course{}
    err := r.DB.QueryRow("SELECT id, title, description, teacher_id, created_at, updated_at FROM courses WHERE id = $1", id).Scan(
        &course.ID, &course.Title, &course.Description, &course.TeacherID, &course.CreatedAt, &course.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("course not found")
        }
        return nil, err
    }
    return course, nil
}

func (r *CourseRepository) CreateCourse(course *models.Course) error {
    err := r.DB.QueryRow(
        "INSERT INTO courses (title, description, teacher_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        course.Title, course.Description, course.TeacherID, course.CreatedAt, course.UpdatedAt).Scan(&course.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *CourseRepository) UpdateCourse(course *models.Course) error {
    _, err := r.DB.Exec(
        "UPDATE courses SET title = $1, description = $2, teacher_id = $3, updated_at = $4 WHERE id = $5",
        course.Title, course.Description, course.TeacherID, course.UpdatedAt, course.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *CourseRepository) DeleteCourse(id int) error {
    _, err := r.DB.Exec("DELETE FROM courses WHERE id = $1", id)
    if err != nil {
        return err
    }
    return nil
}

func (r *CourseRepository) GetAllCourses(filter map[string]string, sort string, page, limit int) ([]*models.Course, error) {
    query := "SELECT id, title, description, teacher_id, created_at, updated_at FROM courses"
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

    courses := []*models.Course{}
    for rows.Next() {
        course := &models.Course{}
        err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.TeacherID, &course.CreatedAt, &course.UpdatedAt)
        if err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }

    return courses, nil
}
