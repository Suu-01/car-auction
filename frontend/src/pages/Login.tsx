import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../services/api'

export default function Login() {
  const [email, setEmail]         = useState('')
  const [password, setPassword]   = useState('')
  const navigate                  = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const res = await login(email, password)
      const { token, role } = res.data
      localStorage.setItem('token', token)
      if (role === 'seller') {
        navigate('/auctions')
      } else {
        navigate('/auctions')
      }
    } catch (err) {
      console.error(err)
      alert('ログインに失敗しました。 メールアドレスとパスワードを確認してください。')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-sm mx-auto p-4 space-y-4">
      <h1 className="text-2xl mb-4">ログイン</h1>
      <input
        type="email" value={email}
        onChange={e=>setEmail(e.target.value)}
        placeholder="Email"
        className="w-full p-2 border mb-2"
      />
      <input
        type="password" value={password}
        onChange={e=>setPassword(e.target.value)}
        placeholder="Password"
        className="w-full p-2 border mb-4"
      />
      <button type="submit" className="w-full py-2 bg-blue-500 text-white rounded">ログイン</button>
      <div className="mt-4 text-center">
       <button
         type="button"
         onClick={() => navigate('/signup')}
         className="w-full py-2 bg-blue-500 text-white rounded"
       >
         会員登録
       </button>
     </div>
    </form>
  )
}
