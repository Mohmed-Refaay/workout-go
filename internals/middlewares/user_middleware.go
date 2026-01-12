package middlewares

import (
	"backend-go/internals/store"
	"backend-go/internals/tokens"
	"backend-go/internals/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey = contextKey("user")

type UserMiddleware struct {
	UserStore store.UserStore
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) (*store.User, error) {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		panic("failed to get the user")
	}
	return user, nil
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth_header := r.Header.Get("Authorization")

		if auth_header == "" {
			r := SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(auth_header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		user, err := um.UserStore.GetUserFromToken(headerParts[1], tokens.ScopeAuth)
		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
			return
		}
		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "token is invalid or expired"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUser(r)

		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
			return
		}
		if user.IsAnonymous() {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you have to be logged in"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
