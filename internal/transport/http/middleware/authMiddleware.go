package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// Protects a route with authentication
func AuthMiddleware(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Auth endpoint is requested!")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := authHeaderParts[1]

		if validateToken(token) {
			original(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func validateToken(accessToken string) bool {
	log.Printf("Validating access token %s", accessToken)
	// Replace this by loading in a private RSA cert for more security
	var mySigningKey = []byte(os.Getenv("AUTH_SECRET"))
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error parsing access token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	return token.Valid
}
