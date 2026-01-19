<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { LogIn, Loader2, BookOpen } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  institution_code: '',
  username: '',
  password: ''
})

const handleLogin = async () => {
  // [FIX] Kirim data persis seperti input user (Case Sensitive)
  const success = await authStore.login(form.value)
  
  if (success) {
    router.push('/dashboard')
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="max-w-md w-full bg-white shadow-xl rounded-2xl p-8 border border-gray-100">
      
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 text-primary mb-4">
          <BookOpen class="w-8 h-8" />
        </div>
        <h1 class="text-2xl font-bold text-gray-900">Mutabaah App</h1>
        <p class="text-gray-500 text-sm mt-1">Masuk untuk melanjutkan aktivitas</p>
      </div>

      <div v-if="authStore.error" class="mb-6 p-4 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100 flex items-center gap-2">
        <span>⚠️</span> {{ authStore.error }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Kode Lembaga</label>
          <input 
            v-model="form.institution_code"
            type="text" 
            placeholder="Contoh: tahfidz01"
            class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition"
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username / NIP</label>
          <input 
            v-model="form.username"
            type="text" 
            class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition"
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input 
            v-model="form.password"
            type="password" 
            class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition"
            required
          />
        </div>

        <button 
          type="submit" 
          :disabled="authStore.loading"
          class="w-full bg-primary hover:bg-blue-700 text-white font-semibold py-3 rounded-lg transition-all shadow-sm hover:shadow-md flex justify-center items-center gap-2"
        >
          <Loader2 v-if="authStore.loading" class="animate-spin w-5 h-5" />
          <span>{{ authStore.loading ? 'Sedang Masuk...' : 'Masuk Sekarang' }}</span>
        </button>
      </form>

      <div class="mt-8 pt-6 border-t border-gray-100 text-center">
        <p class="text-sm text-gray-500 mb-2">Belum punya akun?</p>
        <div class="space-y-2">
            <div class="text-xs text-gray-400">Hubungi Admin Sekolah untuk mendapatkan akun.</div>
            
            <div class="pt-2">
                <router-link to="/register-institution" class="text-gray-400 hover:text-primary text-xs transition">Daftarkan Sekolah Baru</router-link>
            </div>
        </div>
      </div>

    </div>
  </div>
</template>