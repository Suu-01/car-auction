import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getAuction, listBids, placeBid } from '../services/api'
import type { Auction, Bid } from '../services/api'

export default function AuctionDetail() {
  const { id } = useParams<{id: string}>()
  const navigate = useNavigate()

  const [auction, setAuction] = useState<Auction | null>(null)
  const [bids, setBids] = useState<Bid[]>([])
  const [page, _setPage] = useState(1)
  const [size] = useState(10)
  const [amount, setAmount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [role, setRole] = useState<string | null>(null)
  const [currentUserId, setCurrentUserId] = useState<number | null>(null)
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        setRole(payload.role)
        setCurrentUserId(payload.user_id)
      } catch {
        setRole(null)
        setCurrentUserId(null)
      }
    }
  }, [])

  useEffect(() => {
    if (!id) return
    setLoading(true)
    getAuction(+id)
      .then(res => setAuction(res.data))
      .catch(err => setError(err.message))
    listBids(+id, page, size)
      .then(res => setBids(res.data.data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [id, page, size])

  const currentBid = auction
    ? Math.max(auction.start_price, ...(bids.map(b => b.amount)))
    : 0

  const auctionRef = useRef<Auction | null>(null)
  useEffect(() => { auctionRef.current = auction }, [auction])
  useEffect(() => {
    if (!id || currentUserId == null) return
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const wsUrl = `${protocol}://${window.location.host}/ws/auctions/${id}`
    const socket = new WebSocket(wsUrl)
    socket.onmessage = e => {
      const ev: { bid: Bid; new_end_at?: string } = JSON.parse(e.data)
      if (ev.bid.user_id !== currentUserId) {
        setBids(prev => [ev.bid, ...prev])
      }
      if (ev.new_end_at && auctionRef.current) {
        setAuction({ ...auctionRef.current, end_at: ev.new_end_at })
      }
    }
    socket.onerror = () => socket.close()
    return () => socket.close()
  }, [id, currentUserId])

  const handleBid = async () => {
    if (!id || amount <= 0) return
    if (amount <= currentBid) {
      alert('現在の入札価格より低い価格は登録できません')
      return
    }
    try {
      const res = await placeBid(+id, { amount })
      setBids(prev => [res.data, ...prev])
      setAmount(0)
    } catch (err: any) {
      alert('入札失敗: ' + err.message)
    }
  }

  if (loading) return <div className="p-4">ローディング中…</div>
  if (error) return <div className="p-4 text-red-600">エラー: {error}</div>
  if (!auction) return <div className="p-4">オークション情報が見つかりません</div>

  return (
    <div className="p-6 max-w-3xl mx-auto space-y-6">
      <button onClick={() => navigate(-1)} className="text-sm underline">← バック</button>

    <h1 className="text-3xl font-bold">{auction.title}</h1>

    {auction.photo_url
      ? <img src={`http://localhost:8080${auction.photo_url}`} alt={auction.title} className="max-w-[400px] max-h-[300px] w-auto h-auto object-contain rounded shadow" />
      : <div className="w-full h-64 bg-gray-200 flex items-center justify-center rounded">写真なし</div>
    }

    <ul className="grid grid-cols-2 gap-4 text-gray-700">
      <li><strong>車両説明:</strong> {auction.description}</li>
      <li><strong>メーカー名:</strong> {auction.maker || '情報なし'}</li>
      <li><strong>車種名:</strong> {auction.model_name || '情報なし'}</li>
      <li>
        <strong>走行距離:</strong>
        {auction.mileage != null
          ? auction.mileage.toLocaleString() + ' km'
          : '情報なし'}
      </li>
      <li><strong>年式:</strong> {auction.year || '情報なし'}</li>
      <li>
        <strong>開始価格:</strong>
        {auction.start_price != null
          ? auction.start_price.toLocaleString() + '円'
          : '情報なし'}
      </li>
      <li>
        <strong>終了時間:</strong>
        {auction.end_at
          ? new Date(auction.end_at).toLocaleString()
          : '情報なし'}
      </li>
    </ul>

      {role === 'bidder' && (
        <div className="flex items-center space-x-2">
          <input
            type="number"
            value={amount}
            onChange={e => setAmount(+e.target.value)}
            placeholder="入札価格入力"
            className="border p-2 rounded w-32"
          />
          <button onClick={handleBid}
                  className="px-4 py-2 bg-blue-600 text-white rounded">
            入札する
          </button>
        </div>
      )}

      <div className="mt-4">
        <span className="font-semibold">現在の入札価格:</span>{' '}
        <span className="text-lg">{currentBid.toLocaleString()}円</span>
      </div>
    </div>
  )
}

