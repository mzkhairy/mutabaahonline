<script setup>
import { BookOpen } from 'lucide-vue-next'
</script>

<template>
  <div>
    <h2 class="text-lg font-semibold text-primary mb-4">Panel Pengajar</h2>
    
    <div class="bg-surface p-8 rounded-xl border border-gray-100 text-center shadow-sm">
      <div class="mx-auto w-12 h-12 bg-orange-50 text-orange-600 rounded-full flex items-center justify-center mb-4">
        <BookOpen class="w-6 h-6" />
      </div>
      <h3 class="text-xl font-bold text-primary">Kelas Ajar Anda</h3>
      <p class="text-secondary mb-6">Mulai sesi mutabaah atau absen murid.</p>
      
      <router-link 
        to="/guru/classes" 
        class="inline-flex px-6 py-2 bg-primary text-white rounded-lg hover:bg-accent transition font-medium"
      >
        Lihat Daftar Kelas
      </router-link>
    </div>
  </div>
</template>