import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import RegisterAdminView from '../RegisterAdminView.vue'

// [FIX 1] Gunakan vi.hoisted() untuk membuat mock variable
// Ini memastikan variabel mockPost sudah ada SEBELUM vi.mock dijalankan
const { mockPost } = vi.hoisted(() => {
  return { mockPost: vi.fn() }
})

// [FIX 2] Mock API menggunakan variabel dari vi.hoisted
vi.mock('../../services/api', () => ({
  default: { post: mockPost }
}))

// Mock Router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush })
}))

describe('RegisterAdminView.vue', () => {
  it('renders registration form', () => {
    const wrapper = mount(RegisterAdminView, {
      global: { stubs: ['router-link'] }
    })
    expect(wrapper.text()).toContain('Daftarkan Sekolah')
    expect(wrapper.find('input[placeholder="Contoh: Ponpes Al-Hidayah"]').exists()).toBe(true)
  })

  it('submits registration successfully', async () => {
    const wrapper = mount(RegisterAdminView, {
      global: { stubs: ['router-link'] }
    })
    
    // Isi Form
    await wrapper.find('input[placeholder="Contoh: Ponpes Al-Hidayah"]').setValue('Ponpes Test')
    await wrapper.find('input[placeholder="alhidayah01"]').setValue('test01')
    await wrapper.find('input[placeholder="Nama Lengkap"]').setValue('Admin Test')
    await wrapper.find('input[type="email"]').setValue('admin@test.com')
    
    // Cari input password dan isi
    const passwordInputs = wrapper.findAll('input[type="password"]')
    await passwordInputs[0].setValue('password123')
    
    // Isi username (sesuai urutan di template, username ada setelah email)
    // Kita cari input text yang belum terisi atau spesifik
    const textInputs = wrapper.findAll('input[type="text"]')
    // [Tips] Jika sulit mencari by index, cari input terakhir yang kosong atau by placeholder jika ada
    // Di template: Institution(0), Code(1), AdminName(2), Username(3)
    if (textInputs[3]) await textInputs[3].setValue('admintest')

    // Setup Mock Return Value
    mockPost.mockResolvedValue({})

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises() // Tunggu semua proses async selesai
    
    // Assert
    expect(mockPost).toHaveBeenCalled()
    expect(mockPost).toHaveBeenCalledWith('/auth/register-institution', expect.objectContaining({
      institution_code: 'test01',
      email: 'admin@test.com'
    }))
  })
})