// formatters.ts - Utility functions for formatting data
import { TIMESTAMP_FORMAT } from '../constants'

/**
 * Formats a timestamp to a localized string
 */
export function formatTimestamp(timestamp: Date | string): string {
  const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp
  return date.toLocaleString('en-US', TIMESTAMP_FORMAT)
}

/**
 * Formats a currency amount
 */
export function formatCurrency(amount: number): string {
  return amount.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

/**
 * Renders a message with clickable links
 * Extracts URLs from text and converts them to anchor tags
 */
export function renderLinksInText(text: string): (string | JSX.Element)[] {
  const urlRegex = /(https?:\/\/[^\s]+)/g
  const parts = text.split(urlRegex)

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
