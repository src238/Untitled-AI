// config.go - Application configuration and constants
package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Application constants
const (
	DefaultPort            = "8080"
	DefaultLiminalBaseURL  = "https://api.liminal.cash"
	DefaultMaxTokens       = 4096
	ClaudeModel            = "claude-sonnet-4-20250514"

	// Analysis loop timing
	AnalysisInitialDelay   = 0 * time.Second
	AnalysisInterval       = 5 * time.Second
	AnalysisResetDelay     = 60 * time.Second
	AnalysisTimeout        = 60 * time.Second

	// Alert settings
	MaxAlertsStored        = 100
	AlertRetentionHours    = 24

	// Polling intervals
	AlertPollInterval      = 5 * time.Second

	// Product analysis
	MinimumSavings         = 5.0
	TransactionLookbackDays = 7

	// Mock data
	MockTransactionsFile   = "mock_transactions.txt"
)

// Config holds the application configuration
type Config struct {
	AnthropicKey   string
	LiminalBaseURL string
	Port           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (optional - will use system env vars if not found)
	_ = godotenv.Load()

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey == "" {
		log.Fatal("❌ ANTHROPIC_API_KEY environment variable is required")
	}

	liminalBaseURL := os.Getenv("LIMINAL_BASE_URL")
	if liminalBaseURL == "" {
		liminalBaseURL = DefaultLiminalBaseURL
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	return &Config{
		AnthropicKey:   anthropicKey,
		LiminalBaseURL: liminalBaseURL,
		Port:           port,
	}, nil
}
