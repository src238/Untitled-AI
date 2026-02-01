import { useState, useMemo, useEffect } from 'react'

interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  isIncoming: boolean
  merchant?: string
}

interface BudgetPlannerProps {
  balance: number
}

export default function BudgetPlanner({ balance }: BudgetPlannerProps) {
  const [budgetTarget, setBudgetTarget] = useState<string>('')
  const [mockTransactions, setMockTransactions] = useState<Transaction[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isFocused, setIsFocused] = useState(false)

  // Fetch transactions directly from mock_transactions.txt via API
  useEffect(() => {
    const fetchMockTransactions = async () => {
      try {
        const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
        const response = await fetch(`${apiBaseUrl}/api/transactions`)
        if (response.ok) {
          const data = await response.json()
          setMockTransactions(data)
        } else {
          console.error('Failed to fetch mock transactions:', response.status)
        }
      } catch (error) {
        console.error('Error fetching mock transactions:', error)
      } finally {
        setIsLoading(false)
      }
    }

    fetchMockTransactions()
  }, [])

  // Calculate month results - we have the COMPLETE month's data
  const monthResults = useMemo(() => {
    const startingBalance = 0
    const endingBalance = balance // $2547.83
    const profit = endingBalance - startingBalance

    // Calculate totals from ALL transactions (full month Jan 1-30)
    const incomeTransactions = mockTransactions.filter(t => t.isIncoming)
    const expenseTransactions = mockTransactions.filter(t => !t.isIncoming)

    const totalIncome = incomeTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0)
    const totalExpenses = expenseTransactions.reduce((sum, t) => sum + Math.abs(t.amount), 0)

    return {
      startingBalance,
      endingBalance,
      profit,
      totalIncome,
      totalExpenses,
      incomeCount: incomeTransactions.length,
      expenseCount: expenseTransactions.length
    }
  }, [mockTransactions, balance])

  // Calculate budget comparison
  const budgetComparison = useMemo(() => {
    if (!budgetTarget || isNaN(parseFloat(budgetTarget))) {
      return null
    }

    const target = parseFloat(budgetTarget)
    const difference = monthResults.endingBalance - target
    const achievedTarget = difference >= 0

    return {
      target,
      difference,
      achievedTarget
    }
  }, [budgetTarget, monthResults])

  if (isLoading) {
    return (
      <section className="transactions-section">
        <div style={{ padding: '20px', textAlign: 'center' }}>
          <p style={{ fontSize: '16px', color: '#666' }}>Loading...</p>
        </div>
      </section>
    )
  }

  return (
    <>
      <header className="balance-header">
        <h3 className="section-title">Savings Planner</h3>
      </header>

      <section className="transactions-section">
        <div className="transactions-list" style={{ maxHeight: 'none', padding: '0' }}>

          {/* Budget Target Input Card */}
          <div style={{
            padding: '2rem 2rem 1.5rem',
            borderBottom: '1px solid var(--light-border)',
            background: 'var(--white)'
          }}>
            <label htmlFor="budget-input" style={{
              display: 'block',
              fontFamily: 'ABC Marist, sans-serif',
              fontSize: '0.875rem',
              color: 'var(--brown)',
              textTransform: 'uppercase',
              letterSpacing: '0.05em',
              marginBottom: '0.75rem',
              fontWeight: '500'
            }}>
              Set Monthly Target
            </label>
            <div style={{ position: 'relative' }}>
              <span style={{
                position: 'absolute',
                left: '1rem',
                top: '50%',
                transform: 'translateY(-50%)',
                fontSize: '1.5rem',
                fontFamily: 'ABC Marist, sans-serif',
                color: isFocused || budgetTarget ? 'var(--orange)' : 'var(--beige)',
                transition: 'color 0.2s ease',
                pointerEvents: 'none',
                fontWeight: '500'
              }}>$</span>
              <input
                id="budget-input"
                type="number"
                value={budgetTarget}
                onChange={(e) => setBudgetTarget(e.target.value)}
                onFocus={() => setIsFocused(true)}
                onBlur={() => setIsFocused(false)}
                placeholder="2000"
                style={{
                  width: '100%',
                  padding: '1rem 1rem 1rem 2.5rem',
                  fontSize: '1.5rem',
                  fontFamily: 'ABC Marist, sans-serif',
                  fontWeight: '500',
                  border: `2px solid ${isFocused ? 'var(--orange)' : 'var(--light-border)'}`,
                  borderRadius: '8px',
                  background: 'var(--cream)',
                  color: 'var(--black)',
                  transition: 'all 0.2s ease',
                  outline: 'none',
                  letterSpacing: '-0.02em'
                }}
              />
            </div>
          </div>

          {/* Financial Overview Grid */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(3, 1fr)',
            borderBottom: '1px solid var(--light-border)'
          }}>
            {/* Start Balance */}
            <div style={{
              padding: '2rem',
              borderRight: '1px solid var(--light-border)',
              background: 'var(--white)',
              transition: 'all 0.2s ease'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--cream)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'var(--white)'}>
              <div style={{
                fontSize: '0.6875rem',
                color: 'var(--beige)',
                marginBottom: '0.75rem',
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                fontWeight: '500'
              }}>January 1st</div>
              <div style={{
                fontSize: '2rem',
                fontWeight: '500',
                fontFamily: 'ABC Marist, sans-serif',
                color: 'var(--black)',
                letterSpacing: '-0.03em',
                lineHeight: '1'
              }}>${monthResults.startingBalance.toFixed(2)}</div>
              <div style={{
                fontSize: '0.75rem',
                color: 'var(--beige)',
                marginTop: '0.5rem',
                letterSpacing: '0.01em'
              }}>Starting balance</div>
            </div>

            {/* End Balance */}
            <div style={{
              padding: '2rem',
              borderRight: '1px solid var(--light-border)',
              background: 'var(--white)',
              transition: 'all 0.2s ease'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--cream)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'var(--white)'}>
              <div style={{
                fontSize: '0.6875rem',
                color: 'var(--beige)',
                marginBottom: '0.75rem',
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                fontWeight: '500'
              }}>January 30th</div>
              <div style={{
                fontSize: '2rem',
                fontWeight: '500',
                fontFamily: 'ABC Marist, sans-serif',
                color: 'var(--orange)',
                letterSpacing: '-0.03em',
                lineHeight: '1'
              }}>${monthResults.endingBalance.toFixed(2)}</div>
              <div style={{
                fontSize: '0.75rem',
                color: 'var(--beige)',
                marginTop: '0.5rem',
                letterSpacing: '0.01em'
              }}>Ending balance</div>
            </div>

            {/* Net Change */}
            <div style={{
              padding: '2rem',
              background: 'var(--white)',
              transition: 'all 0.2s ease'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--cream)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'var(--white)'}>
              <div style={{
                fontSize: '0.6875rem',
                color: 'var(--beige)',
                marginBottom: '0.75rem',
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                fontWeight: '500'
              }}>Net Change</div>
              <div style={{
                fontSize: '2rem',
                fontWeight: '500',
                fontFamily: 'ABC Marist, sans-serif',
                color: monthResults.profit >= 0 ? 'var(--blue)' : 'var(--orange)',
                letterSpacing: '-0.03em',
                lineHeight: '1'
              }}>
                {monthResults.profit >= 0 ? '+' : ''}${monthResults.profit.toFixed(2)}
              </div>
              <div style={{
                fontSize: '0.75rem',
                color: 'var(--beige)',
                marginTop: '0.5rem',
                letterSpacing: '0.01em'
              }}>Monthly profit</div>
            </div>
          </div>

          {/* Income & Expenses Split */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: '1fr 1fr',
            borderBottom: budgetComparison ? '1px solid var(--light-border)' : 'none'
          }}>
            {/* Income */}
            <div style={{
              padding: '2rem',
              borderRight: '1px solid var(--light-border)',
              background: 'var(--white)',
              transition: 'all 0.2s ease'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--cream)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'var(--white)'}>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                marginBottom: '1rem'
              }}>
                <div style={{
                  fontSize: '0.6875rem',
                  color: 'var(--beige)',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  fontWeight: '500'
                }}>Income</div>
                <div style={{
                  width: '28px',
                  height: '28px',
                  borderRadius: '50%',
                  background: 'var(--blue)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '1rem',
                  color: 'var(--white)'
                }}>↓</div>
              </div>
              <div style={{
                fontSize: '1.75rem',
                fontWeight: '500',
                fontFamily: 'ABC Marist, sans-serif',
                color: 'var(--blue)',
                letterSpacing: '-0.03em',
                marginBottom: '0.5rem'
              }}>+${monthResults.totalIncome.toFixed(2)}</div>
              <div style={{
                fontSize: '0.75rem',
                color: 'var(--beige)',
                letterSpacing: '0.01em'
              }}>{monthResults.incomeCount} transactions</div>
            </div>

            {/* Expenses */}
            <div style={{
              padding: '2rem',
              background: 'var(--white)',
              transition: 'all 0.2s ease'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--cream)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'var(--white)'}>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                marginBottom: '1rem'
              }}>
                <div style={{
                  fontSize: '0.6875rem',
                  color: 'var(--beige)',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  fontWeight: '500'
                }}>Expenses</div>
                <div style={{
                  width: '28px',
                  height: '28px',
                  borderRadius: '50%',
                  background: 'var(--orange)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '1rem',
                  color: 'var(--white)'
                }}>↑</div>
              </div>
              <div style={{
                fontSize: '1.75rem',
                fontWeight: '500',
                fontFamily: 'ABC Marist, sans-serif',
                color: 'var(--orange)',
                letterSpacing: '-0.03em',
                marginBottom: '0.5rem'
              }}>-${monthResults.totalExpenses.toFixed(2)}</div>
              <div style={{
                fontSize: '0.75rem',
                color: 'var(--beige)',
                letterSpacing: '0.01em'
              }}>{monthResults.expenseCount} transactions</div>
            </div>
          </div>

          {/* Budget Achievement Result */}
          {budgetComparison && (
            <div style={{
              padding: '2rem',
              background: budgetComparison.achievedTarget ? 'var(--blue)' : 'var(--orange)',
              color: 'var(--white)',
              position: 'relative',
              overflow: 'hidden'
            }}>
              {/* Decorative background pattern */}
              <div style={{
                position: 'absolute',
                top: 0,
                right: 0,
                width: '200px',
                height: '200px',
                opacity: 0.1,
                pointerEvents: 'none'
              }}>
                <svg width="200" height="200" viewBox="0 0 200 200">
                  <circle cx="100" cy="100" r="80" stroke="currentColor" strokeWidth="2" fill="none" />
                  <circle cx="100" cy="100" r="60" stroke="currentColor" strokeWidth="2" fill="none" />
                  <circle cx="100" cy="100" r="40" stroke="currentColor" strokeWidth="2" fill="none" />
                </svg>
              </div>

              <div style={{ position: 'relative', zIndex: 1 }}>
                <div style={{
                  fontSize: '0.6875rem',
                  color: 'rgba(255, 255, 255, 0.8)',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  marginBottom: '1rem',
                  fontWeight: '500'
                }}>Target: ${budgetComparison.target.toFixed(2)}</div>

                <div style={{
                  fontSize: '2.5rem',
                  fontWeight: '500',
                  fontFamily: 'ABC Marist, sans-serif',
                  letterSpacing: '-0.03em',
                  marginBottom: '0.75rem',
                  lineHeight: '1'
                }}>
                  {budgetComparison.achievedTarget ? '✓ Savings Target Achieved' : '⚠ Target Missed'}
                </div>

                <div style={{
                  fontSize: '1rem',
                  letterSpacing: '0.01em',
                  opacity: 0.95
                }}>
                  ${Math.abs(budgetComparison.difference).toFixed(2)} {budgetComparison.achievedTarget ? 'above' : 'below'} your monthly target
                </div>
              </div>
            </div>
          )}

          {!budgetTarget && (
            <div style={{
              padding: '4rem 2rem',
              textAlign: 'center',
              background: 'var(--white)'
            }}>
              <div style={{
                fontSize: '0.875rem',
                color: 'var(--beige)',
                letterSpacing: '0.01em'
              }}>Set a monthly budget target to track your progress</div>
            </div>
          )}
        </div>
      </section>
    </>
  )
}
