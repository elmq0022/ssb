package router

import (
	"net/http"
)

type AuthFunc func(request *http.Request) error

func WithAuth(handler JSONHandler, auth AuthFunc) JSONHandler {
	return func(request *http.Request) (any, int, error) {
		if err := auth(request); err != nil {
			return nil, http.StatusUnauthorized, err
		}
		return handler(request)
	}
}
