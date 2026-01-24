import { Link } from 'react-router-dom'

function Home() {
  return (
    <div className="text-center py-16 px-8">
      <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
        Welcome to FreeGemini
      </h1>
      <p className="text-xl text-slate-400 mb-8">
        A minimal full-stack application with React, Go, and PostgreSQL
      </p>
      <div className="flex gap-4 justify-center">
        <Link to="/tasks">
          <button className="px-8 py-4 text-lg font-semibold bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:from-indigo-500 hover:to-purple-500 transition-all shadow-lg hover:shadow-indigo-500/25">
            Tasks
          </button>
        </Link>
        <Link to="/emails">
          <button className="px-8 py-4 text-lg font-semibold bg-gradient-to-r from-pink-600 to-rose-600 text-white rounded-lg hover:from-pink-500 hover:to-rose-500 transition-all shadow-lg hover:shadow-pink-500/25">
            Emails
          </button>
        </Link>
      </div>
    </div>
  )
}

export default Home
