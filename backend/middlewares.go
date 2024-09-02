package main

import (
	"net/http"
	"os"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "development" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			// I don't understand CORS enough to know what this is doing halp
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.WriteHeader(http.StatusOK)
				return
			}
			// switch r.Method {
			// case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
			// 	w.Header().Set("Access-Control-Allow-Methods", r.Method)
			// case http.MethodOptions:
			// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			// 	w.WriteHeader(http.StatusOK)
			// 	return
			// }
		}
		next.ServeHTTP(w, r)
	})
}
