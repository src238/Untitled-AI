import React, { useState, useEffect, useRef } from 'react'
import ReactDOM from 'react-dom/client'
import ReactMarkdown from 'react-markdown'
//Mode Components:
import RecentTransactions from './src/components/RecentTransactions'
import RecurringPayments from './src/components/RecurringPayments'
import GraphAnalysis from './src/components/GraphAnalysis'
import BudgetPlanner from './src/components/BudgetPlanner'

import '@liminalcash/nim-chat/styles.css'

import './styles.css'

interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  isIncoming: boolean
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

  const [alertFilter, setAlertFilter] = useState<'all' | 'info' | 'warning' | 'success'>('all')

  const [currentMode, setCurrentMode] = useState('recent-transactions') //either 'recent-transactions', 'graph-analysis', 'recurring-payments' or 'budget-planner'

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
          setAlerts(mappedAlerts.reverse())
        }
      } catch (error) {
        console.error('Failed to fetch alerts:', error)
      }
    }

    fetchAlerts()
    const interval = setInterval(fetchAlerts, 5000)
    return () => clearInterval(interval)
  }, [apiBaseUrl])

  useEffect(() => {
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      console.log('Connected to AI assistant')
      // Initialize conversation with the nim-go-sdk server
      ws.send(JSON.stringify({
        type: 'new_conversation',
        user: 'user'
      }))
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        if (data.type === 'text' && data.content) {
          // Check if we already have this message from streaming chunks
          setMessages(prev => {
            const lastMsg = prev[prev.length - 1]
            // If last message is from assistant and matches this content, skip (already streamed)
            if (lastMsg && lastMsg.role === 'assistant' && lastMsg.content === data.content) {
              return prev
            }
            // Otherwise add it as a new message
            return [...prev, { role: 'assistant', content: data.content }]
          })
        } else if (data.type === 'text_chunk' && data.content) {
          // Handle streaming text chunks
          setMessages(prev => {
            const newMessages = [...prev]
            if (newMessages.length > 0 && newMessages[newMessages.length - 1].role === 'assistant') {
              // Append to last assistant message
              newMessages[newMessages.length - 1].content += data.content
            } else {
              // Create new assistant message
              newMessages.push({ role: 'assistant', content: data.content })
            }
            return newMessages
          })
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
const filteredAlerts = alertFilter === 'all' 
  ? alerts 
  : alerts.filter(alert => alert.type === alertFilter)

const alertCounts = {
  all: alerts.length,
  info: alerts.filter(a => a.type === 'info').length,
  warning: alerts.filter(a => a.type === 'warning').length,
  success: alerts.filter(a => a.type === 'success').length,
}

  return (
    <div className="dashboard-container">
      {/* Left Sidebar - AI Alerts */}
      <aside className="alerts-sidebar">
        <div className="sidebar-header">
          <h2>AI Insights</h2>
          <div className="pulse-indicator"></div>
        </div>
	<div className="alert-filters">
  	  <button
    	    className={`filter-btn ${alertFilter === 'all' ? 'active' : ''}`}
    	    onClick={() => setAlertFilter('all')}
  	  >
    	    All <span className="filter-count">{alertCounts.all}</span>
  	  </button>
  	  <button
    	    className={`filter-btn filter-warning ${alertFilter === 'warning' ? 'active' : ''}`}
    	    onClick={() => setAlertFilter('warning')}
  	  >
    	    Warnings <span className="filter-count">{alertCounts.warning}</span>
  	  </button>
  	  <button
    	    className={`filter-btn filter-success ${alertFilter === 'success' ? 'active' : ''}`}
    	    onClick={() => setAlertFilter('success')}
  	  >
    	    Savings <span className="filter-count">{alertCounts.success}</span>
  	  </button>
  	  <button
    	    className={`filter-btn filter-info ${alertFilter === 'info' ? 'active' : ''}`}
    	    onClick={() => setAlertFilter('info')}
  	  >
    	    Info <span className="filter-count">{alertCounts.info}</span>
  	  </button>
	</div>
        <div className="alerts-container">
          {alerts.length === 0 ? (
            <div className="empty-state">
              <p>No alerts yet. I'll notify you of important insights.</p>
            </div>
          ) : (
            alerts.map(alert => (
              <div key={alert.id} className={`alert-card alert-${alert.type}`}>
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
        <div className="balance-info">
            <span className="balance-label">Current Balance</span>
            <h1 className="balance-amount">
              ${balance.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </h1>
          </div>
        <div className="mode-buttons">
          <button
            className={`action-btn ${currentMode === 'recent-transactions' ? 'active' : 'inactive'}`}
            onClick={() => setCurrentMode('recent-transactions')}
          >
            Recent Transactions
          </button>
          <button
            className={`action-btn ${currentMode === 'graph-analysis' ? 'active' : 'inactive'}`}
            onClick={() => setCurrentMode('graph-analysis')}
          >
            Graph Analysis
          </button>
          <button
            className={`action-btn ${currentMode === 'recurring-payments' ? 'active' : 'inactive'}`}
            onClick={() => setCurrentMode('recurring-payments')}
          >
            Recurring Payments
          </button>
          <button
            className={`action-btn ${currentMode === 'budget-planner' ? 'active' : 'inactive'}`}
            onClick={() => setCurrentMode('budget-planner')}
          >
            Budget Planner
          </button>
        </div>
        {currentMode === 'recent-transactions' && (
          <RecentTransactions transactions={transactions} balance={balance} />
        )}
        {currentMode === 'graph-analysis' && (
          <GraphAnalysis transactions={transactions} balance={balance} />
        )}
        {currentMode === 'recurring-payments' && (
          <RecurringPayments transactions={transactions} balance={balance} />
        )}
        {currentMode === 'budget-planner' && (
          <BudgetPlanner transactions={transactions} balance={balance} />
        )}
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
