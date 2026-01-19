<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '../../stores/auth'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue' 
import { 
  Save, Loader2, CheckCircle, AlertCircle, 
  BookOpen, Calendar, Edit3, ChevronRight 
} from 'lucide-vue-next'

const authStore = useAuthStore()
const loading = ref(false)
const fetchingDetails = ref(false)
const sessions = ref([])
const saving = ref(false)

const selectedSessionId = ref(null)

const inputData = ref({
    attendance: 'HADIR',
    note: '',
    data: {} 
})

// --- FETCH DATA ---
const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.get('/sessions/student')
    const rawSessions = res.data.data || []
    
    sessions.value = rawSessions.map(s => {
      let structure = []
      try {
        structure = Array.isArray(s.template_structure) 
          ? s.template_structure 
          : JSON.parse(s.template_structure || '[]')
      } catch (e) { structure = [] }

      // Filter Field: Backend sudah filter 'student_activity', tapi kita double check
      const fields = structure.filter(f => f.type === 'student_activity' || !f.type) // !f.type fallback

      return { ...s, fields }
    })

    // Hanya pilih dari sesi yang AKTIF (Todo)
    const firstTodo = todoSessions.value.find(s => !s.has_entry) || todoSessions.value[0]
    if (firstTodo) {
        selectSession(firstTodo)
    }

  } catch (err) {
    console.error("Gagal load session:", err)
  } finally {
    loading.value = false
  }
}

// --- LOGIC SELECTION ---
const selectSession = async (session) => {
    selectedSessionId.value = session.id
    
    // 1. Reset Form
    inputData.value = {
        attendance: 'HADIR',
        note: '',
        data: {} 
    }

    // Init keys
    session.fields.forEach(f => {
        if (f.name) inputData.value.data[f.name] = '' 
    })
    
    // 2. Jika sudah ada isian (Edit Mode), fetch detail
    if (session.has_entry) {
        fetchingDetails.value = true
        try {
            const res = await api.get('/mutabaah', { params: { session_id: session.id } })
            const entries = res.data.data.entries || []
            const myEntry = entries.find(e => e.student_id === authStore.user.id) || entries[0]
            
            if (myEntry) {
                inputData.value.attendance = myEntry.attendance
                inputData.value.note = myEntry.note || ''
                // Merge data nilai
                inputData.value.data = {
                    ...inputData.value.data,
                    ...(myEntry.student_activity_data || {})
                }
            }
        } catch (e) {
            console.error(e)
        } finally {
            fetchingDetails.value = false
        }
    }
}

// Hanya tampilkan sesi yang is_student_input_allowed = TRUE
const todoSessions = computed(() => {
    return sessions.value.filter(s => s.is_student_input_allowed)
})

const currentSession = computed(() => {
    return sessions.value.find(s => s.id === selectedSessionId.value)
})

