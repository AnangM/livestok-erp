package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type contextkey string

const RequestIdKey contextkey = "request_id"

func ReqeustId() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// generate request id
			// add request id to context

			id, err := uuid.NewV7()

			if err != nil {
				// non blocking error, log instead
				log.Printf("[RequestId] failed to generate request id: %v", err)
			}

			r = r.WithContext(context.WithValue(r.Context(), RequestIdKey, id))
			next.ServeHTTP(w, r)
		})
	}
}
