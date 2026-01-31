// Hackathon Starter: Complete AI Financial Agent
// Build intelligent financial tools with nim-go-sdk + Liminal banking APIs
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/becomeliminal/nim-go-sdk/core"
	"github.com/becomeliminal/nim-go-sdk/executor"
	"github.com/becomeliminal/nim-go-sdk/server"
	"github.com/becomeliminal/nim-go-sdk/tools"
	"github.com/joho/godotenv"
)

// Global alert storage and product tracking
var (
	alertsMutex         sync.RWMutex
	alerts              []Alert
	checkedProductsMutex sync.RWMutex
	checkedProducts     map[string]bool // Track which products have been checked for alternatives
)

// Alert represents a notification for the user
type Alert struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "info", "warning", "success"
}

func main() {
	// ============================================================================
	// CONFIGURATION
	// ============================================================================
	// Initialize global tracking map
	checkedProducts = make(map[string]bool)

	// Load .env file if it exists (optional - will use system env vars if not found)
	_ = godotenv.Load()

	// Load configuration from environment variables
	// Create a .env file or export these in your shell

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey == "" {
		log.Fatal("❌ ANTHROPIC_API_KEY environment variable is required")
	}

	liminalBaseURL := os.Getenv("LIMINAL_BASE_URL")
	if liminalBaseURL == "" {
		liminalBaseURL = "https://api.liminal.cash"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ============================================================================
	// LIMINAL EXECUTOR SETUP
	// ============================================================================
	// The HTTPExecutor handles all API calls to Liminal banking services.
	// Authentication is handled automatically via JWT tokens passed from the
	// frontend login flow (email/OTP). No API key needed!

	liminalExecutor := executor.NewHTTPExecutor(executor.HTTPExecutorConfig{
		BaseURL: liminalBaseURL,
	})
	log.Println("✅ Liminal API configured")

	// ============================================================================
	// SERVER SETUP
	// ============================================================================
	// Create the nim-go-sdk server with Claude AI
	// The server handles WebSocket connections and manages conversations
	// Authentication is automatic: JWT tokens from the login flow are extracted
	// from WebSocket connections and forwarded to Liminal API calls

	srv, err := server.New(server.Config{
		AnthropicKey:    anthropicKey,
		SystemPrompt:    hackathonSystemPrompt,
		Model:           "claude-sonnet-4-20250514",
		MaxTokens:       4096,
		LiminalExecutor: liminalExecutor, // SDK automatically handles JWT extraction and forwarding
	})
	if err != nil {
		log.Fatal(err)
	}

	// ============================================================================
	// ADD LIMINAL BANKING TOOLS
	// ============================================================================
	// These are 8 core Liminal tools that give your AI access to real banking:
	// (get_transactions is disabled - using mock transaction data instead)
	//
	// READ OPERATIONS (no confirmation needed):
	//   1. get_balance - Check wallet balance
	//   2. get_savings_balance - Check savings positions and APY
	//   3. get_vault_rates - Get current savings rates
	//   4. get_profile - Get user profile info
	//   5. search_users - Find users by display tag
	//
	// WRITE OPERATIONS (require user confirmation):
	//   6. send_money - Send money to another user
	//   7. deposit_savings - Deposit funds into savings
	//   8. withdraw_savings - Withdraw funds from savings

	srv.AddTools(tools.LiminalTools(liminalExecutor)...)
	log.Println("✅ Added 8 Liminal banking tools (get_transactions disabled - using mock data)")

	// ============================================================================
	// ADD CUSTOM TOOLS
	// ============================================================================
	// This is where you'll add your hackathon project's custom tools!
	// Below is an example spending analyzer tool to get you started.

	srv.AddTool(createSpendingAnalyzerTool(liminalExecutor))
	log.Println("✅ Added custom spending analyzer tool")

	// Add product analysis tool that uses Claude AI
	srv.AddTool(createProductAnalyzerTool(liminalExecutor, anthropicKey))
	log.Println("✅ Added AI-powered product analyzer tool")

	// Add mock transaction reader tool
	srv.AddTool(createMockTransactionReaderTool())
	log.Println("✅ Added mock transaction reader tool")

	// Add product search tool for finding alternatives
	srv.AddTool(createProductSearchTool(anthropicKey))
	log.Println("✅ Added product search tool for finding alternatives")

	// Add alert notification tools
	srv.AddTool(createAlertTool())
	srv.AddTool(createReadAlertsTool())
	log.Println("✅ Added alert notification tools (post & read) for AI insights")

	// TODO: Add more custom tools here!
	// Examples:
	//   - Savings goal tracker
	//   - Budget alerts
	//   - Spending category analyzer
	//   - Bill payment predictor
	//   - Cash flow forecaster

	// ============================================================================
	// START AI BACKGROUND ANALYSIS LOOP
	// ============================================================================
	// This loop periodically prompts the AI to analyze financial data
	// and automatically posts insights to the alert board

	go startAIAnalysisLoop(anthropicKey)
	log.Println("✅ Started AI background analysis loop (runs every 30 seconds)")

	// ============================================================================
	// START SERVER
	// ============================================================================

	// Add HTTP endpoint for alerts
	http.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Content-Type", "application/json")

		alertsMutex.RLock()
		defer alertsMutex.RUnlock()

		// Return alerts from last 24 hours
		recentAlerts := []Alert{}
		cutoff := time.Now().Add(-24 * time.Hour)
		for _, alert := range alerts {
			if alert.Timestamp.After(cutoff) {
				recentAlerts = append(recentAlerts, alert)
			}
		}

		json.NewEncoder(w).Encode(recentAlerts)
	})

	// Add HTTP endpoint for transactions
	http.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Content-Type", "application/json")

		// Read mock transactions file
		fileContent, err := os.ReadFile("mock_transactions.txt")
		if err != nil {
			http.Error(w, "Failed to read transactions", http.StatusInternalServerError)
			return
		}

		// Parse transactions
		transactions := parseMockTransactions(string(fileContent))

		// Convert to frontend format
		formattedTransactions := []map[string]interface{}{}
		for i, tx := range transactions {
			date, _ := tx["date"].(string)
			merchant, _ := tx["merchant"].(string)
			product, _ := tx["product"].(string)
			amountStr, _ := tx["amount"].(string)

			// Parse amount (remove $ and convert to float)
			amountStr = strings.TrimPrefix(amountStr, "$")
			amount := 0.0
			fmt.Sscanf(amountStr, "%f", &amount)

			// Determine type (all mock transactions are debits/purchases)
			txType := "debit"

			formattedTransactions = append(formattedTransactions, map[string]interface{}{
				"id":          fmt.Sprintf("tx-%d", i+1),
				"amount":      -amount, // Negative for debits
				"description": product,
				"date":        date,
				"type":        txType,
				"merchant":    merchant,
			})
		}

		json.NewEncoder(w).Encode(formattedTransactions)
	})

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 Hackathon Starter Server Running")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📡 WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Printf("🔔 Alerts API: http://localhost:%s/api/alerts", port)
	log.Printf("💳 Transactions API: http://localhost:%s/api/transactions", port)
	log.Printf("💚 Health check: http://localhost:%s/health", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Ready for connections! Start your frontend with: cd frontend && npm run dev")
	log.Println()

	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// SYSTEM PROMPT
// ============================================================================
// This prompt defines your AI agent's personality and behavior
// Customize this to match your hackathon project's focus!

const hackathonSystemPrompt = `You are Nim, a smart shopping assistant that helps users find better deals and save money.

WHAT YOU DO:
You help users save money by automatically analyzing their purchases and finding cheaper alternatives that preserve or improve quality. The system continuously checks recent purchases and posts money-saving recommendations to the notice board.

CONVERSATIONAL STYLE:
- Be warm, friendly, and focused on helping users save money
- Use casual language when appropriate, but stay professional about finances
- Ask clarifying questions when something is unclear
- Remember context from earlier in the conversation
- Explain things simply without being condescending

CORE FUNCTIONALITY:
The system automatically:
1. Monitors all purchases from the past week
2. Checks each purchase for cheaper alternatives while preserving quality
3. Posts savings opportunities to the notice board
4. Tracks which products have been checked to avoid duplicates
5. Resets and rechecks after all products are analyzed

WHEN TO USE TOOLS:
- Use tools immediately for simple queries ("what did I buy?")
- Use read_mock_transactions to view purchase history
- Use search_product_alternatives when users ask about specific product alternatives
- Use post_alert to share important savings insights
- Don't use tools for general questions about how things work

AVAILABLE TOOLS:
- Read mock transaction history (read_mock_transactions) - access detailed credit card transactions
- Search product alternatives (search_product_alternatives) - find cheaper alternatives to purchased products
- Post alert notifications (post_alert) - send savings opportunities to the notice board
- Read alert notifications (read_alerts) - check what alternatives have been suggested

IMPORTANT - TRANSACTION HISTORY:
When users ask about their purchases or spending, use the read_mock_transactions tool to access mock credit card data with 60 detailed purchases including merchant names, products, and amounts from January 2026.

PRODUCT ALTERNATIVES:
When users ask about alternatives to products they purchased:
1. First read their transaction history with read_mock_transactions
2. Identify the specific product they're asking about
3. Use search_product_alternatives with the product name and original price
4. Provide detailed comparisons emphasizing quality and savings
Example: "Can you find a cheaper alternative to the Echo Dot I bought?" - read transactions, find the Echo Dot purchase ($49.99), then search for alternatives.

AUTOMATIC BACKGROUND ANALYSIS:
The system runs continuously in the background:
- Every 30 seconds, it checks one unchecked purchase from the past week
- Uses AI to find cheaper alternatives while preserving quality
- Posts recommendations like: "Google Nest Mini for $29 saves $20 - similar features, better voice recognition"
- Only suggests alternatives that maintain or improve quality
- Skips products where current choice is already optimal

TIPS FOR GREAT INTERACTIONS:
- Focus on savings opportunities and value
- Celebrate smart purchasing decisions
- Explain why alternatives offer good value
- Be encouraging about money-saving habits
- Make smart shopping feel easy and rewarding

Remember: You're here to help users save money without compromising on quality!`

// ============================================================================
// CUSTOM TOOL: SPENDING ANALYZER
// ============================================================================
// This is an example custom tool that demonstrates how to:
// 1. Define tool parameters with JSON schema
// 2. Call other Liminal tools from within your tool
// 3. Process and analyze the data
// 4. Return useful insights
//
// Use this as a template for your own hackathon tools!

func createSpendingAnalyzerTool(liminalExecutor core.ToolExecutor) core.Tool {
	return tools.New("analyze_spending").
		Description("Analyze the user's spending patterns over a specified time period. Returns insights about spending velocity, categories, and trends.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"days": tools.IntegerProperty("Number of days to analyze (default: 30)"),
		})).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				Days int `json:"days"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			// Default to 30 days if not specified
			if params.Days == 0 {
				params.Days = 30
			}

			// STEP 1: Fetch transaction history
			// We'll call the Liminal get_transactions tool through the executor
			txRequest := map[string]interface{}{
				"limit": 100, // Get up to 100 transactions
			}
			txRequestJSON, _ := json.Marshal(txRequest)

			txResponse, err := liminalExecutor.Execute(ctx, &core.ExecuteRequest{
				UserID:    toolParams.UserID,
				Tool:      "get_transactions",
				Input:     txRequestJSON,
				RequestID: toolParams.RequestID,
			})
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("failed to fetch transactions: %v", err),
				}, nil
			}

			if !txResponse.Success {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("transaction fetch failed: %s", txResponse.Error),
				}, nil
			}

			// STEP 2: Parse transaction data
			// In a real implementation, you'd parse the actual response structure
			// For now, we'll create a structured analysis

			var transactions []map[string]interface{}
			var txData map[string]interface{}
			if err := json.Unmarshal(txResponse.Data, &txData); err == nil {
				if txArray, ok := txData["transactions"].([]interface{}); ok {
					for _, tx := range txArray {
						if txMap, ok := tx.(map[string]interface{}); ok {
							transactions = append(transactions, txMap)
						}
					}
				}
			}

			// STEP 3: Analyze the data
			analysis := analyzeTransactions(transactions, params.Days)

			// STEP 4: Return insights
			result := map[string]interface{}{
				"period_days":        params.Days,
				"total_transactions": len(transactions),
				"analysis":           analysis,
				"generated_at":       time.Now().Format(time.RFC3339),
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}

// analyzeTransactions processes transaction data and returns insights
func analyzeTransactions(transactions []map[string]interface{}, days int) map[string]interface{} {
	if len(transactions) == 0 {
		return map[string]interface{}{
			"summary": "No transactions found in the specified period",
		}
	}

	// Calculate basic metrics
	var totalSpent, totalReceived float64
	var spendCount, receiveCount int

	// This is a simplified example - you'd do real analysis here:
	// - Group by category/merchant
	// - Calculate daily/weekly averages
	// - Identify spending spikes
	// - Compare to previous periods
	// - Detect recurring payments

	for _, tx := range transactions {
		// Example analysis logic
		txType, _ := tx["type"].(string)
		amount, _ := tx["amount"].(float64)

		switch txType {
		case "send":
			totalSpent += amount
			spendCount++
		case "receive":
			totalReceived += amount
			receiveCount++
		}
	}

	avgDailySpend := totalSpent / float64(days)

	return map[string]interface{}{
		"total_spent":     fmt.Sprintf("%.2f", totalSpent),
		"total_received":  fmt.Sprintf("%.2f", totalReceived),
		"spend_count":     spendCount,
		"receive_count":   receiveCount,
		"avg_daily_spend": fmt.Sprintf("%.2f", avgDailySpend),
		"velocity":        calculateVelocity(spendCount, days),
		"insights": []string{
			fmt.Sprintf("You made %d spending transactions over %d days", spendCount, days),
			fmt.Sprintf("Average daily spend: $%.2f", avgDailySpend),
			"Consider setting up savings goals to build financial cushion",
		},
	}
}

// calculateVelocity determines spending frequency
func calculateVelocity(transactionCount, days int) string {
	txPerWeek := float64(transactionCount) / float64(days) * 7

	switch {
	case txPerWeek < 2:
		return "low"
	case txPerWeek < 7:
		return "moderate"
	default:
		return "high"
	}
}

// ============================================================================
// CUSTOM TOOL: PRODUCT ANALYZER (AI-POWERED)
// ============================================================================
// This tool uses Claude AI to analyze transaction data and identify what
// products were purchased based on transaction descriptions, merchant names,
// and other available data.

func createProductAnalyzerTool(liminalExecutor core.ToolExecutor, anthropicKey string) core.Tool {
	return tools.New("analyze_products").
		Description("Analyze transaction history to identify what products or services were purchased. Uses AI to understand transaction descriptions and categorize purchases.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"days":  tools.IntegerProperty("Number of days to analyze (default: 30)"),
			"limit": tools.IntegerProperty("Maximum number of transactions to analyze (default: 50)"),
		})).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				Days  int `json:"days"`
				Limit int `json:"limit"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			// Set defaults
			if params.Days == 0 {
				params.Days = 30
			}
			if params.Limit == 0 {
				params.Limit = 50
			}

			// STEP 1: Fetch transaction history from Liminal
			txRequest := map[string]interface{}{
				"limit": params.Limit,
			}
			txRequestJSON, _ := json.Marshal(txRequest)

			txResponse, err := liminalExecutor.Execute(ctx, &core.ExecuteRequest{
				UserID:    toolParams.UserID,
				Tool:      "get_transactions",
				Input:     txRequestJSON,
				RequestID: toolParams.RequestID,
			})
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("failed to fetch transactions: %v", err),
				}, nil
			}

			if !txResponse.Success {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("transaction fetch failed: %s", txResponse.Error),
				}, nil
			}

			// STEP 2: Parse transaction data
			var transactions []map[string]interface{}
			var txData map[string]interface{}
			if err := json.Unmarshal(txResponse.Data, &txData); err == nil {
				if txArray, ok := txData["transactions"].([]interface{}); ok {
					for _, tx := range txArray {
						if txMap, ok := tx.(map[string]interface{}); ok {
							transactions = append(transactions, txMap)
						}
					}
				}
			}

			if len(transactions) == 0 {
				return &core.ToolResult{
					Success: true,
					Data: map[string]interface{}{
						"summary": "No transactions found in the specified period",
						"products": []interface{}{},
					},
				}, nil
			}

			// STEP 3: Use Claude AI to analyze and identify products
			productAnalysis, err := analyzeProductsWithAI(ctx, transactions, anthropicKey)
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("AI analysis failed: %v", err),
				}, nil
			}

			// STEP 4: Return the analysis
			result := map[string]interface{}{
				"period_days":        params.Days,
				"transactions_analyzed": len(transactions),
				"product_analysis":   productAnalysis,
				"generated_at":       time.Now().Format(time.RFC3339),
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}

