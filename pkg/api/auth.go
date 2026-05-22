package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("change-me")

func hashPassword(password string) string {
	result := sha256.Sum256([]byte(password))
	return hex.EncodeToString(result[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "wrong method"}, http.StatusMethodNotAllowed)
		return
	}

	var buf bytes.Buffer
	var input struct {
		Password string `json:"password"`
	}

	_, err := buf.ReadFrom(r.Body)

	if err = json.Unmarshal(buf.Bytes(), &input); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == input.Password {
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"password_hash": hashPassword(envPassword),
			"exp":           time.Now().Add(8 * time.Hour).Unix(),
		})

		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]any{"token": signedToken}, http.StatusOK)
		return
	}

	writeJSON(w, map[string]string{"error": "wrong password"}, http.StatusUnauthorized)

}

func IsValidToken(tokenStr string) bool {
	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) == 0 {
		return true
	}

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil

	})

	if err != nil || !parsedToken.Valid {
		return false
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		currentHash := hashPassword(pass)
		if tokenHash, exists := claims["password_hash"]; exists {
			return tokenHash == currentHash
		}

		return false

	}
	return false
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenStr string

			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
			}

			if !IsValidToken(tokenStr) {

				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
