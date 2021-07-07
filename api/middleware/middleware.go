package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saranrajv123/levelupblog/api/auth"
	"github.com/saranrajv123/levelupblog/api/responses"
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

func TokenAuthMiddleware() gin.HandlerFunc {
	errList := map[string]string{}

	return func(ctx *gin.Context) {
		err := auth.TokenValid(ctx.Request)
		fmt.Println("err ", err)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  errList,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}
