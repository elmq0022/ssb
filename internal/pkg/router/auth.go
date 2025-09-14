package router

import (
	"context"
	"net/http"
	"ssb/internal/models"
	"ssb/internal/pkg/auth"
	"ssb/internal/repo"
)

type ctxKey string

const userKey ctxKey = "user"

type AuthFunc func(request *http.Request) (string, error)

// TODO: Move to the JWT package?
func NewJWTAuthFunction(jwtConfig *auth.JWTConfig) AuthFunc {
	return func(request *http.Request) (string, error) {
		token := request.Header.Get("Bearer")
		claims, err := jwtConfig.IsValidToken(token)
		if err != nil {
			return "", err
		}
		return claims.Subject, nil
	}
}

// TODO: pass user repository and check that user exists
// i.e. do full auth here.
func WithAuth(handler JSONHandler, auth AuthFunc, ur repo.UserRepository) JSONHandler {
	return func(request *http.Request) (any, int, error) {
		username, err := auth(request)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		user, err := ur.GetByUserName(username)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		ctx := context.WithValue(request.Context(), userKey, user)
		request = request.WithContext(ctx)
		return handler(request)
	}
}

func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userKey).(models.User)
	return user, ok
}
