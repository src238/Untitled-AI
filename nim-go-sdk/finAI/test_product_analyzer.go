package main

import (
	"fmt"
	"strings"
)

// Mock transaction data for testing
var mockTransactions = []map[string]interface{}{
	{
		"type":        "send",
		"amount":      15.99,
		"description": "Monthly subscription",
		"merchant":    "Netflix",
		"memo":        "Streaming service",
	},
	{
		"type":        "send",
		"amount":      45.67,
		"description": "Grocery shopping",
		"merchant":    "Whole Foods",
		"memo":        "Weekly groceries",
	},
	{
		"type":        "send",
		"amount":      12.50,
		"description": "Lunch",
		"merchant":    "Chipotle",
		"memo":        "Burrito bowl",
	},
	{
		"type":        "send",
		"amount":      89.99,
		"description": "Electronics",
		"merchant":    "Best Buy",
		"memo":        "Wireless mouse",
	},
	{
		"type":        "send",
		"amount":      10.99,
		"description": "Music streaming",
		"merchant":    "Spotify",
		"memo":        "Premium subscription",
	},
	{
		"type":        "send",
		"amount":      32.45,
		"description": "Dinner",
		"merchant":    "Local Italian Restaurant",
		"memo":        "Date night",
	},
	{
		"type":        "send",
		"amount":      156.78,
		"description": "Groceries",
		"merchant":    "Trader Joe's",
		"memo":        "Monthly stock up",
	},
	{
		"type":        "send",
		"amount":      8.50,
		"description": "Coffee",
		"merchant":    "Starbucks",
		"memo":        "Morning coffee",
	},
}

// Test function to demonstrate the AI analysis
func testProductAnalyzer() {
	fmt.Println("🧪 Testing Product Analyzer")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Build transaction summary like the real function does
	var txSummary strings.Builder
	txSummary.WriteString("Here are the recent transactions to analyze:\n\n")

	for i, tx := range mockTransactions {
		txSummary.WriteString(fmt.Sprintf("Transaction %d:\n", i+1))

		if txType, ok := tx["type"].(string); ok {
			txSummary.WriteString(fmt.Sprintf("  Type: %s\n", txType))
		}
		if amount, ok := tx["amount"].(float64); ok {
			txSummary.WriteString(fmt.Sprintf("  Amount: $%.2f\n", amount))
		}
		if description, ok := tx["description"].(string); ok {
			txSummary.WriteString(fmt.Sprintf("  Description: %s\n", description))
		}
		if merchant, ok := tx["merchant"].(string); ok {
			txSummary.WriteString(fmt.Sprintf("  Merchant: %s\n", merchant))
		}
		if memo, ok := tx["memo"].(string); ok && memo != "" {
			txSummary.WriteString(fmt.Sprintf("  Memo: %s\n", memo))
		}
		txSummary.WriteString("\n")
	}

	fmt.Println("📊 Transaction Data to be Analyzed:")
	fmt.Println("───────────────────────────────────────────────────")
	fmt.Println(txSummary.String())

	fmt.Println("🤖 This data would be sent to Claude AI for analysis...")
	fmt.Println()
	fmt.Println("📝 Expected Analysis Format:")
	fmt.Println("───────────────────────────────────────────────────")
	fmt.Println(`
**Purchase Categories:**
- Dining & Restaurants: 3 transactions - Mix of lunch, dinner, and coffee
- Groceries: 2 transactions - Weekly and monthly grocery shopping
- Digital Subscriptions: 2 transactions - Entertainment streaming services
- Electronics: 1 transaction - Computer accessory purchase

**Specific Products Identified:**
- Netflix subscription - $15.99
- Spotify Premium - $10.99
- Wireless mouse from Best Buy - $89.99
- Whole Foods groceries - $45.67
- Trader Joe's monthly stock up - $156.78
- Chipotle burrito bowl - $12.50
- Italian Restaurant dinner - $32.45
- Starbucks coffee - $8.50

**Insights:**
- Total spending: $372.87 across 8 transactions
- Subscription services account for $26.98/month
- Dining out represents 3 transactions ($53.45 total)
- Grocery spending is healthy at $202.45 for the period
- Consider evaluating if all subscriptions are being used regularly
- Electronics purchase was a one-time expense for productivity
`)

	fmt.Println("✅ Test Complete!")
	fmt.Println()
	fmt.Println("💡 To test with real data:")
	fmt.Println("   1. Open http://localhost:5173")
	fmt.Println("   2. Login with your Liminal account")
	fmt.Println("   3. Ask: 'What products did I buy recently?'")
	fmt.Println("   4. Or: 'Analyze my purchases from the last 30 days'")
}

func main() {
	testProductAnalyzer()

	// Show the actual tool parameters
	fmt.Println()
	fmt.Println("🔧 Tool Configuration:")
	fmt.Println("───────────────────────────────────────────────────")
	fmt.Println("Tool Name: analyze_products")
	fmt.Println("Parameters:")
	fmt.Println("  - days (optional): Number of days to analyze (default: 30)")
	fmt.Println("  - limit (optional): Max transactions to analyze (default: 50)")
	fmt.Println()
	fmt.Println("AI Model: Claude 3.5 Sonnet (claude-3-5-sonnet-20241022)")
	fmt.Println("Max Tokens: 2048")
	fmt.Println()
}
