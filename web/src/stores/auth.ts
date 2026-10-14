import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as api from '@/api'
import { getToken, setToken, clearToken } from '@/api/client'
import type { UserInfoRes } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(getToken())
  const userInfo = ref<UserInfoRes | null>(null)

  const perms = computed<string[]>(() => userInfo.value?.perms ?? [])

  // 超级管理员 perms 含 "*"，表示拥有全部权限。
  const hasPerm = (perm: string): boolean => {
    const p = perms.value
    return p.includes('*') || p.includes(perm)
  }

  async function login(username: string, password: string) {
    const res = await api.login({ username, password })
    token.value = res.token
    setToken(res.token)
    await fetchInfo()
  }

  async function fetchInfo() {
    userInfo.value = await api.getUserInfo()
  }

  async function logout() {
    try {
      await api.logout()
    } catch {
      // 忽略登出接口错误，本地照常清理
    }
    reset()
  }

  function reset() {
    token.value = ''
    userInfo.value = null
    clearToken()
  }

  return { token, userInfo, perms, hasPerm, login, fetchInfo, logout, reset }
})
