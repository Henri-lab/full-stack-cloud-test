import { useState, useEffect } from 'react'
import api from '../services/api'

function Tasks() {
  const [tasks, setTasks] = useState([])
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchTasks()
  }, [])

  const fetchTasks = async () => {
    try {
      const response = await api.get('/tasks')
      setTasks(response.data)
      setLoading(false)
    } catch (err) {
      setError('Failed to fetch tasks')
      setLoading(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    try {
      await api.post('/tasks', { title, description })
      setTitle('')
      setDescription('')
      fetchTasks()
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to create task')
    }
  }

  const handleDelete = async (id) => {
    try {
      await api.delete(`/tasks/${id}`)
      fetchTasks()
    } catch (err) {
      setError('Failed to delete task')
    }
  }

  if (loading) {
    return <div>Loading...</div>
  }

  return (
    <div className="tasks-container">
      <div className="tasks-header">
        <h1>Tasks</h1>
      </div>

      <div className="form-container">
        <h2>Create New Task</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="title">Title</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="description">Description</label>
            <input
              type="text"
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
          {error && <div className="error-message">{error}</div>}
          <div className="form-actions">
            <button type="submit">Create Task</button>
          </div>
        </form>
      </div>

      <div className="task-list">
        <h2>All Tasks</h2>
        {tasks.length === 0 ? (
          <p>No tasks yet. Create one above!</p>
        ) : (
          <ul>
            {tasks.map((task) => (
              <li key={task.id} className="task-item">
                <div className="task-info">
                  <h3>{task.title}</h3>
                  <p>{task.description}</p>
                  <small>Status: {task.status}</small>
                </div>
                <div className="task-actions">
                  <button onClick={() => handleDelete(task.id)}>Delete</button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

export default Tasks
