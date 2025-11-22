import axios from 'axios'

// Use relative URL to leverage Vite proxy in development
// In production, this will be handled by the backend serving static files
const API_URL = import.meta.env.VITE_API_URL || '/api'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const attendeesAPI = {
  getAll: () => api.get('/attendees'),
  register: (data) => api.post('/attendees', data),
  getCount: () => api.get('/attendees/count'),
}

export const speakersAPI = {
  getAll: () => api.get('/speakers'),
  create: (data) => api.post('/speakers', data),
  update: (id, data) => api.put(`/speakers/${id}`, data),
  delete: (id) => api.delete(`/speakers/${id}`),
}

export const sessionsAPI = {
  getAll: () => api.get('/sessions'),
  create: (data) => api.post('/sessions', data),
  update: (id, data) => api.put(`/sessions/${id}`, data),
  delete: (id) => api.delete(`/sessions/${id}`),
}

export const adminAPI = {
  login: (password) => api.post('/admin/login', { password }),
}