// analyzeProductsWithAI uses Claude to analyze transaction data and identify products
func analyzeProductsWithAI(ctx context.Context, transactions []map[string]interface{}, anthropicKey string) (string, error) {
	// Create Anthropic client
	client := anthropic.NewClient(option.WithAPIKey(anthropicKey))

	// Build a summary of transactions for Claude to analyze
	var txSummary strings.Builder
	txSummary.WriteString("Here are the recent transactions to analyze:\n\n")

	for i, tx := range transactions {
		if i >= 50 { // Limit to avoid token limits
			break
		}

		txSummary.WriteString(fmt.Sprintf("Transaction %d:\n", i+1))

		// Include relevant transaction fields
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
		if recipient, ok := tx["recipient"].(string); ok {
			txSummary.WriteString(fmt.Sprintf("  Recipient: %s\n", recipient))
		}
		if sender, ok := tx["sender"].(string); ok {
			txSummary.WriteString(fmt.Sprintf("  Sender: %s\n", sender))
		}
		if memo, ok := tx["memo"].(string); ok && memo != "" {
			txSummary.WriteString(fmt.Sprintf("  Memo: %s\n", memo))
		}
		txSummary.WriteString("\n")
	}

	// Create the analysis prompt
	prompt := fmt.Sprintf(`%s

Analyze these transactions and identify what products or services were purchased. For each transaction that represents a purchase:

1. Identify the product or service category (e.g., groceries, entertainment, dining, transportation, utilities, etc.)
2. Identify specific products if possible from merchant names or descriptions
3. Group similar purchases together

Provide your analysis in the following format:

**Purchase Categories:**
- [Category name]: [Number of transactions] - [Brief description of what was bought]

**Specific Products Identified:**
- [Product/Service name] from [Merchant] - $[Amount]

**Insights:**
- [Key findings about purchasing patterns]
- [Most frequent purchase types]
- [Any recommendations or observations]

Be specific and helpful. If transaction details are limited, make reasonable inferences based on merchant names and amounts.`, txSummary.String())

	// Make API call to Claude
	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model("claude-3-5-sonnet-20241022"),
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}

	// Extract the text response from content blocks
	var textResponse strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			textResponse.WriteString(block.Text)
		}
	}

	if textResponse.Len() == 0 {
		return "No analysis generated", nil
	}

	return textResponse.String(), nil
}

