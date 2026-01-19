import axios from 'axios'

// Setup koneksi dasar ke Backend Go
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // Sesuai router Go kamu
  headers: {
    'Content-Type': 'application/json'
  }
})

// Otomatis selipkan Token JWT jika sudah login
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default api