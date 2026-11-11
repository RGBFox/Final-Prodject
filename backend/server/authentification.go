package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 8 * time.Hour

// jwtSecret строит ключ подписи из пароля: при смене TODO_PASSWORD
// все ранее выданные токены перестают действовать.
// Сам пароль или его хеш в токен не попадают.
func jwtSecret(pass string) []byte {
	sum := sha256.Sum256([]byte("todo-jwt-secret:" + pass))
	return sum[:]
}

// signinHandler принимает {"password": "..."} и возвращает {"token": "..."}.
func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("метод не поддерживается"))
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("аутентификация не включена"))
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, fmt.Errorf("ошибка формата данных"))
		return
	}

	// Сравнение за постоянное время, чтобы не светить длину/префикс пароля.
	// Обычное сравнение останавливается на первой ошибке. Если настоящий пароль abcdef, то попытка zzzzzz отвергается мгновенно, а abcxxx чуть дольше, потому что три буквы совпали
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(pass)) != 1 {
		writeJSON(w, http.StatusUnauthorized, fmt.Errorf("Неверный пароль"))
		return
	}

	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret(pass))
	if err != nil {
		log.Printf("ошибка созданяя токена: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Errorf("ошибка создания токена"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// auth пропускает запрос дальше, только если пароль не задан
// или в куке token лежит корректный неистёкший JWT.
func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var valid bool
			if cookie, err := r.Cookie("token"); err == nil {
				tok, err := jwt.ParseWithClaims(
					cookie.Value,
					&jwt.RegisteredClaims{},
					func(t *jwt.Token) (any, error) { return jwtSecret(pass), nil },
					jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
					jwt.WithExpirationRequired(),
				)
				valid = err == nil && tok.Valid
			}

			if !valid {
				writeJSON(w, http.StatusUnauthorized, fmt.Errorf("Authentification required"))
				return
			}
		}
		next(w, r)
	})
}
