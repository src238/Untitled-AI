// types.go - Type definitions for the application
package main

import (
	"time"
)

// Alert represents a notification for the user
type Alert struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "info", "warning", "success"
}

// Transaction represents a parsed transaction from mock data
type Transaction struct {
	Date       string
	Merchant   string
	Product    string
	Amount     string
	IsIncoming bool
}

// TransactionAPI represents a transaction in the API response format
type TransactionAPI struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Merchant    string  `json:"merchant"`
	IsIncoming  bool    `json:"isIncoming"`
}

// ToolParams represents common parameters for tool handlers
type AnalyzeSpendingParams struct {
	Days int `json:"days"`
}

type AnalyzeProductsParams struct {
	Days  int `json:"days"`
	Limit int `json:"limit"`
}

type MockTransactionParams struct {
	Format string `json:"format"`
}

type ProductSearchParams struct {
	ProductName    string `json:"product_name"`
	OriginalPrice  string `json:"original_price"`
	SearchCriteria string `json:"search_criteria"`
	MaxPrice       string `json:"max_price"`
	CategoryFilter string `json:"category_filter"`
}

type AlertParams struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ReadAlertsParams struct {
	Hours string `json:"hours"`
	Type  string `json:"type"`
}
