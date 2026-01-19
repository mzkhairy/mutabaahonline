<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Loader2, Building2, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  institution_code: '',
  username: '',
  password: ''
})

const handleLogin = async () => {
  // Login Admin tetap menggunakan endpoint login biasa
  const success = await authStore.login(form.value)
  if (success) {
    router.push('/dashboard') 
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <div class="max-w-md w-full bg-surface shadow-lg rounded-xl p-8 border border-gray-100">
      
      <button @click="router.push('/')" class="text-secondary hover:text-primary mb-6 flex items-center gap-2 text-sm transition-colors">
        <ArrowLeft class="w-4 h-4" /> Kembali ke Login Utama
      </button>

      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-primary/10 text-primary mb-4">
          <Building2 class="w-6 h-6" />
        </div>
        <h1 class="text-2xl font-bold text-accent mb-2">Login Admin Institusi</h1>
        <p class="text-secondary text-sm">Masuk untuk mengelola lembaga Anda</p>
      </div>

      <div v-if="authStore.error" class="mb-4 p-3 bg-red-50 text-danger text-sm rounded-lg border border-red-100">
        {{ authStore.error }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-primary mb-1">Kode Lembaga</label>
          <input v-model="form.institution_code" type="text" class="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:border-primary uppercase" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-primary mb-1">Username Admin</label>
          <input v-model="form.username" type="text" class="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:border-primary" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-primary mb-1">Password</label>
          <input v-model="form.password" type="password" class="w-full px-4 py-2 border border-gray-300 rounded-lg outline-none focus:border-primary" required />
        </div>

        <button type="submit" :disabled="authStore.loading" class="w-full bg-primary hover:bg-accent text-white font-medium py-2.5 rounded-lg transition-colors flex justify-center items-center gap-2">
          <Loader2 v-if="authStore.loading" class="animate-spin w-5 h-5" />
          <span>{{ authStore.loading ? 'Masuk Dashboard' : 'Masuk' }}</span>
        </button>

        <div class="text-center mt-6 pt-4 border-t border-gray-100">
          <p class="text-sm text-secondary">
            Institusi belum terdaftar? 
            <router-link to="/register-institution" class="font-semibold text-primary hover:text-accent transition-colors">
              Buat Institusi Baru
            </router-link>
          </p>
        </div>
      </form>
    </div>
  </div>
</template>