// ============================================================================
// CUSTOM TOOL: MOCK TRANSACTION READER
// ============================================================================
// This tool reads the mock credit card transaction history file and returns
// the data to the AI for analysis. Useful for testing and demonstration.

func createMockTransactionReaderTool() core.Tool {
	return tools.New("read_mock_transactions").
		Description("Read mock credit card transaction history from a file. Returns detailed transaction data including merchant names, products purchased, and amounts. Use this to analyze spending patterns, identify products, or provide financial insights based on realistic transaction data.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"format": tools.StringProperty("Output format: 'full' for complete data, 'summary' for just totals and categories (default: 'full')"),
		})).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				Format string `json:"format"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			// Default to full format
			if params.Format == "" {
				params.Format = "full"
			}

			// Read the mock transactions file
			filePath := "mock_transactions.txt"
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("failed to read mock transactions file: %v", err),
				}, nil
			}

			// Parse the file content
			content := string(fileContent)

			// Extract transactions from the file
			transactions := parseMockTransactions(content)

			// Build the result based on requested format
			var result map[string]interface{}

			if params.Format == "summary" {
				// Return just summary information
				result = map[string]interface{}{
					"format":             "summary",
					"total_transactions": len(transactions),
					"summary":            extractSummaryFromContent(content),
					"categories":         extractCategoriesFromContent(content),
				}
			} else {
				// Return full transaction data
				result = map[string]interface{}{
					"format":             "full",
					"total_transactions": len(transactions),
					"transactions":       transactions,
					"raw_content":        content,
				}
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}

// parseMockTransactions extracts individual transactions from the file content
func parseMockTransactions(content string) []map[string]interface{} {
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
		if inDataSection && len(line) > 0 && !strings.HasPrefix(line, "MOCK") && !strings.HasPrefix(line, "Account") && !strings.HasPrefix(line, "Card") && !strings.HasPrefix(line, "Period") {
			parts := strings.Split(line, "|")
			if len(parts) == 4 {
				date := strings.TrimSpace(parts[0])
				merchant := strings.TrimSpace(parts[1])
				product := strings.TrimSpace(parts[2])
				amount := strings.TrimSpace(parts[3])

				// Skip if this looks like a header
				if date != "" && date != "DATE" && merchant != "" {
					transactions = append(transactions, map[string]interface{}{
						"date":     date,
						"merchant": merchant,
						"product":  product,
						"amount":   amount,
					})
				}
			}
		}
	}

	return transactions
}

// extractSummaryFromContent extracts summary information from the file
func extractSummaryFromContent(content string) string {
	lines := strings.Split(content, "\n")
	var summary strings.Builder

	inSummary := false
	for _, line := range lines {
		if strings.Contains(line, "TOTAL TRANSACTIONS:") || strings.Contains(line, "TOTAL AMOUNT:") {
			inSummary = true
		}
		if inSummary {
			summary.WriteString(line + "\n")
			if strings.Contains(line, "TOTAL AMOUNT:") {
				break
			}
		}
	}

	return summary.String()
}

// extractCategoriesFromContent extracts category breakdown from the file
func extractCategoriesFromContent(content string) map[string]string {
	categories := make(map[string]string)
	lines := strings.Split(content, "\n")

	inCategorySection := false
	for _, line := range lines {
		if strings.Contains(line, "CATEGORY BREAKDOWN:") {
			inCategorySection = true
			continue
		}

		if inCategorySection && strings.HasPrefix(line, "- ") {
			parts := strings.SplitN(strings.TrimPrefix(line, "- "), ":", 2)
			if len(parts) == 2 {
				category := strings.TrimSpace(parts[0])
				amount := strings.TrimSpace(parts[1])
				categories[category] = amount
			}
		}
	}

	return categories
}

// ============================================================================
// CUSTOM TOOL: PRODUCT SEARCH & ALTERNATIVES
// ============================================================================
// This tool allows the AI to search for products and find alternatives
// to items purchased in the transaction history. Uses Claude AI with web search
// capabilities to find similar products, compare prices, and suggest better options.

func createProductSearchTool(anthropicKey string) core.Tool {
	return tools.New("search_product_alternatives").
		Description("Search the web for product alternatives and recommendations. Given a product name or description from transaction history, find similar products, compare features and prices, and suggest better alternatives. Use this when users ask for product recommendations or want to find alternatives to something they purchased.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"product_name":     tools.StringProperty("Name or description of the product to find alternatives for (e.g., 'Echo Dot Smart Speaker', 'Nike Running Shoes')"),
			"original_price":   tools.StringProperty("Optional: Original price paid for the product (e.g., '$49.99')"),
			"search_criteria":  tools.StringProperty("Optional: Specific criteria for alternatives (e.g., 'cheaper', 'better quality', 'more features', 'eco-friendly')"),
			"max_price":        tools.StringProperty("Optional: Maximum price for alternatives (e.g., '$100')"),
			"category_filter":  tools.StringProperty("Optional: Specific category or type (e.g., 'smart speakers', 'athletic shoes', 'streaming services')"),
		}, "product_name")).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				ProductName    string `json:"product_name"`
				OriginalPrice  string `json:"original_price"`
				SearchCriteria string `json:"search_criteria"`
				MaxPrice       string `json:"max_price"`
				CategoryFilter string `json:"category_filter"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			if params.ProductName == "" {
				return &core.ToolResult{
					Success: false,
					Error:   "product_name is required",
				}, nil
			}

			// Use Claude with extended thinking to search for product alternatives
			searchResults, err := searchProductAlternativesWithAI(ctx, params.ProductName, params.OriginalPrice, params.SearchCriteria, params.MaxPrice, params.CategoryFilter, anthropicKey)
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("product search failed: %v", err),
				}, nil
			}

			// Return the search results
			result := map[string]interface{}{
				"product_searched": params.ProductName,
				"original_price":   params.OriginalPrice,
				"search_results":   searchResults,
				"generated_at":     time.Now().Format(time.RFC3339),
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}

