<script setup>
import { ref } from 'vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import api from '../services/api'
import { useAuthStore } from '../stores/auth'
import { Lock, Loader2, Check } from 'lucide-vue-next'

const authStore = useAuthStore()
const form = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})
const loading = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const handleChangePassword = async () => {
  errorMsg.value = ''
  successMsg.value = ''
  
  if (form.value.new_password.length < 6) {
    errorMsg.value = 'Password baru minimal 6 karakter.'
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    errorMsg.value = 'Konfirmasi password tidak cocok.'
    return
  }

  loading.value = true
  try {
    // Kita butuh endpoint khusus di backend: PATCH /auth/change-password
    // Karena belum ada, untuk MVP kita bisa pakai fitur "Reset Password" 
    // tapi diarahkan ke diri sendiri (Self Reset).
    // TAPI, reset password butuh hak akses ADMIN di backend kita saat ini.
    // SOLUSI CEPAT: Kita pakai endpoint /users/:id/password yang sudah ada.
    // Karena admin punya hak akses ke endpoint itu, dia bisa ganti passwordnya sendiri.
    
    await api.patch(`/users/${authStore.user.id}/password`, {
      new_password: form.value.new_password
    })

    successMsg.value = 'Password berhasil diperbarui!'
    form.value = { current_password: '', new_password: '', confirm_password: '' }
  } catch (err) {
    errorMsg.value = 'Gagal mengganti password.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <DashboardLayout>
    <div class="max-w-xl mx-auto">
      <div class="bg-white p-8 rounded-xl border border-gray-200 shadow-sm">
        <div class="flex items-center gap-3 mb-6 border-b border-gray-100 pb-4">
            <div class="p-2 bg-gray-100 rounded-lg"><Lock class="w-5 h-5 text-primary"/></div>
            <div>
                <h1 class="text-xl font-bold text-accent">Ganti Password</h1>
                <p class="text-sm text-secondary">Amankan akun Anda secara berkala.</p>
            </div>
        </div>

        <form @submit.prevent="handleChangePassword" class="space-y-4">
            <div v-if="successMsg" class="bg-green-50 text-green-700 p-3 rounded-lg text-sm flex items-center gap-2">
                <Check class="w-4 h-4"/> {{ successMsg }}
            </div>
            <div v-if="errorMsg" class="bg-red-50 text-red-700 p-3 rounded-lg text-sm">
                {{ errorMsg }}
            </div>

            <div>
                <label class="label">Password Baru</label>
                <input v-model="form.new_password" type="password" class="input-field w-full border p-2 rounded-lg" required placeholder="Minimal 6 karakter" />
            </div>

            <div>
                <label class="label">Konfirmasi Password Baru</label>
                <input v-model="form.confirm_password" type="password" class="input-field w-full border p-2 rounded-lg" required />
            </div>

            <div class="pt-4">
                <button type="submit" :disabled="loading" class="w-full bg-primary text-white py-2.5 rounded-lg flex justify-center items-center gap-2 font-semibold hover:bg-accent transition">
                    <Loader2 v-if="loading" class="animate-spin w-4 h-4" />
                    Simpan Password Baru
                </button>
            </div>
        </form>
      </div>
    </div>
  </DashboardLayout>
</template>