package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/araibayaly/learnify/pkg/learnify/models"
	"github.com/araibayaly/learnify/pkg/learnify/repository"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    UserID int `json:"user_id"`
    jwt.StandardClaims
}

type AuthHandler struct {
    userRepo *repository.UserRepository
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
    return &AuthHandler{
        userRepo: repository.NewUserRepository(db),
    }
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var user models.User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Log received password for debugging (remove in production)
    log.Printf("Received password for signup: %s", user.Password)

    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    err = h.userRepo.CreateUser(&user)
    if err != nil {
        http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if creds.Email == "" || creds.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    log.Printf("Login attempt for email: %s", creds.Email)

    user, err := h.userRepo.GetUserByEmail(creds.Email)
    if err != nil {
        log.Printf("Error retrieving user by email: %v", err)
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    log.Printf("Stored password for %s: %s", user.Email, user.Password)
    log.Printf("Password provided: %s", creds.Password)

    // Compare the provided password with the stored plain text password
    if user.Password != creds.Password {
        log.Printf("Password comparison failed")
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: int(user.ID),
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error logging in", http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "message": "Login successful",
        "token":   tokenString,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
