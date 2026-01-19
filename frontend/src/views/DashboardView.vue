<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import DashboardLayout from '../layouts/DashboardLayout.vue'

// Import Sub-Components (HANYA Admin & Guru)
import DashboardAdmin from './dashboard/DashboardAdmin.vue'
import DashboardGuru from './dashboard/DashboardGuru.vue'
// [FIX] DashboardMurid TIDAK di-import disini agar tidak double layout

const authStore = useAuthStore()
const router = useRouter()
const role = authStore.user?.role

// Helper text role
const roleText = {
  'ADMIN': 'administrator',
  'GURU': 'guru',
  'MURID': 'murid'
}

onMounted(() => {
  // [SOLUSI DOUBLE LAYOUT]
  // Jika role adalah MURID, jangan render DashboardView.
  // Langsung alihkan ke '/mutabaah' yang sudah punya layout sendiri.
  if (role === 'MURID') {
    router.replace('/mutabaah')
  }
})
</script>

<template>
  <div v-if="role === 'MURID'" class="flex h-screen items-center justify-center bg-gray-50">
      <span class="text-gray-500 animate-pulse">Mengalihkan ke halaman utama...</span>
  </div>

  <DashboardLayout v-else>
    <div class="max-w-4xl mx-auto">
      
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-accent">Ahlan wa Sahlan, {{ authStore.user?.name }}!</h1>
        <p class="text-secondary mt-1">
          Selamat datang di dashboard {{ roleText[role] || 'User' }}.
        </p>
      </div>

      <DashboardAdmin v-if="role === 'ADMIN'" />
      <DashboardGuru v-else-if="role === 'GURU'" />
      
      <div v-else class="text-red-500">Role tidak dikenali atau akses ditolak.</div>

    </div>
  </DashboardLayout>
</template>