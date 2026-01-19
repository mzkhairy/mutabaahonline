<script setup>
import { ref, onMounted } from 'vue'
import api from '../../services/api'
import { Users, FileText, Calendar, BookOpen, UserCheck, CheckCircle } from 'lucide-vue-next'

const stats = ref({ guru_active: 0, murid_active: 0, total_users: 0 })
const loading = ref(false)

const fetchStats = async () => {
  loading.value = true
  try {
    const res = await api.get('/users', { params: { limit: 1 } }) 
    if (res.data.data.stats) {
      stats.value = res.data.data.stats
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(fetchStats)
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:justify-between sm:items-end gap-2">
        <div>
            <h2 class="text-xl font-bold text-accent">Panel Administrator</h2>
            <p class="text-secondary text-sm">Ringkasan aktivitas lembaga Anda.</p>
        </div>
        <div class="self-start sm:self-auto">
             <p class="text-xs text-secondary bg-white border border-gray-200 px-3 py-1 rounded-full shadow-sm">
                Total User: <span class="font-bold text-primary">{{ stats.total_registered }}</span>
            </p>
        </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="bg-gradient-to-br from-blue-500 to-blue-600 text-white p-6 rounded-xl shadow-lg flex items-center justify-between relative overflow-hidden">
            <div class="z-10">
                <p class="text-blue-100 text-sm font-medium mb-1">Guru Aktif</p>
                <h3 class="text-3xl font-bold">{{ stats.guru_active }}</h3>
            </div>
            <div class="p-3 bg-white/20 rounded-lg z-10">
                <UserCheck class="w-8 h-8 text-white" />
            </div>
            <div class="absolute -bottom-4 -right-4 w-24 h-24 bg-white/10 rounded-full"></div>
        </div>

        <div class="bg-gradient-to-br from-green-500 to-green-600 text-white p-6 rounded-xl shadow-lg flex items-center justify-between relative overflow-hidden">
            <div class="z-10">
                <p class="text-green-100 text-sm font-medium mb-1">Murid Aktif</p>
                <h3 class="text-3xl font-bold">{{ stats.murid_active }}</h3>
            </div>
            <div class="p-3 bg-white/20 rounded-lg z-10">
                <CheckCircle class="w-8 h-8 text-white" />
            </div>
            <div class="absolute -bottom-4 -right-4 w-24 h-24 bg-white/10 rounded-full"></div>
        </div>
    </div>

    <h3 class="text-sm font-bold text-gray-500 uppercase tracking-wider mt-4">Menu Utama</h3>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      
      <div class="bg-white p-5 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition group cursor-pointer active:scale-95 duration-200" @click="$router.push('/admin/users')">
        <div class="flex items-start gap-4">
            <div class="p-3 bg-blue-50 text-blue-600 rounded-lg group-hover:bg-blue-600 group-hover:text-white transition shrink-0">
                <Users class="w-6 h-6" />
            </div>
            <div>
                <h3 class="font-bold text-primary mb-1 text-base">Kelola Akun</h3>
                <p class="text-sm text-secondary leading-snug">Tambah & atur guru/murid.</p>
            </div>
        </div>
      </div>

      <div class="bg-white p-5 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition group cursor-pointer active:scale-95 duration-200" @click="$router.push('/admin/academic-years')">
        <div class="flex items-start gap-4">
            <div class="p-3 bg-green-50 text-green-600 rounded-lg group-hover:bg-green-600 group-hover:text-white transition shrink-0">
                <BookOpen class="w-6 h-6" />
            </div>
            <div>
                <h3 class="font-bold text-primary mb-1 text-base">Manajemen Kelas</h3>
                <p class="text-sm text-secondary leading-snug">Atur kelas & tahun ajaran.</p>
            </div>
        </div>
      </div>

      <div class="bg-white p-5 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition group cursor-pointer active:scale-95 duration-200" @click="$router.push('/admin/templates')">
        <div class="flex items-start gap-4">
            <div class="p-3 bg-purple-50 text-purple-600 rounded-lg group-hover:bg-purple-600 group-hover:text-white transition shrink-0">
                <FileText class="w-6 h-6" />
            </div>
            <div>
                <h3 class="font-bold text-primary mb-1 text-base">Template Penilaian</h3>
                <p class="text-sm text-secondary leading-snug">Desain format mutabaah.</p>
            </div>
        </div>
      </div>

    </div>
  </div>
</template>