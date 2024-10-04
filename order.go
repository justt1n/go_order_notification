package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	log.Println("body: ", string(body))
	if err != nil || len(body) == 0 {
		// If the body is empty, respond with a default product and current time
		handleEmptyRequest(w)
		return
	}

	// Parse the JSON body into an Order struct
	var order Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Handle empty or missing fields in the order
	if order.OrderID == "" {
		order.OrderID = "Unknown Order ID"
	}
	// Set the created time to the current time if it's missing and convert to GMT+7 if not null
	if order.Created == "" {
		order.Created = getCurrentTimeInGMT7()
	} else {
		order.Created, _ = convertToGMT7(order.Created)
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
	_, err = fmt.Fprint(w, "Order notification received and sent to Discord\n")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
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

func convertToGMT7(timeStr string) (string, error) {
	// Parse the time string
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", err
	}
	t = t.In(location)
	return t.Format(time.RFC3339), nil
}
