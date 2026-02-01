// Hackathon Starter: Complete AI Financial Agent
// Build intelligent financial tools with nim-go-sdk + Liminal banking APIs
package main

import (
	"log"
	"sync"

	"github.com/becomeliminal/nim-go-sdk/executor"
	"github.com/becomeliminal/nim-go-sdk/server"
	"github.com/becomeliminal/nim-go-sdk/tools"
)

// Global state for alerts and product tracking
var (
	alertsMutex          sync.RWMutex
	alerts               []Alert
	checkedProductsMutex sync.RWMutex
	checkedProducts      map[string]bool // Track which products have been checked for alternatives
	checkedLargeTransactions sync.RWMutex    
	largeTransactionsSeen    map[string]bool
)

func main() {
	// Initialize global tracking map
	checkedProducts = make(map[string]bool)
	largeTransactionsSeen = make(map[string]bool)

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Setup Liminal executor
	liminalExecutor := executor.NewHTTPExecutor(executor.HTTPExecutorConfig{
		BaseURL: cfg.LiminalBaseURL,
	})
	log.Println("✅ Liminal API configured")

	// Create nim-go-sdk server
	srv, err := server.New(server.Config{
		AnthropicKey:    cfg.AnthropicKey,
		SystemPrompt:    hackathonSystemPrompt,
		Model:           ClaudeModel,
		MaxTokens:       DefaultMaxTokens,
		LiminalExecutor: liminalExecutor,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add Liminal banking tools
	srv.AddTools(tools.LiminalTools(liminalExecutor)...)
	log.Println("✅ Added 8 Liminal banking tools (get_transactions disabled - using mock data)")

	// Add custom tools
	srv.AddTool(createSpendingAnalyzerTool(liminalExecutor))
	log.Println("✅ Added custom spending analyzer tool")

	srv.AddTool(createProductAnalyzerTool(liminalExecutor, cfg.AnthropicKey))
	log.Println("✅ Added AI-powered product analyzer tool")

	srv.AddTool(createMockTransactionReaderTool())
	log.Println("✅ Added mock transaction reader tool")

	srv.AddTool(createProductSearchTool(cfg.AnthropicKey))
	log.Println("✅ Added product search tool for finding alternatives")

	srv.AddTool(createAlertTool())
	srv.AddTool(createReadAlertsTool())
	log.Println("✅ Added alert notification tools (post & read) for AI insights")

	// Start AI background analysis loop
	go startAIAnalysisLoop(cfg.AnthropicKey)
	log.Println("✅ Started AI background analysis loop (runs every 5 seconds)")

	go startLargeTransactionMonitor()
	log.Println("✅ Started large transaction monitor (checks for $1000+ transactions)")

	// For detecting recurring payments given transactions
	go detectRecurringPayments(cfg.AnthropicKey)

	// Setup HTTP endpoints
	setupHTTPHandlers()

	// Print startup information
	printStartupBanner(cfg.Port)

	// Start server
	if err := srv.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// printStartupBanner displays server information on startup
func printStartupBanner(port string) {
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
}
