<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router' // [BARU] Import Router
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { 
  Calendar, Plus, Edit, Trash2, CheckCircle, Power, Loader2, Search, XCircle, ArrowRight
} from 'lucide-vue-next'

const router = useRouter() // [BARU]
const academicYears = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const searchQuery = ref('')
const processingId = ref(null)

const form = ref({
  id: '',
  name: '',
  start_date: '',
  end_date: ''
})

const fetchAcademicYears = async () => {
  loading.value = true
  try {
    const res = await api.get('/academic-years')
    academicYears.value = res.data.data || []
  } catch (err) { console.error(err) } 
  finally { loading.value = false }
}

const filteredYears = computed(() => {
  if (!searchQuery.value) return academicYears.value
  return academicYears.value.filter(y => y.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const openModal = (year = null) => {
  if (year) {
    isEdit.value = true
    form.value = { 
        ...year,
        // Format tanggal untuk input date HTML (YYYY-MM-DD)
        start_date: year.start_date.split('T')[0],
        end_date: year.end_date.split('T')[0]
    }
  } else {
    isEdit.value = false
    form.value = { id: '', name: '', start_date: '', end_date: '' }
  }
  showModal.value = true
}

const handleSave = async () => {
  try {
    const payload = { ...form.value }
    if (!isEdit.value) payload.is_active = false 
    else {
        const old = academicYears.value.find(y => y.id === form.value.id)
        payload.is_active = old ? old.is_active : false
    }
    
    if (isEdit.value) {
      await api.put(`/academic-years/${form.value.id}`, payload)
    } else {
      await api.post('/academic-years', payload)
    }
    await fetchAcademicYears()
    showModal.value = false
  } catch (err) {
    alert(err.response?.data?.error?.message || 'Gagal menyimpan')
  }
}

const handleDelete = async (id) => {
  if (!confirm('Yakin hapus tahun ajaran ini?')) return
  try {
    await api.delete(`/academic-years/${id}`)
    await fetchAcademicYears()
  } catch (err) { alert('Gagal menghapus') }
}

const handleActivate = async (year) => {
    if (year.is_active) return
    if (!confirm(`Aktifkan tahun ajaran ${year.name}? Tahun ajaran lain akan dinonaktifkan.`)) return

    processingId.value = year.id
    try {
        await api.put(`/academic-years/${year.id}`, { 
            name: year.name,
            // [FIX ERROR 500] Ambil hanya YYYY-MM-DD, buang bagian jam (T00:00:00Z)
            start_date: year.start_date.split('T')[0],
            end_date: year.end_date.split('T')[0],
            is_active: true 
        })
        await fetchAcademicYears()
    } catch (err) {
        alert('Gagal mengaktifkan: ' + (err.response?.data?.error?.message || 'Error Server'))
    } finally {
        processingId.value = null
    }
}

// [BARU] Fungsi Navigasi ke Daftar Kelas
const goToClasses = (yearId) => {
    router.push({ 
        path: '/admin/classes', 
        query: { academic_year_id: yearId } 
    })
}

onMounted(fetchAcademicYears)
</script>

<template>
  <DashboardLayout>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
            <div>
                <h2 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                    <Calendar class="w-5 h-5 text-primary"/> Tahun Ajaran
                </h2>
                <p class="text-sm text-gray-500">Kelola periode akademik sekolah.</p>
            </div>
            <button @click="openModal()" class="bg-primary hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-bold flex items-center gap-2 shadow-sm transition">
                <Plus class="w-4 h-4"/> Tambah Baru
            </button>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
            <div class="p-4 border-b border-gray-100">
                <div class="relative max-w-md">
                    <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"/>
                    <input v-model="searchQuery" type="text" placeholder="Cari tahun ajaran..." class="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-primary/50">
                </div>
            </div>

            <div v-if="loading" class="p-10 text-center"><Loader2 class="w-8 h-8 animate-spin mx-auto text-primary"/></div>

            <div v-else>
                <table class="w-full text-left text-sm hidden sm:table">
                    <thead class="bg-gray-50 text-gray-600 font-semibold border-b">
                        <tr>
                            <th class="px-6 py-3">Nama Periode</th>
                            <th class="px-6 py-3">Durasi</th>
                            <th class="px-6 py-3 text-center">Status</th>
                            <th class="px-6 py-3 text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        <tr v-for="year in filteredYears" :key="year.id" class="hover:bg-gray-50 group">
                            <td class="px-6 py-3 font-medium cursor-pointer text-blue-600 hover:text-blue-800 hover:underline" @click="goToClasses(year.id)">
                                {{ year.name }}
                                <ArrowRight class="inline w-3 h-3 ml-1 opacity-0 group-hover:opacity-100 transition"/>
                            </td>
                            <td class="px-6 py-3 text-gray-500 text-xs">
                                {{ new Date(year.start_date).toLocaleDateString('id') }} - {{ new Date(year.end_date).toLocaleDateString('id') }}
                            </td>
                            <td class="px-6 py-3 text-center">
                                <span v-if="year.is_active" class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold bg-green-100 text-green-700 border border-green-200">
                                    <CheckCircle class="w-3 h-3"/> AKTIF
                                </span>
                                <button v-else @click="handleActivate(year)" :disabled="processingId === year.id" class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-500 border border-gray-200 hover:bg-blue-50 hover:text-blue-600 hover:border-blue-200 transition">
                                    <Loader2 v-if="processingId === year.id" class="w-3 h-3 animate-spin"/>
                                    <Power v-else class="w-3 h-3"/> Aktifkan
                                </button>
                            </td>
                            <td class="px-6 py-3 flex justify-end gap-2">
                                <button @click="openModal(year)" class="p-1.5 text-blue-600 hover:bg-blue-50 rounded"><Edit class="w-4 h-4"/></button>
                                <button @click="handleDelete(year.id)" class="p-1.5 text-red-600 hover:bg-red-50 rounded"><Trash2 class="w-4 h-4"/></button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                
                <div class="sm:hidden divide-y divide-gray-100">
                    <div v-for="year in filteredYears" :key="year.id" class="p-4 space-y-3">
                        <div class="flex justify-between items-start">
                            <div @click="goToClasses(year.id)" class="cursor-pointer">
                                <h3 class="font-bold text-blue-700 flex items-center gap-2">
                                    {{ year.name }} <ArrowRight class="w-4 h-4"/>
                                </h3>
                                <p class="text-xs text-gray-500 mt-1">{{ new Date(year.start_date).toLocaleDateString('id') }} - {{ new Date(year.end_date).toLocaleDateString('id') }}</p>
                            </div>
                            <button @click="openModal(year)" class="text-gray-400 hover:text-blue-600"><Edit class="w-4 h-4"/></button>
                        </div>
                        <div v-if="year.is_active" class="w-full py-2 bg-green-100 text-green-800 text-center rounded-lg text-sm font-bold border border-green-200">
                            Sedang Aktif
                        </div>
                        <button v-else @click="handleActivate(year)" :disabled="processingId === year.id" class="w-full py-2 bg-white border border-gray-300 text-gray-600 hover:bg-blue-50 hover:text-blue-600 rounded-lg text-sm font-bold flex justify-center gap-2">
                            <Loader2 v-if="processingId === year.id" class="w-4 h-4 animate-spin"/>
                            <Power v-else class="w-4 h-4"/> Aktifkan
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm">
            <div class="bg-white rounded-xl shadow-xl w-full max-w-md animate-scale-in">
                <div class="px-6 py-4 border-b border-gray-100 flex justify-between items-center">
                    <h3 class="font-bold text-gray-800">{{ isEdit ? 'Edit Tahun Ajaran' : 'Tahun Ajaran Baru' }}</h3>
                    <button @click="showModal = false"><XCircle class="w-5 h-5 text-gray-400"/></button>
                </div>
                <div class="p-6 space-y-4">
                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Periode</label>
                        <input v-model="form.name" type="text" class="w-full border border-gray-300 rounded-lg p-2.5 text-sm outline-none focus:border-primary">
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Mulai</label>
                            <input v-model="form.start_date" type="date" class="w-full border border-gray-300 rounded-lg p-2.5 text-sm outline-none focus:border-primary">
                        </div>
                        <div>
                            <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Selesai</label>
                            <input v-model="form.end_date" type="date" class="w-full border border-gray-300 rounded-lg p-2.5 text-sm outline-none focus:border-primary">
                        </div>
                    </div>
                </div>
                <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
                    <button @click="showModal = false" class="px-4 py-2 text-gray-600 text-sm font-medium">Batal</button>
                    <button @click="handleSave" class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-bold shadow-sm">Simpan</button>
                </div>
            </div>
        </div>
    </div>
  </DashboardLayout>
</template>