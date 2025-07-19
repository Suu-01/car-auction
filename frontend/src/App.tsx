import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Signup from './pages/Signup'
import AuctionList from './pages/AuctionList'
import AuctionDetail from './pages/AuctionDetail'
import NotFound from './pages/NotFound'
import AuctionCreate from "./pages/AuctionCreate"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* 루트 접속 시 로그인 화면으로 */}
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/login" element={<Login />} />
        <Route path="/auctions/create" element={<AuctionCreate />} />
        <Route path="/auctions" element={<AuctionList />} />
        <Route path="/auctions/:id" element={<AuctionDetail />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  )
}
