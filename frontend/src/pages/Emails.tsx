import { useState, useEffect, useRef, ChangeEvent } from 'react'
import * as OTPAuth from 'otpauth'
import api from '../services/api'
import { AxiosError } from 'axios'

interface EmailFamily {
  id: number
  email: string
  password: string
  code: string
  contact: string
  issue: string
}

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
  status: string  // unknown, live, verify, dead
  meta: EmailMeta
  familys: EmailFamily[]
}

interface EmailImportOption {
  id: number
  name: string
  created_at: string
  count: number
}

function Emails() {
  const [emails, setEmails] = useState<Email[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [imports, setImports] = useState<EmailImportOption[]>([])
  const [selectedImportId, setSelectedImportId] = useState<number | ''>('')
  const [searchTerm, setSearchTerm] = useState('')
  const [copiedField, setCopiedField] = useState<string | null>(null)
  const [totpCodes, setTotpCodes] = useState<{ [key: number]: string }>({})
  const [timeRemaining, setTimeRemaining] = useState(30)
  const [importing, setImporting] = useState(false)
  const [importMessage, setImportMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  // 验证相关状态
  const [verifying, setVerifying] = useState(false)
  const [verifyKey, setVerifyKey] = useState('')
  const [verifyMethod, setVerifyMethod] = useState<'api' | 'smtp'>('smtp')
  const [showKeyInput, setShowKeyInput] = useState(false)
  const [selectedEmails, setSelectedEmails] = useState<Set<number>>(new Set())

  // License Key 相关状态
  const [licenseKey, setLicenseKey] = useState('')
  const [showLicenseInput, setShowLicenseInput] = useState(false)
  const [licenseStatus, setLicenseStatus] = useState<{ valid: boolean; quota_remaining?: number } | null>(null)

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
    const loadImports = async () => {
      try {
        const response = await api.get<EmailImportOption[]>('/emails/imports')
        if (isMounted) {
          setImports(response.data)
        }
      } catch (err) {
        console.error('Failed to load imports', err)
        if (isMounted) {
          setError('Failed to load import history.')
        }
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    loadImports()
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

  const fetchEmails = async (importId?: number) => {
    setLoading(true)
    try {
      const response = await api.get<Email[]>('/emails', {
        params: importId ? { import_id: importId } : undefined,
      })
      setEmails(response.data)
    } catch (err) {
      console.error('Failed to load emails', err)
      setError('Failed to load emails.')
    } finally {
      setLoading(false)
    }
  }

  const fetchImports = async () => {
    try {
      const response = await api.get<EmailImportOption[]>('/emails/imports')
      setImports(response.data)
    } catch (err) {
      console.error('Failed to load imports', err)
      setError('Failed to load import history.')
    }
  }

  const handleImportClick = () => {
    fileInputRef.current?.click()
  }

  const handleDownloadTestJson = () => {
    const link = document.createElement('a')
    link.href = '/test-emails.json'
    link.download = 'test-emails.json'
    document.body.appendChild(link)
    link.click()
    link.remove()
    setImportMessage({ type: 'success', text: 'Test JSON downloaded. Please click "Import JSON" to upload it.' })
  }

  const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setImporting(true)
    setImportMessage(null)

    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await api.post<{ message: string; imported: number; import_id: number; import_name: string }>('/emails/import', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      setImportMessage({ type: 'success', text: `${response.data.message} (${response.data.imported} emails)` })
      await fetchImports()
      setSelectedImportId(response.data.import_id)
      fetchEmails(response.data.import_id)
    } catch (err) {
      const axiosError = err as AxiosError<{ error: string }>
      const errorMsg = axiosError.response?.data?.error || 'Import failed'
      setImportMessage({ type: 'error', text: errorMsg })
      if (errorMsg.includes('License Key') || errorMsg.includes('Key') || errorMsg.includes('额度')) {
        setShowLicenseInput(true)
      }
    } finally {
      setImporting(false)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const handleLoadImport = () => {
    if (!selectedImportId) {
      setImportMessage({ type: 'error', text: 'Please select a saved dataset to load' })
      return
    }
    fetchEmails(selectedImportId)
  }

  // 验证邮箱功能
  const handleVerifyEmails = async () => {
    // 检查 License Key
    if (!licenseKey.trim()) {
      setImportMessage({ type: 'error', text: 'Please enter License Key to use verification feature' })
      setShowLicenseInput(true)
      return
    }

    // API 方法需要 key
    if (verifyMethod === 'api' && !verifyKey.trim()) {
      setImportMessage({ type: 'error', text: 'Please enter verification key for API method' })
      return
    }

    if (selectedEmails.size === 0) {
      setImportMessage({ type: 'error', text: 'Please select emails to verify' })
      return
    }

    setVerifying(true)
    setImportMessage(null)

    const emailsToVerify = Array.from(selectedEmails).map(id => {
      const email = emails.find(e => e.id === id)
      return email?.main
    }).filter(Boolean) as string[]

    try {
      const payload: { mail: string[]; method: string; key?: string } = {
        mail: emailsToVerify,
        method: verifyMethod
      }

      // 只有 API 方法才需要 key
      if (verifyMethod === 'api') {
        payload.key = verifyKey
      }

      // 添加 License Key 到请求头
      const config = {
        headers: {
          'X-License-Key': licenseKey
        }
      }

      const response = await api.post<{ results: Array<{ email: string; status: string; error?: string }>; total: number; method: string }>('/emails/verify', payload, config)

      const successCount = response.data.results.filter(r => r.status !== 'error').length
      const methodName = verifyMethod === 'smtp' ? 'SMTP' : 'API'
      setImportMessage({ type: 'success', text: `Verified ${successCount}/${response.data.total} emails successfully using ${methodName}` })

      // 更新本地邮箱状态
      setEmails(prevEmails => prevEmails.map(email => {
        const result = response.data.results.find(r => r.email === email.main)
        if (result) {
          return { ...email, status: result.status }
        }
        return email
      }))

      // 清空选择
      setSelectedEmails(new Set())

      // 刷新 License Key 状态
      await checkLicenseKey()
    } catch (err) {
      const axiosError = err as AxiosError<{ error: string }>
      const errorMsg = axiosError.response?.data?.error || 'Verification failed'
      setImportMessage({ type: 'error', text: errorMsg })

      // 如果是 License Key 相关错误，显示输入框
      if (errorMsg.includes('License Key') || errorMsg.includes('额度')) {
        setShowLicenseInput(true)
      }
    } finally {
      setVerifying(false)
    }
  }

  // 检查 License Key
  const checkLicenseKey = async () => {
    if (!licenseKey.trim()) {
      setLicenseStatus(null)
      return
    }

    try {
      const response = await api.post<{ key: any; quota_remaining: number }>('/keys/check', {
        key_code: licenseKey
      })
      setLicenseStatus({
        valid: true,
        quota_remaining: response.data.quota_remaining
      })
      setImportMessage({ type: 'success', text: `License Key valid! Remaining quota: ${response.data.quota_remaining}` })
    } catch (err) {
      const axiosError = err as AxiosError<{ error: string }>
      setLicenseStatus({ valid: false })
      setImportMessage({ type: 'error', text: axiosError.response?.data?.error || 'Invalid License Key' })
    }
  }

  // 保存 License Key 到 localStorage
  const saveLicenseKey = () => {
    if (licenseKey.trim()) {
      localStorage.setItem('license_key', licenseKey)
      checkLicenseKey()
    }
  }

  // 从 localStorage 加载 License Key
  useEffect(() => {
    const savedKey = localStorage.getItem('license_key')
    if (savedKey) {
      setLicenseKey(savedKey)
      // 不自动检查，让用户手动触发
    }
  }, [])

  // 切换邮箱选择
  const toggleEmailSelection = (id: number) => {
    setSelectedEmails(prev => {
      const newSet = new Set(prev)
      if (newSet.has(id)) {
        newSet.delete(id)
      } else {
        newSet.add(id)
      }
      return newSet
    })
  }

  // 全选/取消全选
  const toggleSelectAll = () => {
    if (selectedEmails.size === filteredEmails.length) {
      setSelectedEmails(new Set())
    } else {
      setSelectedEmails(new Set(filteredEmails.map(e => e.id)))
    }
  }

  // 批量复制选中的邮箱地址
  const handleCopySelectedEmails = () => {
    if (selectedEmails.size === 0) {
      setImportMessage({ type: 'error', text: 'Please select emails to copy' })
      return
    }

    const selectedEmailAddresses = Array.from(selectedEmails)
      .map(id => emails.find(e => e.id === id)?.main)
      .filter(Boolean)
      .join('\n')

    navigator.clipboard.writeText(selectedEmailAddresses).then(() => {
      setImportMessage({ type: 'success', text: `Copied ${selectedEmails.size} email addresses to clipboard` })
      setTimeout(() => setImportMessage(null), 3000)
    }).catch(() => {
      setImportMessage({ type: 'error', text: 'Failed to copy to clipboard' })
    })
  }

  // 获取状态徽章
  const getVerifyStatusBadge = (status: string) => {
    switch (status) {
      case 'live':
        return <span className="px-2 py-1 text-xs font-semibold rounded bg-green-500/20 text-green-400 border border-green-500/50">Live</span>
      case 'verify':
        return <span className="px-2 py-1 text-xs font-semibold rounded bg-yellow-500/20 text-yellow-400 border border-yellow-500/50">Verify</span>
      case 'dead':
        return <span className="px-2 py-1 text-xs font-semibold rounded bg-red-500/20 text-red-400 border border-red-500/50">Dead</span>
      default:
        return <span className="px-2 py-1 text-xs font-semibold rounded bg-gray-500/20 text-gray-400 border border-gray-500/50">Unknown</span>
    }
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

      {/* Search, Import and Verify */}
      <div className="mb-6 space-y-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <input
            type="text"
            placeholder="Search emails..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full max-w-md px-4 py-3 bg-slate-800 border-2 border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:shadow-lg focus:shadow-indigo-500/20 transition-all"
          />
          <input
            type="file"
            ref={fileInputRef}
            accept=".json"
            onChange={handleFileChange}
            className="hidden"
          />
          <button
            onClick={handleImportClick}
            disabled={importing}
            className="px-6 py-3 bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-semibold rounded-lg hover:from-emerald-500 hover:to-teal-500 transition-all shadow-lg hover:shadow-emerald-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {importing ? 'Importing...' : 'Import JSON'}
          </button>
          <button
            onClick={handleDownloadTestJson}
            className="px-6 py-3 bg-gradient-to-r from-slate-600 to-slate-700 text-white font-semibold rounded-lg hover:from-slate-500 hover:to-slate-600 transition-all shadow-lg hover:shadow-slate-500/25"
          >
            Download Test JSON
          </button>
          <button
            onClick={() => setShowKeyInput(!showKeyInput)}
            className="px-6 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 text-white font-semibold rounded-lg hover:from-blue-500 hover:to-cyan-500 transition-all shadow-lg hover:shadow-blue-500/25"
          >
            {showKeyInput ? 'Hide Verify' : 'Verify Emails'}
          </button>
        </div>
        <div className="flex flex-col sm:flex-row gap-4">
          <select
            value={selectedImportId}
            onChange={(e) => setSelectedImportId(e.target.value ? Number(e.target.value) : '')}
            className="w-full max-w-md px-4 py-3 bg-slate-800 border-2 border-slate-700 rounded-lg text-white focus:outline-none focus:border-indigo-500 focus:shadow-lg focus:shadow-indigo-500/20 transition-all"
          >
            <option value="">Select saved dataset...</option>
            {imports.map(item => (
              <option key={item.id} value={item.id}>
                {item.name} ({item.count}) - {item.created_at ? new Date(item.created_at).toLocaleString() : ''}
              </option>
            ))}
          </select>
          <button
            onClick={handleLoadImport}
            disabled={!selectedImportId}
            className="px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg hover:from-indigo-500 hover:to-purple-500 transition-all shadow-lg hover:shadow-indigo-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Load Dataset
          </button>
        </div>

        {/* Verification Options */}
        {showKeyInput && (
          <div className="flex flex-col gap-4 p-4 bg-slate-800/50 border border-slate-700 rounded-lg">
            {/* License Key Input */}
            <div className="flex flex-col gap-2 p-3 bg-yellow-500/10 border border-yellow-500/30 rounded-lg">
              <div className="flex items-center gap-2">
                <svg className="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <span className="text-yellow-400 font-semibold">License Key Required</span>
              </div>
              <div className="flex flex-col sm:flex-row gap-2">
                <input
                  type="text"
                  placeholder="Enter your License Key (e.g., XXXX-XXXX-XXXX-XXXX)"
                  value={licenseKey}
                  onChange={(e) => setLicenseKey(e.target.value)}
                  className="flex-1 px-4 py-2 bg-slate-800 border-2 border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-yellow-500 focus:shadow-lg focus:shadow-yellow-500/20 transition-all font-mono text-sm"
                />
                <button
                  onClick={saveLicenseKey}
                  disabled={!licenseKey.trim()}
                  className="px-4 py-2 bg-gradient-to-r from-yellow-600 to-orange-600 text-white font-semibold rounded-lg hover:from-yellow-500 hover:to-orange-500 transition-all shadow-lg hover:shadow-yellow-500/25 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
                >
                  Check Key
                </button>
                <button
                  onClick={() => window.location.href = '/payment'}
                  className="px-4 py-2 bg-gradient-to-r from-green-600 to-emerald-600 text-white font-semibold rounded-lg hover:from-green-500 hover:to-emerald-500 transition-all shadow-lg hover:shadow-green-500/25 whitespace-nowrap"
                >
                  Buy Key
                </button>
              </div>
              {licenseStatus && (
                <div className={`text-sm ${licenseStatus.valid ? 'text-green-400' : 'text-red-400'}`}>
                  {licenseStatus.valid ? (
                    <span>✓ Valid License Key - Remaining quota: {licenseStatus.quota_remaining}</span>
                  ) : (
                    <span>✗ Invalid License Key</span>
                  )}
                </div>
              )}
            </div>

            {/* Method Selection */}
            <div className="flex gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="verifyMethod"
                  value="smtp"
                  checked={verifyMethod === 'smtp'}
                  onChange={(e) => setVerifyMethod(e.target.value as 'smtp')}
                  className="w-4 h-4 text-blue-600 bg-slate-700 border-slate-600 focus:ring-blue-500"
                />
                <span className="text-white font-medium">SMTP Verification</span>
                <span className="text-xs text-slate-400">(No key needed, slower but reliable)</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="verifyMethod"
                  value="api"
                  checked={verifyMethod === 'api'}
                  onChange={(e) => setVerifyMethod(e.target.value as 'api')}
                  className="w-4 h-4 text-blue-600 bg-slate-700 border-slate-600 focus:ring-blue-500"
                />
                <span className="text-white font-medium">API Verification</span>
                <span className="text-xs text-slate-400">(Requires key from gmailver.com)</span>
              </label>
            </div>

            {/* Key Input (only for API method) */}
            {verifyMethod === 'api' && (
              <div className="flex flex-col sm:flex-row gap-4">
                <input
                  type="text"
                  placeholder="Enter verification key from gmailver.com..."
                  value={verifyKey}
                  onChange={(e) => setVerifyKey(e.target.value)}
                  className="flex-1 px-4 py-3 bg-slate-800 border-2 border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-blue-500 focus:shadow-lg focus:shadow-blue-500/20 transition-all"
                />
              </div>
            )}

            {/* Verify Button */}
            <div className="flex justify-end">
              <button
                onClick={handleVerifyEmails}
                disabled={verifying || selectedEmails.size === 0 || (verifyMethod === 'api' && !verifyKey.trim()) || !licenseKey.trim()}
                className="px-6 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 text-white font-semibold rounded-lg hover:from-blue-500 hover:to-cyan-500 transition-all shadow-lg hover:shadow-blue-500/25 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
              >
                {verifying ? 'Verifying...' : `Verify ${selectedEmails.size} Email${selectedEmails.size !== 1 ? 's' : ''}`}
              </button>
            </div>
          </div>
        )}
      </div>

      {importMessage && (
        <div className={`mb-6 p-4 rounded-lg ${importMessage.type === 'success' ? 'bg-green-500/20 border border-green-500/50 text-green-400' : 'bg-red-500/20 border border-red-500/50 text-red-400'}`}>
          {importMessage.text}
        </div>
      )}

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
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={selectedEmails.size === filteredEmails.length && filteredEmails.length > 0}
                    onChange={toggleSelectAll}
                    className="w-4 h-4 rounded border-slate-600 bg-slate-700 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-slate-800"
                  />
                  {selectedEmails.size > 0 && (
                    <button
                      onClick={handleCopySelectedEmails}
                      className="px-2 py-1 text-xs font-semibold rounded bg-gradient-to-r from-purple-600 to-pink-600 text-white hover:from-purple-500 hover:to-pink-500 transition-all shadow-lg hover:shadow-purple-500/25"
                      title="Copy selected email addresses"
                    >
                      Copy ({selectedEmails.size})
                    </button>
                  )}
                </div>
              </th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">#</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Main Email</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Password</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Deputy Email</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">2FA Key</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">2FA Code</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Verify Status</th>
              <th className="px-4 py-4 text-left text-xs font-semibold text-indigo-400 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {filteredEmails.map((email, index) => (
              <tr
                key={email.id}
                className={`transition-colors ${email.meta.banned ? 'bg-red-500/5 hover:bg-red-500/10' : 'hover:bg-indigo-500/5'}`}
              >
                <td className="px-4 py-4">
                  <input
                    type="checkbox"
                    checked={selectedEmails.has(email.id)}
                    onChange={() => toggleEmailSelection(email.id)}
                    className="w-4 h-4 rounded border-slate-600 bg-slate-700 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-slate-800"
                  />
                </td>
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
                <td className="px-4 py-4">{getVerifyStatusBadge(email.status)}</td>
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
