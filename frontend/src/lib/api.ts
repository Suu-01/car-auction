import axios from 'axios'

const api = axios.create({
  baseURL: '/',            // vite proxy를 쓰고 있다면 슬래시만
})

// 요청 전마다 로컬 스토리지에서 토큰 꺼내서 헤더에 붙이기
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  // config.headers 가 undefined 가 아닐 때만
  if (token && config.headers) {
    // 기존 헤더 객체에 Authorization 속성만 추가
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default api