// --- SAVE ---
const handleSave = async () => {
  if (!selectedSessionId.value) return
  
  saving.value = true
  try {
    await api.post('/mutabaah', {
      session_id: selectedSessionId.value,
      student_id: authStore.user.id,
      attendance: inputData.value.attendance,
      note: inputData.value.note,
      student_activity_data: inputData.value.data,
      status: 'LANJUT'
    })

    alert(`Alhamdulillah, laporan berhasil tersimpan!`)
    
    // Update local state agar UI refresh
    if (currentSession.value) {
        currentSession.value.has_entry = true
    }
    
  } catch (err) {
    console.error(err)
    alert(err.response?.data?.error?.message || 'Gagal menyimpan laporan.')
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <DashboardLayout>
    <div class="max-w-5xl mx-auto space-y-8 pb-10">
        
        <div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm flex flex-col md:flex-row justify-between items-center gap-4">
            <div>
                <h1 class="text-xl font-bold text-gray-800">Mutabaah Siswa</h1>
                <p class="text-sm text-gray-500">
                    Ahlan wa Sahlan, <span class="font-semibold text-primary">{{ authStore.user?.name }}</span>
                </p>
            </div>
            <div v-if="todoSessions.length > 0" class="bg-orange-50 text-orange-700 px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 border border-orange-100">
                <AlertCircle class="w-4 h-4"/>
                <span>Ada {{ todoSessions.length }} sesi aktif.</span>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 shadow-md overflow-hidden relative">
            
            <div v-if="fetchingDetails" class="absolute inset-0 bg-white/80 z-20 flex items-center justify-center backdrop-blur-sm">
                <div class="flex flex-col items-center gap-2 text-primary font-medium">
                    <Loader2 class="w-8 h-8 animate-spin"/>
                    <span>Mengambil data isian...</span>
                </div>
            </div>

            <div class="bg-gray-50 px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                <h2 class="font-bold text-gray-800 flex items-center gap-2">
                    <Edit3 class="w-5 h-5 text-primary"/> 
                    Form Input Laporan
                </h2>
                <div v-if="currentSession" class="text-xs text-gray-500 bg-white px-2 py-1 rounded border border-gray-200">
                    Template: {{ currentSession.template_name }}
                </div>
            </div>

            <div class="p-6">
                <div v-if="loading" class="py-10 text-center text-gray-400 flex flex-col items-center gap-2">
                    <Loader2 class="w-8 h-8 animate-spin"/> Memuat data...
                </div>

                <div v-else-if="!currentSession" class="py-12 text-center text-gray-400 bg-gray-50/50 rounded-lg border border-dashed border-gray-300">
                    <p>Tidak ada sesi yang dipilih.</p>
                    <p class="text-sm mt-1">Silakan pilih sesi dari daftar di bawah.</p>
                </div>

                <div v-else class="space-y-6">
                    <div class="flex items-start gap-4 mb-6 bg-blue-50 p-4 rounded-lg border border-blue-100">
                        <div class="bg-white p-2 rounded text-center min-w-[60px] shadow-sm">
                            <div class="text-xl font-bold text-primary">{{ new Date(currentSession.date).getDate() }}</div>
                            <div class="text-[10px] uppercase text-gray-500">{{ new Date(currentSession.date).toLocaleString('id', { month: 'short' }) }}</div>
                        </div>
                        <div>
                            <h3 class="font-bold text-lg text-gray-800">{{ currentSession.name }}</h3>
                            <p class="text-sm text-gray-600">{{ currentSession.class_name }}</p>
                            <span v-if="currentSession.has_entry" class="mt-1 inline-flex items-center gap-1 text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded font-bold">
                               <CheckCircle class="w-3 h-3"/> Data Tersimpan
                            </span>
                        </div>
                    </div>

                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-2">Status Kehadiran</label>
                        <div class="flex gap-4">
                            <label class="flex items-center gap-2 cursor-pointer">
                                <input type="radio" v-model="inputData.attendance" value="HADIR" class="w-4 h-4 text-primary focus:ring-primary">
                                <span class="text-sm font-medium">Hadir</span>
                            </label>
                            <label class="flex items-center gap-2 cursor-pointer">
                                <input type="radio" v-model="inputData.attendance" value="IZIN" class="w-4 h-4 text-yellow-500 focus:ring-yellow-500">
                                <span class="text-sm font-medium">Izin</span>
                            </label>
                            <label class="flex items-center gap-2 cursor-pointer">
                                <input type="radio" v-model="inputData.attendance" value="SAKIT" class="w-4 h-4 text-red-500 focus:ring-red-500">
                                <span class="text-sm font-medium">Sakit</span>
                            </label>
                        </div>
                    </div>

                    <div v-if="inputData.attendance === 'HADIR'" class="bg-gray-50 p-4 rounded-lg border border-gray-200">
                        <h3 class="text-sm font-bold text-gray-700 mb-4 border-b pb-2">Capaian Personal</h3>
                        
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div v-for="(field, idx) in currentSession.fields" :key="idx">
                                <label class="block text-xs font-bold text-gray-600 uppercase mb-1.5">
                                    {{ field.name }}
                                </label>
                                
                                <select v-if="field.options" 
                                    v-model="inputData.data[field.name]"
                                    class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-primary/50 bg-white">
                                    <option value="">- Pilih -</option>
                                    <option v-for="opt in field.options" :key="opt" :value="opt">{{ opt }}</option>
                                </select>

                                <input v-else 
                                    v-model="inputData.data[field.name]"
                                    :type="field.type === 'number' ? 'number' : 'text'"
                                    class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-primary/50"
                                    :placeholder="field.name" />
                            </div>
                            
                            <div v-if="currentSession.fields.length === 0" class="col-span-2 text-xs text-gray-400 italic">
                                Tidak ada kolom penilaian khusus di template ini.
                            </div>
                        </div>
                    </div>

                    <div>
                        <label class="block text-xs font-bold text-gray-500 uppercase mb-1.5">Catatan (Opsional)</label>
                        <textarea v-model="inputData.note" rows="2" class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-primary/50" placeholder="Tulis catatan..."></textarea>
                    </div>

                    <div class="pt-4 border-t border-gray-100 flex justify-end">
                        <button @click="handleSave" :disabled="saving" class="bg-primary hover:bg-blue-700 text-white px-8 py-3 rounded-xl shadow-lg font-bold transition flex items-center gap-2 disabled:opacity-70">
                            <Loader2 v-if="saving" class="w-5 h-5 animate-spin"/>
                            <Save v-else class="w-5 h-5"/>
                            {{ saving ? 'Menyimpan...' : 'Simpan Laporan' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div>
            <h3 class="font-bold text-gray-800 mb-4 flex items-center gap-2">
                <BookOpen class="w-5 h-5 text-gray-500"/> Daftar Sesi (Pilih untuk Mengisi)
            </h3>

            <div v-if="todoSessions.length === 0" class="text-center py-8 bg-gray-50 rounded-xl border border-dashed border-gray-300 text-gray-400">
                <CheckCircle class="w-8 h-8 mx-auto mb-2 text-green-400"/>
                <p>Tidak ada sesi aktif saat ini.</p>
            </div>

            <div v-else class="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
                <div v-for="session in todoSessions" :key="session.id" 
                     @click="selectSession(session)"
                     class="cursor-pointer border rounded-xl p-4 transition-all duration-200 group relative overflow-hidden"
                     :class="selectedSessionId === session.id 
                        ? 'bg-blue-50 border-blue-500 ring-1 ring-blue-500 shadow-sm' 
                        : 'bg-white border-gray-200 hover:border-blue-300 hover:shadow-md'">
                    
                    <div class="flex justify-between items-start">
                        <div>
                            <div class="text-xs text-gray-500 mb-1 flex items-center gap-1">
                                <Calendar class="w-3 h-3"/> {{ new Date(session.date).toLocaleDateString('id-ID', { dateStyle: 'long' }) }}
                            </div>
                            <h4 class="font-bold text-gray-800 group-hover:text-primary transition">{{ session.name }}</h4>
                        </div>
                        <div v-if="session.has_entry" class="text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded font-bold">Sudah Diisi</div>
                        <div v-else class="text-[10px] bg-orange-100 text-orange-700 px-2 py-0.5 rounded font-bold">Belum Diisi</div>
                    </div>
                    
                    <div class="mt-3 flex justify-between items-end">
                        <span class="text-xs text-gray-400 truncate max-w-[150px]">{{ session.template_name }}</span>
                        <ChevronRight class="w-4 h-4 text-gray-300 group-hover:text-primary group-hover:translate-x-1 transition-all"/>
                    </div>
                </div>
            </div>
        </div>

    </div>
  </DashboardLayout>
</template>