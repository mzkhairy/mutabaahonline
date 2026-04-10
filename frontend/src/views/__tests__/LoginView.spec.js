import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginView from '../LoginView.vue' // File ini sejajar, jadi path ini BENAR
import { createPinia, setActivePinia } from 'pinia'

// Mock Router & Auth Store
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush })
}))

describe('LoginView.vue', () => {
  it('renders login form completely', () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: {
        stubs: ['router-link'] // [FIX] Hilangkan warning router-link
      }
    })
    
    expect(wrapper.text()).toContain('Mutabaah App')
    const inputs = wrapper.findAll('input')
    expect(inputs.length).toBe(3)
    expect(wrapper.find('button').text()).toContain('Masuk Sekarang')
  })

  it('validates input handling', async () => {
    setActivePinia(createPinia())
    const wrapper = mount(LoginView, {
      global: { stubs: ['router-link'] }
    })

    await wrapper.find('input[placeholder="Contoh: tahfidz01"]').setValue('tahfidz01')
    await wrapper.find('input[type="password"]').setValue('rahasia')
    
    expect(wrapper.vm.form.institution_code).toBe('tahfidz01')
    expect(wrapper.vm.form.password).toBe('rahasia')
  })
})