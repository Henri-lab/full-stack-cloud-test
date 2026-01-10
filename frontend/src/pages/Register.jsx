import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'

function Register() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      await api.post('/auth/register', { username, email, password })
      setSuccess('Registration successful! Redirecting to login...')
      setTimeout(() => navigate('/login'), 2000)
    } catch (err) {
      setError(err.response?.data?.message || 'Registration failed')
    }
  }

  return (
    <div className="form-container">
      <h2>Register</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength="6"
          />
        </div>
        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}
        <div className="form-actions">
          <button type="submit">Register</button>
        </div>
      </form>
    </div>
  )
}

export default Register
