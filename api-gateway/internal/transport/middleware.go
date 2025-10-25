// Package transport is a nice package
package transport

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const (
	customerID ctxKey = "customerID"
)
 
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), customerID, uuid.NewString())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