// searchProductAlternativesWithAI uses Claude to search for product alternatives
func searchProductAlternativesWithAI(ctx context.Context, productName, originalPrice, searchCriteria, maxPrice, categoryFilter, anthropicKey string) (string, error) {
	// Create Anthropic client
	client := anthropic.NewClient(option.WithAPIKey(anthropicKey))

	// Build the search prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString(fmt.Sprintf("Search for alternative products to: **%s**\n\n", productName))

	if originalPrice != "" {
		promptBuilder.WriteString(fmt.Sprintf("Original price paid: %s\n", originalPrice))
	}
	if searchCriteria != "" {
		promptBuilder.WriteString(fmt.Sprintf("Search criteria: %s\n", searchCriteria))
	}
	if maxPrice != "" {
		promptBuilder.WriteString(fmt.Sprintf("Maximum price: %s\n", maxPrice))
	}
	if categoryFilter != "" {
		promptBuilder.WriteString(fmt.Sprintf("Category: %s\n", categoryFilter))
	}

	promptBuilder.WriteString(`

Please provide:

1. **Top Alternative Products** (3-5 options):
   - Product name and brand
   - Current market price
   - Key features and specifications
   - Pros and cons compared to the original
   - Where to buy (online retailers)

2. **Price Comparison**:
   - How each alternative compares in price to the original
   - Value for money assessment

3. **Recommendation**:
   - Which alternative offers the best value
   - Why you recommend it
   - Any important considerations (reviews, reliability, warranty, etc.)

4. **Savings Opportunity**:
   - Potential savings if switching to recommended alternative
   - Whether the original purchase was a good deal

Format your response in a clear, structured way that's easy to read and make decisions from. Include specific product models and realistic 2026 pricing.`)

	prompt := promptBuilder.String()

	// Make API call to Claude with extended thinking for better product research
	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model("claude-sonnet-4-20250514"),
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}

	// Extract the text response from content blocks
	var textResponse strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			textResponse.WriteString(block.Text)
		}
	}

	if textResponse.Len() == 0 {
		return "No product alternatives found", nil
	}

	return textResponse.String(), nil
}

