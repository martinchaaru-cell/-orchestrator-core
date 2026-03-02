import React, { useState } from 'react'
import axios from 'axios'
import styles from '../styles/Dashboard.module.css'

export default function TaskSubmit({ onTaskSubmitted, apiUrl }) {
  const [taskType, setTaskType] = useState('hello_world')
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const task = {
        type: taskType,
        params: {
          message: message || 'test',
          target: message || '192.168.1.0/24',
        },
      }

      const response = await axios.post(`${apiUrl}/task`, task)
      onTaskSubmitted(response.data)
      setMessage('')
    } catch (err) {
      setError('Failed to submit task')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      {error && <div className={styles.error}>{error}</div>}
      
      <div className={styles.formGroup}>
        <label>Task Type:</label>
        <select value={taskType} onChange={(e) => setTaskType(e.target.value)} disabled={loading}>
          <option value="hello_world">Hello World</option>
          <option value="echo_task">Echo Task</option>
          <option value="nmap_scan">Nmap Scan</option>
        </select>
      </div>

      <div className={styles.formGroup}>
        <label>Parameter:</label>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="message or target"
          disabled={loading}
        />
      </div>

      <button type="submit" disabled={loading} className={styles.button}>
        {loading ? 'Submitting...' : '📤 Submit'}
      </button>
    </form>
  )
}
