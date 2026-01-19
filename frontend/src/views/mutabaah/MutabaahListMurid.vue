<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../../services/api'
import DashboardLayout from '../../layouts/DashboardLayout.vue'
import { Printer, Loader2, BookOpen } from 'lucide-vue-next'

const loading = ref(false)
const sessions = ref([])
const enrolledClasses = ref([]) // Menyimpan daftar kelas murid
const selectedClassId = ref('') // Menyimpan ID kelas terpilih untuk filter

// 1. Ambil Daftar Kelas Murid
const fetchClasses = async () => {
    try {
        const res = await api.get('/student/classes')
        enrolledClasses.value = res.data.data || []
        
        // Auto-select kelas pertama (biasanya yang terbaru/aktif)
        if (enrolledClasses.value.length > 0) {
            selectedClassId.value = enrolledClasses.value[0].id
        }
    } catch(e) { 
        console.error("Gagal memuat daftar kelas", e) 
    }
}

// 2. Ambil Data Mutabaah berdasarkan Kelas Terpilih
const fetchData = async () => {
  if (!selectedClassId.value) return 

  loading.value = true
  try {
    const [resSessions, resEntries] = await Promise.all([
        // [PENTING] Kirim parameter class_id agar backend memfilter sesi
        api.get('/sessions/student', { 
            params: { mode: 'report', class_id: selectedClassId.value } 
        }),
        api.get('/mutabaah') // Entri tetap diambil semua untuk mapping
    ])

    const rawSessions = resSessions.data.data || []
    const entries = resEntries.data.data || []

    sessions.value = rawSessions.map(s => {
      let structure = []
      try { 
          structure = Array.isArray(s.template_structure) ? s.template_structure : JSON.parse(s.template_structure || '[]') 
      } catch (e) { structure = [] }
      
      // Sort: Aktivitas kelas (ceklis guru) di atas, input murid di bawah
      const fields = structure.sort((a, b) => {
          if (a.type === 'class_activity') return -1
          return 0
      })

      // Cari entri milik sesi ini
      const existingEntry = entries.find(e => String(e.session_id) === String(s.id))
      
      const values = {}
      fields.forEach(f => {
        if (f.type === 'class_activity') values[f.name] = s.class_activity_data?.[f.name] || '-'
        else values[f.name] = existingEntry?.student_activity_data?.[f.name] || '-'
      })

      return { 
          ...s, 
          fields, 
          entry: existingEntry || { attendance: 'BELUM', note: '-' }, 
          values, 
          has_entry: !!existingEntry 
      }
    })

  } catch (err) { 
      console.error(err)
      alert("Gagal memuat data mutabaah")
  } finally { 
      loading.value = false 
  }
}

// Refresh data saat dropdown berubah
watch(selectedClassId, () => {
    fetchData()
})

onMounted(async () => {
    await fetchClasses()
})
</script>

