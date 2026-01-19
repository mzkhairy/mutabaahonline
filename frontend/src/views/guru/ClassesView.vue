<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { BookOpen, ChevronRight, Users, Loader2 } from 'lucide-vue-next'

const router = useRouter()
const classes = ref([])
const academicYears = ref([])
const selectedAyId = ref('') // ID Tahun Ajaran Terpilih
const loading = ref(false)

// 1. Ambil Data Tahun Ajaran dulu
const fetchAcademicYears = async () => {
    try {
        const res = await api.get('/academic-years')
        academicYears.value = res.data.data || []
        
        // Default select yang aktif
        const active = academicYears.value.find(ay => ay.is_active)
        if (active) selectedAyId.value = active.id
        else if (academicYears.value.length > 0) selectedAyId.value = academicYears.value[0].id
        
        // Baru fetch kelas
        fetchClasses()
    } catch (e) { console.error(e) }
}

// 2. Fetch Kelas dengan Filter
const fetchClasses = async () => {
  if(!selectedAyId.value) return
  loading.value = true
  try {
    // Kirim parameter academic_year_id
    const res = await api.get('/guru/classes', {
        params: { academic_year_id: selectedAyId.value }
    })
    classes.value = res.data.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

// 3. Watch jika dropdown berubah
watch(selectedAyId, () => {
    fetchClasses()
})

onMounted(fetchAcademicYears)
</script>

<template>
  <DashboardLayout>
    <div class="max-w-4xl mx-auto space-y-6">
      
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
            <h1 class="text-2xl font-bold text-gray-800">Kelas Saya</h1>
            <p class="text-gray-500 text-sm">Pilih kelas untuk mulai mengelola kegiatan.</p>
        </div>
        
        <div class="relative w-full sm:w-auto">
            <select v-model="selectedAyId" class="w-full sm:w-48 appearance-none bg-white border border-gray-300 text-gray-700 py-2 pl-3 pr-8 rounded-lg text-sm font-medium focus:outline-none focus:ring-2 focus:ring-primary/50 cursor-pointer shadow-sm">
                <option v-for="ay in academicYears" :key="ay.id" :value="ay.id">
                    {{ ay.name }} {{ ay.is_active ? '(Aktif)' : '' }}
                </option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-500">
                <svg class="h-4 w-4 fill-current" viewBox="0 0 20 20"><path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"/></svg>
            </div>
        </div>
      </div>

      <div v-if="loading" class="py-12 flex justify-center text-gray-400">
        <Loader2 class="w-8 h-8 animate-spin" />
      </div>

      <div v-else-if="classes.length === 0" class="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
        <p class="text-gray-500 font-medium">Tidak ada kelas di tahun ajaran ini.</p>
      </div>

      <div v-else class="grid gap-4">
        <div v-for="cls in classes" :key="cls.id" 
             @click="router.push(`/guru/classes/${cls.id}`)"
             class="bg-white p-5 rounded-xl border border-gray-200 shadow-sm hover:border-primary/50 hover:shadow-md transition cursor-pointer group flex items-center justify-between">
          
          <div class="flex items-center gap-4">
            <div class="bg-blue-50 p-3 rounded-lg text-primary group-hover:bg-primary group-hover:text-white transition">
              <BookOpen class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-bold text-lg text-gray-800 group-hover:text-primary transition">{{ cls.name }}</h3>
              <div class="flex items-center gap-4 text-xs text-gray-500 mt-1">
                <span class="flex items-center gap-1"><Users class="w-3 h-3"/> {{ cls.student_count }} Murid</span>
                <span class="bg-gray-100 px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide">Level {{ cls.level }}</span>
              </div>
            </div>
          </div>

          <ChevronRight class="w-5 h-5 text-gray-300 group-hover:text-primary transition" />
        </div>
      </div>

    </div>
  </DashboardLayout>
</template>