// useWebSocket.ts - Custom hook for WebSocket connection management
import { useEffect, useRef, useState } from 'react'
import { Message, WebSocketMessage, Alert } from '../types'
import { formatTimestamp } from '../utils/formatters.tsx'

interface UseWebSocketProps {
  wsUrl: string
  onBalanceUpdate?: (balance: number) => void
  onTransactionUpdate?: (transaction: any) => void
  onAlertReceived?: (alert: Alert) => void
}

interface UseWebSocketReturn {
  messages: Message[]
  isConnected: boolean
  sendMessage: (content: string) => void
  wsRef: React.MutableRefObject<WebSocket | null>
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>
}

export function useWebSocket({
  wsUrl,
  onBalanceUpdate,
  onTransactionUpdate,
  onAlertReceived
}: UseWebSocketProps): UseWebSocketReturn {
  const [messages, setMessages] = useState<Message[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

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
        const data: WebSocketMessage = JSON.parse(event.data)
        handleWebSocketMessage(data)
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

  const handleWebSocketMessage = (data: WebSocketMessage) => {
    if (data.type === 'text' && data.content) {
      handleTextMessage(data.content)
    } else if (data.type === 'text_chunk' && data.content) {
      handleTextChunk(data.content)
    } else if (data.type === 'alert' && data.message) {
      handleAlert(data)
    } else if (data.type === 'balance_update' && data.balance !== undefined) {
      onBalanceUpdate?.(data.balance)
    } else if (data.type === 'transaction_update' && data.transaction) {
      onTransactionUpdate?.(data.transaction)
    }
  }

  const handleTextMessage = (content: string) => {
    setMessages(prev => {
      const lastMsg = prev[prev.length - 1]
      // If last message is from assistant and matches this content, skip (already streamed)
      if (lastMsg && lastMsg.role === 'assistant' && lastMsg.content === content) {
        return prev
      }
      // Otherwise add it as a new message
      return [...prev, { role: 'assistant', content }]
    })
  }

  const handleTextChunk = (content: string) => {
    setMessages(prev => {
      const newMessages = [...prev]
      if (newMessages.length > 0 && newMessages[newMessages.length - 1].role === 'assistant') {
        // Append to last assistant message
        newMessages[newMessages.length - 1].content += content
      } else {
        // Create new assistant message
        newMessages.push({ role: 'assistant', content })
      }
      return newMessages
    })
  }

  const handleAlert = (data: WebSocketMessage) => {
    if (!data.message) return

    const newAlert: Alert = {
      id: Date.now().toString(),
      message: data.message,
      timestamp: formatTimestamp(new Date()),
      type: data.severity as Alert['type'] || 'info'
    }
    onAlertReceived?.(newAlert)
  }

  const sendMessage = (content: string) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return

    const userMessage: Message = { role: 'user', content }
    setMessages(prev => [...prev, userMessage])

    wsRef.current.send(JSON.stringify({
      type: 'message',
      content
    }))
  }

  return {
    messages,
    isConnected,
    sendMessage,
    wsRef,
    setMessages
  }
}
