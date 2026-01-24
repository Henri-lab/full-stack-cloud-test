import { useState, FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'
import { AxiosError } from 'axios'

interface LoginProps {
  setUser: (user: { token: string } | null) => void
}

interface LoginResponse {
  token: string
}

interface ErrorResponse {
  error: string
  remaining_attempts?: number
}

function Login({ setUser }: LoginProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [remainingAttempts, setRemainingAttempts] = useState<number | null>(null)
  const [isBlocked, setIsBlocked] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')

    if (isBlocked) {
      setError('Too many failed attempts. Please try again later.')
      return
    }

    try {
      const response = await api.post<LoginResponse>('/auth/login', { email, password })
      const { token } = response.data

      localStorage.setItem('token', token)
      setUser({ token })
      setRemainingAttempts(null)
      navigate('/tasks')
    } catch (err) {
      const axiosError = err as AxiosError<ErrorResponse>

      if (axiosError.response?.status === 429) {
        setIsBlocked(true)
        setError('Too many failed login attempts. Please try again in 15 minutes.')
      } else {
        const remaining = axiosError.response?.data?.remaining_attempts
        if (remaining !== undefined) {
          setRemainingAttempts(remaining)
        }
        setError(axiosError.response?.data?.error || 'Login failed')
      }
    }
  }

  return (
    <div className="max-w-md mx-auto mt-8 p-8 bg-slate-800 rounded-xl shadow-2xl border border-slate-700">
      <h2 className="text-2xl font-bold mb-6 text-center bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
        Login
      </h2>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-slate-300 mb-2">
            Email
          </label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={isBlocked}
            className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all disabled:opacity-50"
          />
        </div>
        <div>
          <label htmlFor="password" className="block text-sm font-medium text-slate-300 mb-2">
            Password
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={isBlocked}
            className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all disabled:opacity-50"
          />
        </div>
        {error && (
          <div className="p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-sm">
            {error}
            {remainingAttempts !== null && remainingAttempts > 0 && (
              <p className="mt-1 text-xs">
                Remaining attempts: {remainingAttempts}
              </p>
            )}
          </div>
        )}
        <button
          type="submit"
          disabled={isBlocked}
          className="w-full py-3 px-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg hover:from-indigo-500 hover:to-purple-500 transition-all shadow-lg hover:shadow-indigo-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isBlocked ? 'Temporarily Blocked' : 'Login'}
        </button>
      </form>
    </div>
  )
}

export default Login
