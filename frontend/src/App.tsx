import { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Tasks from './pages/Tasks'
import Emails from './pages/Emails'

interface User {
  token: string
}

function App() {
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      setUser({ token })
    }
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    setUser(null)
  }

  return (
    <Router>
      <div className="min-h-screen flex flex-col bg-slate-900">
        <nav className="bg-slate-800 px-8 py-4 shadow-lg">
          <div className="max-w-7xl mx-auto flex justify-between items-center">
            <Link to="/" className="text-2xl font-bold text-indigo-500 hover:text-indigo-400 transition-colors">
              FreeGemini
            </Link>
            <ul className="flex gap-8 list-none">
              <li>
                <Link to="/" className="text-slate-200 hover:text-indigo-400 transition-colors">
                  Home
                </Link>
              </li>
              {user ? (
                <>
                  <li>
                    <Link to="/tasks" className="text-slate-200 hover:text-indigo-400 transition-colors">
                      Tasks
                    </Link>
                  </li>
                  <li>
                    <Link to="/emails" className="text-slate-200 hover:text-indigo-400 transition-colors">
                      Emails
                    </Link>
                  </li>
                  <li>
                    <button
                      onClick={handleLogout}
                      className="text-slate-200 hover:text-indigo-400 transition-colors bg-transparent border-none cursor-pointer"
                    >
                      Logout
                    </button>
                  </li>
                </>
              ) : (
                <>
                  <li>
                    <Link to="/login" className="text-slate-200 hover:text-indigo-400 transition-colors">
                      Login
                    </Link>
                  </li>
                  <li>
                    <Link to="/register" className="text-slate-200 hover:text-indigo-400 transition-colors">
                      Register
                    </Link>
                  </li>
                </>
              )}
            </ul>
          </div>
        </nav>

        <main className="flex-1 max-w-7xl w-full mx-auto p-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login setUser={setUser} />} />
            <Route path="/register" element={<Register />} />
            <Route path="/tasks" element={<Tasks />} />
            <Route path="/emails" element={<Emails />} />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App
