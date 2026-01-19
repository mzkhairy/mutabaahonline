<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { Loader2, UserPlus, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  institution_code: '', // Kode Sekolah (Wajib ada di DB)
  name: '',
  username: '',
  email: '',
  password: '',
  role: 'MURID' // Default role
})

const handleRegister = async () => {
  const success = await authStore.register(form.value)
  if (success) {
    alert('Pendaftaran berhasil! Silakan login.')
    router.push('/') // Redirect ke halaman login
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4 py-10">
    <div class="max-w-md w-full bg-surface shadow-lg rounded-xl p-8 border border-gray-100">
      
      <button @click="router.push('/')" class="text-secondary hover:text-primary mb-6 flex items-center gap-2 text-sm transition-colors">
        <ArrowLeft class="w-4 h-4" /> Kembali ke Login
      </button>

      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-primary/10 text-primary mb-4">
          <UserPlus class="w-6 h-6" />
        </div>
        <h1 class="text-2xl font-bold text-accent">Daftar Akun Baru</h1>
        <p class="text-secondary text-sm">Bergabung dengan lembaga pendidikan Anda</p>
      </div>

      <div v-if="authStore.error" class="mb-4 p-3 bg-red-50 text-danger text-sm rounded-lg border border-red-100">
        {{ authStore.error }}
      </div>

      <form @submit.prevent="handleRegister" class="space-y-4">
        
        <div>
          <label class="block text-sm font-medium text-primary mb-1">Kode Lembaga</label>
          <input 
            v-model="form.institution_code"
            type="text" 
            placeholder="Contoh: TAHFIDZ01"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-primary bg-white uppercase"
            required
          />
          <p class="text-xs text-secondary mt-1">*Minta kode ini kepada Admin sekolah Anda</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-primary mb-1">Nama Lengkap</label>
          <input 
            v-model="form.name"
            type="text" 
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-primary bg-white"
            required
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-primary mb-1">Username</label>
            <input 
              v-model="form.username"
              type="text" 
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-primary bg-white"
              required
            />
          </div>
           <div>
            <label class="block text-sm font-medium text-primary mb-1">Email</label>
            <input 
              v-model="form.email"
              type="email" 
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-primary bg-white"
              required
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-primary mb-1">Password</label>
          <input 
            v-model="form.password"
            type="password" 
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all text-primary bg-white"
            required
            minlength="6"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-primary mb-2">Daftar Sebagai</label>
          <div class="grid grid-cols-2 gap-4">
            <label 
              class="cursor-pointer border rounded-lg p-3 flex items-center justify-center gap-2 transition-all"
              :class="form.role === 'MURID' ? 'bg-primary text-white border-primary' : 'bg-white text-secondary border-gray-200 hover:bg-gray-50'"
            >
              <input type="radio" v-model="form.role" value="MURID" class="hidden" />
              <span class="font-medium text-sm">Murid</span>
            </label>
            <label 
              class="cursor-pointer border rounded-lg p-3 flex items-center justify-center gap-2 transition-all"
              :class="form.role === 'GURU' ? 'bg-primary text-white border-primary' : 'bg-white text-secondary border-gray-200 hover:bg-gray-50'"
            >
              <input type="radio" v-model="form.role" value="GURU" class="hidden" />
              <span class="font-medium text-sm">Guru</span>
            </label>
          </div>
        </div>

        <button 
          type="submit" 
          :disabled="authStore.loading"
          class="w-full bg-primary hover:bg-accent text-white font-medium py-2.5 rounded-lg transition-colors flex justify-center items-center gap-2 mt-6"
        >
          <Loader2 v-if="authStore.loading" class="animate-spin w-5 h-5" />
          <span>{{ authStore.loading ? 'Mendaftar...' : 'Daftar Sekarang' }}</span>
        </button>

      </form>

    </div>
  </div>
</template>