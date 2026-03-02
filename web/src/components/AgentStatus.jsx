import React from 'react'
import styles from '../styles/Dashboard.module.css'

export default function AgentStatus({ agents }) {
  if (!agents || agents.length === 0) {
    return <p className={styles.empty}>No agents connected</p>
  }

  return (
    <div className={styles.agentList}>
      {agents.map((agent) => (
        <div key={agent.id} className={styles.agentItem}>
          <strong>{agent.id.substring(0, 12)}</strong>
          <p>Status: {agent.status}</p>
        </div>
      ))}
    </div>
  )
}
