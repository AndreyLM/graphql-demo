package auth

import (
	"context"
	"net/http"

	graphql_demo "github.com/andreylm/graphql-demo"
)

// Middleware - authentication middleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			// AUTH LOGIC
			ctx = context.WithValue(ctx, graphql_demo.UserIDCtxKey, auth)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
