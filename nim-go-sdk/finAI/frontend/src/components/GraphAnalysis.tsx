import { useState } from 'react'

interface Transaction {
  id: string
  amount: number
  description: string
  date: string
  isIncoming: boolean
  merchant?: string
}

interface RecentTransactionsProps {
  transactions: Transaction[]
  balance: number
}

export default function GraphAnalysis({ transactions, balance }: RecentTransactionsProps) {
  
  return (
    <>
      <section className="transactions-section">
        <p>hopefully i work????</p>
      </section>
    </>
  )
}