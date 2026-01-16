package middleware

import (
	"log"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func () {
			//it allows the program to manage the behaviour of a panicking goroutine
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)

				//Return 500 to client
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}