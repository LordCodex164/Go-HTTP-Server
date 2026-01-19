package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/LordCodex164/httpserver/internal/handlers"
	"github.com/LordCodex164/httpserver/internal/middleware"
	"github.com/LordCodex164/httpserver/internal/logger"
	"github.com/LordCodex164/httpserver/internal/config"
	"golang.org/x/time/rate"
	"github.com/lpernett/godotenv"
)

func main() {
	//create a multiplexer 

	mux := http.NewServeMux()

	structuredLogger := logger.New()

	cfg := config.Load()


	//register handlers 
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/api/v1/users", handlers.Users)
	mux.HandleFunc("/metrics", handlers.Metrics)
	mux.HandleFunc("/panic", handlers.Panic)
	mux.HandleFunc("/slow", handlers.Slow) //this handler takes 10 seconds

	rateLimiter := middleware.NewRateLimiter(rate.Limit(10), 20)

	go rateLimiter.Cleanup() // Run cleanup in background

	//building the middleware chain recovery => request id => logger => handler
	handler := middleware.Recovery(
		middleware.RequestID((middleware.Timeout(5 * time.Second)(middleware.Logger(mux)))),
	)

	server := http.Server{
		Addr: cfg.Server.Addr(),
		Handler: handler,
		ReadTimeout: cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout: cfg.Server.IdleTimeout,
	}

	//channel to list for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	//start server in a goroutine
	go func() {
    //start server
	structuredLogger.Info("Server Starting", map[string]interface{}{
		"addr": cfg.Server.Addr(),
		"request_timeout": "5",
		"note": "please press ctrl + c to stop",
	})
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello, World!")

	}() 

	//wait for interrupt signal
	<-stop
	structuredLogger.Info("\n Shutting Down gracefully", map[string]interface{}{})

	//create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // maximum shutdown timeout window of 30secs

	defer cancel()

	//attempt graceful shutdown

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server stopped gracefully")
}