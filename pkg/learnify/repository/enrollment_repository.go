package repository

import (
	"database/sql"
	"errors"

	"github.com/araibayaly/learnify/pkg/learnify/models"
)

type EnrollmentRepository struct {
    DB *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
    return &EnrollmentRepository{DB: db}
}

func (r *EnrollmentRepository) GetEnrollmentByID(id int64) (*models.Enrollment, error) {
    enrollment := &models.Enrollment{}
    err := r.DB.QueryRow("SELECT id, student_id, course_id, enrolled_at FROM enrollments WHERE id = $1", id).Scan(
        &enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.EnrolledAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("enrollment not found")
        }
        return nil, err
    }
    return enrollment, nil
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment *models.Enrollment) error {
    err := r.DB.QueryRow(
        "INSERT INTO enrollments (student_id, course_id, enrolled_at) VALUES ($1, $2, $3) RETURNING id",
        enrollment.StudentID, enrollment.CourseID, enrollment.EnrolledAt).Scan(&enrollment.ID)
    if err != nil {
        return err
    }
    return nil
}
