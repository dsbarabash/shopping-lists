package rest

import (
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func UserIdentity(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request)) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success": false, "error": "Unauthorized"}`))
		return
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Cfg.Secret), nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success": false, "error": "Invalid JWT token"}`))
		return
	}
	handler(w, r)
}
