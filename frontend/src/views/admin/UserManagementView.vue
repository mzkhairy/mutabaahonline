<script setup>
import { ref, onMounted, watch } from 'vue'
import { debounce } from 'lodash' 
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import api from '../../services/api'
import { 
  Users, Search, UserCheck, UserX, KeyRound, 
  Loader2, Plus, X, Edit2, Upload, FileDown, CheckCircle, AlertTriangle 
} from 'lucide-vue-next'

const loading = ref(false)
const users = ref([])
const stats = ref({}) 
const activeTab = ref('GURU') 
const searchQuery = ref('')
const processingId = ref(null)

// Tambahkan State untuk Manual Add
const showAddModal = ref(false)
const adding = ref(false)
const addForm = ref({ name: '', username: '', role: 'MURID', serial_number: '', password: '' })

// --- STATE MODAL EDIT USER (Manual) ---
const showEditModal = ref(false)
const editing = ref(false)
const editForm = ref({ id: '', name: '', username: '', serial_number: '' })

// --- STATE MODAL IMPORT CSV ---
const showImportModal = ref(false)
const importing = ref(false)
const importFile = ref(null)
const importResult = ref(null) // Menyimpan hasil { success, skipped, details }

const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.get('/users', {
      params: { role: activeTab.value, q: searchQuery.value }
    })
    users.value = res.data.data.users || []
    if (res.data.data.stats) stats.value = res.data.data.stats
  } catch (err) { console.error(err) } 
  finally { loading.value = false }
}

const handleSearch = debounce(() => fetchData(), 500)

// --- ACTIONS USER ---

const handleAddUser = async () => {
    if (!addForm.value.name || !addForm.value.username) return alert("Nama dan Username wajib diisi")
    adding.value = true
    try {
        await api.post('/users', addForm.value)
        alert("User berhasil ditambahkan")
        showAddModal.value = false
        // Reset form
        addForm.value = { name: '', username: '', role: 'MURID', serial_number: '', password: '' }
        fetchData()
    } catch(err) {
        alert("Gagal tambah user: " + (err.response?.data?.error?.message || 'Error'))
    } finally {
        adding.value = false
    }
}

const toggleStatus = async (user) => {
  const newStatus = user.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE'
  if (!confirm(`Ubah status ${user.name} jadi ${newStatus}?`)) return
  processingId.value = user.id
  try {
    await api.patch(`/users/${user.id}/status`, { status: newStatus })
    user.status = newStatus
    fetchData() 
  } catch (err) { alert("Gagal update status") } 
  finally { processingId.value = null }
}

const resetPassword = async (user) => {
  const newPass = prompt(`Password baru untuk ${user.name}:`, "12345678")
  if (!newPass) return
  try {
    await api.patch(`/users/${user.id}/password`, { new_password: newPass })
    alert("Password direset.")
  } catch (err) { alert("Gagal reset password") }
}

// --- EDIT USER MANUAL ---
const openEdit = (user) => {
    editForm.value = {
        id: user.id,
        name: user.name,
        username: user.username,
        serial_number: user.serial_number || ''
    }
    showEditModal.value = true
}

const handleUpdateUser = async () => {
    editing.value = true
    try {
        await api.put(`/users/${editForm.value.id}`, {
            name: editForm.value.name,
            username: editForm.value.username,
            serial_number: editForm.value.serial_number
        })
        alert('Data user diperbarui')
        showEditModal.value = false
        fetchData()
    } catch(err) {
        alert('Gagal update user: ' + (err.response?.data?.error?.message || 'Error'))
    } finally {
        editing.value = false
    }
}

// --- IMPORT CSV ---
const handleFileSelect = (e) => {
    const file = e.target.files[0]
    if (file && file.type === 'text/csv') {
        importFile.value = file
    } else {
        alert("Mohon upload file .csv")
        e.target.value = null
    }
}

