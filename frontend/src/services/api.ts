import axios, { AxiosInstance } from 'axios'

// Resolve base URL from environment; fall back to relative path for dev
const baseURL =
  import.meta.env?.VITE_API_URL && import.meta.env.VITE_API_URL.length > 0
    ? `${import.meta.env.VITE_API_URL}/api/v1`
    : '/api/v1'

const api: AxiosInstance = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token to requests if it exists
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers = config.headers ?? {}
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Handle response errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
