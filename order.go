package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Product represents each product sold in the order
type Product struct {
	ProductID   int               `json:"product_id"`
	ProductName string            `json:"product_name"`
	Quantity    int               `json:"quantity"`
	UserData    map[string]string `json:"user_data"`
	KeyIDsSold  []int             `json:"key_ids_sold"`
}

// Order represents the incoming order notification
type Order struct {
	OrderID      string    `json:"order_id"`
	Created      string    `json:"created"`
	ProductsSold []Product `json:"products_sold"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Set a limit for the maximum allowed body size (20MB)
	if r.Method == "CONNECT" {
		http.Error(w, "CONNECT method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body2, err2 := io.ReadAll(r.Body)
	if err2 != nil || len(body2) == 0 {
		log.Printf("Error reading body: %v", err2)
		http.Error(w, "can't read body", http.StatusBadRequest)
		handleEmptyRequest(w)
		return
	}
	// Log the request
	log.Printf("Request: %s %s", r.Method, r.URL)
	log.Printf("Request body: %s", body2)

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB limit

	// Parse the JSON body into an Order struct
	var order Order
	err := json.Unmarshal(body2, &order)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Handle empty or missing fields in the order
	if order.OrderID == "" {
		order.OrderID = "Unknown Order ID"
	}
	if order.Created == "" {
		order.Created = getCurrentTimeInGMT7()
	}
	for i, product := range order.ProductsSold {
		if product.ProductName == "" {
			order.ProductsSold[i].ProductName = "New Product"
		}
		if product.Quantity == 0 {
			order.ProductsSold[i].Quantity = 1 // Default to 1 if quantity is missing
		}
	}

	// Send the order details to Discord
	err = sendToDiscord(order)
	if err != nil {
		http.Error(w, "Failed to send to Discord", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	_, err2 = fmt.Fprint(w, "Order notification received and sent to Discord\n")
	if err2 != nil {
		return
	}
}

func handleEmptyRequest(w http.ResponseWriter) {
	// Prepare a default order if the request body is empty
	defaultOrder := Order{
		OrderID: "Unknown Order ID",
		Created: getCurrentTimeInGMT7(),
		ProductsSold: []Product{
			{
				ProductID:   0,
				ProductName: "New Product",
				Quantity:    1,
				UserData:    make(map[string]string), // Empty map for user data
				KeyIDsSold:  []int{},
			},
		},
	}

	// Send the default order to Discord
	err := sendToDiscord(defaultOrder)
	if err != nil {
		http.Error(w, "Failed to send default order to Discord", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, "Empty request handled with default order data and sent to Discord\n")
	if err != nil {
		return
	}
}
