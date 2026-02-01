import { useState, useEffect } from 'react'

interface RecurringPayment {
  id: string
  merchant: string
  product: string
  amount: number
  frequency: string
  lastSeen: string
  occurrences: number
}

interface RecurringPaymentsProps {
  transactions: any[]
  balance: number
}

export default function RecurringPayments({ }: RecurringPaymentsProps) {
  const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const [payments, setPayments] = useState<RecurringPayment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchRecurringPayments = async () => {
      try {
        const response = await fetch(`${apiBaseUrl}/api/recurring-payments`)
        if (response.ok) {
          const data = await response.json()
          setPayments(data)
          setLoading(false)
        } else if (response.status === 202) {
          setTimeout(fetchRecurringPayments, 3000)
          return
        } else {
          setError('Failed to fetch recurring payments')
          setLoading(false)
        }
      } catch (err) {
        setError('Could not connect to server')
        setLoading(false)
      }
    }

    fetchRecurringPayments()
  }, [apiBaseUrl])

  const monthlyTotal = payments.reduce((sum, p) => {
    switch (p.frequency) {
      case 'Weekly': return sum + p.amount * 4.33
      case 'Bi-weekly': return sum + p.amount * 2.17
      case 'Monthly': return sum + p.amount
      case 'Annual': return sum + p.amount / 12
      default: return sum + p.amount
    }
  }, 0)

  const formatDate = (dateStr: string) => {
    const [year, month, day] = dateStr.split('-')
    return `${month}/${day}/${year.slice(2)}`
  }

  if (loading) {
    return (
      <div className="recurring-loading">
        <div className="loading-spinner"></div>
        <p>AI is analyzing your transactions...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="recurring-error">
        <p>⚠ {error}</p>
      </div>
    )
  }

  return (
    <>
      <header className="balance-header">
        <div className="balance-info">
          <span className="balance-label">Est. Monthly Recurring</span>
          <h1 className="balance-amount">${monthlyTotal.toFixed(2)}</h1>
        </div>
        <div className="balance-actions">
          <span className="recurring-count-badge">{payments.length} active</span>
        </div>
      </header>

      <section className="transactions-section" style={{ marginTop: '3rem' }}>
        <h3 className="section-title">Recurring Payments</h3>
        <div className="transactions-list">
          {payments.length === 0 ? (
            <p className="empty-transactions">No recurring payments detected.</p>
          ) : (
            payments.map(payment => (
              <div key={payment.id} className="transaction-item recurring-item">
                <div className="transaction-icon">
                  🔄
                </div>
                <div className="transaction-details">
                  <div className="recurring-main-row">
                    <span className="transaction-description">{payment.merchant}</span>
                    <span className={`recurring-freq-badge freq-${payment.frequency.toLowerCase()}`}>
                      {payment.frequency}
                    </span>
                  </div>
                  <div className="recurring-sub-row">
                    <span className="recurring-sub-text">{payment.occurrences}x detected</span>
                    <span className="recurring-sub-divider">·</span>
                    <span className="recurring-sub-text">Last: {formatDate(payment.lastSeen)}</span>
                  </div>
                </div>
                <span className="transaction-amount outgoing">
                  ${payment.amount.toFixed(2)}
                </span>
              </div>
            ))
          )}
        </div>
      </section>
    </>
  )
}