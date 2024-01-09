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
	UserId string
}

type CookieKey string

const SECRET_KEY = "supersecretkey"
const USER_ID_KEY CookieKey = "UserId"

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Token")
		if errors.Is(http.ErrNoCookie, err) {
			newUUID := uuid.NewString()
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ShortenerClaims{
				RegisteredClaims: jwt.RegisteredClaims{},
				UserId:           newUUID,
			})
			strJwt, err := newToken.SignedString([]byte(SECRET_KEY))
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
			ctx := context.WithValue(r.Context(), USER_ID_KEY, newUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims := &ShortenerClaims{}

		_, err = jwt.ParseWithClaims(token.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), USER_ID_KEY, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
