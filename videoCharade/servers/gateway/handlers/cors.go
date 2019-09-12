package handlers

import (
	"net/http"
)

const headerAuthorization = "Authorization"

//Cors struct that wraps around a mux
type Cors struct {
	handler http.Handler
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", headerAuthorization)
	w.Header().Set("Access-Control-Max-Age", "Access-Control-Max-Age: 600")
	if r.Method == "OPTIONS" {
		if r.Header.Get("Access-Control-Request-Method") == "" || r.Header.Get("Origin") == "" {
			http.Error(w, "Bad CORS Pre-flight request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	c.handler.ServeHTTP(w, r)
}

//NewCors constructs a new Cors struct that wraps around a mux
func NewCors(handlerToWrap http.Handler) *Cors {
	return &Cors{handlerToWrap}
}
