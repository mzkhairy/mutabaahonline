<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { Plus, Users, Search, Loader2, Calendar, BookOpen, Trash2, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const classes = ref([])
const loading = ref(false)
const filterYearId = ref('')

const showModal = ref(false)
const form = ref({
  name: '',
  level: '',
  teacher_id: '',
  academic_year_id: ''
})

const teachers = ref([])
const academicYears = ref([])

const fetchClasses = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterYearId.value) params.academic_year_id = filterYearId.value
    const res = await api.get('/admin/classes', { params }) 
    classes.value = res.data.data || []
  } catch (err) { console.error(err) } 
  finally { loading.value = false }
}

const fetchMasterData = async () => {
    try {
        const [resT, resY] = await Promise.all([
            api.get('/users', { params: { role: 'GURU' } }),
            api.get('/academic-years')
        ])
        
        // [FIX] Ambil array users dari dalam object stats/users
        teachers.value = resT.data.data?.users || []
        
        academicYears.value = resY.data.data || []
        
        // Logic default filter tahun ajaran
        if (!filterYearId.value) {
            const active = academicYears.value.find(y => y.is_active)
            if (active) {
                filterYearId.value = active.id
                form.value.academic_year_id = active.id
            } else if (academicYears.value.length > 0) {
                filterYearId.value = academicYears.value[0].id
            }
        }
        
        fetchClasses()
    } catch (e) { console.error(e) }
}

const handleCreate = async () => {
    try {
        // Jika sedang memfilter tahun, gunakan tahun tersebut
        if (filterYearId.value) {
            form.value.academic_year_id = filterYearId.value
        }
        
        await api.post('/classes', form.value)
        showModal.value = false
        
        // Reset form
        form.value.name = ''
        form.value.level = ''
        form.value.teacher_id = ''
        
        fetchClasses()
    } catch (err) {
        alert('Gagal membuat kelas')
    }
}

const handleDelete = async (id) => {
    if(!confirm("Hapus kelas ini?")) return
    try {
        await api.delete(`/classes/${id}`)
        fetchClasses()
    } catch(e) { alert("Gagal menghapus") }
}

onMounted(() => {
    // Ambil filter dari URL jika ada
    if (route.query.academic_year_id) {
        filterYearId.value = route.query.academic_year_id
    }
    fetchMasterData()
})
</script>

<template>
  <DashboardLayout>
    <div class="space-y-6">
        
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
            <div class="flex items-center gap-3">
                <button @click="router.push('/admin/academic-years')" class="p-2 bg-gray-100 hover:bg-gray-200 rounded-full transition" title="Kembali ke Tahun Ajaran">
                    <ArrowLeft class="w-5 h-5 text-gray-600"/>
                </button>
                <div>
                    <h2 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                        <BookOpen class="w-5 h-5 text-primary"/> Manajemen Kelas
                    </h2>
                    <p class="text-sm text-gray-500">Kelola rombongan belajar.</p>
                </div>
            </div>
            
            <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                <div class="relative">
                    <select v-model="filterYearId" @change="fetchClasses" class="w-full sm:w-48 appearance-none bg-gray-50 border border-gray-300 text-gray-700 py-2 pl-3 pr-8 rounded-lg text-sm font-bold focus:outline-none cursor-pointer">
                        <option value="" disabled>Pilih Tahun Ajaran</option>
                        <option v-for="ay in academicYears" :key="ay.id" :value="ay.id">
                            {{ ay.name }} {{ ay.is_active ? '(Aktif)' : '' }}
                        </option>
                    </select>
                </div>

                <button @click="showModal = true" class="bg-primary hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-bold flex items-center justify-center gap-2 shadow-sm transition">
                    <Plus class="w-4 h-4"/> Buat Kelas
                </button>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden p-4">
            <div v-if="loading" class="py-10 text-center"><Loader2 class="w-8 h-8 animate-spin mx-auto text-primary"/></div>
            
            <div v-else-if="classes.length === 0" class="text-center py-10 text-gray-400 border border-dashed rounded-lg">
                Belum ada data kelas pada tahun ajaran ini.
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <div v-for="cls in classes" :key="cls.id" class="border border-gray-200 rounded-xl p-4 hover:border-primary/50 transition bg-gray-50/50">
                    <div class="flex justify-between items-start mb-2">
                        <div>
                            <h3 class="font-bold text-gray-800 text-lg">{{ cls.name }}</h3>
                            <span class="text-xs font-bold px-2 py-0.5 bg-white border rounded text-gray-500">Level {{ cls.level }}</span>
                        </div>
                        <button @click="handleDelete(cls.id)" class="text-red-400 hover:text-red-600"><Trash2 class="w-4 h-4"/></button>
                    </div>
                    
                    <div class="space-y-1 text-sm text-gray-600 mt-3">
                        <div class="flex items-center gap-2">
                            <Users class="w-4 h-4 text-gray-400"/>
                            <span>{{ cls.teacher_name || 'Belum ada Wali' }}</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <Calendar class="w-4 h-4 text-gray-400"/>
                            <span>{{ cls.academic_year_name || '-' }}</span>
                        </div>
                    </div>

                    <div class="mt-4 pt-3 border-t border-gray-200 flex justify-between items-center">
                        <span class="text-xs text-gray-500 font-medium">{{ cls.student_count }} Murid</span>
                        <router-link :to="`/admin/classes/${cls.id}`" class="text-primary text-sm font-bold hover:underline">Detail &rarr;</router-link>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm">
            <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
                <h3 class="font-bold text-lg mb-4">Buat Kelas Baru</h3>
                <form @submit.prevent="handleCreate" class="space-y-4">
                    
                    <div v-if="filterYearId" class="p-3 bg-blue-50 border border-blue-100 rounded text-sm text-blue-800">
                        <span class="font-bold block text-xs uppercase mb-1">Tahun Ajaran</span>
                        {{ academicYears.find(y => y.id === filterYearId)?.name }}
                    </div>
                    
                    <div v-else>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Tahun Ajaran</label>
                        <select v-model="form.academic_year_id" class="w-full border p-2 rounded-lg" required>
                            <option v-for="y in academicYears" :key="y.id" :value="y.id">{{ y.name }}</option>
                        </select>
                    </div>

                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Kelas</label>
                        <input v-model="form.name" type="text" placeholder="Contoh: 1A Ikhwan" class="w-full border p-2 rounded-lg" required>
                    </div>
                    
                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Level / Tingkat</label>
                        <input v-model="form.level" type="text" placeholder="Contoh: 1 / X / A" class="w-full border p-2 rounded-lg" required>
                    </div>

                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Wali Kelas</label>
                        <select v-model="form.teacher_id" class="w-full border p-2 rounded-lg">
                            <option value="">-- Pilih Guru --</option>
                            <option v-for="t in teachers" :key="t.id" :value="t.id">{{ t.name }}</option>
                        </select>
                    </div>

                    <div class="pt-2 flex justify-end gap-2">
                        <button type="button" @click="showModal = false" class="px-4 py-2 text-gray-600">Batal</button>
                        <button type="submit" class="px-4 py-2 bg-primary text-white rounded-lg font-bold">Simpan</button>
                    </div>
                </form>
            </div>
        </div>

    </div>
  </DashboardLayout>
</template>