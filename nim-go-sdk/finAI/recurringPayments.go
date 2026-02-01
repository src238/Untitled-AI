// recurring_payments.go - AI-powered recurring payment detection
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// RecurringPayment represents a detected recurring payment
type RecurringPayment struct {
	ID          string  `json:"id"`
	Merchant    string  `json:"merchant"`
	Product     string  `json:"product"`
	Amount      float64 `json:"amount"`
	Frequency   string  `json:"frequency"` // e.g. "Monthly", "Weekly", "Bi-weekly"
	LastSeen    string  `json:"lastSeen"`
	Occurrences int     `json:"occurrences"`
}

var (
	recurringPayments      []RecurringPayment
	recurringPaymentsMutex sync.RWMutex
	recurringDetected      bool
)

// detectRecurringPayments runs AI detection once and caches the result
func detectRecurringPayments(anthropicKey string) {
	client := anthropic.NewClient(option.WithAPIKey(anthropicKey))

	transactions, err := readMockTransactions()
	if err != nil {
		log.Printf("❌ Recurring payments: failed to read transactions: %v", err)
		return
	}

	// Only use outgoing transactions
	var outgoing []Transaction
	for _, tx := range transactions {
		if !tx.IsIncoming {
			outgoing = append(outgoing, tx)
		}
	}

	if len(outgoing) == 0 {
		log.Printf("ℹ️  Recurring payments: no outgoing transactions to analyze")
		return
	}

	log.Printf("🔄 Detecting recurring payments from %d outgoing transactions...", len(outgoing))

	result, err := getRecurringPaymentsAI(&client, outgoing)
	if err != nil {
		log.Printf("❌ Recurring payments AI error: %v", err)
		return
	}

	recurringPaymentsMutex.Lock()
	recurringPayments = result
	recurringDetected = true
	recurringPaymentsMutex.Unlock()

	log.Printf("✅ Recurring payments detected: %d", len(result))
}

// getRecurringPaymentsAI sends transaction data to Claude and parses the response
func getRecurringPaymentsAI(client *anthropic.Client, transactions []Transaction) ([]RecurringPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AnalysisTimeout)
	defer cancel()

	prompt := buildRecurringPaymentsPrompt(transactions)

	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(ClaudeModel),
		MaxTokens: 2000,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, err
	}

	// Extract text response
	var responseText strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			responseText.WriteString(block.Text)
		}
	}

	return parseRecurringPaymentsResponse(responseText.String())
}

// buildRecurringPaymentsPrompt creates the prompt with all transaction data
func buildRecurringPaymentsPrompt(transactions []Transaction) string {
	var txList strings.Builder
	for _, tx := range transactions {
		txList.WriteString(fmt.Sprintf("- %s | %s | %s | %s\n", tx.Date, tx.Merchant, tx.Product, tx.Amount))
	}

	return fmt.Sprintf(`You are analyzing a list of financial transactions to detect recurring payments.
A recurring payment is one where the same merchant appears multiple times, especially at regular intervals (weekly, bi-weekly, monthly, etc.), or is a known subscription service.

TRANSACTIONS:
%s

TASK:
Identify all recurring or subscription-like payments. Group by merchant. For each, determine the frequency and typical amount.

RESPOND ONLY with a valid JSON array, no preamble, no markdown, no explanation. Each object must have these exact fields:
[
  {
    "id": "rp-1",
    "merchant": "Netflix",
    "product": "Premium Plan Monthly Subscription",
    "amount": 22.99,
    "frequency": "Monthly",
    "lastSeen": "2026-01-29",
    "occurrences": 1
  }
]

RULES:
- "id" must be "rp-1", "rp-2", etc.
- "amount" must be a number (no $ sign)
- "frequency" must be one of: "Weekly", "Bi-weekly", "Monthly", "Annual", "Irregular"
- "lastSeen" must be the most recent date that merchant appeared
- "occurrences" is how many times that merchant appears in the data
- Include known subscription services even if they only appear once (Netflix, Spotify, Hulu, Disney+, etc.)
- Include merchants that appear 2+ times even if the interval is irregular
- Do NOT include one-off purchases from merchants that only appear once and are not subscription services`, txList.String())
}

// parseRecurringPaymentsResponse parses the JSON array from Claude's response
func parseRecurringPaymentsResponse(response string) ([]RecurringPayment, error) {
	// Strip any markdown code fences if present
	cleaned := strings.TrimSpace(response)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var payments []RecurringPayment
	if err := json.Unmarshal([]byte(cleaned), &payments); err != nil {
		return nil, fmt.Errorf("failed to parse recurring payments JSON: %w\nraw response: %s", err, cleaned)
	}

	return payments, nil
}

// getRecurringPayments returns the cached recurring payments (thread-safe)
func getRecurringPayments() []RecurringPayment {
	recurringPaymentsMutex.RLock()
	defer recurringPaymentsMutex.RUnlock()
	return recurringPayments
}

// isRecurringDetected returns whether detection has run yet
func isRecurringDetected() bool {
	recurringPaymentsMutex.RLock()
	defer recurringPaymentsMutex.RUnlock()
	return recurringDetected
}
