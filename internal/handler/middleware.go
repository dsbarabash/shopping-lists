package handler

import (
	"fmt"
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
	splitToken := strings.Split(authHeader, "Bearer")
	if len(splitToken) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success": false, "error": "Unauthorized"}`))
		return
	}
	tokenString := splitToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil || token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"success": false, "error": "Invalid JWT token"}`))
		return
	}
	handler(w, r)

}
