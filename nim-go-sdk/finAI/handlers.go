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
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	alertsMutex.RLock()
	defer alertsMutex.RUnlock()

	// Return alerts from last 24 hours
	recentAlerts := filterRecentAlerts(alerts, AlertRetentionHours)
	json.NewEncoder(w).Encode(recentAlerts)
}

// handleTransactions returns all transactions from the mock data file
func handleTransactions(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	// Read and parse transactions
	transactions, err := readMockTransactions()
	if err != nil {
		http.Error(w, "Failed to read transactions", http.StatusInternalServerError)
		return
	}

	// Convert to API format
	formattedTransactions := formatTransactionsForAPI(transactions)
	json.NewEncoder(w).Encode(formattedTransactions)
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
		fmt.Sscanf(amountStr, "%f", &amount)

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
