import { useState } from 'react'
import { signup } from '../services/api'
import { useNavigate } from 'react-router-dom'

export default function Signup() {
  const [email, setEmail]     = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole]         = useState<'bidder'|'seller'>('bidder')
  const navigate = useNavigate()
  const [error, setError]       = useState<string|null>(null)

   const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await signup({ email, password, role })
      navigate("/login")
    } catch (err: any) {
      setError("登録されているメールアドレスです。")
    }
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-sm mx-auto p-4 space-y-4">
      <h1 className="text-2xl mb-4">会員登録</h1>
      {error && <div className="mb-2 text-red-600">{error}</div>}
      <input
        type="email"
        value={email}
        onChange={e => setEmail(e.target.value)}
        placeholder="Email"
        required
        className="w-full p-2 border rounded"
      />
      <input
        type="password"
        value={password}
        onChange={e => setPassword(e.target.value)}
        placeholder="Password"
        required
        className="w-full p-2 border rounded"
      />
      <fieldset className="mb-4">
        <legend className="mb-1">役割</legend>
        <label className="mr-4">
          <input
            type="radio"
            value="bidder"
            checked={role==="bidder"}
            onChange={()=>setRole("bidder")}
          /> 入札者
        </label>
        <label>
          <input
            type="radio"
            value="seller"
            checked={role==="seller"}
            onChange={()=>setRole("seller")}
          /> 出品者
        </label>
      </fieldset>
      <button type="submit" className="w-full py-2 bg-green-500 text-white rounded">
        登録する
      </button>
      <button
        type="button"
        onClick={() => navigate('/login')}
        className="w-full mt-2 py-2 bg-blue-600 text-white rounded"
      >
        ログインへ
      </button>
    </form>
  )
}

