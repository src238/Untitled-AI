// tools.go - Custom tool creation functions
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/becomeliminal/nim-go-sdk/core"
	"github.com/becomeliminal/nim-go-sdk/tools"
)
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
			txRequestJSON, err := json.Marshal(txRequest)
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("Error marshaling transaction request: %v", err),
				}, nil
			}

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
		txType, ok := tx["type"].(string)
		if !ok {
			continue // Skip transactions with invalid type
		}
		amount, ok := tx["amount"].(float64)
		if !ok {
			continue // Skip transactions with invalid amount
		}

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
			txRequestJSON, err := json.Marshal(txRequest)
			if err != nil {
				return &core.ToolResult{
					Success: false,
					Error:   fmt.Sprintf("Error marshaling transaction request: %v", err),
				}, nil
			}

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
				if _, err := fmt.Sscanf(params.Hours, "%d", &hours); err != nil {
					// Log error but use default value
					log.Printf("Warning: Failed to parse hours '%s': %v, using default 24", params.Hours, err)
				}
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
