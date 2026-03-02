import React, { useState, useEffect } from 'react'
import axios from 'axios'
import TaskSubmit from './TaskSubmit'
import TaskList from './TaskList'
import AgentStatus from './AgentStatus'
import styles from '../styles/Dashboard.module.css'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export default function Dashboard() {
  const [tasks, setTasks] = useState([])
  const [agents, setAgents] = useState([])
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [tasksRes, agentsRes] = await Promise.all([
          axios.get(`${API_URL}/tasks`),
          axios.get(`${API_URL}/agents`),
        ])
        setTasks(tasksRes.data.tasks || [])
        setAgents(agentsRes.data.agents || [])
        setError(null)
      } catch (err) {
        setError('Failed to fetch data')
      }
    }

    fetchData()
    const interval = setInterval(fetchData, 2000)
    return () => clearInterval(interval)
  }, [])

  const handleTaskSubmitted = (newTask) => {
    setTasks([newTask, ...tasks])
  }

  return (
    <div className={styles.dashboard}>
      <div className={styles.grid}>
        <section className={styles.section}>
          <h2>📋 Submit Task</h2>
          <TaskSubmit onTaskSubmitted={handleTaskSubmitted} apiUrl={API_URL} />
        </section>

        <section className={styles.section}>
          <h2>📊 Tasks ({tasks.length})</h2>
          {error && <div className={styles.error}>{error}</div>}
          <TaskList tasks={tasks} />
        </section>

        <section className={styles.section}>
          <h2>🤖 Agents ({agents.length})</h2>
          <AgentStatus agents={agents} />
        </section>
      </div>
    </div>
  )
}
