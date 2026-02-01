// constants.ts - Application constants and configuration
export const API_CONFIG = {
  DEFAULT_WS_URL: 'ws://localhost:8080/ws',
  DEFAULT_API_URL: 'http://localhost:8080',
  ALERTS_POLL_INTERVAL: 5000, // 5 seconds
} as const

export const DEFAULT_BALANCE = 2547.83

export const INITIAL_MESSAGE = {
  role: 'assistant' as const,
  content: "Hello! I'm your AI financial assistant. How can I help you today?"
}

export const TIMESTAMP_FORMAT = {
  year: 'numeric' as const,
  month: '2-digit' as const,
  day: '2-digit' as const,
  hour: '2-digit' as const,
  minute: '2-digit' as const,
}
