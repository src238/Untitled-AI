// analysis.go - Background AI analysis for finding product alternatives
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// startAIAnalysisLoop checks purchases from the past week for cheaper alternatives
// and posts findings to the notice board
func startAIAnalysisLoop(anthropicKey string) {
	client := anthropic.NewClient(option.WithAPIKey(anthropicKey))
	clientPtr := &client

	// Wait before first analysis
	time.Sleep(AnalysisInitialDelay)

	for {
		analyzeNextProduct(clientPtr)
		time.Sleep(AnalysisInterval)
	}
}

// analyzeNextProduct finds and analyzes one unchecked product
func analyzeNextProduct(client *anthropic.Client) {
	// Read and filter transactions
	transactions, err := readMockTransactions()
	if err != nil {
		log.Printf("❌ Failed to read transactions: %v", err)
		return
	}

	recentTransactions := filterRecentTransactions(transactions, TransactionLookbackDays)
	if len(recentTransactions) == 0 {
		log.Printf("ℹ️  No recent transactions to analyze")
		return
	}

	// Find unchecked product
	productToCheck := findUncheckedProduct(recentTransactions)
	if productToCheck == nil {
		// All products checked, reset and wait
		resetCheckedProducts()
		log.Printf("✅ All recent purchases checked for alternatives, resetting...")
		time.Sleep(AnalysisResetDelay)
		return
	}

	// Analyze and post alternative if found
	analyzeAndPostAlternative(client, productToCheck)
}

// findUncheckedProduct finds the first unchecked product in transactions
func findUncheckedProduct(transactions []Transaction) *Transaction {
	checkedProductsMutex.Lock()
	defer checkedProductsMutex.Unlock()

	for _, tx := range transactions {
		if tx.Product != "" && !checkedProducts[tx.Product] {
			checkedProducts[tx.Product] = true
			return &tx
		}
	}

	return nil
}

// resetCheckedProducts clears the checked products map
func resetCheckedProducts() {
	checkedProductsMutex.Lock()
	defer checkedProductsMutex.Unlock()
	checkedProducts = make(map[string]bool)
}

// analyzeAndPostAlternative uses AI to find alternatives and posts if savings meet minimum
func analyzeAndPostAlternative(client *anthropic.Client, tx *Transaction) {
	log.Printf("🔍 Checking for cheaper alternatives: %s ($%s)", tx.Product, tx.Amount)

	// Get AI recommendation
	recommendation, err := getAIRecommendation(client, tx)
	if err != nil {
		log.Printf("❌ AI analysis error: %v", err)
		return
	}

	// Validate and post if meets criteria
	if shouldPostRecommendation(recommendation, tx.Product) {
		postAlternativeAlert(recommendation)
	}
}

// getAIRecommendation asks Claude to find a cheaper alternative
func getAIRecommendation(client *anthropic.Client, tx *Transaction) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AnalysisTimeout)
	defer cancel()

	prompt := buildAlternativePrompt(tx)

	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(ClaudeModel),
		MaxTokens: 200,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	// Extract text response
	var responseText strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			responseText.WriteString(block.Text)
		}
	}

	return strings.TrimSpace(responseText.String()), nil
}

// buildAlternativePrompt creates the AI prompt for finding alternatives
func buildAlternativePrompt(tx *Transaction) string {
	return fmt.Sprintf(`You are analyzing a purchase to find cheaper alternatives while preserving quality.

PURCHASE DETAILS:
- Product: %s
- Price Paid: %s
- Merchant: %s
- Date: %s

TASK:
Find a cheaper alternative that maintains or improves quality. You MUST save at least $%.2f to recommend an alternative.

FORMAT YOUR RESPONSE EXACTLY LIKE THIS (no preamble):
"[Original Product] ($[original price]) - Alternative: [New Product] ($[new price]) - Save: $[difference] - Buy: [URL]"

EXAMPLES:
- "Echo Dot 5th Gen ($49.99) - Alternative: Google Nest Mini ($29.99) - Save: $20.00 - Buy: https://store.google.com/product/google_nest_mini"
- "Nike Air Max ($139.99) - Alternative: Adidas Ultraboost ($120.00) - Save: $19.99 - Buy: https://www.adidas.com/us/ultraboost"
- "Whole Foods Groceries ($127.83) - Alternative: Trader Joe's Organic Mix ($95.00) - Save: $32.83 - Buy: https://www.traderjoes.com"

CRITICAL RULES:
- MUST save at least $%.2f or respond: "Not enough savings (under $%.0f)"
- MUST include a real, working purchase link (Amazon, official store, major retailer)
- Use exact format with " - " separators (no vertical pipes)
- Show all prices with $ and two decimal places
- Alternative must maintain or improve quality
- URL should be direct product page when possible
- Focus on 2026 realistic pricing and real retailers`,
		tx.Product, tx.Amount, tx.Merchant, tx.Date,
		MinimumSavings, MinimumSavings, MinimumSavings)
}

// shouldPostRecommendation checks if recommendation meets posting criteria
func shouldPostRecommendation(recommendation, product string) bool {
	if recommendation == "" {
		return false
	}

	recLower := strings.ToLower(recommendation)

	// Skip if insufficient savings
	if strings.Contains(recLower, "optimal") ||
		strings.Contains(recLower, "not enough savings") ||
		strings.Contains(recLower, "under $5") {
		log.Printf("✓ No better alternative for: %s (insufficient savings or optimal)", product)
		return false
	}

	// Extract and validate savings amount
	if !strings.Contains(recommendation, "Save: $") {
		return false
	}

	savings := extractSavingsAmount(recommendation)
	if savings < MinimumSavings {
		log.Printf("✓ Savings too low for %s: $%.2f (minimum $%.2f)", product, savings, MinimumSavings)
		return false
	}

	return true
}

// extractSavingsAmount parses the savings amount from a recommendation
func extractSavingsAmount(recommendation string) float64 {
	saveIdx := strings.Index(recommendation, "Save: $")
	if saveIdx == -1 {
		return 0
	}

	saveStr := recommendation[saveIdx+7:] // Skip "Save: $"
	endIdx := strings.IndexAny(saveStr, " -")
	if endIdx != -1 {
		saveStr = saveStr[:endIdx]
	}

	var savings float64
	fmt.Sscanf(saveStr, "%f", &savings)
	return savings
}

// postAlternativeAlert creates and stores an alert for a product alternative
func postAlternativeAlert(message string) {
	alert := Alert{
		ID:        fmt.Sprintf("alt-%d", time.Now().UnixNano()),
		Message:   message,
		Timestamp: time.Now(),
		Type:      "success",
	}

	alertsMutex.Lock()
	alerts = append(alerts, alert)
	if len(alerts) > MaxAlertsStored {
		alerts = alerts[len(alerts)-MaxAlertsStored:]
	}
	alertsMutex.Unlock()

	log.Printf("💡 Posted alternative: %s", message)
}
