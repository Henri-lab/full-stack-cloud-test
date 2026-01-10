import { Link } from 'react-router-dom'

function Home() {
  return (
    <div className="hero">
      <h1>Welcome to FullStack App</h1>
      <p>A minimal full-stack application with React, Go, and PostgreSQL</p>
      <Link to="/tasks">
        <button className="cta-button">Get Started</button>
      </Link>
    </div>
  )
}

export default Home
