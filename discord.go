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
	message += "📦 **Order Received**:\n"
	message += fmt.Sprintf("🔹 **Order ID**: `%s`\n", order.OrderID)

	// Ensure Created At is not empty
	if order.Created != "" {
		message += fmt.Sprintf("**Created At**: `%s`\n", order.Created)
	} else {
		message += "**Created At**: `Unknown`\n" // Fallback in case the date is empty
	}

	for i, product := range order.ProductsSold {
		// Add a horizontal line between products if there's more than one
		if i > 0 {
			message += "\n───────────────\n" // Separator for readability between products
		}

		// Format product details
		message += fmt.Sprintf("🔹 **Product Name**: *%s* (ID: `%d`), \n**Quantity**: `%d`\n", product.ProductName, product.ProductID, product.Quantity)

		// Add user data (if present)
		if len(product.UserData) > 0 {
			message += "🔹 **User Data**:\n"
			for key, value := range product.UserData {
				message += fmt.Sprintf("> • **%s**: `%s`\n", key, value) // Indented bullet points for user data
			}
		}

		// Add key IDs sold
		message += fmt.Sprintf("🔑 **Key IDs Sold**: `%v`\n", product.KeyIDsSold)
	}
	message += "\n───────────────\n"
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

	// Check for any error in response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message to Discord, status: %s", resp.Status)
	}

	return nil
}

func sendToDebugDiscord(order Order) error {
	// Load .env file for webhook URL
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webhookUrl := os.Getenv("DISCORD_DEBUG_WEBHOOK_URL")
	if webhookUrl == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL is not set in .env")
	}

	// Create the Discord message
	var message string
	message += "📦 **Order Received**:\n"
	message += fmt.Sprintf("**Order ID**: `%s`\n", order.OrderID)

	// Ensure Created At is not empty
	if order.Created != "" {
		message += fmt.Sprintf("**Created At**: `%s`\n", order.Created)
	} else {
		message += "**Created At**: `Unknown`\n" // Fallback in case the date is empty
	}

	for i, product := range order.ProductsSold {
		// Add a horizontal line between products if there's more than one
		if i > 0 {
			message += "\n───────────────\n" // Separator for readability between products
		}

		// Format product details
		message += fmt.Sprintf("🔹 **Product Name**: *%s* (ID: `%d`), \n**Quantity**: `%d`\n", product.ProductName, product.ProductID, product.Quantity)

		// Add user data (if present)
		if len(product.UserData) > 0 {
			message += "🔹 **User Data**:\n"
			for key, value := range product.UserData {
				message += fmt.Sprintf("> • **%s**: `%s`\n", key, value) // Indented bullet points for user data
			}
		}

		// Add key IDs sold
		message += fmt.Sprintf("🔑 **Key IDs Sold**: `%v`\n", product.KeyIDsSold)
	}
	message += "\n───────────────\n"
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

	// Check for any error in response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message to Discord, status: %s", resp.Status)
	}

	return nil
}
