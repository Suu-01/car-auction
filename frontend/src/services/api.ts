import axios from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'

const BASE_URL = import.meta.env.DEV
  ? '' 
  : import.meta.env.VITE_API_BASE_URL

export const api = axios.create({
  baseURL: BASE_URL,     
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers = config.headers ?? {}
    // @ts-ignore
    config.headers.Authorization = `Bearer ${token}`
    console.log('[API DEBUG] token=', token, '→', config.method, config.url)
  }
  return config
})

export interface AuthResponse {
  token: string
  role: 'seller' | 'bidder'
}

export interface LoginRequest {
  email: string
  password: string
}

export interface SignupRequest {
  email: string
  password: string
  role: 'bidder' | 'seller'
}

export interface Auction {
  id: number
  title: string
  description: string
  start_price: number
  created_at: string
  end_at: string
  maker: string       
  model_name: string   
  mileage: number      
  year: number         
  photo_url: string
  seller_id: number
}

export interface Bid {
  id: number
  amount: number
  auction_id: number
  user_id: number
  created_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  size: number
  total_count: number
}

export const signup = (req: SignupRequest) =>
  api.post<AuthResponse>('/api/users/signup', req)

export const login = (email: string, password: string) =>
  api.post<AuthResponse>('/api/users/login', { email, password })

export const listAuctions = (
  page = 1,
  size = 10,
  title = ''
) =>
  api.get<PaginatedResponse<Auction>>(
    `/api/auctions?page=${page}&size=${size}&title=${encodeURIComponent(title)}`
  )

export const getAuction = (id: number) =>
  api.get<Auction>(`/api/auctions/${id}`)

export const listBids = (id: number, page = 1, size = 10) =>
  api.get<PaginatedResponse<Bid>>(
    `/api/auctions/${id}/bids?page=${page}&size=${size}`
  )
  
export interface CreateAuctionRequest {
  title: string
  description: string
  start_price: number
  end_at: string
  maker: string
  model_name: string
  mileage: number
  year: number
  photo_url: string
}

export const createAuction = (req: CreateAuctionRequest) =>
  api.post<Auction>('/api/auctions', req)

export interface CreateBidReq {
  amount: number
}
export const placeBid = (auctionId: number, req: CreateBidReq) =>
  api.post<Bid>(`/api/auctions/${auctionId}/bids`, req)

export function deleteAuction(id: number) {
  const token = localStorage.getItem("token")
  return axios.delete(`/api/auctions/${id}`, {
    headers: { Authorization: `Bearer ${token}` }
  })
}