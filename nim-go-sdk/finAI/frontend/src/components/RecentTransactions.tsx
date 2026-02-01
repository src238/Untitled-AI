import { useState } from 'react'

interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  isIncoming: boolean
  merchant?: string
}

type FilterType = 'all' | 'send' | 'receive'

interface RecentTransactionsProps {
  transactions: Transaction[]
  balance: number
}

export default function RecentTransactions({ transactions, balance }: RecentTransactionsProps) {
  const [activeFilter, setActiveFilter] = useState<FilterType>('all')

  const filteredTransactions = transactions.filter(tx => {
    if (activeFilter === 'receive') return tx.isIncoming
    if (activeFilter === 'send') return !tx.isIncoming
    return true
  })

  return (
    <>
      <header className="balance-header">
        <h3 className="section-title">Recent Transactions</h3>
        <div className="balance-actions">
          <button
            className={`action-btn action-all ${activeFilter === 'all' ? 'active' : 'inactive'}`}
            onClick={() => setActiveFilter('all')}
          >
            All
          </button>
          <button
            className={`action-btn action-send ${activeFilter === 'send' ? 'active' : 'inactive'}`}
            onClick={() => setActiveFilter('send')}
          >
            Send
          </button>
          <button
            className={`action-btn action-receive ${activeFilter === 'receive' ? 'active' : 'inactive'}`}
            onClick={() => setActiveFilter('receive')}
          >
            Receive
          </button>
        </div>
      </header>

      <section className="transactions-section">
        <div className="transactions-list">
          {filteredTransactions.length === 0 ? (
            <p className="empty-transactions">No {activeFilter !== 'all' ? activeFilter : ''} transactions yet.</p>
          ) : (
            filteredTransactions.map(tx => (
              <div key={tx.id} className="transaction-item">
                <div className="transaction-icon">
                  {tx.isIncoming ? '↑' : '↓'}
                </div>
                <div className="transaction-details">
                  <span className="transaction-description">{tx.description}</span>
                  <span className="transaction-date">{tx.date}</span>
                </div>
                <span className={`transaction-amount ${tx.isIncoming ? 'incoming' : 'outgoing'}`}>
                  {tx.isIncoming ? '+' : '-'}${Math.abs(tx.amount).toFixed(2)}
                </span>
              </div>
            ))
          )}
        </div>
      </section>
    </>
  )
}