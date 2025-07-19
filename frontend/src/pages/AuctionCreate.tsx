import { useState } from "react"
import { createAuction, api } from "../services/api"
import { useNavigate } from "react-router-dom"

export default function AuctionCreate() {
  const [title, setTitle]         = useState('')
  const [desc,  setDesc]          = useState('')
  const [price, setPrice]         = useState(0)
  const [maker, setMaker]         = useState('')
  const [modelName, setModelName] = useState('')
  const [mileage, setMileage]     = useState(0)
  const [year, setYear]           = useState(2020)
  const [photoUrl, setPhotoUrl]      = useState('')
  const [endAt, setEndAt]         = useState('')
  const navigate                  = useNavigate()

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
  const file = e.currentTarget.files?.[0]
  if (!file) return

  const form = new FormData()
  form.append('file', file)

  try {
    const res = await api.post<{ url: string }>('/api/upload', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    setPhotoUrl(res.data.url)
  } catch {
    alert('写真のアップロードに失敗しました。')
  }
}

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const res = await createAuction({
        title,
        description: desc,
        start_price: price,
        maker,
        model_name: modelName,
        mileage,
        year,
        photo_url: photoUrl,
        end_at: new Date(endAt).toISOString(),
      })
      navigate(`/auctions/${res.data.id}`)
    } catch (err: any) {
      alert('出品失敗: '+ err.message)
    }
  }

  return (
    <div className="p-6 max-w-3xl mx-auto space-y-6">
     <button
       onClick={() => navigate(-1)}
       className="text-sm text-gray-600 underline">← バック</button>

      <h1 className="text-2xl font-bold mb-4">車両出品</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block mb-1 font-medium">タイトル</label>
          <input
            type="text"
            placeholder="タイトルを入力してください"
            value={title}
            onChange={e => setTitle(e.target.value)}
            required
            className="w-full p-2 border rounded"
          />
        </div>

        <div>
          <label className="block mb-1 font-medium">車両説明</label>
          <textarea
            placeholder="車両説明を入力してください"
            value={desc}
            onChange={e => setDesc(e.target.value)}
            rows={4}
            required
            className="w-full p-2 border rounded resize-none"
          />
        </div>

        <div>
          <label className="block mb-1 font-medium">開始価格 (円)</label>
          <input
            type="number"
            placeholder="ex: 1,000,000"
            value={price}
            onChange={e => setPrice(+e.target.value)}
            min={0}
            required
            className="w-full p-2 border rounded"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block mb-1 font-medium">メーカー名</label>
            <input
              type="text"
              placeholder="ex: TOYOTA"
              value={maker}
              onChange={e => setMaker(e.target.value)}
              required
              className="w-full p-2 border rounded"
            />
          </div>
          <div>
            <label className="block mb-1 font-medium">車種名</label>
            <input
              type="text"
              placeholder="ex: Prius"
              value={modelName}
              onChange={e => setModelName(e.target.value)}
              required
              className="w-full p-2 border rounded"
            />
          </div>
          <div>
            <label className="block mb-1 font-medium">走行距離 (km)</label>
            <input
              type="number"
              placeholder="예: 50,000"
              value={mileage}
              onChange={e => setMileage(+e.target.value)}
              min={0}
              required
              className="w-full p-2 border rounded"
            />
          </div>
          <div>
            <label className="block mb-1 font-medium">年式 (年)</label>
            <input
              type="number"
              placeholder={`1900 ~ ${new Date().getFullYear()}`}
              value={year}
              onChange={e => setYear(+e.target.value)}
              min={1900}
              max={new Date().getFullYear()}
              required
              className="w-full p-2 border rounded"
            />
          </div>
        </div>

        <div>
          <label className="block mb-1 font-medium">写真選択</label>
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className="block"
          />
          {photoUrl && (
            <img
              src={photoUrl}
              alt="preview"
              className="mt-2 max-h-40 object-contain rounded border"
            />
          )}
        </div>

        <div>
          <label className="block mb-1 font-medium">締切日時</label>
          <input
            type="datetime-local"
            value={endAt}
            onChange={e => setEndAt(e.target.value)}
            required
            className="w-full p-2 border rounded"
          />
        </div>

        <button
          type="submit"
          className="w-full py-2 bg-green-600 text-white font-semibold rounded hover:bg-green-700"
        >
          出品する
        </button>
      </form>
    </div>
  )
}
