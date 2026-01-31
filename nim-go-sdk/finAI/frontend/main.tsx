import React, { useState, useEffect, useRef } from 'react'
import ReactDOM from 'react-dom/client'
import ReactMarkdown from 'react-markdown'
import '@liminalcash/nim-chat/styles.css'
import './styles.css'

interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  type: 'debit' | 'credit'
  merchant?: string
}

interface Alert {
  id: string
  message: string
  timestamp: string
  type: 'info' | 'warning' | 'success'
}

interface Message {
  role: 'user' | 'assistant'
  content: string
}

function App() {
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
  const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const [balance, setBalance] = useState<number>(2547.83)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [alerts, setAlerts] = useState<Alert[]>([])

  const [messages, setMessages] = useState<Message[]>([
    { role: 'assistant', content: 'Hello! I\'m your AI financial assistant. How can I help you today?' }
  ])
  const [inputValue, setInputValue] = useState('')
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  // Fetch transactions on mount
  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const response = await fetch(`${apiBaseUrl}/api/transactions`)
        if (response.ok) {
          const data = await response.json()
          setTransactions(data)
        }
      } catch (error) {
        console.error('Failed to fetch transactions:', error)
      }
    }

    fetchTransactions()
  }, [apiBaseUrl])

  // Poll alerts API every 5 seconds
  useEffect(() => {
    const fetchAlerts = async () => {
      try {
        const response = await fetch(`${apiBaseUrl}/api/alerts`)
        if (response.ok) {
          const data = await response.json()
          // Map API response to Alert format
          const mappedAlerts = data.map((alert: any) => ({
            id: alert.id,
            message: alert.message,
            timestamp: new Date(alert.timestamp).toLocaleString('en-US', {
              year: 'numeric',
              month: '2-digit',
              day: '2-digit',
              hour: '2-digit',
              minute: '2-digit'
            }),
            type: alert.type
          }))
          // Reverse to show newest alerts at the top
          setAlerts(mappedAlerts.reverse())
        }
      } catch (error) {
        console.error('Failed to fetch alerts:', error)
      }
    }

    // Initial fetch
    fetchAlerts()

    // Poll every 5 seconds
    const interval = setInterval(fetchAlerts, 5000)

    return () => clearInterval(interval)
  }, [apiBaseUrl])

  useEffect(() => {
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      console.log('Connected to AI assistant')
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        if (data.type === 'message' && data.content) {
          setMessages(prev => [...prev, { role: 'assistant', content: data.content }])
        } else if (data.type === 'alert') {
          const newAlert: Alert = {
            id: Date.now().toString(),
            message: data.message,
            timestamp: new Date().toLocaleString('en-US', {
              year: 'numeric',
              month: '2-digit',
              day: '2-digit',
              hour: '2-digit',
              minute: '2-digit'
            }),
            type: data.severity || 'info'
          }
          setAlerts(prev => [newAlert, ...prev])
        } else if (data.type === 'balance_update') {
          setBalance(data.balance)
        } else if (data.type === 'transaction_update') {
          setTransactions(prev => [data.transaction, ...prev])
        }
      } catch (e) {
        console.error('Failed to parse message:', e)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      setIsConnected(false)
    }

    ws.onclose = () => {
      setIsConnected(false)
      console.log('Disconnected from AI assistant')
    }

    return () => {
      ws.close()
    }
  }, [wsUrl])

  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault()
    if (!inputValue.trim() || !wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return

    const userMessage: Message = { role: 'user', content: inputValue }
    setMessages(prev => [...prev, userMessage])

    wsRef.current.send(JSON.stringify({
      type: 'message',
      content: inputValue
    }))

    setInputValue('')
  }

  // Helper function to render alert message with clickable links
  const renderAlertMessage = (message: string) => {
    const urlRegex = /(https?:\/\/[^\s]+)/g
    const parts = message.split(urlRegex)

    return parts.map((part, index) => {
      if (part.match(urlRegex)) {
        return (
          <a
            key={index}
            href={part}
            target="_blank"
            rel="noopener noreferrer"
            className="alert-link"
          >
            {part}
          </a>
        )
      }
      return <span key={index}>{part}</span>
    })
  }

  return (
    <div className="dashboard-container">
      {/* Left Sidebar - AI Alerts */}
      <aside className="alerts-sidebar">
        <div className="sidebar-header">
          <h2>AI Insights</h2>
          <div className="pulse-indicator"></div>
        </div>
        <div className="alerts-container">
          {alerts.length === 0 ? (
            <div className="empty-state">
              <p>No alerts yet. I'll notify you of important insights.</p>
            </div>
          ) : (
            alerts.map(alert => (
              <div key={alert.id} className={`alert-card alert-${alert.type}`}>
                <div className="alert-icon">
                  {alert.type === 'warning' && '⚠'}
                  {alert.type === 'success' && '✓'}
                  {alert.type === 'info' && 'ℹ'}
                </div>
                <div className="alert-content">
                  <p className="alert-message">{renderAlertMessage(alert.message)}</p>
                  <time className="alert-time">{alert.timestamp}</time>
                </div>
              </div>
            ))
          )}
        </div>
      </aside>

      {/* Center - Balance & Transactions */}
      <main className="main-content">
        <header className="balance-header">
          <div className="balance-info">
            <span className="balance-label">Current Balance</span>
            <h1 className="balance-amount">
              ${balance.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </h1>
          </div>
          <div className="balance-actions">
            <button className="action-btn action-send">Send</button>
            <button className="action-btn action-receive">Receive</button>
          </div>
        </header>

        <section className="transactions-section">
          <h3 className="section-title">Recent Transactions</h3>
          <div className="transactions-list">
            {transactions.map(tx => (
              <div key={tx.id} className="transaction-item">
                <div className="transaction-icon">
                  {tx.type === 'credit' ? '↑' : '↓'}
                </div>
                <div className="transaction-details">
                  <span className="transaction-description">{tx.description}</span>
                  <span className="transaction-date">{tx.date}</span>
                </div>
                <span className={`transaction-amount ${tx.type === 'credit' ? 'credit' : 'debit'}`}>
                  {tx.type === 'credit' ? '+' : '-'}${Math.abs(tx.amount).toFixed(2)}
                </span>
              </div>
            ))}
          </div>
        </section>
      </main>

      {/* Right Sidebar - Chat */}
      <aside className="chat-sidebar">
        <div className="chat-header">
          <h2>AI Assistant</h2>
          <div className={`connection-status ${isConnected ? 'connected' : 'disconnected'}`}>
            {isConnected ? '● Online' : '○ Offline'}
          </div>
        </div>
        <div className="chat-messages">
          {messages.map((msg, idx) => (
            <div key={idx} className={`chat-message ${msg.role}`}>
              <div className="message-avatar">
                {msg.role === 'assistant' ? 'AI' : 'You'}
              </div>
              <div className="message-bubble">
                <ReactMarkdown>{msg.content}</ReactMarkdown>
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
        <form className="chat-input-form" onSubmit={sendMessage}>
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Ask me anything..."
            className="chat-input"
            disabled={!isConnected}
          />
          <button
            type="submit"
            className="chat-send-btn"
            disabled={!isConnected || !inputValue.trim()}
          >
            →
          </button>
        </form>
      </aside>
    </div>
  )
}

ReactDOM.createRoot(document.getElementById('root')!).render(<App />)
