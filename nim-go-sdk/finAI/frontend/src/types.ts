// types.ts - TypeScript type definitions
export interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  isIncoming: boolean
  merchant?: string
}

export interface Alert {
  id: string
  message: string
  timestamp: string
  type: 'info' | 'warning' | 'success'
}

export interface Message {
  role: 'user' | 'assistant'
  content: string
}

export interface WebSocketMessage {
  type: string
  content?: string
  message?: string
  severity?: string
  balance?: number
  transaction?: Transaction
}
