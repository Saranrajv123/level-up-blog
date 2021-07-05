package middleware

import (
	"errors"
	"github.com/saranrajv123/levelupblog/api/auth"
	"github.com/saranrajv123/levelupblog/api/responses"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(write http.ResponseWriter, req *http.Request) {
		write.Header().Set("Content-Type", "application/json")
		next(write, req)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(write http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(write, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(write, r)
	}
}
