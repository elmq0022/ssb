package router

import (
	"context"
	"net/http"
)

type ctxKey string

const userKey ctxKey = "username"

type AuthFunc func(request *http.Request) (string, error)

func WithAuth(handler JSONHandler, auth AuthFunc) JSONHandler {
	return func(request *http.Request) (any, int, error) {
		username, err := auth(request)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		ctx := context.WithValue(request.Context(), userKey, username)
		request = request.WithContext(ctx)
		return handler(request)
	}
}
