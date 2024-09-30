package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func sendToDiscord(order Order) error {
	// Load .env file for webhook URL
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webhookUrl := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookUrl == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL is not set in .env")
	}

	// Create the Discord message
	var message string
	message += fmt.Sprintf("📦 **Order Received**:\n")
	message += fmt.Sprintf("Order ID: %s\n", order.OrderID)
	message += fmt.Sprintf("Created At: %s\n", order.Created)

	for _, product := range order.ProductsSold {
		message += fmt.Sprintf("🔹 **Product Name**: %s (ID: %d), Quantity: %d\n", product.ProductName, product.ProductID, product.Quantity)
		message += "User Data:\n"
		for key, value := range product.UserData {
			message += fmt.Sprintf("- %s: %s\n", key, value)
		}
		message += fmt.Sprintf("🔑 **Key IDs Sold**: %v\n", product.KeyIDsSold)
	}

	// Prepare the payload for Discord
	discordPayload := map[string]string{
		"content": message,
	}
	payloadBytes, err := json.Marshal(discordPayload)
	if err != nil {
		return err
	}

	// Send the request to Discord webhook
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message to Discord, status: %s", resp.Status)
	}

	return nil
}
