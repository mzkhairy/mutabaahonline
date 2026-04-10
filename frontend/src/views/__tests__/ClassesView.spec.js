import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
// [FIX 1] Perbaiki Path Import (tambahkan /admin/)
import ClassesView from '../admin/ClassesView.vue' 

// Mock Layout
vi.mock('../../layouts/DashboardLayout.vue', () => ({
  default: { template: '<div><slot /></div>' }
}))

// [FIX 2] Gunakan vi.hoisted untuk mock API agar tidak kena error "initialization"
const { mockApi } = vi.hoisted(() => {
  return { 
    mockApi: { 
      get: vi.fn(), 
      post: vi.fn(),
      delete: vi.fn() 
    } 
  }
})

vi.mock('../../services/api', () => ({ default: mockApi }))

// Mock Router (Karena ClassesView menggunakan useRoute & useRouter)
vi.mock('vue-router', () => ({
  useRoute: () => ({ query: {} }),
  useRouter: () => ({ push: vi.fn() })
}))

describe('ClassesView.vue', () => {
  it('displays list of classes', async () => {
    // Setup Mock Responses untuk endpoint yang dipanggil saat onMounted
    mockApi.get.mockImplementation((url) => {
      // 1. Mock List Kelas
      if (url === '/admin/classes') {
        return Promise.resolve({ 
          data: { 
            data: [
              { id: 'c1', name: 'Kelas 7A', level: '7', student_count: 20, teacher_name: 'Guru Budi', academic_year_name: '2025/2026' }
            ] 
          } 
        })
      }
      // 2. Mock List Guru (fetchMasterData)
      if (url === '/users') {
        return Promise.resolve({ data: { data: { users: [] } } })
      } 
      // 3. Mock Tahun Ajaran (fetchMasterData)
      if (url === '/academic-years') {
        return Promise.resolve({ data: { data: [{ id: 'ay1', name: '2025/2026', is_active: true }] } })
      }
      
      return Promise.resolve({ data: { data: [] } })
    })

    const wrapper = mount(ClassesView, {
        global: { stubs: ['router-link'] }
    })
    
    // Tunggu semua promise fetch selesai
    await flushPromises()

    // Assert render berhasil
    expect(wrapper.text()).toContain('Manajemen Kelas')
    expect(wrapper.text()).toContain('Kelas 7A')
    expect(wrapper.text()).toContain('20 Murid')
  })
})