package middleware

import (
	"VK_HR/pkg/sessionrepo"
	"VK_HR/pkg/userrepo"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("checking authorize for query...")

		if _, ok := middleware.AuthHandlers[NewPair(r.URL.Path, r.Method)]; !ok {
			log.Info("no need authorize")
			next.ServeHTTP(w, r)
			return
		}

		log.Info("need authorize")

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
				log.WithFields(log.Fields{
					"path":   r.URL.Path,
					"method": r.Method,
					"err":    err,
				}).Error("query called panic")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (middleware *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
			"host":   r.Host,
		}).Info("handling query...")

		now := time.Now()
		next.ServeHTTP(w, r)
		dur := time.Since(now)

		log.Infof("took %d milliseconds", dur.Milliseconds())
	})
}

func (middleware *Middleware) AddAuthHandler(path, method string) {
	middleware.AuthHandlers[NewPair(path, method)] = struct{}{}
}

func (middleware *Middleware) jsonHTTPError(w http.ResponseWriter, err error, status int) {
	log.WithFields(log.Fields{
		"err":    err.Error(),
		"status": status,
	}).Error("http error")

	http.Error(w, fmt.Sprintf("{\"err\":\"%s\"}", err.Error()), status)
}
