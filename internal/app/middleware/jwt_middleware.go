package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ShortenerClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type CookieKey string

const SecretKey string = "supersecretkey"
const UserIDKey CookieKey = "UserId"

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Token")
		if errors.Is(http.ErrNoCookie, err) {
			newUUID := uuid.NewString()
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ShortenerClaims{
				RegisteredClaims: jwt.RegisteredClaims{},
				UserID:           newUUID,
			})
			strJwt, err := newToken.SignedString([]byte(SecretKey))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "Token",
				Value: strJwt,
			})
			if r.RequestURI == "/api/user/urls" {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, newUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims := &ShortenerClaims{}

		_, err = jwt.ParseWithClaims(token.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
