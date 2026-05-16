<template>
  <div class="auth-page">
    <el-card class="auth-card">
      <h2>登录 Twitta</h2>
      <el-form :model="form" @submit.prevent="submit" label-position="top">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-button type="primary" native-type="submit" :loading="loading" style="width:100%">登录</el-button>
      </el-form>
      <div class="links">
        <router-link to="/register">还没有账号？立即注册</router-link>
      </div>

      <!-- 登录记录 -->
      <div v-if="records.length" class="records">
        <div class="records-title">最近登录记录</div>
        <div v-for="r in records" :key="r.time" class="record-row">
          <span class="record-user">{{ r.username }}</span>
          <span class="record-time">{{ r.time }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

interface LoginRecord { username: string; time: string }

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const form = reactive({ username: '', password: '' })
const records = ref<LoginRecord[]>([])

function loadRecords() {
  try {
    records.value = JSON.parse(localStorage.getItem('login_records') || '[]')
  } catch { records.value = [] }
}

function saveRecord(username: string) {
  const list: LoginRecord[] = JSON.parse(localStorage.getItem('login_records') || '[]')
  list.unshift({ username, time: new Date().toLocaleString('zh-CN') })
  localStorage.setItem('login_records', JSON.stringify(list.slice(0, 5)))
}

async function submit() {
  loading.value = true
  try {
    const res = await authApi.login({ ...form, platform: 'twitta' })
    saveRecord(form.username)
    auth.setAuth(res)
    router.push('/home')
  } finally {
    loading.value = false
  }
}

onMounted(loadRecords)
</script>

<style scoped>
.auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f7fa; }
.auth-card { width: 380px; }
h2 { text-align: center; margin-bottom: 24px; color: #1da1f2; }
.links { margin-top: 16px; text-align: center; font-size: 13px; }
.records { margin-top: 20px; border-top: 1px solid #f0f0f0; padding-top: 12px; }
.records-title { font-size: 12px; color: #999; margin-bottom: 8px; }
.record-row { display: flex; justify-content: space-between; font-size: 13px; padding: 4px 0; }
.record-user { color: #333; font-weight: 500; }
.record-time { color: #bbb; font-size: 12px; }
</style>
