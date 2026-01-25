import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'

interface Product {
  type: string
  name: string
  price: number
  quota_amount: number
  features: string[]
}

interface Order {
  id: number
  order_no: string
  amount: number
  product_type: string
  quota_amount: number
  status: string
  expired_at: string
  created_at: string
}

interface LicenseKey {
  id: number
  key_code: string
  product_type: string
  quota_total: number
  quota_used: number
  status: string
  activated_at: string
  created_at: string
}

const Payment: React.FC = () => {
  const navigate = useNavigate()
  const [products, setProducts] = useState<Product[]>([])
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null)
  const [myKeys, setMyKeys] = useState<LicenseKey[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [activeTab, setActiveTab] = useState<'buy' | 'keys'>('buy')
  const [activeKey, setActiveKey] = useState(() => localStorage.getItem('license_key') || '')

  // 模拟支付状态
  const [paymentStep, setPaymentStep] = useState<'select' | 'pay' | 'success'>('select')

  useEffect(() => {
    fetchProducts()
    fetchMyKeys()
  }, [])

  const fetchProducts = async () => {
    try {
      const response = await api.get('/payments/products')
      setProducts(response.data.products)
    } catch (err: any) {
      setError(err.response?.data?.error || '获取产品列表失败')
    }
  }

  const fetchMyKeys = async () => {
    try {
      const response = await api.get('/keys')
      setMyKeys(response.data.keys)
    } catch (err: any) {
      console.error('获取密钥列表失败:', err)
    }
  }

  const handleSelectProduct = (product: Product) => {
    setSelectedProduct(product)
  }

  const handleCreateOrder = async () => {
    if (!selectedProduct) return

    setLoading(true)
    setError('')

    try {
      const response = await api.post('/payments/orders', {
        product_type: selectedProduct.type,
      })
      setCurrentOrder(response.data.order)
      setPaymentStep('pay')
    } catch (err: any) {
      setError(err.response?.data?.error || '创建订单失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSimulatePayment = async (method: 'alipay' | 'wechat') => {
    if (!currentOrder) return

    setLoading(true)
    setError('')

    try {
      // 模拟支付回调
      const response = await api.post<{ license_key: { key_code: string } }>('/payments/notify', {
        order_no: currentOrder.order_no,
        transaction_id: `TXN${Date.now()}`,
        payment_method: method,
      })

      setSuccess('支付成功！License Key 已生成')
      setPaymentStep('success')
      if (response.data?.license_key?.key_code) {
        localStorage.setItem('license_key', response.data.license_key.key_code)
        setActiveKey(response.data.license_key.key_code)
      }

      // 刷新密钥列表
      await fetchMyKeys()

      // 3秒后切换到密钥列表
      setTimeout(() => {
        setActiveTab('keys')
        setPaymentStep('select')
        setCurrentOrder(null)
        setSelectedProduct(null)
      }, 3000)
    } catch (err: any) {
      setError(err.response?.data?.error || '支付失败')
    } finally {
      setLoading(false)
    }
  }

  const applyKey = (keyCode: string) => {
    localStorage.setItem('license_key', keyCode)
    setActiveKey(keyCode)
    setSuccess('已设置为当前使用的 License Key')
  }

  const formatPrice = (price: number) => {
    return (price / 100).toFixed(2)
  }

  const getStatusBadge = (status: string) => {
    const badges: { [key: string]: string } = {
      active: 'bg-green-100 text-green-800',
      exhausted: 'bg-red-100 text-red-800',
      revoked: 'bg-gray-100 text-gray-800',
    }
    const labels: { [key: string]: string } = {
      active: '可用',
      exhausted: '已用尽',
      revoked: '已撤销',
    }
    return (
      <span className={`px-2 py-1 rounded text-xs ${badges[status] || 'bg-gray-100 text-gray-800'}`}>
        {labels[status] || status}
      </span>
    )
  }

  const getProductTypeBadge = (type: string) => {
    const badges: { [key: string]: string } = {
      basic: 'bg-blue-100 text-blue-800',
      pro: 'bg-purple-100 text-purple-800',
      enterprise: 'bg-yellow-100 text-yellow-800',
    }
    const labels: { [key: string]: string } = {
      basic: '基础版',
      pro: '专业版',
      enterprise: '企业版',
    }
    return (
      <span className={`px-2 py-1 rounded text-xs ${badges[type] || 'bg-gray-100 text-gray-800'}`}>
        {labels[type] || type}
      </span>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">License Key 管理</h1>
          <p className="mt-2 text-gray-600">购买 License Key 以使用高级功能</p>
        </div>

        {/* Tabs */}
        <div className="mb-6 border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            <button
              onClick={() => setActiveTab('buy')}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'buy'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              购买 Key
            </button>
            <button
              onClick={() => setActiveTab('keys')}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'keys'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              我的 Key ({myKeys.length})
            </button>
          </nav>
        </div>

        {/* Error/Success Messages */}
        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}
        {success && (
          <div className="mb-4 bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
            {success}
          </div>
        )}

        {/* Buy Tab */}
        {activeTab === 'buy' && (
          <>
            {paymentStep === 'select' && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {products.map((product) => (
                  <div
                    key={product.type}
                    className={`bg-white rounded-lg shadow-md p-6 cursor-pointer transition-all ${
                      selectedProduct?.type === product.type
                        ? 'ring-2 ring-blue-500 transform scale-105'
                        : 'hover:shadow-lg'
                    }`}
                    onClick={() => handleSelectProduct(product)}
                  >
                    <div className="text-center">
                      <h3 className="text-xl font-bold text-gray-900 mb-2">{product.name}</h3>
                      <div className="text-3xl font-bold text-blue-600 mb-4">
                        ¥{formatPrice(product.price)}
                      </div>
                      <div className="text-gray-600 mb-4">
                        {product.quota_amount} 次验证额度
                      </div>
                      <ul className="text-left space-y-2 mb-6">
                        {product.features.map((feature, index) => (
                          <li key={index} className="flex items-center text-sm text-gray-700">
                            <svg
                              className="w-4 h-4 mr-2 text-green-500"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M5 13l4 4L19 7"
                              />
                            </svg>
                            {feature}
                          </li>
                        ))}
                      </ul>
                      {selectedProduct?.type === product.type && (
                        <button
                          onClick={handleCreateOrder}
                          disabled={loading}
                          className="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 disabled:bg-gray-400"
                        >
                          {loading ? '创建中...' : '立即购买'}
                        </button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {paymentStep === 'pay' && currentOrder && (
              <div className="max-w-2xl mx-auto bg-white rounded-lg shadow-md p-8">
                <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
                  扫码支付
                </h2>

                <div className="mb-6 p-4 bg-gray-50 rounded">
                  <div className="flex justify-between mb-2">
                    <span className="text-gray-600">订单号:</span>
                    <span className="font-mono">{currentOrder.order_no}</span>
                  </div>
                  <div className="flex justify-between mb-2">
                    <span className="text-gray-600">产品:</span>
                    <span>{getProductTypeBadge(currentOrder.product_type)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">金额:</span>
                    <span className="text-2xl font-bold text-blue-600">
                      ¥{formatPrice(currentOrder.amount)}
                    </span>
                  </div>
                </div>

                {/* 收款码展示区域 */}
                <div className="mb-6">
                  <div className="grid grid-cols-2 gap-4">
                    {/* 支付宝收款码 */}
                    <div className="text-center">
                      <div className="bg-gray-100 h-64 flex items-center justify-center rounded mb-2">
                        <div className="text-center">
                          <svg
                            className="w-32 h-32 mx-auto text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={1}
                              d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"
                            />
                          </svg>
                          <p className="text-sm text-gray-500 mt-2">支付宝收款码</p>
                        </div>
                      </div>
                      <button
                        onClick={() => handleSimulatePayment('alipay')}
                        disabled={loading}
                        className="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 disabled:bg-gray-400"
                      >
                        {loading ? '处理中...' : '模拟支付宝支付'}
                      </button>
                    </div>

                    {/* 微信收款码 */}
                    <div className="text-center">
                      <div className="bg-gray-100 h-64 flex items-center justify-center rounded mb-2">
                        <div className="text-center">
                          <svg
                            className="w-32 h-32 mx-auto text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={1}
                              d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"
                            />
                          </svg>
                          <p className="text-sm text-gray-500 mt-2">微信收款码</p>
                        </div>
                      </div>
                      <button
                        onClick={() => handleSimulatePayment('wechat')}
                        disabled={loading}
                        className="w-full bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 disabled:bg-gray-400"
                      >
                        {loading ? '处理中...' : '模拟微信支付'}
                      </button>
                    </div>
                  </div>
                </div>

                <div className="text-center text-sm text-gray-500">
                  <p>请在 15 分钟内完成支付</p>
                  <p className="mt-2">
                    实际使用时，这里会显示真实的收款码图片
                  </p>
                </div>

                <button
                  onClick={() => {
                    setPaymentStep('select')
                    setCurrentOrder(null)
                    setSelectedProduct(null)
                  }}
                  className="mt-4 w-full text-gray-600 hover:text-gray-800"
                >
                  取消订单
                </button>
              </div>
            )}

            {paymentStep === 'success' && (
              <div className="max-w-md mx-auto bg-white rounded-lg shadow-md p-8 text-center">
                <div className="mb-4">
                  <svg
                    className="w-16 h-16 mx-auto text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">支付成功！</h2>
                <p className="text-gray-600 mb-4">License Key 已生成，请在"我的 Key"中查看</p>
                <p className="text-sm text-gray-500">3 秒后自动跳转...</p>
              </div>
            )}
          </>
        )}

        {/* Keys Tab */}
        {activeTab === 'keys' && (
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            {myKeys.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                <p>暂无 License Key</p>
                <button
                  onClick={() => setActiveTab('buy')}
                  className="mt-4 text-blue-600 hover:text-blue-700"
                >
                  立即购买
                </button>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Key Code
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        产品类型
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        额度
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        状态
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        创建时间
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {myKeys.map((key) => (
                      <tr key={key.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">
                            {key.key_code}
                          </code>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getProductTypeBadge(key.product_type)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {key.quota_used} / {key.quota_total}
                          <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                            <div
                              className="bg-blue-600 h-2 rounded-full"
                              style={{
                                width: `${(key.quota_used / key.quota_total) * 100}%`,
                              }}
                            ></div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getStatusBadge(key.status)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {new Date(key.created_at).toLocaleString('zh-CN')}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <button
                            onClick={() => applyKey(key.key_code)}
                            className={`px-3 py-1 rounded text-xs font-medium ${
                              activeKey === key.key_code
                                ? 'bg-green-100 text-green-700'
                                : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
                            }`}
                          >
                            {activeKey === key.key_code ? '已使用' : '设为当前'}
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default Payment
