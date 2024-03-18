package middleware

import (
	"VK_HR/pkg/sessionrepo"
	"VK_HR/pkg/userrepo"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type IMiddleware interface {
	Auth(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	Logging(next http.Handler) http.Handler
	AddAuthHandler(path, method string)
}

type Middleware struct {
	SessionManager sessionrepo.SessionPacker
	UsersRepo      userrepo.UsersRepository
	AuthHandlers   map[Pair]struct{}
}

func NewMiddleware(packer sessionrepo.SessionPacker, repository userrepo.UsersRepository) *Middleware {
	return &Middleware{
		SessionManager: packer,
		UsersRepo:      repository,
		AuthHandlers:   make(map[Pair]struct{}),
	}
}

func (middleware *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := middleware.AuthHandlers[NewPair(r.URL.Path, r.Method)]; !ok {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		s, err := middleware.SessionManager.Unpack(token)
		if err != nil {
			middleware.jsonHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		role, err := middleware.UsersRepo.SelectRole(r.Context(), s.Sub)
		if err != nil {
			middleware.jsonHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user_role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (middleware *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//

		next.ServeHTTP(w, r)
	})
}

func (middleware *Middleware) AddAuthHandler(path, method string) {
	middleware.AuthHandlers[NewPair(path, method)] = struct{}{}
}

func (middleware *Middleware) jsonHTTPError(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf("{\"err\":\"%s\"}", err.Error()), status)
}
