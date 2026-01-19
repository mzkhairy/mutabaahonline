<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api' 
import { 
  Building2, User, Lock, Loader2, ArrowLeft, School, CheckCircle, Mail 
} from 'lucide-vue-next'

const router = useRouter()
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const form = ref({
  institution_name: '',
  institution_code: '',
  admin_name: '',
  username: '',
  email: '',
  password: ''
})

const handleRegister = async () => {
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    // [FIX] Kirim kode lembaga apa adanya (tanpa toUpperCase)
    await api.post('/auth/register-institution', {
        institution_name: form.value.institution_name,
        institution_code: form.value.institution_code, 
        name: form.value.admin_name,
        username: form.value.username,
        email: form.value.email,
        password: form.value.password,
        role: 'ADMIN'
    })

    successMsg.value = 'Pendaftaran berhasil! Silakan masuk.'
    
    setTimeout(() => {
        router.push('/')
    }, 2000)

  } catch (err) {
    console.error(err)
    errorMsg.value = err.response?.data?.error?.message || 'Gagal mendaftar. Coba kode lain atau periksa koneksi.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-10">
    <div class="w-full max-w-md bg-white shadow-xl rounded-2xl p-6 sm:p-8 border border-gray-100">
      
      <div class="mb-8">
        <button @click="router.push('/')" class="text-gray-400 hover:text-gray-600 flex items-center gap-1 text-sm font-medium mb-4 transition">
            <ArrowLeft class="w-4 h-4" /> Kembali
        </button>
        <div class="text-center">
            <div class="inline-flex items-center justify-center w-14 h-14 rounded-full bg-primary/10 text-primary mb-3">
            <School class="w-7 h-7" />
            </div>
            <h1 class="text-2xl font-bold text-gray-900">Daftarkan Sekolah</h1>
            <p class="text-gray-500 text-sm mt-1">Buat akun lembaga & admin baru.</p>
        </div>
      </div>

      <div v-if="errorMsg" class="mb-6 p-4 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100 flex items-start gap-2">
        <span class="mt-0.5">⚠️</span> <span>{{ errorMsg }}</span>
      </div>

      <div v-if="successMsg" class="mb-6 p-4 bg-green-50 text-green-700 text-sm rounded-lg border border-green-100 flex items-center gap-2">
        <CheckCircle class="w-5 h-5" /> <span>{{ successMsg }}</span>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-4">
        
        <div class="space-y-4 border-b border-gray-100 pb-4">
            <div>
                <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Lembaga / Sekolah</label>
                <div class="relative">
                    <Building2 class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input 
                        v-model="form.institution_name"
                        type="text" 
                        placeholder="Contoh: Ponpes Al-Hidayah"
                        class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm"
                        required
                    />
                </div>
            </div>

            <div>
                <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Kode Lembaga (Unik)</label>
                <div class="relative">
                    <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 font-bold text-xs">#</span>
                    <input 
                        v-model="form.institution_code"
                        type="text" 
                        placeholder="alhidayah01"
                        class="w-full pl-8 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm font-medium"
                        required
                    />
                </div>
                <p class="text-[10px] text-gray-400 mt-1">*Kode ini akan digunakan guru & murid untuk login.</p>
            </div>
        </div>

        <div class="space-y-4 pt-2">
            <div>
                <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Admin</label>
                <div class="relative">
                    <User class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input 
                        v-model="form.admin_name"
                        type="text" 
                        placeholder="Nama Lengkap"
                        class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm"
                        required
                    />
                </div>
            </div>

            <div>
                <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Email</label>
                <div class="relative">
                    <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input 
                        v-model="form.email"
                        type="email" 
                        placeholder="admin@sekolah.com"
                        class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm"
                        required
                    />
                </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Username</label>
                    <input 
                        v-model="form.username"
                        type="text" 
                        class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm"
                        required
                    />
                </div>
                <div>
                    <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Password</label>
                    <input 
                        v-model="form.password"
                        type="password" 
                        class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition text-sm"
                        required
                    />
                </div>
            </div>
        </div>

        <div class="pt-4">
            <button 
                type="submit" 
                :disabled="loading"
                class="w-full bg-primary hover:bg-blue-700 text-white font-bold py-3 rounded-lg transition-all shadow-sm hover:shadow-md flex justify-center items-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed"
            >
                <Loader2 v-if="loading" class="animate-spin w-5 h-5" />
                <span>{{ loading ? 'Memproses...' : 'Daftarkan Sekolah' }}</span>
            </button>
        </div>

        <div class="text-center mt-4">
            <p class="text-xs text-gray-500">
                Sudah punya akun? 
                <router-link to="/" class="text-primary font-bold hover:underline">Masuk disini</router-link>
            </p>
        </div>

      </form>

    </div>
  </div>
</template>