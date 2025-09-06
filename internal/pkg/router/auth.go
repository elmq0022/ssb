package router

import (
	"net/http"
)

type AuthFunc func(request *http.Request) (bool, error)

func WithAuth(handler JSONHandler, auth AuthFunc) JSONHandler {
	return func(request *http.Request) (any, int, error) {
		match, err := auth(request)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if !match {
			return nil, http.StatusUnauthorized, nil
		}
		return handler(request)
	}
}
