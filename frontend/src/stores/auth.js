import { defineStore } from 'pinia'
import api from '../services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    mustChangePassword: (state) => state.user?.must_change_password === true
  },
  actions: {
    async login(payload) {
      this.loading = true
      this.error = null
      try {
        // payload = { institution_code, username, password }
        const res = await api.post('/auth/login', payload)
        
        // Ambil data dari response backend Go
        // Struktur response backend: { data: { access_token, user: {...} } }
        const { access_token, user } = res.data.data
        
        // Simpan ke State & LocalStorage
        this.token = access_token
        this.user = user
        localStorage.setItem('token', access_token)
        localStorage.setItem('user', JSON.stringify(user))
        
        return true // Login sukses
      } catch (err) {
        console.error(err)
        this.error = err.response?.data?.error?.message || 'Gagal login ke server'
        return false // Login gagal
      } finally {
        this.loading = false
      }
    },
    _formatError(err) {
      const resp = err.response?.data?.error
      if (!resp) return 'Terjadi kesalahan jaringan'

      // Jika ada details (error validasi form)
      if (resp.details) {
        // Jika details berupa string panjang dari Go Validator
        if (typeof resp.details === 'string') {
          if (resp.details.includes('Password') && resp.details.includes('min')) {
            return 'Password terlalu pendek (minimal 8 karakter)'
          }
          if (resp.details.includes('email')) {
            return 'Format email tidak valid'
          }
          if (resp.details.includes('required')) {
            return 'Mohon lengkapi semua kolom yang wajib diisi'
          }
          return resp.details // Fallback ke pesan asli jika tidak dikenal
        }
        // Jika details berupa object (validation errors yang sudah diparse backend)
        if (typeof resp.details === 'object') {
          const firstKey = Object.keys(resp.details)[0]
          return `${firstKey}: ${resp.details[firstKey]}`
        }
      }

      // Jika tidak ada details, pakai message utama
      return resp.message || 'Terjadi kesalahan pada server'
    },

    async register(payload) {
      this.loading = true; this.error = null
      try {
        await api.post('/auth/register', payload)
        return true
      } catch (err) {
        console.error(err)
        this.error = this._formatError(err) // <--- Pakai Helper tadi
        return false
      } finally {
        this.loading = false
      }
    },

    async registerInstitution(payload) {
      this.loading = true; this.error = null
      try {
        await api.post('/auth/register-institution', payload)
        return true
      } catch (err) {
        console.error(err)
        this.error = this._formatError(err) // <--- Pakai Helper tadi
        return false
      } finally {
        this.loading = false
      }
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      // Kita akan redirect di komponen nanti
    }
  }
})