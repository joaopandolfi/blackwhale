package handlers

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type GeneralHandle func(http.Handler) http.Handler

// Chain applies middlewares to a http.HandlerFunc
// @handler
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