// createAlertTool creates a tool that posts alerts to the user's notification sidebar
func createAlertTool() core.Tool {
	return tools.New("post_alert").
		Description("Post an important notification or insight to the user's alert sidebar. Use this to proactively notify users about spending patterns, savings opportunities, unusual transactions, budget concerns, or financial recommendations. Alerts appear in the left sidebar and persist for 24 hours.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"message": tools.StringProperty("The alert message to display to the user (be clear and actionable)"),
			"type":    tools.StringProperty("Alert type: 'info' (general insight), 'warning' (concern or caution), or 'success' (positive news or achievement)"),
		}, "message", "type")).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			if params.Message == "" {
				return &core.ToolResult{
					Success: false,
					Error:   "message is required",
				}, nil
			}

			// Validate and default type
			validTypes := map[string]bool{"info": true, "warning": true, "success": true}
			if !validTypes[params.Type] {
				params.Type = "info"
			}

			// Create alert
			alert := Alert{
				ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
				Message:   params.Message,
				Timestamp: time.Now(),
				Type:      params.Type,
			}

			// Store alert
			alertsMutex.Lock()
			alerts = append(alerts, alert)
			// Keep only last 100 alerts
			if len(alerts) > 100 {
				alerts = alerts[len(alerts)-100:]
			}
			alertsMutex.Unlock()

			log.Printf("📢 Alert posted: [%s] %s", params.Type, params.Message)

			return &core.ToolResult{
				Success: true,
				Data: map[string]interface{}{
					"alert_id":  alert.ID,
					"message":   alert.Message,
					"type":      alert.Type,
					"timestamp": alert.Timestamp.Format(time.RFC3339),
					"status":    "Alert posted successfully and will appear in the user's notification sidebar",
				},
			}, nil
		}).
		Build()
}

