import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
// [FIX] Path diperbaiki agar mengarah ke folder mutabaah
import MutabaahListMurid from '../mutabaah/MutabaahListMurid.vue'
import { createPinia, setActivePinia } from 'pinia'

// 1. Mock Layout & Komponen Global
vi.mock('../../layouts/DashboardLayout.vue', () => ({
  default: { template: '<div><slot /></div>' }
}))

// 2. Setup Mock API dengan vi.hoisted
const { mockApi } = vi.hoisted(() => {
  return { 
    mockApi: { get: vi.fn() } 
  }
})

vi.mock('../../services/api', () => ({ default: mockApi }))

describe('MutabaahListMurid.vue', () => {
  // Setup data dummy
  const mockClasses = [
    { id: 'cls-1', name: 'Kelas 7A', academic_year_name: '2025/2026' }
  ]

  const mockSessions = [
    { 
      id: 'sess-1', 
      name: 'Setoran Subuh', 
      date: '2025-01-27', 
      template_name: 'Hafalan',
      template_structure: JSON.stringify([{ type: 'student_activity', name: 'Surah' }]),
      class_activity_data: {}
    }
  ]

  const mockEntries = [
    { 
      session_id: 'sess-1', 
      student_id: 'std-1', 
      attendance: 'HADIR', 
      note: 'Lancar',
      student_activity_data: { 'Surah': 'Al-Mulk' }
    }
  ]

  it('loads class dropdown on mount', async () => {
    setActivePinia(createPinia())

    // Mock API responses
    mockApi.get.mockImplementation((url) => {
      if (url === '/student/classes') return Promise.resolve({ data: { data: mockClasses } })
      return Promise.resolve({ data: { data: [] } })
    })

    const wrapper = mount(MutabaahListMurid, {
      global: { stubs: ['router-link'] }
    })

    await flushPromises()

    // Cek Judul
    expect(wrapper.text()).toContain('Mutabaah Saya')

    // Cek Dropdown terisi
    const options = wrapper.findAll('option')
    expect(options.length).toBe(1)
    expect(options[0].text()).toContain('Kelas 7A')
  })

  it('displays report table when data is loaded', async () => {
    setActivePinia(createPinia())

    mockApi.get.mockImplementation((url) => {
      if (url === '/student/classes') return Promise.resolve({ data: { data: mockClasses } })
      if (url === '/sessions/student') return Promise.resolve({ data: { data: mockSessions } })
      if (url === '/mutabaah') return Promise.resolve({ data: { data: mockEntries } })
      return Promise.resolve({ data: { data: [] } })
    })

    const wrapper = mount(MutabaahListMurid, {
      global: { stubs: ['router-link'] }
    })

    await flushPromises()

    // Cek Tabel Muncul & Data Rendered
    expect(wrapper.text()).toContain('Setoran Subuh') // Nama Sesi
    expect(wrapper.text()).toContain('Hafalan')       // Nama Template
    expect(wrapper.text()).toContain('Al-Mulk')       // Data Nilai
    expect(wrapper.text()).toContain('Lancar')        // Catatan
  })

  it('shows empty state when no data', async () => {
    setActivePinia(createPinia())

    mockApi.get.mockImplementation((url) => {
      if (url === '/student/classes') return Promise.resolve({ data: { data: mockClasses } })
      // Return array kosong untuk sesi
      if (url === '/sessions/student') return Promise.resolve({ data: { data: [] } })
      return Promise.resolve({ data: { data: [] } })
    })

    const wrapper = mount(MutabaahListMurid, {
      global: { stubs: ['router-link'] }
    })

    await flushPromises()

    // Pastikan pesan kosong muncul
    expect(wrapper.text()).toContain('Tidak ada data mutabaah')
  })
})