import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
// [FIX 1] Perbaiki Path Import (tambahkan /admin/)
import UserManagementView from '../admin/UserManagementView.vue'

// Mock Layout
vi.mock('../../layouts/DashboardLayout.vue', () => ({
  default: { template: '<div><slot /></div>' }
}))

// [FIX 2] Gunakan vi.hoisted() untuk Mock API
const { mockGet, mockPost } = vi.hoisted(() => {
  return {
    mockGet: vi.fn(),
    mockPost: vi.fn()
  }
})

vi.mock('../../services/api', () => ({
  default: {
    get: mockGet,
    post: mockPost,
    put: vi.fn(),   // Tambahan mock method lain yg dipakai komponen
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

describe('UserManagementView.vue', () => {
  it('renders user list and stats', async () => {
    // Mock Response Data
    mockGet.mockResolvedValue({
      data: {
        data: {
          users: [
            { id: '1', name: 'Guru Budi', role: 'GURU', status: 'ACTIVE', username: 'guru_budi' }
          ],
          stats: { 
            total_registered: 100,
            guru_active: 10,
            murid_active: 90
          }
        }
      }
    })

    const wrapper = mount(UserManagementView, {
      global: { stubs: ['router-link'] }
    })
    
    await flushPromises()

    // Assert UI
    expect(wrapper.text()).toContain('Kelola Akun')
    expect(wrapper.text()).toContain('Guru Budi')
    expect(wrapper.text()).toContain('Total User')
  })

  it('opens add user modal', async () => {
    // Mock Data Kosong agar render berhasil
    mockGet.mockResolvedValue({ data: { data: { users: [], stats: {} } } })

    const wrapper = mount(UserManagementView, {
        global: { stubs: ['router-link'] }
    })
    await flushPromises()

    // Cari tombol "Tambah User" & Klik
    const buttons = wrapper.findAll('button')
    const addButton = buttons.find(b => b.text().includes('Tambah User'))
    
    expect(addButton).toBeDefined()
    await addButton.trigger('click')

    // Cek Modal Muncul
    expect(wrapper.text()).toContain('Tambah User Baru')
  })
})