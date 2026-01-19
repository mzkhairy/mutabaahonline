<script setup>
import { ref, onMounted } from 'vue'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import api from '../../services/api'
import { FileText, List, Loader2 } from 'lucide-vue-next'

const templates = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.get('/templates')
    templates.value = res.data.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchData())
</script>

<template>
  <DashboardLayout>
    <div class="max-w-5xl mx-auto">
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-accent">Daftar Template Penilaian</h1>
        <p class="text-secondary text-sm">Referensi format penilaian yang tersedia.</p>
      </div>

      <div v-if="loading" class="py-12 flex justify-center"><Loader2 class="animate-spin w-8 h-8 text-primary" /></div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="t in templates" :key="t.id" class="bg-white p-5 rounded-xl border border-gray-200 hover:shadow-md transition">
          <div class="flex items-start justify-between mb-3">
            <div class="p-2 bg-purple-50 text-purple-600 rounded-lg">
              <FileText class="w-6 h-6" />
            </div>
          </div>
          <h4 class="font-bold text-lg text-accent mb-2">{{ t.name }}</h4>
          
          <div class="bg-gray-50 rounded-lg p-3 space-y-2">
            <p class="text-xs font-bold text-secondary uppercase mb-1">Kolom Penilaian:</p>
            <div v-for="(item, i) in t.structure" :key="i" class="text-sm flex items-center gap-2 text-secondary">
              <List class="w-3 h-3 text-gray-400 shrink-0" /> 
              <span class="truncate flex-1">{{ item.name }}</span>
              <span class="text-[10px] bg-white border px-1.5 py-0.5 rounded text-gray-500">
                {{ item.type === 'class_activity' ? 'Kelas' : 'Personal' }}
              </span>
            </div>
          </div>
        </div>

        <div v-if="templates.length === 0" class="col-span-full text-center py-12 text-secondary bg-gray-50 rounded-xl border-dashed border">
          Belum ada template tersedia. Hubungi Admin.
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>