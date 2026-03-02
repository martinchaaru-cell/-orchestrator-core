import React from 'react'
import Dashboard from '../components/Dashboard'
import styles from '../styles/Dashboard.module.css'

export default function Home() {
  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>🚀 Modular Security Orchestrator</h1>
        <p>Distribute security tasks, execute in parallel</p>
      </header>
      <Dashboard />
    </div>
  )
}
