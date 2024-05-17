package repository

import (
	"database/sql"
	"errors"

	"github.com/araibayaly/learnify/pkg/learnify/models"
)

// UserRepository handles the CRUD operations for the User model.
type UserRepository struct {
    DB *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    err := r.DB.QueryRow("SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE email = $1", email).Scan(
        &user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return user, nil
}

// CreateUser creates a new user.
func (r *UserRepository) CreateUser(user *models.User) error {
    err := r.DB.QueryRow(
        "INSERT INTO users (first_name, last_name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
        user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
    if err != nil {
        return err
    }
    return nil
}
