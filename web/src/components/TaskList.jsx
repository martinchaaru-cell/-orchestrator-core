import React from 'react'
import styles from '../styles/Dashboard.module.css'

export default function TaskList({ tasks }) {
  if (!tasks || tasks.length === 0) {
    return <p className={styles.empty}>No tasks yet</p>
  }

  return (
    <div className={styles.taskList}>
      {tasks.map((task) => (
        <div key={task.id} className={`${styles.taskItem} ${styles[task.status]}`}>
          <div className={styles.taskHeader}>
            <strong>{task.type}</strong>
            <span className={styles.status}>{task.status}</span>
          </div>
          {task.result && <pre className={styles.taskResult}>{task.result.substring(0, 200)}</pre>}
        </div>
      ))}
    </div>
  )
}
