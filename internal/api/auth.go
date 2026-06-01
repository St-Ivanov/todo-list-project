package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	Pass = ""
)

// Authentication verification function
func handlerAuth(w http.ResponseWriter, r *http.Request) {
	var (
		req models.AuthRequest
		buf bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if len(Pass) > 0 && Pass == req.Password {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(8 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := jwtToken.SignedString([]byte(Pass))
		if err != nil {
			writeJson(w, map[string]string{"error": "Invalid password"}, http.StatusUnauthorized)
			return
		}

		writeJson(w, map[string]string{"token": signedToken}, http.StatusOK)
	}
}

// User authorization verification function
func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(Pass) > 0 {
			var jwtCookie string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}

			token, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("the encryption method does not match")
				}
				return []byte(Pass), nil
			}, jwt.WithExpirationRequired())
			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