<template>
  <DashboardLayout>
    <div class="max-w-[1400px] mx-auto pb-10 font-sans px-2 sm:px-4">
      
      <div class="bg-white p-4 rounded-lg border border-gray-200 shadow-sm mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 print:hidden">
        <div>
            <h1 class="text-lg font-bold text-gray-800">Mutabaah Saya</h1>
            <p class="text-xs sm:text-sm text-gray-500">Rekapitulasi kegiatan dan ibadah harian.</p>
        </div>

        <div class="flex gap-2 w-full sm:w-auto">
            <div class="relative w-full sm:w-72">
                <select v-model="selectedClassId" class="w-full appearance-none bg-blue-50 border border-blue-200 text-blue-800 py-2 pl-3 pr-8 rounded-lg text-xs sm:text-sm font-bold focus:outline-none focus:ring-2 focus:ring-primary/50 cursor-pointer shadow-sm">
                    <option v-for="cls in enrolledClasses" :key="cls.id" :value="cls.id">
                        {{ cls.name }} - {{ cls.academic_year_name }}
                    </option>
                </select>
                <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-blue-600">
                    <BookOpen class="w-4 h-4" />
                </div>
            </div>

            <button onclick="window.print()" class="flex items-center justify-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-900 text-white rounded text-sm font-bold shadow-sm transition active:scale-95">
                <Printer class="w-4 h-4"/> <span class="hidden sm:inline">Cetak</span>
            </button>
        </div>
      </div>

      <div class="hidden print:block mb-6 text-center border-b-2 border-black pb-4">
          <h2 class="text-xl font-bold uppercase tracking-wide mb-2">Laporan Mutabaah Yaumiyah</h2>
          <div class="flex justify-center gap-8 text-sm font-medium" v-if="selectedClassId">
              <span>Kelas: {{ enrolledClasses.find(c => c.id === selectedClassId)?.name }}</span>
              <span>Tahun Ajaran: {{ enrolledClasses.find(c => c.id === selectedClassId)?.academic_year_name }}</span>
          </div>
      </div>

      <div class="bg-white border border-gray-400 print:border-none shadow-sm overflow-hidden rounded-lg">
        
        <div v-if="loading" class="p-10 text-center text-gray-400 print:hidden">
            <Loader2 class="w-6 h-6 animate-spin mx-auto mb-2"/> Memuat data...
        </div>

        <div v-else-if="sessions.length === 0" class="p-10 text-center text-gray-400 italic">
            Tidak ada data mutabaah untuk kelas ini.
        </div>

        <div v-else class="overflow-x-auto w-full">
            <table class="w-full text-[10px] sm:text-[11px] border-collapse text-left font-medium min-w-[600px] sm:min-w-full">
                <thead class="bg-gray-100 text-gray-900 font-bold uppercase border-b-2 border-gray-500 print:bg-gray-200">
                    <tr>
                        <th class="border border-gray-400 px-2 py-2 w-8 text-center bg-gray-200">No</th>
                        <th class="border border-gray-400 px-2 py-2 w-24 bg-gray-50">Hari / Tanggal</th>
                        <th class="border border-gray-400 px-2 py-2 w-32 sm:w-48 bg-gray-50">Materi & Kegiatan</th>
                        <th class="border border-gray-400 px-2 py-2 min-w-[200px] bg-white text-center">
                            Laporan Capaian
                        </th>
                        <th class="border border-gray-400 px-2 py-2 w-32 sm:w-48 bg-gray-50">Catatan</th>
                    </tr>
                </thead>

                <tbody>
                    <tr v-for="(session, idx) in sessions" :key="session.id" class="print:break-inside-avoid hover:bg-gray-50">
                        <td class="border border-gray-400 px-1 py-1.5 text-center font-mono text-gray-500">{{ idx + 1 }}</td>

                        <td class="border border-gray-400 px-2 py-1.5 align-top">
                            <div class="font-bold text-black">{{ new Date(session.date).toLocaleDateString('id-ID', {day:'numeric', month:'numeric', year:'2-digit'}) }}</div>
                            <div class="text-[9px] uppercase text-gray-600">{{ new Date(session.date).toLocaleDateString('id-ID', {weekday:'long'}) }}</div>
                        </td>

                        <td class="border border-gray-400 px-2 py-1.5 align-top">
                            <div class="font-bold text-black leading-tight mb-1">{{ session.name }}</div>
                            <div class="text-[9px] text-gray-600 border border-gray-300 rounded px-1 inline-block bg-gray-50">
                                {{ session.template_name }}
                            </div>
                        </td>

                        <td class="border border-gray-400 px-2 py-1.5 align-top bg-white">
                            <div v-if="!session.has_entry" class="h-full flex items-center justify-center text-gray-300 text-[10px] italic py-2">(Belum diisi)</div>

                            <div v-else-if="session.entry.attendance === 'HADIR'" class="flex flex-col gap-1">
                                <div v-for="field in session.fields" :key="field.name" 
                                     class="flex justify-between items-center border-b border-dotted border-gray-300 last:border-0 pb-0.5">
                                    <span class="text-[9px] font-bold uppercase" :class="field.type === 'class_activity' ? 'text-blue-700' : 'text-gray-600'">
                                        {{ field.name }}
                                    </span>
                                    <span class="font-mono font-bold text-black">{{ session.values[field.name] }}</span>
                                </div>
                                <div v-if="session.fields.length === 0" class="text-center italic text-gray-400 text-[9px]">- Kosong -</div>
                            </div>

                            <div v-else class="h-full flex items-center justify-center py-2">
                                <span class="font-bold px-3 py-1 rounded text-[10px] border uppercase tracking-wider"
                                    :class="{
                                        'bg-yellow-50 text-yellow-700 border-yellow-300': session.entry.attendance === 'IZIN',
                                        'bg-red-50 text-red-700 border-red-300': session.entry.attendance === 'SAKIT',
                                        'bg-gray-100 text-gray-700 border-gray-300': session.entry.attendance === 'ALPA'
                                    }">
                                    {{ session.entry.attendance }}
                                </span>
                            </div>
                        </td>

                        <td class="border border-gray-400 px-2 py-1.5 align-top italic text-gray-600 text-[10px]">
                            {{ session.entry.note !== '-' ? session.entry.note : '' }}
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    
        <div class="hidden print:flex mt-8 justify-between px-10 text-xs text-black">
            <div class="text-center">
                <p>Orang Tua Wali</p>
                <div class="h-16"></div>
                <p>( ........................... )</p>
            </div>
            <div class="text-center">
                <p>Guru Pendamping</p>
                <div class="h-16"></div>
                <p>( ........................... )</p>
            </div>
        </div>

      </div>
    </div>
  </DashboardLayout>
</template>