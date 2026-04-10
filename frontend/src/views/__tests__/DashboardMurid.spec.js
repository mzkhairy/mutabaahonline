import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
// [FIX 1] Perbaiki Path Import (tambahkan /dashboard/)
import DashboardMurid from '../dashboard/DashboardMurid.vue'
import { createPinia, setActivePinia } from 'pinia'

// [FIX 2] Gunakan vi.hoisted() untuk Mock API
const { mockGet, mockPost } = vi.hoisted(() => {
  return {
    mockGet: vi.fn(),
    mockPost: vi.fn()
  }
})

// Setup Mock API
vi.mock('../../services/api', () => ({
  default: {
    get: mockGet,
    post: mockPost
  }
}))

// Mock Layout Dashboard
vi.mock('../../layouts/DashboardLayout.vue', () => ({
  default: { template: '<div><slot /></div>' }
}))

describe('DashboardMurid.vue', () => {
  it('renders "Mutabaah Siswa" title', async () => {
    setActivePinia(createPinia())
    
    // Mock Response Data Sesi Kosong
    mockGet.mockResolvedValue({ data: { data: [] } })

    const wrapper = mount(DashboardMurid, {
      global: { 
        stubs: ['router-link'] // [FIX 3] Stub router-link
      }
    })
    
    await flushPromises()

    expect(wrapper.text()).toContain('Mutabaah Siswa')
    expect(wrapper.text()).toContain('Ahlan wa Sahlan')
  })

  it('shows session list correctly', async () => {
    setActivePinia(createPinia())

    // Mock Response: Ada 1 Sesi Aktif
    const mockSession = [{
      id: 'sess-1',
      name: 'Setoran Subuh',
      date: '2025-01-27',
      is_student_input_allowed: true,
      has_entry: false,
      template_name: 'Hafalan Rutin'
    }]
    mockGet.mockResolvedValue({ data: { data: mockSession } })

    const wrapper = mount(DashboardMurid, {
      global: { 
        stubs: ['router-link'] 
      }
    })
    
    await flushPromises()

    // Cek apakah kartu sesi muncul
    expect(wrapper.text()).toContain('Setoran Subuh')
    expect(wrapper.text()).toContain('Belum Diisi')
  })
})