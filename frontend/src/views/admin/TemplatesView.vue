<script setup>
import { ref, onMounted } from 'vue'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { FileText, Plus, Trash2, Loader2, CheckCircle, GripVertical } from 'lucide-vue-next'

const templates = ref([])
const loading = ref(false)
const showModal = ref(false)

const form = ref({
  name: '',
  structure: [] 
})

// [FIX] Default type kegiatan kelas
const newItem = ref({ name: '', type: 'class_activity' }) 

const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await api.get('/templates')
    templates.value = res.data.data || []
  } catch (err) { console.error(err) } 
  finally { loading.value = false }
}

const addField = () => {
  if (!newItem.value.name) return
  form.value.structure.push({ 
      name: newItem.value.name, 
      type: newItem.value.type 
  })
  newItem.value.name = ''
}

const removeField = (index) => {
  form.value.structure.splice(index, 1)
}

const handleCreate = async () => {
  if (form.value.structure.length === 0) return alert('Tambahkan minimal 1 kolom penilaian')
  try {
    // Kirim structure sebagai Array Object (Go akan handle konversi JSONB)
    const payload = {
        name: form.value.name,
        structure: form.value.structure 
    }
    await api.post('/templates', payload)
    showModal.value = false
    form.value.name = ''
    form.value.structure = []
    fetchTemplates()
  } catch (err) { 
      console.error(err)
      alert('Gagal membuat template') 
  }
}

const handleDelete = async (id) => {
    if(!confirm('Hapus template ini?')) return
    try {
        await api.delete(`/templates/${id}`)
        fetchTemplates()
    } catch(e) { alert('Gagal menghapus') }
}

onMounted(fetchTemplates)
</script>

<template>
  <DashboardLayout>
    <div class="space-y-6">
        <div class="flex justify-between items-center bg-white p-4 rounded-xl shadow-sm border border-gray-200">
            <div>
                <h2 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                    <FileText class="w-5 h-5 text-primary"/> Template Penilaian
                </h2>
                <p class="text-sm text-gray-500">Buat format penilaian untuk kegiatan mutabaah.</p>
            </div>
            <button @click="showModal = true" class="bg-primary hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-bold flex items-center gap-2 shadow-sm transition">
                <Plus class="w-4 h-4"/> Buat Template
            </button>
        </div>

        <div v-if="loading" class="py-10 text-center"><Loader2 class="w-8 h-8 animate-spin mx-auto text-primary"/></div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="t in templates" :key="t.id" class="bg-white border border-gray-200 rounded-xl p-5 hover:shadow-md transition group">
                <div class="flex justify-between items-start mb-4">
                    <h3 class="font-bold text-gray-800 text-lg">{{ t.name }}</h3>
                    <button @click="handleDelete(t.id)" class="text-gray-300 hover:text-red-500"><Trash2 class="w-4 h-4"/></button>
                </div>
                <div class="space-y-2">
                    <div v-for="(item, idx) in (t.structure || []).slice(0, 5)" :key="idx" class="flex items-center gap-2 text-sm text-gray-600">
                        <CheckCircle class="w-3 h-3 text-green-500"/>
                        <span>{{ item.name }}</span>
                        <span class="text-[10px] px-1.5 py-0.5 rounded bg-gray-100 uppercase font-bold text-gray-500">
                            {{ item.type === 'class_activity' ? 'Kelas' : 'Individu' }}
                        </span>
                    </div>
                    <div v-if="t.structure && t.structure.length > 5" class="text-xs text-gray-400 italic pl-5">
                        + {{ t.structure.length - 5 }} kolom lainnya
                    </div>
                </div>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm">
            <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl p-6 h-[90vh] flex flex-col">
                <h3 class="font-bold text-lg mb-4">Buat Template Baru</h3>
                
                <div class="mb-4">
                    <label class="block text-xs font-bold text-gray-500 uppercase mb-1">Nama Template</label>
                    <input v-model="form.name" type="text" placeholder="Contoh: Mutabaah Harian Tahfidz" class="w-full border p-2 rounded-lg" required>
                </div>

                <div class="flex-1 overflow-y-auto border border-gray-100 rounded-lg bg-gray-50 p-4 mb-4 space-y-2">
                    <div v-if="form.structure.length === 0" class="text-center text-gray-400 py-10 italic">Belum ada kolom penilaian.</div>
                    
                    <div v-for="(field, idx) in form.structure" :key="idx" class="flex items-center justify-between bg-white p-3 rounded shadow-sm border border-gray-200">
                        <div class="flex items-center gap-3">
                            <GripVertical class="w-4 h-4 text-gray-300 cursor-move"/>
                            <div>
                                <div class="font-bold text-gray-800">{{ field.name }}</div>
                                <div class="text-xs text-blue-600 font-medium bg-blue-50 inline-block px-1 rounded mt-0.5">
                                    {{ field.type === 'class_activity' ? 'Kegiatan Kelas' : 'Kegiatan Individu' }}
                                </div>
                            </div>
                        </div>
                        <button @click="removeField(idx)" class="text-red-400 hover:text-red-600"><Trash2 class="w-4 h-4"/></button>
                    </div>
                </div>

                <div class="border-t pt-4">
                    <label class="block text-xs font-bold text-gray-500 uppercase mb-2">Tambah Kolom Baru</label>
                    <div class="flex flex-col sm:flex-row gap-2 mb-4">
                        <input v-model="newItem.name" type="text" placeholder="Nama Kegiatan (Misal: Shalat Dhuha)" class="flex-1 border p-2 rounded-lg text-sm">
                        
                        <select v-model="newItem.type" class="border p-2 rounded-lg text-sm bg-white cursor-pointer">
                            <option value="class_activity">Kegiatan Kelas</option>
                            <option value="student_activity">Kegiatan Individu</option>
                        </select>

                        <button @click="addField" class="bg-gray-800 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-black">Tambah</button>
                    </div>

                    <div class="flex justify-end gap-2">
                        <button @click="showModal = false" class="px-4 py-2 text-gray-600">Batal</button>
                        <button @click="handleCreate" class="px-4 py-2 bg-primary text-white rounded-lg font-bold">Simpan Template</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </DashboardLayout>
</template>