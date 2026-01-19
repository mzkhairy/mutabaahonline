<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { 
  ArrowLeft, Save, Calendar, Users, 
  BookOpen, Loader2, CheckCircle, Edit2, Lock, X, RefreshCw, AlertTriangle
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const sessionId = route.params.id

const loading = ref(false)
const saving = ref(false)
const refreshingStudents = ref(false)

const session = ref(null)
const students = ref([])
const templateFields = ref([])

const classData = ref({}) 
const studentEntries = ref({}) 

const showEditModal = ref(false)
const editForm = ref({ name: '', date: '' })

// --- HELPER DATE FORMAT ---
const formattedDate = computed(() => {
    if (!session.value?.date) return '-'
    const date = new Date(session.value.date)
    return date.toLocaleDateString('id-ID', { 
        weekday: 'long', 
        day: 'numeric', 
        month: 'long', 
        year: 'numeric' 
    })
})

// --- FETCH DATA ---
const fetchData = async () => {
  loading.value = true
  try {
    const resSession = await api.get(`/sessions/${sessionId}`)
    session.value = resSession.data.data
    
    const rawDate = new Date(session.value.date).toISOString().split('T')[0]
    editForm.value = { name: session.value.name, date: rawDate }

    let rawStructure = session.value.template_structure
    if (!Array.isArray(rawStructure)) {
        try { rawStructure = JSON.parse(rawStructure) } catch (e) { rawStructure = [] }
    }
    templateFields.value = rawStructure || []
    
    classData.value = session.value.class_activity_data || {}

    const classId = session.value.class_id
    const resClassStudents = await api.get(`/classes/${classId}/students`)
    const allStudents = resClassStudents.data.data || []
    
    students.value = allStudents.map(s => ({
        id: s.id,
        name: s.name,
        serial_number: s.serial_number
    }))

    await fetchStudentEntries()

  } catch (err) {
    console.error(err)
    alert("Gagal memuat data sesi.")
  } finally {
    loading.value = false
  }
}

const fetchStudentEntries = async () => {
    try {
        const resEntries = await api.get('/mutabaah', { params: { session_id: sessionId } })
        const existingEntries = resEntries.data.data || []

        students.value.forEach(student => {
            const entry = existingEntries.find(e => e.student_id === student.id)
            studentEntries.value[student.id] = {
                attendance: entry?.attendance || 'HADIR',
                note: entry?.note || '',
                data: entry?.student_activity_data || {} 
            }
        })
    } catch (e) { console.error(e) }
}

const handleRefreshStudents = async () => {
    refreshingStudents.value = true
    await fetchStudentEntries()
    refreshingStudents.value = false
}

const classFields = computed(() => templateFields.value.filter(f => f.type === 'class_activity'))
const personalFields = computed(() => templateFields.value.filter(f => f.type === 'student_activity'))

const canSave = computed(() => {
    if (session.value?.is_student_input_allowed) return false
    return true
})

const handleUpdateInfo = async () => {
    try {
        await api.put(`/sessions/${sessionId}`, {
            name: editForm.value.name,
            date: editForm.value.date
        })
        session.value.name = editForm.value.name
        session.value.date = editForm.value.date
        showEditModal.value = false
    } catch(err) { alert("Gagal update info.") }
}

const handleLockSession = async () => {
    if(!confirm("Kunci sesi ini agar input siswa dimatikan?")) return
    try {
        await api.put(`/sessions/${sessionId}`, { is_student_input_allowed: false })
        session.value.is_student_input_allowed = false
    } catch(err) { alert("Gagal mengunci.") }
}

const markAllPresent = () => {
    if(!confirm("Set semua siswa menjadi HADIR?")) return
    Object.keys(studentEntries.value).forEach(id => {
        studentEntries.value[id].attendance = 'HADIR'
    })
}

const handleSaveAll = async () => {
  saving.value = true
  try {
    await api.put(`/sessions/${sessionId}`, { class_activity_data: classData.value })

    const promises = Object.keys(studentEntries.value).map(studentId => {
        const entry = studentEntries.value[studentId]
        return api.post('/mutabaah', {
            session_id: sessionId,
            student_id: studentId,
            attendance: entry.attendance,
            status: 'LANJUT',
            note: entry.note,
            student_activity_data: entry.data
        })
    })

    await Promise.all(promises)
    alert("Berhasil menyimpan data!")
    await fetchStudentEntries()
  } catch (err) {
    console.error(err)
    alert("Gagal menyimpan.")
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <DashboardLayout>
    <div v-if="loading" class="flex h-64 items-center justify-center">
      <Loader2 class="w-8 h-8 animate-spin text-primary" />
    </div>

    <div v-else class="max-w-6xl mx-auto space-y-4 sm:space-y-6 pb-32"> 
      
      <div class="sticky top-0 bg-white/95 backdrop-blur-sm z-30 py-3 px-4 border-b border-gray-200 shadow-sm -mx-4 sm:mx-0 sm:rounded-xl sm:border sm:px-6">
        <div class="flex flex-col sm:flex-row justify-between sm:items-center gap-3">
            <div class="flex items-start gap-3">
                <button @click="router.back()" class="p-1.5 hover:bg-gray-100 rounded-full transition shrink-0 mt-1 sm:mt-0">
                    <ArrowLeft class="w-5 h-5 text-gray-600" />
                </button>
                <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2">
                        <h1 class="text-lg sm:text-2xl font-bold text-gray-800 truncate leading-tight">{{ session?.name }}</h1>
                        <button @click="showEditModal = true" class="text-gray-400 hover:text-primary shrink-0"><Edit2 class="w-4 h-4"/></button>
                    </div>
                    <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs sm:text-sm text-gray-500 mt-1">
                        <span class="flex items-center gap-1"><Calendar class="w-3 h-3"/> {{ formattedDate }}</span>
                        <span class="flex items-center gap-1"><BookOpen class="w-3 h-3"/> {{ session?.template_name }}</span>
                    </div>
                </div>
            </div>

            <div class="hidden md:flex items-center gap-2">
                <span v-if="session?.is_student_input_allowed" class="text-xs text-orange-600 bg-orange-100 px-2 py-1 rounded font-bold flex items-center gap-1 border border-orange-200">
                    <AlertTriangle class="w-3 h-3"/> Input Siswa Aktif
                </span>
                
                <button v-if="session?.is_student_input_allowed" @click="handleLockSession" class="text-sm bg-white border border-orange-300 text-orange-600 hover:bg-orange-50 flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-bold transition">
                    <Lock class="w-3 h-3" /> Kunci
                </button>
                
                <button @click="handleSaveAll" :disabled="saving || !canSave" class="text-sm bg-primary text-white flex items-center gap-1.5 px-4 py-2 rounded-lg shadow-sm font-bold transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-700 active:scale-95">
                    <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
                    <Save v-else class="w-4 h-4" />
                    <span>{{ saving ? 'Menyimpan...' : 'Simpan' }}</span>
                </button>
            </div>
        </div>
      </div>

      <div v-if="!canSave" class="mx-4 sm:mx-0 bg-orange-50 border border-orange-200 text-orange-800 px-4 py-3 rounded-lg flex items-start gap-3 text-sm shadow-sm">
        <Lock class="w-5 h-5 mt-0.5 shrink-0 text-orange-500" />
        <div>
            <p class="font-bold">Mode Read-Only (Hanya Baca)</p>
            <p class="text-xs mt-1 text-orange-700">Sesi ini sedang terbuka untuk siswa. Anda harus mengunci sesi terlebih dahulu untuk melakukan penilaian final dan menyimpan data.</p>
        </div>
      </div>

      <div class="bg-white mx-4 sm:mx-0 rounded-xl shadow-sm border border-gray-200 p-4 sm:p-6">
        <div class="flex items-center gap-2 mb-4 border-b border-gray-100 pb-2">
            <div class="bg-blue-100 p-1.5 rounded-md text-primary"><BookOpen class="w-4 h-4"/></div>
            <h2 class="font-bold text-gray-800 text-sm sm:text-base">Aktivitas Kelas</h2>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div v-for="(field, idx) in classFields" :key="idx">
                <label class="block text-xs font-bold text-gray-500 mb-1.5 uppercase tracking-wide">{{ field.name }}</label>
                <input v-model="classData[field.name]" type="text" class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition text-sm bg-gray-50 focus:bg-white" placeholder="Isi keterangan..." />
            </div>
            <div v-if="classFields.length === 0" class="text-xs text-gray-400 italic">Tidak ada parameter aktivitas kelas.</div>
        </div>
      </div>

      <div class="bg-white mx-4 sm:mx-0 rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        
        <div class="p-4 sm:p-6 border-b border-gray-200 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 bg-gray-50/50">
             <div class="flex items-center gap-2">
                <div class="bg-green-100 p-1.5 rounded-md text-green-700"><Users class="w-4 h-4"/></div>
                <h2 class="font-bold text-gray-800 text-sm sm:text-base">Input Nilai Siswa</h2>
                <span class="text-xs text-gray-500 bg-gray-200 px-2 py-0.5 rounded-full font-bold">{{ students.length }}</span>
            </div>
            
            <div class="flex gap-2 w-full sm:w-auto">
                <button @click="handleRefreshStudents" :disabled="refreshingStudents" class="flex-1 sm:flex-none justify-center text-xs flex items-center gap-1 bg-white text-gray-600 px-3 py-2 rounded-lg border border-gray-300 hover:bg-gray-100 transition font-medium active:scale-95">
                    <Loader2 v-if="refreshingStudents" class="w-3 h-3 animate-spin"/>
                    <RefreshCw v-else class="w-3 h-3"/>
                    Refresh
                </button>

                <button @click="markAllPresent" class="flex-1 sm:flex-none justify-center text-xs flex items-center gap-1 bg-green-50 text-green-700 px-3 py-2 rounded-lg border border-green-200 hover:bg-green-100 transition font-medium active:scale-95">
                    <CheckCircle class="w-3 h-3"/> Hadir Semua
                </button>
            </div>
        </div>

        <div class="hidden md:block overflow-x-auto">
          <table class="w-full text-sm text-left">
            <thead class="bg-gray-50 text-gray-500 border-b font-semibold uppercase text-xs">
              <tr>
                <th class="px-4 py-3 w-10 text-center">#</th>
                <th class="px-4 py-3 min-w-[200px]">Nama Siswa</th>
                <th class="px-4 py-3 w-[120px]">Kehadiran</th>
                <th v-for="(field, idx) in personalFields" :key="idx" class="px-4 py-3 min-w-[120px]">{{ field.name }}</th>
                <th class="px-4 py-3 min-w-[200px]">Catatan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="(student, index) in students" :key="student.id" class="hover:bg-blue-50/30 transition group">
                <td class="px-4 py-3 text-center text-gray-400 font-mono text-xs">{{ index + 1 }}</td>
                <td class="px-4 py-3 font-medium text-gray-800">
                    {{ student.name }}
                    <div class="text-[10px] text-gray-400 font-mono">{{ student.serial_number || '' }}</div>
                </td>
                
                <td class="px-4 py-3">
                    <select v-model="studentEntries[student.id].attendance" 
                        class="w-full text-xs border-gray-300 rounded focus:ring-primary focus:border-primary p-1.5 font-bold cursor-pointer"
                        :class="{
                            'bg-green-100 text-green-700 border-green-200': studentEntries[student.id].attendance === 'HADIR', 
                            'bg-red-100 text-red-700 border-red-200': studentEntries[student.id].attendance === 'ALPA', 
                            'bg-yellow-100 text-yellow-700 border-yellow-200': ['SAKIT','IZIN'].includes(studentEntries[student.id].attendance)
                        }">
                        <option value="HADIR">Hadir</option>
                        <option value="SAKIT">Sakit</option>
                        <option value="IZIN">Izin</option>
                        <option value="ALPA">Alpa</option>
                    </select>
                </td>

                <td v-for="(field, idx) in personalFields" :key="idx" class="px-4 py-3">
                    <div :class="{'opacity-40 pointer-events-none grayscale': studentEntries[student.id].attendance !== 'HADIR'}">
                        <input v-model="studentEntries[student.id].data[field.name]" type="text" class="w-full text-xs border-gray-300 rounded p-1.5 focus:ring-1 focus:ring-primary outline-none transition text-center font-mono" placeholder="-" />
                    </div>
                </td>

                <td class="px-4 py-3">
                    <input type="text" v-model="studentEntries[student.id].note" class="w-full text-xs border-transparent hover:border-gray-200 focus:border-gray-300 rounded p-1.5 focus:ring-0 outline-none transition bg-transparent focus:bg-white truncate" placeholder="Tulis..." />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="md:hidden bg-gray-50 p-3 space-y-3">
            <div v-for="(student, index) in students" :key="student.id" class="bg-white rounded-xl border border-gray-200 shadow-sm p-4">
                
                <div class="flex justify-between items-start mb-4 border-b border-gray-100 pb-3">
                    <div>
                        <h3 class="font-bold text-gray-800 text-sm">{{ student.name }}</h3>
                        <p class="text-[10px] text-gray-400 font-mono">{{ student.serial_number || 'No ID' }}</p>
                    </div>
                    <select v-model="studentEntries[student.id].attendance" 
                        class="text-xs border-none rounded-lg p-1.5 font-bold cursor-pointer ring-1 ring-inset focus:ring-2 focus:ring-primary outline-none"
                        :class="{
                            'bg-green-50 text-green-700 ring-green-200': studentEntries[student.id].attendance === 'HADIR', 
                            'bg-red-50 text-red-700 ring-red-200': studentEntries[student.id].attendance === 'ALPA', 
                            'bg-yellow-50 text-yellow-700 ring-yellow-200': ['SAKIT','IZIN'].includes(studentEntries[student.id].attendance)
                        }">
                        <option value="HADIR">Hadir</option>
                        <option value="SAKIT">Sakit</option>
                        <option value="IZIN">Izin</option>
                        <option value="ALPA">Alpa</option>
                    </select>
                </div>

                <div class="space-y-4 transition-opacity duration-200" :class="{'opacity-50 pointer-events-none': studentEntries[student.id].attendance !== 'HADIR'}">
                    
                    <div v-for="(field, idx) in personalFields" :key="idx">
                        <label class="block text-[10px] font-bold text-gray-400 uppercase mb-1 truncate">{{ field.name }}</label>
                        <input 
                            v-model="studentEntries[student.id].data[field.name]" 
                            type="text" 
                            class="w-full text-sm border border-gray-300 rounded-lg p-2.5 focus:border-primary focus:ring-1 focus:ring-primary outline-none font-medium text-gray-700 transition placeholder-gray-300" 
                            placeholder="..." 
                        />
                    </div>
                    
                    <div>
                        <label class="block text-[10px] font-bold text-gray-400 uppercase mb-1">Catatan</label>
                        <input 
                            type="text" 
                            v-model="studentEntries[student.id].note" 
                            class="w-full text-sm border border-gray-300 rounded-lg p-2.5 focus:border-primary focus:ring-1 focus:ring-primary outline-none font-medium text-gray-700 transition placeholder-gray-300" 
                            placeholder="Tulis catatan..." 
                        />
                    </div>
                </div>

            </div>
        </div>

      </div>

      <div class="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 p-4 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)] z-40 flex items-center gap-3">
            <button v-if="session?.is_student_input_allowed" @click="handleLockSession" class="bg-orange-50 text-orange-600 border border-orange-200 px-4 py-3 rounded-xl font-bold flex items-center justify-center gap-2 transition active:scale-95 shadow-sm">
                <Lock class="w-5 h-5" />
            </button>
            
            <button @click="handleSaveAll" :disabled="saving || !canSave" class="flex-1 bg-primary text-white px-6 py-3 rounded-xl shadow-lg font-bold transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-blue-700 active:scale-95 flex items-center justify-center gap-2">
                <Loader2 v-if="saving" class="w-5 h-5 animate-spin" />
                <Save v-else class="w-5 h-5" />
                <span>{{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
            </button>
      </div>

      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
        <div class="bg-white rounded-xl w-full max-w-sm shadow-xl overflow-hidden animate-scale-in">
            <div class="bg-gray-50 px-4 py-3 border-b border-gray-100 flex justify-between items-center">
                <h3 class="font-bold text-gray-800">Edit Sesi</h3>
                <button @click="showEditModal = false" class="text-gray-400 hover:text-gray-600"><X class="w-5 h-5"/></button>
            </div>
            <form @submit.prevent="handleUpdateInfo" class="p-4 space-y-4">
                <div>
                    <label class="block text-xs font-bold text-gray-500 mb-1 uppercase">Nama Sesi</label>
                    <input v-model="editForm.name" type="text" class="w-full border border-gray-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-primary/20 outline-none" required />
                </div>
                <div>
                    <label class="block text-xs font-bold text-gray-500 mb-1 uppercase">Tanggal</label>
                    <input v-model="editForm.date" type="date" class="w-full border border-gray-300 rounded-lg p-2.5 text-sm focus:ring-2 focus:ring-primary/20 outline-none" required />
                </div>
                <div class="pt-2 flex gap-2">
                    <button type="button" @click="showEditModal = false" class="flex-1 py-2 text-sm text-gray-600 font-medium hover:bg-gray-100 rounded-lg">Batal</button>
                    <button type="submit" class="flex-1 bg-primary text-white py-2 rounded-lg font-bold text-sm shadow-sm hover:bg-blue-700">Simpan</button>
                </div>
            </form>
        </div>
      </div>

    </div>
  </DashboardLayout>
</template>