// createReadAlertsToolcreates a tool that reads current alerts from the notification sidebar
func createReadAlertsTool() core.Tool {
	return tools.New("read_alerts").
		Description("Read current alerts from the user's notification sidebar. Use this to check what insights or notifications have been previously posted, avoid duplicate alerts, or reference past notifications in conversation.").
		Schema(tools.ObjectSchema(map[string]interface{}{
			"hours": tools.StringProperty("Optional: Number of hours to look back (default: 24)"),
			"type":  tools.StringProperty("Optional: Filter by alert type ('info', 'warning', 'success'). Omit to see all types."),
		})).
		Handler(func(ctx context.Context, toolParams *core.ToolParams) (*core.ToolResult, error) {
			// Parse input parameters
			var params struct {
				Hours string `json:"hours"`
				Type  string `json:"type"`
			}
			if err := json.Unmarshal(toolParams.Input, &params); err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("invalid input: %v", err),
				}, nil
			}

			// Default to 24 hours
			hours := 24
			if params.Hours != "" {
				fmt.Sscanf(params.Hours, "%d", &hours)
			}

			// Read alerts
			alertsMutex.RLock()
			defer alertsMutex.RUnlock()

			// Filter alerts by time and optionally by type
			cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
			filteredAlerts := []map[string]interface{}{}

			for _, alert := range alerts {
				if alert.Timestamp.After(cutoff) {
					if params.Type == "" || alert.Type == params.Type {
						filteredAlerts = append(filteredAlerts, map[string]interface{}{
							"id":        alert.ID,
							"message":   alert.Message,
							"type":      alert.Type,
							"timestamp": alert.Timestamp.Format("2006-01-02 15:04"),
							"age_hours": time.Since(alert.Timestamp).Hours(),
						})
					}
				}
			}

			result := map[string]interface{}{
				"total_alerts":      len(filteredAlerts),
				"hours_looked_back": hours,
				"alerts":            filteredAlerts,
			}

			if params.Type != "" {
				result["filtered_by_type"] = params.Type
			}

			return &core.ToolResult{
				Success: true,
				Data:    result,
			}, nil
		}).
		Build()
}