const handleImport = async () => {
    if (!importFile.value) return
    importing.value = true
    importResult.value = null

    const formData = new FormData()
    formData.append('file', importFile.value)

    try {
        const res = await api.post('/users/import', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
        importResult.value = res.data.data
        // Refresh data setelah import
        fetchData()
    } catch (err) {
        alert("Gagal import: " + (err.response?.data?.error?.message || 'Error Server'))
    } finally {
        importing.value = false
        importFile.value = null // Reset input
    }
}

const downloadTemplate = () => {
    const csvContent = "data:text/csv;charset=utf-8,Nama Lengkap,Nomor Induk,Peran\nBudi Santoso,12345,MURID\nUstadz Ahmad,G001,GURU"
    const encodedUri = encodeURI(csvContent)
    const link = document.createElement("a")
    link.setAttribute("href", encodedUri)
    link.setAttribute("download", "template_import_user.csv")
    document.body.appendChild(link)
    link.click()
}

watch(activeTab, () => { searchQuery.value = ''; fetchData() })
onMounted(fetchData)
</script>

<template>
  <DashboardLayout>
    <div class="max-w-6xl mx-auto">
      
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end mb-8 gap-4">
        <div>
            <h1 class="text-2xl font-bold text-accent">Kelola Akun</h1>
            <p class="text-secondary text-sm">Manajemen data Guru dan Murid</p>
        </div>
        <div class="flex gap-2">
             <button @click="showImportModal = true; importResult = null" class="btn-secondary flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-gray-50 bg-white shadow-sm">
                <Upload class="w-4 h-4" /> Import CSV
             </button>
             <button @click="showAddModal = true" class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2 shadow-sm transition">
                <Plus class="w-4 h-4" /> Tambah User
             </button>
             </div>
      </div>

      <div v-if="showAddModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white p-6 rounded-xl w-full max-w-md shadow-xl relative">
            <button @click="showAddModal = false" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600"><X class="w-5 h-5"/></button>
            <h3 class="font-bold text-lg mb-4 text-primary">Tambah User Baru</h3>
            
            <form @submit.prevent="handleAddUser" class="space-y-4">
                <div>
                    <label class="label">Nama Lengkap</label>
                    <input v-model="addForm.name" type="text" class="input-field w-full border p-2 rounded-lg" required />
                </div>
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label class="label">Peran</label>
                        <select v-model="addForm.role" class="input-field w-full border p-2 rounded-lg bg-white">
                            <option value="MURID">Murid</option>
                            <option value="GURU">Guru</option>
                        </select>
                    </div>
                    <div>
                         <label class="label">Username</label>
                         <input v-model="addForm.username" type="text" class="input-field w-full border p-2 rounded-lg" required />
                    </div>
                </div>
                <div>
                    <label class="label">Nomor Induk (NIP/NIS)</label>
                    <input v-model="addForm.serial_number" type="text" class="input-field w-full border p-2 rounded-lg" placeholder="Opsional (Disarankan)" />
                </div>
                 <div>
                    <label class="label">Password Awal</label>
                    <input v-model="addForm.password" type="text" class="input-field w-full border p-2 rounded-lg" placeholder="Default: 123456" />
                </div>
                
                <div class="pt-2">
                    <button type="submit" :disabled="adding" class="w-full bg-primary text-white py-2 rounded-lg font-medium flex justify-center gap-2">
                        <Loader2 v-if="adding" class="animate-spin w-4 h-4" /> Simpan User
                    </button>
                </div>
            </form>
        </div>
      </div>

      <div class="mb-6 bg-blue-50 border border-blue-100 p-4 rounded-lg flex items-start gap-3">
         <div class="bg-blue-100 p-1 rounded text-blue-600"><KeyRound class="w-4 h-4"/></div>
         <div>
            <h4 class="text-sm font-bold text-blue-800">Informasi Keamanan</h4>
            <p class="text-xs text-blue-600 mt-1">
                Password default untuk user baru (Import/Manual) adalah <b>123456</b>. 
                User akan dipaksa mengganti password saat login pertama kali.
            </p>
         </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div class="bg-white p-4 rounded-xl border border-gray-200 shadow-sm flex items-center gap-4">
          <div class="p-3 bg-blue-50 text-primary rounded-lg"><UserCheck class="w-6 h-6" /></div>
          <div><p class="text-xs text-secondary">Guru Aktif</p><p class="text-2xl font-bold text-accent">{{ stats.guru_active || 0 }}</p></div>
        </div>
        <div class="bg-white p-4 rounded-xl border border-gray-200 shadow-sm flex items-center gap-4">
          <div class="p-3 bg-green-50 text-green-600 rounded-lg"><Users class="w-6 h-6" /></div>
          <div><p class="text-xs text-secondary">Murid Aktif</p><p class="text-2xl font-bold text-accent">{{ stats.murid_active || 0 }}</p></div>
        </div>
        <div class="bg-white p-4 rounded-xl border border-gray-200 shadow-sm flex items-center gap-4">
          <div class="p-3 bg-gray-100 text-gray-600 rounded-lg"><Users class="w-6 h-6" /></div>
          <div><p class="text-xs text-secondary">Total User</p><p class="text-2xl font-bold text-accent">{{ stats.total_registered || 0 }}</p></div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <div class="flex flex-col md:flex-row justify-between items-center p-4 border-b border-gray-100 gap-4">
          <div class="flex bg-gray-100 p-1 rounded-lg">
            <button v-for="role in ['GURU', 'MURID']" :key="role" @click="activeTab = role"
              :class="['px-4 py-2 text-sm font-medium rounded-md transition', activeTab === role ? 'bg-white shadow text-primary' : 'text-secondary hover:text-primary']">
              Data {{ role === 'GURU' ? 'Guru' : 'Murid' }}
            </button>
          </div>
          <div class="relative w-full md:w-64">
            <Search class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
            <input v-model="searchQuery" @input="handleSearch" type="text" placeholder="Cari nama / nomor induk..." 
              class="w-full pl-9 pr-4 py-2 text-sm border border-gray-300 rounded-lg outline-none focus:border-primary">
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-50 text-secondary border-b border-gray-200">
              <tr>
                <th class="px-6 py-3 font-semibold">Identitas</th>
                <th class="px-6 py-3 font-semibold">Username</th>
                <th class="px-6 py-3 font-semibold">Status</th>
                <th class="px-6 py-3 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="loading" class="animate-pulse"><td colspan="4" class="p-8 text-center text-gray-400">Memuat data...</td></tr>
              
              <tr v-else v-for="user in users" :key="user.id" class="hover:bg-gray-50/50 transition">
                <td class="px-6 py-4">
                    <div class="font-medium text-primary">{{ user.name }}</div>
                    <div v-if="user.serial_number" class="text-xs text-blue-600 bg-blue-50 border border-blue-100 px-1.5 py-0.5 rounded inline-block mt-1 font-mono">
                        {{ user.serial_number }}
                    </div>
                    <div v-else class="text-xs text-gray-400 italic mt-1 bg-gray-100 px-1 rounded inline-block">No ID</div>
                </td>
                <td class="px-6 py-4 text-secondary font-mono text-xs">{{ user.username }}</td>
                <td class="px-6 py-4">
                  <span :class="['inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-bold', user.status === 'ACTIVE' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700']">
                    {{ user.status === 'ACTIVE' ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right flex justify-end gap-2">
                  <button @click="openEdit(user)" class="p-2 border rounded hover:bg-gray-100 text-blue-600" title="Edit Data">
                    <Edit2 class="w-4 h-4" />
                  </button>
                  <button @click="resetPassword(user)" class="p-2 border rounded hover:bg-gray-100 text-secondary" title="Reset Password"><KeyRound class="w-4 h-4" /></button>
                  <button @click="toggleStatus(user)" class="p-2 border rounded hover:bg-gray-100" :class="user.status === 'ACTIVE' ? 'text-red-600' : 'text-green-600'">
                    <UserX v-if="user.status === 'ACTIVE'" class="w-4 h-4" />
                    <UserCheck v-else class="w-4 h-4" />
                  </button>
                </td>
              </tr>
              <tr v-if="!loading && users.length === 0">
                <td colspan="4" class="p-12 text-center text-secondary">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="showImportModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
         <div class="bg-white p-6 rounded-xl w-full max-w-lg shadow-xl relative max-h-[90vh] overflow-y-auto">
            <button @click="showImportModal = false" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600"><X class="w-5 h-5"/></button>
            
            <h3 class="font-bold text-lg mb-2 text-primary">Import Data User</h3>
            <p class="text-sm text-secondary mb-6">Upload file CSV untuk menambahkan Guru atau Murid secara massal.</p>

            <div v-if="importResult" class="mb-6 bg-gray-50 p-4 rounded-lg border border-gray-200">
                <h4 class="font-bold text-sm mb-2 flex items-center gap-2">
                    <CheckCircle class="w-4 h-4 text-green-600"/> Hasil Import:
                </h4>
                <div class="grid grid-cols-3 gap-2 text-center mb-4">
                    <div class="bg-white p-2 rounded border border-gray-200">
                        <div class="text-lg font-bold text-primary">{{ importResult.total_processed }}</div>
                        <div class="text-[10px] text-secondary">Diproses</div>
                    </div>
                    <div class="bg-green-50 p-2 rounded border border-green-200">
                        <div class="text-lg font-bold text-green-700">{{ importResult.success_count }}</div>
                        <div class="text-[10px] text-green-700">Sukses</div>
                    </div>
                    <div class="bg-yellow-50 p-2 rounded border border-yellow-200">
                        <div class="text-lg font-bold text-yellow-700">{{ importResult.skipped_count }}</div>
                        <div class="text-[10px] text-yellow-700">Dilewati</div>
                    </div>
                </div>
                
                <div v-if="importResult.skipped_details.length > 0">
                    <p class="text-xs font-bold text-secondary mb-1 flex items-center gap-1"><AlertTriangle class="w-3 h-3"/> Detail Skip:</p>
                    <ul class="text-[10px] text-secondary max-h-32 overflow-y-auto list-disc list-inside space-y-1 bg-white p-2 rounded border">
                        <li v-for="(msg, i) in importResult.skipped_details" :key="i">{{ msg }}</li>
                    </ul>
                </div>
                <button @click="importResult = null" class="mt-4 w-full py-2 bg-gray-200 text-secondary text-xs rounded hover:bg-gray-300">Upload Lagi</button>
            </div>

            <div v-else>
                <div class="border-2 border-dashed border-gray-300 rounded-xl p-8 text-center hover:bg-gray-50 transition cursor-pointer relative"
                    @dragover.prevent @drop.prevent="handleFileSelect">
                    <input type="file" accept=".csv" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" @change="handleFileSelect">
                    <Upload class="w-8 h-8 text-gray-400 mx-auto mb-2" />
                    <p class="text-sm font-medium text-gray-700" v-if="!importFile">Klik atau tarik file CSV ke sini</p>
                    <p class="text-sm font-bold text-primary" v-else>{{ importFile.name }}</p>
                </div>

                <div class="mt-4 flex justify-between items-center">
                    <button @click="downloadTemplate" class="text-xs text-blue-600 flex items-center gap-1 hover:underline">
                        <FileDown class="w-3 h-3" /> Download Template CSV
                    </button>
                </div>

                <button @click="handleImport" :disabled="!importFile || importing" class="w-full mt-6 bg-primary text-white py-2.5 rounded-lg font-medium flex justify-center gap-2 disabled:opacity-50">
                    <Loader2 v-if="importing" class="animate-spin w-4 h-4" /> {{ importing ? 'Mengupload...' : 'Proses Import' }}
                </button>
            </div>
         </div>
      </div>

      <div v-if="showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white p-6 rounded-xl w-full max-w-md shadow-xl relative">
            <button @click="showEditModal = false" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600"><X class="w-5 h-5"/></button>
            <h3 class="font-bold text-lg mb-4 text-primary">Edit Data User</h3>
            
            <form @submit.prevent="handleUpdateUser" class="space-y-4">
                <div>
                    <label class="label">Nama Lengkap</label>
                    <input v-model="editForm.name" type="text" class="input-field w-full border p-2 rounded-lg" required />
                </div>
                <div>
                    <label class="label">Username</label>
                    <input v-model="editForm.username" type="text" class="input-field w-full border p-2 rounded-lg" required />
                </div>
                <div>
                    <label class="label">Nomor Induk (NIP/NIS)</label>
                    <input v-model="editForm.serial_number" type="text" class="input-field w-full border p-2 rounded-lg" placeholder="Wajib diisi unik" />
                    <p class="text-xs text-secondary mt-1">Nomor ini akan digunakan untuk login awal & identifikasi.</p>
                </div>
                <div class="pt-2">
                    <button type="submit" :disabled="editing" class="w-full bg-primary text-white py-2 rounded-lg font-medium flex justify-center gap-2">
                        <Loader2 v-if="editing" class="animate-spin w-4 h-4" /> Simpan Perubahan
                    </button>
                </div>
            </form>
        </div>
      </div>

    </div>
  </DashboardLayout>
</template>