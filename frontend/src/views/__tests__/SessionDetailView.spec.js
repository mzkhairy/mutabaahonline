import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
// [FIX 1] Perbaiki Path Import (tambahkan /guru/)
import SessionDetailView from '../guru/SessionDetailView.vue'

// Mock Router
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { id: 'sess-1' } }),
  useRouter: () => ({ back: vi.fn() })
}))

// Mock Layout
vi.mock('../../layouts/DashboardLayout.vue', () => ({
  default: { template: '<div><slot /></div>' }
}))

// [FIX 2] Gunakan vi.hoisted() untuk Mock API
const { mockApi } = vi.hoisted(() => {
  return { 
    mockApi: { 
      get: vi.fn(), 
      put: vi.fn(), 
      post: vi.fn() 
    } 
  }
})

vi.mock('../../services/api', () => ({ default: mockApi }))

describe('SessionDetailView.vue', () => {
  it('renders session details', async () => {
    // Setup Mock Responses
    mockApi.get.mockImplementation((url) => {
      // 1. Get Session Detail
      if (url.includes('/sessions/sess-1')) {
        return Promise.resolve({ 
          data: { 
            data: { 
              id: 'sess-1', 
              name: 'Setoran Subuh', 
              date: '2025-01-27',
              template_name: 'Hafalan',
              template_structure: JSON.stringify([{ type: 'student_activity', name: 'Hafalan' }]),
              class_activity_data: {},
              is_student_input_allowed: true 
            } 
          } 
        })
      }
      // 2. Get Students in Class
      if (url.includes('/students')) {
        return Promise.resolve({ 
          data: { data: [{ id: 's1', name: 'Ahmad', serial_number: '123' }] } 
        })
      }
      // 3. Get Existing Mutabaah Entries
      if (url.includes('/mutabaah')) {
        return Promise.resolve({ data: { data: [] } })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(SessionDetailView, {
        global: { stubs: ['router-link'] }
    })
    
    await flushPromises()

    // Cek Judul Sesi
    expect(wrapper.text()).toContain('Setoran Subuh')
    
    // Cek Tombol "Simpan" muncul
    expect(wrapper.text()).toContain('Simpan')
    
    // Cek Tombol "Kunci" muncul (karena input siswa aktif)
    expect(wrapper.text()).toContain('Kunci') 
    
    // Cek Nama Siswa muncul di tabel
    expect(wrapper.text()).toContain('Ahmad')
  })
})