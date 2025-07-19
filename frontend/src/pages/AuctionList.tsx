import { useState, useEffect } from 'react'
import { listAuctions, deleteAuction } from '../services/api'
import type { Auction } from '../services/api'
import { useNavigate } from 'react-router-dom'

export default function AuctionList() {
  const [auctions, setAuctions] = useState<Auction[]>([])
  const [loading,  setLoading]  = useState(true)
  const [error,    setError]    = useState<string|null>(null)
  const [role,     setRole]     = useState<string|null>(null)
  const [currentUserId, setUserId] = useState<number|null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        setUserId(payload.user_id)
        setRole(payload.role)
      } catch {
        setUserId(null)
        setRole(null)
      }
    }
  }, [])

  useEffect(() => {
    listAuctions(1, 10, '')   // page=1, size=10
      .then(res => {
        // res.data === { data: Auction[], page, size, total_count }
        setAuctions(res.data.data)
      })
      .catch(err => {
        setError(err.message)
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  const handleDelete = async (id: number) => {
    if (!window.confirm('本当にこのオークションを削除しますか？')) return
    try {
      await deleteAuction(id)
      setAuctions(prev => prev.filter(a => a.id !== id))
    } catch (err: any) {
      alert('削除に失敗しました: ' + err.message)
    }
  }
  
  const handleLogout = () => {
    localStorage.removeItem('token')
    navigate('/login')
  }


  if (loading) return <div className="p-4">ローディング中…</div>
  if (error)   return <div className="p-4 text-red-600">エラー: {error}</div>

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">オークションリスト</h1>
        <div className="flex space-x-2">
          {role === 'seller' && (
            <button
              onClick={() => navigate('/auctions/create')}
              className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
            >
              出品登録
            </button>
          )}
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition"
          >
            ログアウト
          </button>
        </div>
      </div>

      <ul className="space-y-2">
        {auctions.map(a => (
          <li
            key={a.id}
            className="flex justify-between items-center p-4 rounded hover:bg-gray-800 transition-colors"
          >
            <div
              className="flex-1 cursor-pointer"
              onClick={() => navigate(`/auctions/${a.id}`)}
            >
              <h2 className="text-xl font-semibold text-white">{a.title}</h2>
              <p className="text-sm text-gray-400">
                開始価格: {a.start_price.toLocaleString()}円 ・ 終了時間: {new Date(a.end_at).toLocaleString()}
              </p>
            </div>

            {role === 'seller' && currentUserId === a.seller_id && (
              <button
                onClick={e => {
                  e.stopPropagation()
                  handleDelete(a.id)
                }}
                className="
                  px-3 py-1
                  bg-red-600 hover:bg-red-700
                  text-sm text-white rounded
                  transition
                "
              >
                削除
              </button>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}