import { useState, useEffect } from 'react'
import * as OTPAuth from 'otpauth'
import api from '../services/api'

interface EmailMeta {
  banned: boolean
  created_at: string
  updated_at: string
  price: number
  sold: boolean
  need_repair: boolean
  from?: string
}

interface Email {
  id: number
  main: string
  password: string
  deputy: string
  key_2FA: string
  meta: EmailMeta
}

function Emails() {
  const [emails, setEmails] = useState<Email[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [copiedField, setCopiedField] = useState<string | null>(null)
  const [totpCodes, setTotpCodes] = useState<{ [key: number]: string }>({})
  const [timeRemaining, setTimeRemaining] = useState(30)

  // Generate TOTP code for a specific email
  const generateTOTP = (secret: string): string => {
    try {
      const totp = new OTPAuth.TOTP({
        issuer: 'FreeGemini',
        label: 'Email',
        algorithm: 'SHA1',
        digits: 6,
        period: 30,
        secret: secret.toUpperCase().replace(/\s/g, ''),
      })
      return totp.generate()
    } catch (error) {
      console.error('Error generating TOTP:', error)
      return 'ERROR'
    }
  }

  // Show TOTP code for specific email
  const showTOTP = (index: number, secret: string) => {
    const code = generateTOTP(secret)
    setTotpCodes(prev => ({ ...prev, [index]: code }))
  }

  const filteredEmails = emails.filter(email =>
    email.main.toLowerCase().includes(searchTerm.toLowerCase()) ||
    email.deputy.toLowerCase().includes(searchTerm.toLowerCase())
  )

  useEffect(() => {
    let isMounted = true
    const fetchEmails = async () => {
      try {
        const response = await api.get<Email[]>('/emails')
        if (isMounted) {
          setEmails(response.data)
        }
      } catch (err) {
        console.error('Failed to load emails', err)
        if (isMounted) {
          setError('Failed to load emails.')
        }
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    fetchEmails()
    return () => {
      isMounted = false
    }
  }, [])

  // Update time remaining and refresh codes
  useEffect(() => {
    const interval = setInterval(() => {
      const now = Math.floor(Date.now() / 1000)
      const remaining = 30 - (now % 30)
      setTimeRemaining(remaining)

      // Auto-refresh displayed TOTP codes
      if (remaining === 30) {
        setTotpCodes(prev => {
          const updated: { [key: number]: string } = {}
          Object.keys(prev).forEach(key => {
            const id = parseInt(key, 10)
            const email = filteredEmails.find(item => item.id === id)
            if (email) {
              updated[id] = generateTOTP(email.key_2FA)
            }
          })
          return updated
        })
      }
    }, 1000)

    return () => clearInterval(interval)
  }, [filteredEmails])

  const copyToClipboard = (text: string, fieldId: string) => {
    if (!text) {
      return
    }
    navigator.clipboard.writeText(text)
    setCopiedField(fieldId)
    setTimeout(() => setCopiedField(null), 1500)
  }

  const getStatusBadge = (meta: EmailMeta) => {
    if (meta.banned) return <span className="px-3 py-1 text-xs font-semibold rounded-full bg-gradient-to-r from-red-500 to-pink-500 text-white">Banned</span>
    if (meta.need_repair) return <span className="px-3 py-1 text-xs font-semibold rounded-full bg-gradient-to-r from-yellow-500 to-orange-500 text-slate-900">Need Repair</span>
    if (meta.sold) return <span className="px-3 py-1 text-xs font-semibold rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 text-white">Sold</span>
    return <span className="px-3 py-1 text-xs font-semibold rounded-full bg-gradient-to-r from-green-500 to-emerald-500 text-white">Active</span>
  }

  return (
    <div className="max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8">
        <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
          Email Accounts
        </h1>
        <div className="flex gap-4">
          <div className="bg-gradient-to-br from-slate-800 to-slate-900 border border-slate-700 rounded-xl px-6 py-4 text-center shadow-lg">
            <span className="block text-2xl font-bold text-pink-500">{emails.length}</span>
            <span className="text-xs text-slate-400 uppercase tracking-wider">Total</span>
          </div>
          <div className="bg-gradient-to-br from-slate-800 to-slate-900 border border-slate-700 rounded-xl px-6 py-4 text-center shadow-lg">
            <span className="block text-2xl font-bold text-green-500">{emails.filter(e => !e.meta.banned).length}</span>
            <span className="text-xs text-slate-400 uppercase tracking-wider">Active</span>
          </div>
          <div className="bg-gradient-to-br from-slate-800 to-slate-900 border border-slate-700 rounded-xl px-6 py-4 text-center shadow-lg">
            <span className="block text-2xl font-bold text-red-500">{emails.filter(e => e.meta.banned).length}</span>
            <span className="text-xs text-slate-400 uppercase tracking-wider">Banned</span>
          </div>
        </div>
      </div>

      {/* Search */}
      <div className="mb-6">
        <input
          type="text"
          placeholder="Search emails..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full max-w-md px-4 py-3 bg-slate-800 border-2 border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:shadow-lg focus:shadow-indigo-500/20 transition-all"
        />
      </div>

      {loading && (
        <div className="mb-6 text-slate-400">Loading emails...</div>
      )}
      {error && (
        <div className="mb-6 text-red-400">{error}</div>
      )}

      {/* Table */}
      <div className="overflow-x-auto rounded-xl shadow-2xl">
        <table className="w-full bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl overflow-hidden">
          <thead className="bg-gradient-to-r from-slate-700 to-slate-800">
            <tr>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">#</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Main Email</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Password</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Deputy Email</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">2FA Key</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">2FA Code</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {filteredEmails.map((email, index) => (
              <tr
                key={email.id}
                className={`transition-colors ${email.meta.banned ? 'bg-red-500/5 hover:bg-red-500/10' : 'hover:bg-indigo-500/5'}`}
              >
                <td className="px-4 py-4 text-indigo-400 font-bold">{index + 1}</td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="text-white font-medium">{email.main}</span>
                    <button
                      onClick={() => copyToClipboard(email.main, `main-${email.id}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `main-${email.id}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `main-${email.id}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <code className="px-2 py-1 bg-slate-700 rounded text-pink-400 font-mono text-sm">{email.password}</code>
                    <button
                      onClick={() => copyToClipboard(email.password, `pwd-${email.id}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `pwd-${email.id}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `pwd-${email.id}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="text-slate-400 text-sm">{email.deputy}</span>
                    <button
                      onClick={() => copyToClipboard(email.deputy, `deputy-${email.id}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `deputy-${email.id}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `deputy-${email.id}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <code className="px-2 py-1 bg-slate-700 rounded text-pink-400 font-mono text-sm">{email.key_2FA}</code>
                    <button
                      onClick={() => copyToClipboard(email.key_2FA, `key-${email.id}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `key-${email.id}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `key-${email.id}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    {totpCodes[email.id] ? (
                      <>
                        <code className="px-3 py-2 bg-gradient-to-r from-cyan-600 to-blue-600 rounded text-white font-mono text-lg font-bold tracking-wider">
                          {totpCodes[email.id]}
                        </code>
                        <div className="flex flex-col items-center">
                          <div className="text-xs text-slate-400">{timeRemaining}s</div>
                          <div className="w-12 h-1 bg-slate-700 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all"
                              style={{ width: `${(timeRemaining / 30) * 100}%` }}
                            />
                          </div>
                        </div>
                        <button
                          onClick={() => copyToClipboard(totpCodes[email.id], `totp-${email.id}`)}
                          className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                            copiedField === `totp-${email.id}`
                              ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                              : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                          }`}
                        >
                          {copiedField === `totp-${email.id}` ? 'Copied!' : 'Copy'}
                        </button>
                      </>
                    ) : (
                      <button
                        onClick={() => showTOTP(email.id, email.key_2FA)}
                        className="px-3 py-1 text-xs font-semibold rounded bg-gradient-to-r from-cyan-600 to-blue-600 text-white hover:from-cyan-500 hover:to-blue-500 transition-all"
                      >
                        Generate
                      </button>
                    )}
                  </div>
                </td>
                <td className="px-4 py-4">{getStatusBadge(email.meta)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {filteredEmails.length === 0 && (
        <div className="text-center py-12 text-slate-400">
          No emails found matching "{searchTerm}"
        </div>
      )}
    </div>
  )
}

export default Emails
