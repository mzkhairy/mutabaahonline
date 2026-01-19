<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'
import { ShieldCheck, Loader2, LogOut, Mail } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const form = ref({
  username: authStore.user?.username || '', 
  email: '', // [BARU]
  password: '',
  confirmPassword: ''
})

const errorMsg = ref('')

const handleSetup = async () => {
  errorMsg.value = ''
  
  if (form.value.password !== form.value.confirmPassword) {
    errorMsg.value = 'Konfirmasi password tidak sesuai.'
    return
  }

  loading.value = true
  try {
    await api.post('/auth/setup-account', {
      username: form.value.username,
      password: form.value.password,
      email: form.value.email // [BARU]
    })

    // Update state lokal
    authStore.user.must_change_password = false
    authStore.user.username = form.value.username
    authStore.user.email = form.value.email
    localStorage.setItem('user', JSON.stringify(authStore.user))

    alert('Akun berhasil diaktifkan! Selamat datang.')
    router.push('/dashboard')
  } catch (err) {
    errorMsg.value = err.response?.data?.error?.message || 'Gagal memproses data.'
  } finally {
    loading.value = false
  }
}

const logout = () => {
    authStore.logout()
    router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
    <div class="bg-white max-w-md w-full p-8 rounded-2xl shadow-xl border border-gray-100">
      
      <div class="text-center mb-6">
        <div class="w-16 h-16 bg-blue-50 text-primary rounded-full flex items-center justify-center mx-auto mb-4">
          <ShieldCheck class="w-8 h-8" />
        </div>
        <h1 class="text-2xl font-bold text-accent">Aktivasi Akun Baru</h1>
        <p class="text-secondary text-sm mt-2">
          Demi keamanan, silakan lengkapi data akun Anda sebelum melanjutkan.
        </p>
      </div>

      <div v-if="errorMsg" class="mb-4 p-3 bg-red-50 text-red-600 text-xs rounded-lg border border-red-100">
        {{ errorMsg }}
      </div>

      <form @submit.prevent="handleSetup" class="space-y-4">
        <div>
          <label class="label">Username Baru</label>
          <input v-model="form.username" type="text" class="input-field w-full border p-2 rounded-lg" required minlength="4" />
          <p class="text-[10px] text-gray-400 mt-1">Gunakan huruf dan angka tanpa spasi (Min. 4 karakter).</p>
        </div>

        <div>
          <label class="label">Email Pemulihan</label>
          <div class="relative">
             <Mail class="w-4 h-4 absolute left-3 top-3 text-gray-400" />
             <input v-model="form.email" type="email" placeholder="contoh@email.com" class="input-field w-full border p-2 pl-9 rounded-lg" required />
          </div>
          <p class="text-[10px] text-gray-400 mt-1">Digunakan jika Anda lupa password.</p>
        </div>

        <div>
          <label class="label">Password Baru</label>
          <input v-model="form.password" type="password" class="input-field w-full border p-2 rounded-lg" required minlength="6" />
        </div>

        <div>
          <label class="label">Ulangi Password</label>
          <input v-model="form.confirmPassword" type="password" class="input-field w-full border p-2 rounded-lg" required />
        </div>

        <button type="submit" :disabled="loading" class="w-full btn-primary py-3 rounded-xl font-bold bg-primary text-white hover:bg-accent transition flex justify-center items-center gap-2 mt-4">
           <Loader2 v-if="loading" class="animate-spin w-5 h-5"/> Simpan & Lanjutkan
        </button>
      </form>

      <button @click="logout" class="mt-6 w-full text-center text-xs text-gray-400 hover:text-red-500 flex justify-center items-center gap-1 transition">
        <LogOut class="w-3 h-3" /> Batalkan & Keluar
      </button>

    </div>
  </div>
</template>