package middleware

import (
	"net/http"
)

// AddCors to a request
func AddCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			w.WriteHeader(200)
			w.Write([]byte(""))
			return
		}

		h.ServeHTTP(w, r)
	})
}
