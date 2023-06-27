package api

import "net/http"

type corsMiddleware struct {
	handler http.Handler
}

func (m *corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PATCH, PUT, HEAD, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	m.handler.ServeHTTP(w, r)
}

// CorsHandler is for gorrilla middleware
func CorsHandler() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &corsMiddleware{
			handler: h,
		}
	}
}
