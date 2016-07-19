package web

import (
	"net/http"
	"net/http/httptest"

	"github.com/NYTimes/gziphandler"
)

// MiddlewareFunc defines a function that returns a middleware. A Middleware can call next to continue the chain, or not if it is not appropriate
type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

// MiddlewareChain is a list of middlewares to be applied. The first element will be called first on a request.
type MiddlewareChain []MiddlewareFunc

func (c MiddlewareChain) Extend(middlewares ...MiddlewareFunc) MiddlewareChain {
	newC := make(MiddlewareChain, 0, len(c)+len(middlewares))
	for _, m := range c {
		newC = append(newC, m)
	}
	for _, m := range middlewares {
		newC = append(newC, m)
	}
	return newC
}

// Build creates a single function that can be called to create concrete handlers from a middleware chain
func (c MiddlewareChain) Build() func(http.HandlerFunc) http.HandlerFunc {
	return func(root http.HandlerFunc) http.HandlerFunc {
		chain := root
		for i := len(c) - 1; i >= 0; i-- {
			chain = c[i](chain)
		}
		return chain
	}
}

func gZip(next http.HandlerFunc) http.HandlerFunc {
	return gziphandler.GzipHandler(next).ServeHTTP
}

func readAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func writeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func jsonp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		next(rec, r)
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
	}
}
