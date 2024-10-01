package main

import (
	"fmt"
	"log"
	"net/http"
)

// Main function to start the server
func main() {
	// Create a new ServeMux (router)
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/api/order", handlePost)
	mux.HandleFunc("/", checkStatus)

	// Wrap the router with logging middleware
	loggedMux := LoggingMiddleware(mux)

	// Start the HTTP server
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", loggedMux))
}

// Health check handler
func checkStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running"))
}

// LoggingMiddleware to log incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the HTTP method and the requested path
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		log.Printf("Request body: %s", r.Body)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
