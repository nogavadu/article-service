package middlewares

import (
	"context"
	"github.com/nogavadu/articles-service/internal/lib/api/request"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
)

const authTokenKey = "authorization"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.GetAuthToken(r)
		if err != nil {
			response.Err(w, r, "invalid auth token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authTokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
