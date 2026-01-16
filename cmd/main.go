package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LordCodex164/httpserver/internal/handlers"
	"github.com/LordCodex164/httpserver/internal/middleware"
)

func main() {
	//create a multiplexer 

	mux := http.NewServeMux()

	//register handlers 
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/api/v1/users", handlers.Users)

	handler := middleware.Logger(mux)

	server := http.Server{
		Addr: ":8080",
		Handler: handler,
	}
	//start server 
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello, World!")
}