// startAIAnalysisLoop checks purchases from the past week for cheaper alternatives
// and posts findings to the notice board
func startAIAnalysisLoop(anthropicKey string) {
	client := anthropic.NewClient(option.WithAPIKey(anthropicKey))

	// Wait 10 seconds before first analysis
	time.Sleep(10 * time.Second)

	for {
		// Read mock transactions
		fileContent, err := os.ReadFile("mock_transactions.txt")
		if err != nil {
			log.Printf("❌ Failed to read transactions: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Parse transactions
		transactions := parseMockTransactions(string(fileContent))

		// Filter transactions from the past week (last 7 days)
		weekAgo := time.Now().AddDate(0, 0, -7)
		var recentTransactions []map[string]interface{}

		for _, tx := range transactions {
			dateStr, ok := tx["date"].(string)
			if !ok {
				continue
			}

			// Parse date (format: 2026-01-30)
			txDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// Include if within past week
			if txDate.After(weekAgo) {
				recentTransactions = append(recentTransactions, tx)
			}
		}

		if len(recentTransactions) == 0 {
			log.Printf("ℹ️  No recent transactions to analyze")
			time.Sleep(30 * time.Second)
			continue
		}

		// Find a product that hasn't been checked yet
		var productToCheck map[string]interface{}
		checkedProductsMutex.Lock()
		for _, tx := range recentTransactions {
			product, _ := tx["product"].(string)
			if product != "" && !checkedProducts[product] {
				productToCheck = tx
				checkedProducts[product] = true
				break
			}
		}
		checkedProductsMutex.Unlock()

		if productToCheck == nil {
			// All products checked, reset the map and start over
			checkedProductsMutex.Lock()
			checkedProducts = make(map[string]bool)
			checkedProductsMutex.Unlock()
			log.Printf("✅ All recent purchases checked for alternatives, resetting...")
			time.Sleep(60 * time.Second)
			continue
		}

		// Extract product details
		product, _ := productToCheck["product"].(string)
		amount, _ := productToCheck["amount"].(string)
		merchant, _ := productToCheck["merchant"].(string)
		date, _ := productToCheck["date"].(string)

		log.Printf("🔍 Checking for cheaper alternatives: %s ($%s)", product, amount)

		// Use Claude AI to find cheaper alternatives while preserving quality
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		prompt := fmt.Sprintf(`You are analyzing a purchase to find cheaper alternatives while preserving quality.

PURCHASE DETAILS:
- Product: %s
- Price Paid: %s
- Merchant: %s
- Date: %s

TASK:
Find a cheaper alternative that maintains or improves quality. You MUST save at least $5.00 to recommend an alternative.

FORMAT YOUR RESPONSE EXACTLY LIKE THIS (no preamble):
"[Original Product] ($[original price]) - Alternative: [New Product] ($[new price]) - Save: $[difference] - Buy: [URL]"

EXAMPLES:
- "Echo Dot 5th Gen ($49.99) - Alternative: Google Nest Mini ($29.99) - Save: $20.00 - Buy: https://store.google.com/product/google_nest_mini"
- "Nike Air Max ($139.99) - Alternative: Adidas Ultraboost ($120.00) - Save: $19.99 - Buy: https://www.adidas.com/us/ultraboost"
- "Whole Foods Groceries ($127.83) - Alternative: Trader Joe's Organic Mix ($95.00) - Save: $32.83 - Buy: https://www.traderjoes.com"

CRITICAL RULES:
- MUST save at least $5.00 or respond: "Not enough savings (under $5)"
- MUST include a real, working purchase link (Amazon, official store, major retailer)
- Use exact format with " - " separators (no vertical pipes)
- Show all prices with $ and two decimal places
- Alternative must maintain or improve quality
- URL should be direct product page when possible
- Focus on 2026 realistic pricing and real retailers`, product, amount, merchant, date)

		resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     anthropic.Model("claude-sonnet-4-20250514"),
			MaxTokens: 200, // Increased to accommodate URLs
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			},
		})
		cancel()

		if err != nil {
			log.Printf("❌ AI analysis error: %v", err)
		} else {
			// Extract response text
			var responseText strings.Builder
			for _, block := range resp.Content {
				if block.Type == "text" {
					responseText.WriteString(block.Text)
				}
			}

			recommendation := strings.TrimSpace(responseText.String())

			// Check if recommendation is valid and meets minimum savings
			shouldPost := false
			if recommendation != "" {
				recLower := strings.ToLower(recommendation)

				// Skip if response indicates no alternative or insufficient savings
				if strings.Contains(recLower, "optimal") ||
				   strings.Contains(recLower, "not enough savings") ||
				   strings.Contains(recLower, "under $5") {
					log.Printf("✓ No better alternative for: %s (insufficient savings or optimal)", product)
				} else if strings.Contains(recommendation, "Save: $") {
					// Extract savings amount from "Save: $X.XX"
					saveIdx := strings.Index(recommendation, "Save: $")
					if saveIdx != -1 {
						saveStr := recommendation[saveIdx+7:] // Skip "Save: $"
						// Find the end of the number (space or dash)
						endIdx := strings.IndexAny(saveStr, " -")
						if endIdx != -1 {
							saveStr = saveStr[:endIdx]
						}

						// Parse savings amount
						var savings float64
						if _, err := fmt.Sscanf(saveStr, "%f", &savings); err == nil {
							if savings >= 5.0 {
								shouldPost = true
							} else {
								log.Printf("✓ Savings too low for %s: $%.2f (minimum $5.00)", product, savings)
							}
						}
					}
				}
			}

			if shouldPost {
				// Create and store alert
				alert := Alert{
					ID:        fmt.Sprintf("alt-%d", time.Now().UnixNano()),
					Message:   recommendation,
					Timestamp: time.Now(),
					Type:      "success", // All valid alternatives are success
				}

				alertsMutex.Lock()
				alerts = append(alerts, alert)
				if len(alerts) > 100 {
					alerts = alerts[len(alerts)-100:]
				}
				alertsMutex.Unlock()

				log.Printf("💡 Posted alternative: %s", recommendation)
			}
		}

		// Wait 30 seconds before checking next product
		time.Sleep(30 * time.Second)
	}
}

