import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type LoginResponse } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<LoginResponse | null>(
    JSON.parse(localStorage.getItem('user') || 'null')
  )

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(data: LoginResponse) {
    token.value = data.token
    user.value = data
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data))
  }

  function logout() {
    authApi.logout().catch(() => {})
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
  }

  return { token, user, isLoggedIn, setAuth, logout }
})
