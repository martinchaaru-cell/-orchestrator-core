import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
})

export const submitTask = (type, params) => {
  return api.post('/task', { type, params })
}

export const getTask = (id) => {
  return api.get(`/task/${id}`)
}

export const listTasks = () => {
  return api.get('/tasks')
}

export const listAgents = () => {
  return api.get('/agents')
}

export default api
