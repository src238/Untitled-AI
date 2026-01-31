# FinAI Usage Guide

## New Feature: AI-Powered Product Analysis

FinAI now includes an intelligent product analyzer that uses Claude AI to identify what products and services were purchased from your transaction history.

### How It Works

1. The tool fetches your recent transactions from Liminal
2. Extracts relevant data (merchant names, descriptions, amounts, memos)
3. Sends this data to Claude AI for intelligent analysis
4. Returns categorized insights about your purchases

### Using the Product Analyzer

Simply ask Nim in natural language:

**Basic queries:**
- "What products did I buy recently?"
- "Analyze my purchases from the last 30 days"
- "Show me what I've been spending money on"
- "What did I buy this month?"

**With parameters:**
- "Analyze my last 20 transactions for products" (limits to 20 transactions)
- "What products did I purchase in the last 7 days?" (analyzes 7 day window)

### What You'll Get

The AI analysis includes:

**Purchase Categories:**
- Breakdown by category (groceries, dining, entertainment, etc.)
- Number of transactions per category
- Description of what was purchased

**Specific Products:**
- Identified products/services
- Merchant names
- Transaction amounts

**Insights:**
- Purchasing patterns
- Most frequent purchase types
- Recommendations and observations

### Technical Details

**Tool Name:** `analyze_products`

**Parameters:**
- `days` (optional): Number of days to analyze (default: 30)
- `limit` (optional): Max transactions to analyze (default: 50)

**Example Response:**
```json
{
  "period_days": 30,
  "transactions_analyzed": 45,
  "product_analysis": "**Purchase Categories:**\n- Dining: 12 transactions - Restaurant meals and food delivery\n- Groceries: 8 transactions - Weekly shopping at supermarkets\n...",
  "generated_at": "2026-01-31T16:30:00Z"
}
```

### How This Differs from `analyze_spending`

- **`analyze_spending`**: Provides numerical analysis (amounts, velocity, counts)
- **`analyze_products`**: Uses AI to understand *what* was purchased and provides semantic insights

### Architecture

The product analyzer:
1. Uses the Liminal `get_transactions` tool to fetch data
2. Creates a Claude API client with your Anthropic key
3. Formats transaction data into a structured prompt
4. Calls Claude 3.5 Sonnet for analysis
5. Returns the AI-generated insights

### Performance

- Analyzes up to 50 transactions per request (configurable)
- Uses Claude 3.5 Sonnet (fast and accurate)
- Typical response time: 2-4 seconds
- Token usage: ~1000-2000 tokens per analysis

### Privacy & Security

- All transaction data is processed via Anthropic's API
- No transaction data is stored beyond the analysis
- Analysis is ephemeral and user-specific
- Follows Anthropic's data handling policies

## Example Conversations

**User:** "What products have I been buying?"

**Nim:** *[Calls analyze_products tool]*

"Based on your recent transactions, here's what you've been purchasing:

**Purchase Categories:**
- Dining & Food Delivery: 8 transactions - Mix of restaurants and food delivery services
- Digital Services: 5 transactions - Streaming subscriptions and software
- Groceries: 4 transactions - Weekly shopping trips
- Transportation: 3 transactions - Ride-sharing services

**Specific Products:**
- Netflix subscription - $15.99
- Spotify Premium - $10.99
- Uber rides (multiple) - $87.50 total
- Whole Foods groceries - $156.32

**Insights:**
- Your dining spending has increased 20% compared to typical patterns
- You have 3 active subscription services
- Transportation costs are primarily weekend rides
- Consider meal planning to reduce dining expenses"

---

**User:** "Analyze my last 10 purchases"

**Nim:** *[Calls analyze_products with limit=10]*

"Looking at your last 10 transactions, here's what I found:

[AI-generated analysis of the specific transactions...]"

## Troubleshooting

### "AI analysis failed" error
- Check your `ANTHROPIC_API_KEY` is valid
- Ensure you have API credits remaining
- Verify network connectivity

### Limited transaction data
- The tool can only analyze transactions returned by Liminal
- Some P2P transfers may have limited merchant info
- AI makes reasonable inferences when data is sparse

### Want more detailed analysis?
- Increase the `limit` parameter for more transactions
- Use a longer time window with the `days` parameter
- Combine with `analyze_spending` for numerical context
