import { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import './App.css'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Tasks from './pages/Tasks'

function App() {
  const [user, setUser] = useState(null)

  useEffect(() => {
    // Check if user is logged in
    const token = localStorage.getItem('token')
    if (token) {
      // Verify token and get user info
      // This is a placeholder - implement actual verification
      setUser({ token })
    }
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    setUser(null)
  }

  return (
    <Router>
      <div className="app">
        <nav className="navbar">
          <div className="nav-container">
            <Link to="/" className="nav-logo">FullStack App</Link>
            <ul className="nav-menu">
              <li className="nav-item">
                <Link to="/" className="nav-link">Home</Link>
              </li>
              {user ? (
                <>
                  <li className="nav-item">
                    <Link to="/tasks" className="nav-link">Tasks</Link>
                  </li>
                  <li className="nav-item">
                    <button onClick={handleLogout} className="nav-link btn-link">Logout</button>
                  </li>
                </>
              ) : (
                <>
                  <li className="nav-item">
                    <Link to="/login" className="nav-link">Login</Link>
                  </li>
                  <li className="nav-item">
                    <Link to="/register" className="nav-link">Register</Link>
                  </li>
                </>
              )}
            </ul>
          </div>
        </nav>

        <main className="main-content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login setUser={setUser} />} />
            <Route path="/register" element={<Register />} />
            <Route path="/tasks" element={<Tasks />} />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App
