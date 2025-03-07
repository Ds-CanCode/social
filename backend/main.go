package main

import (
	"fmt"
	"log"
	"net/http"

	db "funcs/internal/database"
	mux "funcs/rootes"
)

func main() {
	err := db.SetupDb()
	if err != nil {
		log.Println(err)
		return
	}
	config := http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(mux.Rootes()),
	}
	fmt.Println("Server started on http://localhost:8080")
	log.Println(config.ListenAndServe())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigin := "http://localhost:3000" // ✅ Explicitly allow frontend domain

		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin) // ✅ Only allow specific origin
			w.Header().Set("Access-Control-Allow-Credentials", "true")   // ✅ Required for cookies
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
