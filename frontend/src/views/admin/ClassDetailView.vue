<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { debounce } from 'lodash'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import api from '../../services/api'
import { Users, UserPlus, Loader2, Trash2, ArrowLeft, Search, Edit3 } from 'lucide-vue-next'

const route = useRoute()
const classId = route.params.id

const classInfo = ref({})
const students = ref([])
const teachers = ref([])
const loading = ref(false)

// Edit Class Info
const showEditClass = ref(false)
const submittingClass = ref(false)

// Add Student Logic
const showAddModal = ref(false)
const studentSearchQuery = ref('')
const searchResults = ref([])
const searching = ref(false)
const submittingStudent = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    // 1. Get Class Detail (Backend support required: GET /classes/:id)
    const resClass = await api.get(`/classes/${classId}`)
    classInfo.value = resClass.data.data

    // 2. Get Students
    const resStudents = await api.get(`/classes/${classId}/students`)
    students.value = resStudents.data.data || []
    
    // 3. Pre-load teachers just in case user wants to edit
    const resTeachers = await api.get('/users', { params: { role: 'GURU', status: 'ACTIVE' } })
    teachers.value = resTeachers.data.data.users || []

  } catch (err) { console.error(err) }
  finally { loading.value = false }
}

const updateClass = async () => {
  submittingClass.value = true
  try {
    // Kita update nama, level, dan teacher_id
    await api.put(`/classes/${classId}`, {
      name: classInfo.value.name,
      level: classInfo.value.level,
      teacher_id: classInfo.value.teacher_id,
      academic_year_id: classInfo.value.academic_year_id // Tetap kirim biar ga error validasi
    })
    alert("Data kelas diupdate")
    showEditClass.value = false
    fetchData() // Refresh nama guru
  } catch (err) { alert("Gagal update kelas") }
  finally { submittingClass.value = false }
}

// Pencarian Murid Live (Debounced)
const searchStudents = debounce(async (query) => {
  if (!query) { searchResults.value = []; return }
  searching.value = true
  try {
    const res = await api.get('/users', { params: { role: 'MURID', status: 'ACTIVE', q: query } })
    searchResults.value = res.data.data.users || []
  } catch (err) { console.error(err) }
  finally { searching.value = false }
}, 300)

watch(studentSearchQuery, (val) => searchStudents(val))

const addStudent = async (studentId) => {
  submittingStudent.value = true
  try {
    await api.post(`/classes/${classId}/students`, { student_id: studentId })
    showAddModal.value = false
    studentSearchQuery.value = ''
    searchResults.value = []
    fetchData()
  } catch (err) { alert("Gagal/Sudah terdaftar") }
  finally { submittingStudent.value = false }
}

const removeStudent = async (studentId) => {
  if(!confirm("Keluarkan murid ini dari kelas?")) return
  try {
    await api.delete(`/classes/${classId}/students/${studentId}`)
    fetchData()
  } catch(err) { alert("Gagal menghapus") }
}

onMounted(fetchData)
</script>

<template>
  <DashboardLayout>
    <div class="max-w-4xl mx-auto">
      <button @click="$router.back()" class="flex items-center text-sm text-secondary hover:text-primary mb-4">
        <ArrowLeft class="w-4 h-4 mr-1" /> Kembali
      </button>

      <div v-if="classInfo.id" class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm mb-6 flex justify-between items-start">
        <div>
           <div class="flex items-center gap-2 mb-1">
             <h1 class="text-2xl font-bold text-accent">{{ classInfo.name }}</h1>
             <span class="bg-gray-100 text-xs px-2 py-1 rounded font-mono">Level {{ classInfo.level }}</span>
           </div>
           <p class="text-secondary text-sm flex items-center gap-1">
             <Users class="w-4 h-4" /> Wali Kelas: 
             <span class="font-semibold text-primary">{{ classInfo.teacher_name || 'Belum ditentukan' }}</span>
           </p>
        </div>
        <button @click="showEditClass = !showEditClass" class="text-sm text-blue-600 hover:underline flex items-center gap-1">
          <Edit3 class="w-4 h-4" /> Edit Info
        </button>
      </div>

      <div v-if="showEditClass" class="bg-gray-50 p-4 rounded-xl border border-blue-100 mb-6 animate-fade-in">
        <form @submit.prevent="updateClass" class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div><label class="label">Nama Kelas</label><input v-model="classInfo.name" class="input-field w-full" /></div>
          <div><label class="label">Wali Kelas</label>
            <select v-model="classInfo.teacher_id" class="input-field w-full bg-white">
              <option v-for="t in teachers" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </div>
          <div class="flex items-end"><button type="submit" class="bg-primary text-white w-full py-2 rounded-lg">Update</button></div>
        </form>
      </div>

      <div class="flex justify-between items-center mb-4">
        <h3 class="font-bold text-lg text-primary">Daftar Murid ({{ students.length }})</h3>
        <button @click="showAddModal = true" class="bg-primary text-white px-4 py-2 rounded-lg text-sm flex items-center gap-2 hover:bg-accent transition">
          <UserPlus class="w-4 h-4" /> Tambah Murid
        </button>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <table class="w-full text-left text-sm">
          <thead class="bg-gray-50 border-b border-gray-200 text-secondary">
            <tr>
              <th class="px-6 py-3">Nama Lengkap</th>
              <th class="px-6 py-3">No. Induk</th>
              <th class="px-6 py-3">Status</th>
              <th class="px-6 py-3 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="s in students" :key="s.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-primary">{{ s.name }}</td>
              <td class="px-6 py-4 text-secondary">{{ s.serial_number || '-' }}</td>
              <td class="px-6 py-4"><span class="px-2 py-1 rounded text-xs font-bold bg-green-50 text-green-700">ACTIVE</span></td>
              <td class="px-6 py-4 text-right">
                <button @click="removeStudent(s.id)" class="text-red-500 hover:bg-red-50 p-1 rounded transition"><Trash2 class="w-4 h-4" /></button>
              </td>
            </tr>
            <tr v-if="students.length === 0"><td colspan="4" class="p-8 text-center text-secondary">Kelas kosong.</td></tr>
          </tbody>
        </table>
      </div>

      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white w-full max-w-lg rounded-xl p-6 shadow-xl">
          <div class="flex justify-between mb-4">
            <h3 class="font-bold text-lg">Cari & Tambah Murid</h3>
            <button @click="showAddModal = false" class="text-gray-400 hover:text-gray-600">✕</button>
          </div>
          
          <div class="relative mb-4">
            <Search class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
            <input 
              v-model="studentSearchQuery" 
              type="text" 
              placeholder="Ketik nama murid..." 
              class="w-full pl-9 pr-4 py-2 border border-gray-300 rounded-lg outline-none focus:border-primary"
              autofocus
            />
          </div>

          <div class="max-h-60 overflow-y-auto border border-gray-100 rounded-lg">
            <div v-if="searching" class="p-4 text-center text-secondary">Mencari...</div>
            <div v-else-if="searchResults.length === 0 && studentSearchQuery" class="p-4 text-center text-secondary">Tidak ditemukan.</div>
            
            <div 
              v-for="res in searchResults" :key="res.id" 
              @click="addStudent(res.id)"
              class="p-3 hover:bg-blue-50 cursor-pointer border-b last:border-0 flex justify-between items-center"
            >
              <div>
                <p class="font-medium text-primary">{{ res.name }}</p>
                <p class="text-xs text-secondary">{{ res.username }}</p>
              </div>
              <Plus class="w-4 h-4 text-primary" />
            </div>
          </div>
        </div>
      </div>

    </div>
  </DashboardLayout>
</template>