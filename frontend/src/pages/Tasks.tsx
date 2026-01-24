import { useState, useEffect, FormEvent } from 'react'
import api from '../services/api'
import { AxiosError } from 'axios'

interface Task {
  id: number
  title: string
  description: string
  status: string
}

interface ErrorResponse {
  message: string
}

function Tasks() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchTasks()
  }, [])

  const fetchTasks = async () => {
    try {
      const response = await api.get<Task[]>('/tasks')
      setTasks(response.data)
      setLoading(false)
    } catch {
      setError('Failed to fetch tasks')
      setLoading(false)
    }
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      await api.post('/tasks', { title, description })
      setTitle('')
      setDescription('')
      fetchTasks()
    } catch (err) {
      const axiosError = err as AxiosError<ErrorResponse>
      setError(axiosError.response?.data?.message || 'Failed to create task')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/tasks/${id}`)
      fetchTasks()
    } catch {
      setError('Failed to delete task')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
          Tasks
        </h1>
      </div>

      <div className="bg-slate-800 rounded-xl p-6 shadow-2xl border border-slate-700 mb-8">
        <h2 className="text-xl font-semibold text-white mb-4">Create New Task</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-slate-300 mb-2">
              Title
            </label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            />
          </div>
          <div>
            <label htmlFor="description" className="block text-sm font-medium text-slate-300 mb-2">
              Description
            </label>
            <input
              type="text"
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-4 py-3 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            />
          </div>
          {error && (
            <div className="p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-sm">
              {error}
            </div>
          )}
          <button
            type="submit"
            className="w-full py-3 px-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg hover:from-indigo-500 hover:to-purple-500 transition-all shadow-lg hover:shadow-indigo-500/25"
          >
            Create Task
          </button>
        </form>
      </div>

      <div className="bg-slate-800 rounded-xl p-6 shadow-2xl border border-slate-700">
        <h2 className="text-xl font-semibold text-white mb-4">All Tasks</h2>
        {tasks.length === 0 ? (
          <p className="text-slate-400 text-center py-8">No tasks yet. Create one above!</p>
        ) : (
          <ul className="space-y-4">
            {tasks.map((task) => (
              <li
                key={task.id}
                className="p-4 bg-slate-700/50 rounded-lg border border-slate-600 flex justify-between items-center hover:bg-slate-700 transition-colors"
              >
                <div>
                  <h3 className="text-lg font-medium text-white">{task.title}</h3>
                  <p className="text-slate-400 text-sm">{task.description}</p>
                  <span className="inline-block mt-2 px-3 py-1 text-xs font-medium bg-indigo-500/20 text-indigo-400 rounded-full">
                    {task.status}
                  </span>
                </div>
                <button
                  onClick={() => handleDelete(task.id)}
                  className="px-4 py-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors border border-red-500/50"
                >
                  Delete
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

export default Tasks
