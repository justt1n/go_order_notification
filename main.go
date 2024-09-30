package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define a route to handle POST requests
	http.HandleFunc("/api/order", handlePost)
	http.HandleFunc("/", checkStatus)
	// Start the HTTP server
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
