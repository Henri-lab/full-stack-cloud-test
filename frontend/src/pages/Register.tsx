import { useState, FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'
import { AxiosError } from 'axios'

interface ErrorResponse {
  error: string
}

// Password validation
const validatePassword = (password: string): { valid: boolean; message: string } => {
  if (password.length < 8) {
    return { valid: false, message: 'Password must be at least 8 characters' }
  }
  if (!/[A-Z]/.test(password)) {
    return { valid: false, message: 'Password must contain at least one uppercase letter' }
  }
  if (!/[a-z]/.test(password)) {
    return { valid: false, message: 'Password must contain at least one lowercase letter' }
  }
  if (!/[0-9]/.test(password)) {
    return { valid: false, message: 'Password must contain at least one number' }
  }
  return { valid: true, message: '' }
}

function Register() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [passwordStrength, setPasswordStrength] = useState<string[]>([])
  const navigate = useNavigate()

  const checkPasswordStrength = (pwd: string) => {
    const checks: string[] = []
    if (pwd.length >= 8) checks.push('length')
    if (/[A-Z]/.test(pwd)) checks.push('uppercase')
    if (/[a-z]/.test(pwd)) checks.push('lowercase')
    if (/[0-9]/.test(pwd)) checks.push('number')
    setPasswordStrength(checks)
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const pwd = e.target.value
    setPassword(pwd)
    checkPasswordStrength(pwd)
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    // Client-side validation
    const validation = validatePassword(password)
    if (!validation.valid) {
      setError(validation.message)
      return
    }

    try {
      await api.post('/auth/register', { username, email, password })
      setSuccess('Registration successful! Redirecting to login...')
      setTimeout(() => navigate('/login'), 2000)
    } catch (err) {
      const axiosError = err as AxiosError<ErrorResponse>
      setError(axiosError.response?.data?.error || 'Registration failed')
    }
  }

  return (
    <div className="max-w-md mx-auto mt-8 p-8 bg-slate-800 rounded-xl shadow-2xl border border-slate-700">
      <h2 className="text-2xl font-bold mb-6 text-center bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
        Register
      </h2>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="username" className="block text-sm font-medium text-slate-300 mb-2">
            Username
          </label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            minLength={3}
            maxLength={50}
            className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          />
        </div>
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
            className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
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
            onChange={handlePasswordChange}
            required
            minLength={8}
            className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          />
          {/* Password strength indicator */}
          {password && (
            <div className="mt-3 space-y-2">
              <div className="flex gap-1">
                {[1, 2, 3, 4].map((i) => (
                  <div
                    key={i}
                    className={`h-1 flex-1 rounded-full transition-all ${
                      passwordStrength.length >= i
                        ? passwordStrength.length === 4
                          ? 'bg-green-500'
                          : passwordStrength.length >= 3
                          ? 'bg-yellow-500'
                          : 'bg-red-500'
                        : 'bg-slate-600'
                    }`}
                  />
                ))}
              </div>
              <div className="text-xs space-y-1">
                <p className={passwordStrength.includes('length') ? 'text-green-400' : 'text-slate-500'}>
                  {passwordStrength.includes('length') ? '✓' : '○'} At least 8 characters
                </p>
                <p className={passwordStrength.includes('uppercase') ? 'text-green-400' : 'text-slate-500'}>
                  {passwordStrength.includes('uppercase') ? '✓' : '○'} One uppercase letter
                </p>
                <p className={passwordStrength.includes('lowercase') ? 'text-green-400' : 'text-slate-500'}>
                  {passwordStrength.includes('lowercase') ? '✓' : '○'} One lowercase letter
                </p>
                <p className={passwordStrength.includes('number') ? 'text-green-400' : 'text-slate-500'}>
                  {passwordStrength.includes('number') ? '✓' : '○'} One number
                </p>
              </div>
            </div>
          )}
        </div>
        {error && (
          <div className="p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}
        {success && (
          <div className="p-3 bg-green-500/20 border border-green-500/50 rounded-lg text-green-400 text-sm">
            {success}
          </div>
        )}
        <button
          type="submit"
          disabled={passwordStrength.length < 4}
          className="w-full py-3 px-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg hover:from-indigo-500 hover:to-purple-500 transition-all shadow-lg hover:shadow-indigo-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Register
        </button>
      </form>
    </div>
  )
}

export default Register
