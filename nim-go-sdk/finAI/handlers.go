// handlers.go - HTTP request handlers
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// setupHTTPHandlers registers all HTTP endpoints
func setupHTTPHandlers() {
	http.HandleFunc("/api/alerts", handleAlerts)
	http.HandleFunc("/api/transactions", handleTransactions)
}

// handleAlerts returns recent alerts from the last 24 hours
func handleAlerts(w http.ResponseWriter, r *http.Request) {
	// Enable CORS - restrictive origin for production
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:8080" {
		if origin == "" {
			origin = "http://localhost:5173"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	alertsMutex.RLock()
	defer alertsMutex.RUnlock()

	// Return alerts from last 24 hours
	recentAlerts := filterRecentAlerts(alerts, AlertRetentionHours)
	if err := json.NewEncoder(w).Encode(recentAlerts); err != nil {
		http.Error(w, "Failed to encode alerts", http.StatusInternalServerError)
		return
	}
}

// handleTransactions returns all transactions from the mock data file
func handleTransactions(w http.ResponseWriter, r *http.Request) {
	// Enable CORS - restrictive origin for production
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:8080" {
		if origin == "" {
			origin = "http://localhost:5173"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read and parse transactions
	transactions, err := readMockTransactions()
	if err != nil {
		http.Error(w, "Failed to read transactions", http.StatusInternalServerError)
		return
	}

	// Convert to API format
	formattedTransactions := formatTransactionsForAPI(transactions)
	if err := json.NewEncoder(w).Encode(formattedTransactions); err != nil {
		http.Error(w, "Failed to encode transactions", http.StatusInternalServerError)
		return
	}
}

// filterRecentAlerts returns alerts from the last N hours
func filterRecentAlerts(allAlerts []Alert, hours int) []Alert {
	recentAlerts := []Alert{}
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)

	for _, alert := range allAlerts {
		if alert.Timestamp.After(cutoff) {
			recentAlerts = append(recentAlerts, alert)
		}
	}

	return recentAlerts
}

// formatTransactionsForAPI converts transactions to API response format
func formatTransactionsForAPI(transactions []Transaction) []TransactionAPI {
	formatted := []TransactionAPI{}

	for i, tx := range transactions {
		// Parse amount (remove $ and convert to float)
		amountStr := strings.TrimPrefix(tx.Amount, "$")
		amount := 0.0
		if _, err := fmt.Sscanf(amountStr, "%f", &amount); err != nil {
			// Log error but continue with 0.0 as default
			fmt.Printf("Warning: Failed to parse amount '%s': %v\n", amountStr, err)
		}

		formatted = append(formatted, TransactionAPI{
			ID:          fmt.Sprintf("tx-%d", i+1),
			Amount:      -amount, // Negative for debits
			Description: tx.Product,
			Date:        tx.Date,
			Type:        "debit",
			Merchant:    tx.Merchant,
		})
	}

	return formatted
}