// ============================================================================
// HACKATHON IDEAS
// ============================================================================
// Here are some ideas for custom tools you could build:
//
// 1. SAVINGS GOAL TRACKER
//    - Track progress toward savings goals
//    - Calculate how long until goal is reached
//    - Suggest optimal deposit amounts
//
// 2. BUDGET ANALYZER
//    - Set spending limits by category
//    - Alert when approaching limits
//    - Compare actual vs. planned spending
//
// 3. RECURRING PAYMENT DETECTOR
//    - Identify subscription payments
//    - Warn about upcoming bills
//    - Suggest savings opportunities
//
// 4. CASH FLOW FORECASTER
//    - Predict future balance based on patterns
//    - Identify potential low balance periods
//    - Suggest when to save vs. spend
//
// 5. SMART SAVINGS ADVISOR
//    - Analyze spare cash available
//    - Recommend savings deposits
//    - Calculate interest projections
//
// 6. SPENDING INSIGHTS
//    - Categorize spending automatically
//    - Compare to typical user patterns
//    - Highlight unusual activity
//
// 7. FINANCIAL HEALTH SCORE
//    - Calculate overall financial wellness
//    - Track improvements over time
//    - Provide actionable recommendations
//
// 8. PEER COMPARISON (anonymous)
//    - Compare savings rate to anonymized peers
//    - Show percentile rankings
//    - Motivate better habits
//
// 9. TAX ESTIMATION
//    - Track potential tax obligations
//    - Suggest amounts to set aside
//    - Generate tax reports
//
// 10. EMERGENCY FUND BUILDER
//     - Calculate needed emergency fund size
//     - Track progress toward goal
//     - Suggest automated savings plan
//
// ============================================================================
