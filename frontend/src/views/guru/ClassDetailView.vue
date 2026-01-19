<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import api from '../../services/api'
import { 
  ArrowLeft, Plus, Users, ChevronRight, Loader2, FileText, X, CheckCircle, FileBarChart 
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const classId = route.params.id

const loading = ref(false)
const classInfo = ref({})
const sessions = ref([])
const students = ref([]) // [BARU] Data murid
const templates = ref([])

// Tab Nav
const activeTab = ref('SESSIONS') // SESSIONS | STUDENTS

// Modal Create Session
const showModal = ref(false)
const submitting = ref(false)
const form = ref({
  name: '',
  date: new Date().toISOString().substring(0, 10),
  template_id: '',
  is_student_input_allowed: true
})

const fetchData = async () => {
  loading.value = true
  try {
    const resClass = await api.get(`/classes/${classId}`)
    classInfo.value = resClass.data.data

    const resSessions = await api.get(`/sessions`, { params: { class_id: classId } })
    sessions.value = resSessions.data.data || []

    const resTpl = await api.get('/templates')
    templates.value = resTpl.data.data || []
    
    // [BARU] Ambil Data Murid
    const resStudents = await api.get(`/classes/${classId}/students`)
    students.value = resStudents.data.data || []

    // Auto-select first template
    if (templates.value.length > 0) form.value.template_id = templates.value[0].id
  } catch (err) { console.error(err) } 
  finally { loading.value = false }
}

const handleCreateSession = async () => {
  submitting.value = true
  try {
    const payload = {
      class_id: classId,
      template_id: form.value.template_id,
      name: form.value.name,
      date: form.value.date,
      is_student_input_allowed: form.value.is_student_input_allowed,
      class_activity_data: {} 
    }
    await api.post('/sessions', payload)
    alert('Sesi berhasil dibuat!')
    showModal.value = false
    form.value.name = ''
    fetchData()
  } catch (err) { alert('Gagal membuat sesi') } 
  finally { submitting.value = false }
}

onMounted(fetchData)
</script>

<template>
  <DashboardLayout>
    <div class="max-w-5xl mx-auto">
      <button @click="router.push('/guru/classes')" class="flex items-center text-sm text-secondary hover:text-primary mb-6 transition">
        <ArrowLeft class="w-4 h-4 mr-1" /> Kembali
      </button>

      <div v-if="classInfo.id" class="bg-gradient-to-r from-primary to-blue-600 rounded-2xl p-6 text-white shadow-lg mb-8">
        <div class="flex items-center gap-2 mb-2 opacity-90">
            <Users class="w-4 h-4" />
            <span class="text-sm font-medium">Level {{ classInfo.level }}</span>
        </div>
        <h1 class="text-3xl font-bold mb-1">{{ classInfo.name }}</h1>
        <p class="opacity-80 text-sm">{{ classInfo.student_count || 0 }} Murid Terdaftar</p>
      </div>

      <div class="flex border-b border-gray-200 mb-6">
          <button @click="activeTab = 'SESSIONS'" :class="['px-6 py-3 text-sm font-bold border-b-2 transition', activeTab === 'SESSIONS' ? 'border-primary text-primary' : 'border-transparent text-secondary hover:text-primary']">
              Daftar Pertemuan
          </button>
          <button @click="activeTab = 'STUDENTS'" :class="['px-6 py-3 text-sm font-bold border-b-2 transition', activeTab === 'STUDENTS' ? 'border-primary text-primary' : 'border-transparent text-secondary hover:text-primary']">
              Daftar Murid & Rekap
          </button>
      </div>

      <div v-if="loading" class="py-8 flex justify-center"><Loader2 class="animate-spin w-8 h-8 text-primary" /></div>

      <div v-else-if="activeTab === 'SESSIONS'">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-lg font-bold text-gray-700">Riwayat Sesi</h2>
            <button @click="showModal = true" class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-accent transition flex items-center gap-2 shadow-sm text-sm font-medium">
            <Plus class="w-4 h-4" /> Sesi Baru
            </button>
        </div>

        <div class="space-y-3">
            <div v-for="ses in sessions" :key="ses.id" @click="router.push(`/guru/sessions/${ses.id}`)"
            class="bg-white p-4 rounded-xl border border-gray-200 hover:border-primary/50 hover:shadow-md transition cursor-pointer flex items-center justify-between group">
            <div class="flex items-center gap-4">
                <div class="bg-blue-50 text-primary p-3 rounded-lg font-bold text-center min-w-[60px]">
                    <div class="text-xs uppercase opacity-70">{{ new Date(ses.date).toLocaleString('id-ID', { month: 'short' }) }}</div>
                    <div class="text-lg leading-none">{{ new Date(ses.date).getDate() }}</div>
                </div>
                <div>
                    <h4 class="font-bold text-accent group-hover:text-primary transition">{{ ses.name }}</h4>
                    <p class="text-xs text-secondary mt-1 flex items-center gap-1">
                    <FileText class="w-3 h-3" /> Template: {{ ses.template_name || 'Unknown' }}
                    </p>
                </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-300 group-hover:text-primary transition" />
            </div>
            <div v-if="sessions.length === 0" class="text-center py-10 bg-gray-50 text-secondary border border-dashed rounded-xl">Belum ada sesi.</div>
        </div>
      </div>

      <div v-else-if="activeTab === 'STUDENTS'">
          <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
              <table class="w-full text-left text-sm">
                  <thead class="bg-gray-50 border-b border-gray-200 text-gray-600">
                      <tr>
                          <th class="px-6 py-4">Nama Murid</th>
                          <th class="px-6 py-4">Nomor Induk</th>
                          <th class="px-6 py-4 text-right">Aksi</th>
                      </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-100">
                      <tr v-for="s in students" :key="s.id" class="hover:bg-gray-50 transition">
                          <td class="px-6 py-4 font-bold text-primary">{{ s.name }}</td>
                          <td class="px-6 py-4 text-secondary font-mono">{{ s.serial_number || '-' }}</td>
                          <td class="px-6 py-4 text-right">
                              <button @click="router.push(`/guru/students/${s.id}/report`)" class="inline-flex items-center gap-1 px-3 py-1.5 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded-lg text-xs font-bold transition">
                                  <FileBarChart class="w-3 h-3"/> Lihat Rekap
                              </button>
                          </td>
                      </tr>
                      <tr v-if="students.length === 0">
                          <td colspan="3" class="p-8 text-center text-gray-400">Belum ada murid di kelas ini.</td>
                      </tr>
                  </tbody>
              </table>
          </div>
      </div>

      <div v-if="showModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white p-6 rounded-xl w-full max-w-2xl shadow-xl relative max-h-[90vh] overflow-y-auto">
          <button @click="showModal = false" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600"><X class="w-5 h-5"/></button>
          
          <h3 class="font-bold text-lg mb-6 text-primary">Buat Sesi Pertemuan Baru</h3>
          
          <form @submit.prevent="handleCreateSession" class="space-y-6">
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="label">Judul Sesi</label>
                    <input v-model="form.name" type="text" placeholder="Contoh: Halaqoh 1" class="input-field w-full border p-2 rounded-lg" required />
                </div>
                <div>
                    <label class="label">Tanggal</label>
                    <input v-model="form.date" type="date" class="input-field w-full border p-2 rounded-lg" required />
                </div>
            </div>

            <div>
              <label class="label mb-2 block">Pilih Template Penilaian</label>
              <div class="flex gap-3 overflow-x-auto pb-4 snap-x">
                <div 
                    v-for="t in templates" :key="t.id"
                    @click="form.template_id = t.id"
                    :class="['min-w-[200px] w-[200px] p-4 rounded-xl border-2 cursor-pointer transition snap-center flex-shrink-0', 
                             form.template_id === t.id ? 'border-primary bg-blue-50 ring-1 ring-primary' : 'border-gray-200 hover:border-blue-300 bg-white']"
                >
                    <div class="flex justify-between items-start mb-2">
                        <span class="font-bold text-sm text-primary line-clamp-1">{{ t.name }}</span>
                        <CheckCircle v-if="form.template_id === t.id" class="w-4 h-4 text-primary"/>
                    </div>
                    <div class="space-y-1">
                        <div v-for="(item, idx) in t.structure.slice(0, 3)" :key="idx" class="text-[10px] text-secondary flex items-center gap-1">
                            <span class="w-1 h-1 rounded-full bg-gray-400"></span> {{ item.name }}
                        </div>
                        <div v-if="t.structure.length > 3" class="text-[10px] text-gray-400 italic">+ {{ t.structure.length - 3 }} lainnya</div>
                    </div>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-2">
                <input type="checkbox" v-model="form.is_student_input_allowed" id="allowInput" class="w-4 h-4 text-primary rounded border-gray-300 focus:ring-primary">
                <label for="allowInput" class="text-sm text-secondary cursor-pointer">Izinkan murid mengisi capaian mandiri</label>
            </div>

            <button type="submit" :disabled="submitting" class="w-full bg-primary text-white py-2.5 rounded-lg font-medium flex justify-center gap-2">
                <Loader2 v-if="submitting" class="animate-spin w-4 h-4" /> Mulai Sesi
            </button>
          </form>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>