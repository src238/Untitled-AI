# FinAI - AI-Powered Financial Assistant

Built with nim-go-sdk and Liminal banking APIs.

## Quick Start

### Prerequisites

- **Go 1.21+** installed
- **Node.js 18+** installed
- **Anthropic API key** from [console.anthropic.com](https://console.anthropic.com/)

### Setup

1. **Configure environment:**
```bash
cd nim-go-sdk/finAI
cp .env.example .env
# Edit .env and add your ANTHROPIC_API_KEY
```

2. **Start the backend:**
```bash
go mod tidy
go run main.go
```

You should see:
```
✅ Liminal API configured
✅ Added 9 Liminal banking tools
✅ Added custom spending analyzer tool
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚀 Hackathon Starter Server Running
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📡 WebSocket endpoint: ws://localhost:8080/ws
💚 Health check: http://localhost:8080/health
```

3. **Start the frontend (in a new terminal):**
```bash
cd frontend
npm install
npm run dev
```

Your browser will open to `http://localhost:5173` with the chat interface.

### Usage

1. Click the chat bubble
2. Login with your email (you'll receive an OTP code)
3. Enter the code to authenticate
4. Start chatting! Try:
   - "What's my balance?"
   - "Show me my recent transactions"
   - "Analyze my spending over the last 30 days"

## Features

### Built-in Banking Tools

- Check wallet and savings balances
- View transaction history
- Send money to other users
- Manage savings deposits and withdrawals
- Search for users
- View savings rates

### Custom Analytics

- **Spending pattern analysis** - Analyze spending velocity and trends
- **AI-Powered Product Analysis** - Uses Claude AI to identify what products/services were purchased from transaction data
- Transaction insights
- Financial recommendations

## Project Structure

```
finAI/
├── main.go              # Go backend server with AI agent
├── frontend/            # React chat interface
│   ├── main.tsx         # App entry point
│   ├── index.html       # Landing page
│   ├── styles.css       # Styling
│   └── fonts/           # Custom fonts
├── go.mod               # Go dependencies
├── .env.example         # Environment template
└── README.md            # This file
```

## Customization

Edit `main.go` to:
- Add custom tools and financial analysis features
- Customize the AI personality via `hackathonSystemPrompt`
- Configure banking operations

Edit `frontend/main.tsx` to customize the UI.

## License

MIT License
