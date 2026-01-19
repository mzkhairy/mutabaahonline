<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { 
  LayoutDashboard, BookOpen, LogOut, Menu, X, Users, 
  FileText, Calendar, Lock 
} from 'lucide-vue-next'

const sidebarOpen = ref(false)
const authStore = useAuthStore()
const router = useRouter()

const toggleSidebar = () => sidebarOpen.value = !sidebarOpen.value

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

const menuItems = computed(() => {
  const role = authStore.user?.role
  let items = []

  if (role === 'ADMIN') {
    items = [
      { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
      { name: 'Kelola Akun', icon: Users, path: '/admin/users' },
      { name: 'Tahun Ajaran', icon: Calendar, path: '/admin/academic-years' },
      { name: 'Template Mutabaah', icon: FileText, path: '/admin/templates' },
    ]
  } else if (role === 'GURU') {
    items = [
      { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
      { name: 'Kelas Ajar', icon: BookOpen, path: '/guru/classes' },
      { name: 'Referensi Template', icon: FileText, path: '/guru/templates' },
    ]
  } else {
    items = [
      { name: 'Dashboard', icon: LayoutDashboard, path: '/mutabaah' },
      { name: 'Mutabaah Saya', icon: BookOpen, path: '/mutabaah-sheet' },
    ]
  }

  items.push({ name: 'Ganti Password', icon: Lock, path: '/profile' })
  return items
})
</script>

<template>
  <div class="min-h-screen bg-background flex flex-col md:flex-row">
    <header class="md:hidden bg-surface border-b border-gray-200 px-4 py-3 flex justify-between items-center sticky top-0 z-30 print:hidden shadow-sm">
      <div class="font-bold text-lg text-primary flex items-center gap-2">
         <BookOpen class="w-6 h-6"/> Mutabaah App
      </div>
      <button @click="toggleSidebar" class="text-primary hover:bg-gray-100 p-2 rounded-lg transition-colors focus:ring-2 focus:ring-primary/20">
        <Menu v-if="!sidebarOpen" class="w-6 h-6" />
        <X v-else class="w-6 h-6" />
      </button>
    </header>

    <aside :class="['fixed md:static inset-y-0 left-0 z-40 w-64 bg-surface border-r border-gray-200 transform transition-transform duration-300 ease-in-out md:translate-x-0 print:hidden shadow-xl md:shadow-none', sidebarOpen ? 'translate-x-0' : '-translate-x-full']">
      <div class="h-full flex flex-col bg-white">
        <div class="hidden md:flex items-center justify-center h-16 border-b border-gray-200 font-bold text-xl text-primary gap-2">
          <BookOpen class="w-6 h-6"/> Mutabaah App
        </div>
        
        <nav class="flex-1 p-4 space-y-2 overflow-y-auto">
          <router-link v-for="item in menuItems" :key="item.path" :to="item.path" active-class="bg-primary/5 text-primary font-semibold border-r-2 border-primary" class="flex items-center gap-3 px-4 py-3 rounded-lg text-secondary hover:bg-gray-50 hover:text-primary transition-colors" @click="sidebarOpen = false">
            <component :is="item.icon" class="w-5 h-5 shrink-0" />
            <span class="truncate">{{ item.name }}</span>
          </router-link>
        </nav>

        <div class="p-4 border-t border-gray-200 bg-gray-50/50">
          <div class="mb-3 px-2">
            <p class="text-sm font-semibold text-accent truncate">{{ authStore.user?.name || 'User' }}</p>
            <p class="text-xs text-secondary capitalize">{{ authStore.user?.role?.toLowerCase() || 'Guest' }}</p>
          </div>
          <button @click="handleLogout" class="w-full flex items-center gap-3 px-4 py-2 text-danger hover:bg-red-50 rounded-lg transition-colors text-sm font-medium">
            <LogOut class="w-4 h-4" /> Keluar
          </button>
        </div>
      </div>
    </aside>

    <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-black/50 z-30 md:hidden print:hidden backdrop-blur-sm"></div>

    <main class="flex-1 p-4 md:p-8 overflow-y-auto h-[calc(100vh-64px)] md:h-screen bg-gray-50">
      <slot />
    </main>
  </div>
</template>