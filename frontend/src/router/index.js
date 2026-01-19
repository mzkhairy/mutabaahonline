import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// --- AUTH & SETUP VIEWS ---
import LoginView from '../views/LoginView.vue'
import RegisterAdminView from '../views/RegisterAdminView.vue'
import FirstLoginView from '../views/FirstLoginView.vue'

// --- DASHBOARD ---
import DashboardView from '../views/DashboardView.vue'
import ProfileView from '../views/ProfileView.vue'

// --- ADMIN VIEWS ---
import UserManagementView from '../views/admin/UserManagementView.vue'
import AcademicYearsView from '../views/admin/AcademicYearsView.vue'
import TemplatesView from '../views/admin/TemplatesView.vue'
import ClassesView from '../views/admin/ClassesView.vue'
import ClassDetailView from '../views/admin/ClassDetailView.vue'

// --- GURU VIEWS ---
import GuruClassesView from '../views/guru/ClassesView.vue'
import GuruClassDetailView from '../views/guru/ClassDetailView.vue'
import GuruTemplatesView from '../views/guru/TemplatesView.vue'
import GuruSessionDetailView from '../views/guru/SessionDetailView.vue'

// --- MURID VIEWS ---
import DashboardMurid from '../views/dashboard/DashboardMurid.vue'
import MutabaahListMurid from '../views/mutabaah/MutabaahListMurid.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // --- PUBLIC ROUTES ---
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register-institution',
      name: 'register-institution',
      component: RegisterAdminView
    },

    // --- SETUP ACCOUNT (FORCE REDIRECT) ---
    {
      path: '/setup-account',
      name: 'setup-account',
      component: FirstLoginView,
      meta: { requiresAuth: true }
    },

    // --- DASHBOARD --- Protected Routes
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },

    // --- ADMIN ROUTES ---
    {
      path: '/admin/users',
      name: 'admin-users',
      component: UserManagementView,
      meta: { requiresAuth: true, role: 'ADMIN' }
    },
    {
      path: '/admin/academic-years',
      name: 'admin-academic-years',
      component: AcademicYearsView,
      meta: { requiresAuth: true, role: 'ADMIN' }
    },
    {
      path: '/admin/templates',
      name: 'admin-templates',
      component: TemplatesView,
      meta: { requiresAuth: true, role: 'ADMIN' }
    },
    {
      path: '/admin/classes',
      name: 'admin-classes',
      component: ClassesView,
      meta: { requiresAuth: true, role: 'ADMIN' }
    },
    {
      path: '/admin/classes/:id',
      name: 'admin-class-detail',
      component: ClassDetailView,
      meta: { requiresAuth: true, role: 'ADMIN' }
    },

    // --- GURU ROUTES ---
    {
      path: '/guru/classes',
      name: 'guru-classes',
      component: GuruClassesView,
      meta: { requiresAuth: true, role: 'GURU' }
    },
    {
      path: '/guru/classes/:id',
      name: 'guru-class-detail',
      component: GuruClassDetailView,
      meta: { requiresAuth: true, role: 'GURU' }
    },
    {
      path: '/guru/students/:studentId/report',
      name: 'guru-student-report',
      component: () => import('../views/guru/StudentReportView.vue'),
      meta: { requiresAuth: true, role: 'GURU' }
    },
    {
      path: '/guru/templates',
      name: 'guru-templates',
      component: GuruTemplatesView,
      meta: { requiresAuth: true, role: 'GURU' }
    },
    {
      path: '/guru/sessions/:id',
      name: 'guru-session-detail',
      component: GuruSessionDetailView,
      meta: { requiresAuth: true, role: 'GURU' }
    },

    // --- MURID ROUTES ---
    {
      path: '/mutabaah',
      name: 'murid-dashboard',
      component: DashboardMurid,
      meta: { requiresAuth: true, role: 'MURID' }
    },
    { 
      path: '/mutabaah-sheet', 
      name: 'murid-mutabaah-sheet', 
      component: MutabaahListMurid, 
      meta: { requiresAuth: true, role: 'MURID' } 
    },
  ]
})

// --- NAVIGATION GUARDS ---
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated
  // Pastikan akses properti getter state dengan benar (tergantung definisi di store)
  // Di store auth.js sebelumnya: getters: { mustChangePassword: ... }
  const mustChange = authStore.mustChangePassword 

  // 1. Cek Login untuk Halaman Protected
  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: 'login' })
  }

  // 2. Logic Force Setup Account (Jika User Wajib Ganti Password)
  if (isAuthenticated && mustChange) {
    // Arahkan ke 'setup-account' agar user mengisi Username & Password baru
    // sesuai backend endpoint /auth/setup-account
    if (to.name !== 'setup-account') {
      return next({ name: 'setup-account' })
    }
  }

  // 3. Cegah User yang SUDAH Setup masuk kembali ke halaman Setup
  if (isAuthenticated && !mustChange && to.name === 'setup-account') {
    return next({ name: 'dashboard' })
  }

  // 4. Redirect User Login ke Dashboard
  if (isAuthenticated && to.name === 'login') {
    if (!mustChange) {
        const role = authStore.user?.role
        if (role === 'MURID') {
            return next({ name: 'murid-dashboard' })
        } else {
            return next({ name: 'dashboard' })
        }
    }
  }

  // 5. Cek Role Access
  if (to.meta.role && authStore.user?.role !== to.meta.role) {
    // Redirect unauthorized role ke dashboard masing-masing
    return next({ name: 'dashboard' })
  }

  next()
})

export default router