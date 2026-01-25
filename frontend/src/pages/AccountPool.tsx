import { useEffect, useState } from 'react'
import api from '../services/api'

interface FamilyInfo {
  capacity: number
  used: number
}

interface AccountListItem {
  id: number
  type: string
  main: string
  status: string
  source: string
  family?: FamilyInfo
}

interface Subscription {
  id: number
  plan: string
  expires_at: string
  status: string
}

interface Credentials {
  id: number
  main: string
  password: string
  key_2FA: string
}

function AccountPool() {
  const [accounts, setAccounts] = useState<AccountListItem[]>([])
  const [subscription, setSubscription] = useState<Subscription | null>(null)
  const [credentials, setCredentials] = useState<Credentials | null>(null)
  const [loading, setLoading] = useState(true)
  const [message, setMessage] = useState<string | null>(null)

  const loadSubscription = async () => {
    try {
      const response = await api.get<{ subscription: Subscription | null }>('/subscriptions/me')
      setSubscription(response.data.subscription)
    } catch (err) {
      setMessage('Failed to load subscription.')
    }
  }

  const loadAccounts = async () => {
    setLoading(true)
    try {
      const response = await api.get<AccountListItem[]>('/accounts')
      setAccounts(response.data)
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to load accounts.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadSubscription()
    loadAccounts()
  }, [])

  const handleActivateSubscription = async () => {
    try {
      const response = await api.post<{ subscription: Subscription }>('/subscriptions/activate', {
        plan: 'monthly',
        duration_days: 30,
      })
      setSubscription(response.data.subscription)
      setMessage('Subscription activated.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to activate subscription.')
    }
  }

  const claimTemporary = async (accountId: number) => {
    try {
      const response = await api.post('/accounts/temporary/claim', { account_id: accountId })
      setMessage(response.data.message || 'Account claimed.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to claim account.')
    }
  }

  const releaseTemporary = async (accountId: number) => {
    try {
      const response = await api.post('/accounts/temporary/release', { account_id: accountId })
      setMessage(response.data.message || 'Account released.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to release account.')
    }
  }

  const purchaseExclusive = async (accountId: number) => {
    try {
      const response = await api.post('/accounts/exclusive/purchase', { account_id: accountId })
      setMessage(response.data.message || 'Purchase successful.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to purchase account.')
    }
  }

  const bindFamily = async (accountId: number) => {
    const memberEmail = window.prompt('Enter your Gmail for binding:')
    if (!memberEmail) return
    const memberPassword = window.prompt('Enter your Gmail password:')
    if (!memberPassword) return

    try {
      const response = await api.post('/accounts/family/bind', {
        account_id: accountId,
        member_email: memberEmail,
        member_password: memberPassword,
      })
      setMessage(response.data.message || 'Family binding created.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to bind family account.')
    }
  }

  const unbindFamily = async (accountId: number) => {
    try {
      const response = await api.post('/accounts/family/unbind', { account_id: accountId })
      setMessage(response.data.message || 'Family binding removed.')
      loadAccounts()
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to unbind family account.')
    }
  }

  const loadCredentials = async (accountId: number) => {
    try {
      const response = await api.get<Credentials>(`/accounts/exclusive/${accountId}/credentials`)
      setCredentials(response.data)
    } catch (err: any) {
      setMessage(err.response?.data?.error || 'Failed to load credentials.')
    }
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-white">Account Pool</h1>
        {!subscription && (
          <button
            onClick={handleActivateSubscription}
            className="px-4 py-2 rounded bg-indigo-600 text-white hover:bg-indigo-500"
          >
            Activate 30-day Subscription
          </button>
        )}
      </div>

      {subscription && (
        <div className="text-slate-300">
          Subscription: {subscription.plan} (expires {new Date(subscription.expires_at).toLocaleString()})
        </div>
      )}

      {message && (
        <div className="p-3 rounded bg-slate-800 text-slate-200 border border-slate-700">
          {message}
        </div>
      )}

      {credentials && (
        <div className="p-4 rounded bg-slate-800 border border-slate-700 text-slate-200">
          <div className="font-semibold mb-2">Exclusive Credentials</div>
          <div>Main: {credentials.main}</div>
          <div>Password: {credentials.password}</div>
          <div>2FA: {credentials.key_2FA}</div>
        </div>
      )}

      {loading ? (
        <div className="text-slate-400">Loading accounts...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {accounts.map(account => (
            <div key={account.id} className="p-4 rounded bg-slate-800 border border-slate-700 text-slate-200">
              <div className="flex items-center justify-between">
                <div className="font-semibold">{account.main}</div>
                <span className="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300">{account.type}</span>
              </div>
              <div className="text-sm text-slate-400 mt-1">Status: {account.status}</div>
              {account.family && (
                <div className="text-sm text-slate-400 mt-1">
                  Family: {account.family.used}/{account.family.capacity}
                </div>
              )}
              <div className="mt-3 flex flex-wrap gap-2">
                {account.type === 'temporary' && (
                  <>
                    <button
                      onClick={() => claimTemporary(account.id)}
                      className="px-3 py-1 rounded bg-emerald-600 text-white text-sm hover:bg-emerald-500"
                    >
                      Claim
                    </button>
                    <button
                      onClick={() => releaseTemporary(account.id)}
                      className="px-3 py-1 rounded bg-slate-600 text-white text-sm hover:bg-slate-500"
                    >
                      Release
                    </button>
                  </>
                )}
                {account.type === 'exclusive' && (
                  <>
                    <button
                      onClick={() => purchaseExclusive(account.id)}
                      className="px-3 py-1 rounded bg-indigo-600 text-white text-sm hover:bg-indigo-500"
                    >
                      Purchase
                    </button>
                    <button
                      onClick={() => loadCredentials(account.id)}
                      className="px-3 py-1 rounded bg-slate-600 text-white text-sm hover:bg-slate-500"
                    >
                      View Credentials
                    </button>
                  </>
                )}
                {account.type === 'family' && (
                  <>
                    <button
                      onClick={() => bindFamily(account.id)}
                      className="px-3 py-1 rounded bg-cyan-600 text-white text-sm hover:bg-cyan-500"
                    >
                      Bind
                    </button>
                    <button
                      onClick={() => unbindFamily(account.id)}
                      className="px-3 py-1 rounded bg-slate-600 text-white text-sm hover:bg-slate-500"
                    >
                      Unbind
                    </button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default AccountPool
