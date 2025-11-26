// Package main is the entry point for the bank API server
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"server/internal/handler"
	"server/internal/storage"
)

// main initializes and starts the HTTP server with graceful shutdown support
func main() {
	// Initialize SQLite database
	db, err := storage.New("bank.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repositories
	accountRepo := storage.NewAccountRepository(db)
	userRepo := storage.NewUserRepository(db)
	tokenRepo := storage.NewTokenRepository(db)

	// Create a new chi router for handling HTTP requests
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(chimiddleware.StripSlashes)

	// Register all routes with repositories
	handler.Routes(r, accountRepo, userRepo, tokenRepo)

	// Configure the HTTP server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a separate goroutine
	go func() {
		log.Println("API server started on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("\nReceived signal: %v, starting graceful shutdown...", sig)

	// Create a context with 10-second timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server gracefully shut down")
}
