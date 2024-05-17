package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    UserID int64 `json:"user_id"`
    jwt.StandardClaims
}

func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        if tokenStr == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        r.Header.Set("UserID", fmt.Sprintf("%d", claims.UserID))
        next.ServeHTTP(w, r)
    })
}
