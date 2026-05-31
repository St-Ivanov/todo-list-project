package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// Authentication verification function
func handlerAuth(w http.ResponseWriter, r *http.Request) {
	var (
		req models.AuthRequest
		buf bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	password := os.Getenv("TODO_PASSWORD")

	if len(password) > 0 && password == req.Password {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(8 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := jwtToken.SignedString([]byte(password))
		if err != nil {
			writeJson(w, map[string]string{"error": "Invalid password"})
			return
		}

		writeJson(w, map[string]string{"token": signedToken})
	}
}

// User authorization verification function
func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if len(password) > 0 {
			var jwtCookie string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}

			token, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("the encryption method does not match")
				}
				return []byte(password), nil
			}, jwt.WithExpirationRequired())
			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
