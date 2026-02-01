// transactions.go - Transaction parsing and utility functions
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// readMockTransactions reads and parses the mock transactions file
func readMockTransactions() ([]Transaction, error) {
	fileContent, err := os.ReadFile(MockTransactionsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read mock transactions file: %w", err)
	}

	rawTransactions := parseMockTransactionsRaw(string(fileContent))
	transactions := make([]Transaction, len(rawTransactions))
	for i, raw := range rawTransactions {
		transactions[i] = Transaction{
			Date:       raw["date"].(string),
			Merchant:   raw["merchant"].(string),
			Product:    raw["product"].(string),
			Amount:     raw["amount"].(string),
			IsIncoming: raw["incoming"].(bool), //True if "T", else False
		}
	}
	return transactions, nil
}

// parseMockTransactionsRaw extracts individual transactions from the file content
func parseMockTransactionsRaw(content string) []map[string]interface{} {
	var transactions []map[string]interface{}

	lines := strings.Split(content, "\n")
	inDataSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header lines and find the start of transaction data
		if strings.Contains(line, "DATE") && strings.Contains(line, "MERCHANT") {
			inDataSection = true
			continue
		}

		// Stop at the end marker
		if strings.HasPrefix(line, "=======") && inDataSection {
			break
		}

		// Skip separator lines
		if strings.HasPrefix(line, "----") {
			continue
		}

		// Parse transaction lines
		if inDataSection && len(line) > 0 && !isHeaderLine(line) {
			if tx := parseTransactionLine(line); tx != nil {
				transactions = append(transactions, map[string]interface{}{
					"date":     tx.Date,
					"merchant": tx.Merchant,
					"product":  tx.Product,
					"amount":   tx.Amount,
					"incoming": tx.IsIncoming,
				})
			}
		}
	}

	return transactions
}

// isHeaderLine checks if a line is a header or separator
func isHeaderLine(line string) bool {
	headerIndicators := []string{"MOCK", "Account", "Card", "Period", "DATE"}
	for _, indicator := range headerIndicators {
		if strings.HasPrefix(line, indicator) {
			return true
		}
	}
	return false
}

// parseTransactionLine parses a single transaction line
func parseTransactionLine(line string) *Transaction {
	parts := strings.Split(line, "|")
	if len(parts) != 5 {
		return nil
	}

	date := strings.TrimSpace(parts[0])
	merchant := strings.TrimSpace(parts[1])
	product := strings.TrimSpace(parts[2])
	amount := strings.TrimSpace(parts[3])
	isIncoming := strings.TrimSpace(parts[4])

	// Skip if this looks like a header
	if date == "" || date == "DATE" || merchant == "" {
		return nil
	}

	return &Transaction{
		Date:       date,
		Merchant:   merchant,
		Product:    product,
		Amount:     amount,
		IsIncoming: isIncoming == "T",
	}
}

// filterRecentTransactions returns transactions from the past N days
func filterRecentTransactions(transactions []Transaction, days int) []Transaction {
	weekAgo := time.Now().AddDate(0, 0, -days)
	var recent []Transaction

	for _, tx := range transactions {
		// Parse date (format: 2026-01-30)
		txDate, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			continue
		}

		// Include if within lookback period
		if txDate.After(weekAgo) {
			recent = append(recent, tx)
		}
	}

	return recent
}

// These functions are in tools.go to avoid duplication
