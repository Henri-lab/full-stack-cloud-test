import { useState } from 'react'
import emailsData from '../resource/emails.json'

interface EmailMeta {
  banned: boolean
  created_at: string
  updated_at: string
  price: number
  sold: boolean
  need_repair: boolean
}

interface Email {
  main: string
  password: string
  deputy: string
  key_2FA: string
  meta: EmailMeta
}

function Emails() {
  const [emails] = useState<Email[]>(emailsData.emails)
  const [searchTerm, setSearchTerm] = useState('')
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const filteredEmails = emails.filter(email =>
    email.main.toLowerCase().includes(searchTerm.toLowerCase()) ||
    email.deputy.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const copyToClipboard = (text: string, fieldId: string) => {
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
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {filteredEmails.map((email, index) => (
              <tr
                key={index}
                className={`transition-colors ${email.meta.banned ? 'bg-red-500/5 hover:bg-red-500/10' : 'hover:bg-indigo-500/5'}`}
              >
                <td className="px-4 py-4 text-indigo-400 font-bold">{index + 1}</td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="text-white font-medium">{email.main}</span>
                    <button
                      onClick={() => copyToClipboard(email.main, `main-${index}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `main-${index}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `main-${index}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <code className="px-2 py-1 bg-slate-700 rounded text-pink-400 font-mono text-sm">{email.password}</code>
                    <button
                      onClick={() => copyToClipboard(email.password, `pwd-${index}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `pwd-${index}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `pwd-${index}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="text-slate-400 text-sm">{email.deputy}</span>
                    <button
                      onClick={() => copyToClipboard(email.deputy, `deputy-${index}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `deputy-${index}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `deputy-${index}` ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center gap-2 flex-wrap">
                    <code className="px-2 py-1 bg-slate-700 rounded text-pink-400 font-mono text-sm">{email.key_2FA}</code>
                    <button
                      onClick={() => copyToClipboard(email.key_2FA, `key-${index}`)}
                      className={`px-2 py-1 text-xs font-semibold rounded transition-all ${
                        copiedField === `key-${index}`
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white'
                          : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:scale-105'
                      }`}
                    >
                      {copiedField === `key-${index}` ? 'Copied!' : 'Copy'}
                    </button>
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
