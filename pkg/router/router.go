package router

import (
	"VK_HR/pkg/middleware"
	"net/http"
)

type Pair struct {
	path   string
	method string
}

func NewPair(path, method string) Pair {
	return Pair{
		path:   path,
		method: method,
	}
}

type Router struct {
	handlersMap map[Pair]http.HandlerFunc
	middle      middleware.IMiddleware
}

func NewRouter(middle middleware.IMiddleware) *Router {
	return &Router{
		handlersMap: make(map[Pair]http.HandlerFunc),
		middle:      middle,
	}
}

func (router *Router) HandleFunc(pair Pair, handlerFunc http.HandlerFunc) {
	router.handlersMap[pair] = handlerFunc
}

func (router *Router) HandleFuncWithAuth(p Pair, handlerFunc http.HandlerFunc) {
	router.handlersMap[p] = handlerFunc
	router.middle.AddAuthHandler(p.path, p.method)
}

func (router *Router) PackInMiddleware() http.Handler {
	return router.middle.Auth(
		router.middle.Logging(
			router.middle.RecoverPanic(router),
		),
	)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pair := NewPair(r.URL.Path, r.Method)

	handler, exist := router.handlersMap[pair]
	if !exist {
		http.Error(w, newNoHandlerError(pair).Error(), http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}
