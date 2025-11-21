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
	"server/internal/bank"
)

// main initializes and starts the HTTP server with graceful shutdown support
// Setup steps:
//   1. Create a new Bank instance (initializes with default accounts)
//   2. Initialize chi router and apply middleware
//   3. Register all API routes
//   4. Start HTTP server with graceful shutdown on SIGINT/SIGTERM
//
// Graceful Shutdown:
//   - Listens for SIGINT (Ctrl+C) or SIGTERM signals
//   - Waits up to 10 seconds for in-flight requests to complete
//   - Forcefully closes remaining connections after timeout
func main() {
	// Create a new Bank instance with default accounts (IDs: "1", "2")
	b := bank.New()
	
	// Create a new chi router for handling HTTP requests
	r := chi.NewRouter()
	
	// Apply global middleware
	// StripSlashes: automatically redirects /path/ to /path (removes trailing slashes)
	r.Use(chimiddleware.StripSlashes)
	
	// Register all account-related routes
	// Routes registered: GET/POST /account, POST /account/deposit, POST /account/withdraw
	handler.Routes(r, b)

	// Configure the HTTP server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    15 * time.Second,  // Timeout for reading request body
		WriteTimeout:   15 * time.Second,  // Timeout for writing response
		IdleTimeout:    60 * time.Second,  // Timeout for idle connections
		MaxHeaderBytes: 1 << 20,           // Max header size: 1MB
	}

	// Create a channel to receive OS signals (interrupt, terminate)
	// This allows us to detect when the user wants to shutdown the server
	sigChan := make(chan os.Signal, 1)
	
	// Register for SIGINT (Ctrl+C) and SIGTERM signals
	// SIGINT: sent when user presses Ctrl+C
	// SIGTERM: sent by process managers (docker, systemd, etc.)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a separate goroutine to allow graceful shutdown handling
	// This prevents ListenAndServe from blocking the main goroutine
	go func() {
		log.Println("API server started on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// ErrServerClosed is expected during graceful shutdown, so we ignore it
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	// This blocks until SIGINT or SIGTERM is received
	sig := <-sigChan
	log.Printf("\nReceived signal: %v, starting graceful shutdown...", sig)

	// Create a context with 10-second timeout for graceful shutdown
	// This gives in-flight requests 10 seconds to complete
	// After 10 seconds, the server will forcefully close connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	// Stops accepting new connections and waits for existing ones to finish
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server gracefully shut down")
}

