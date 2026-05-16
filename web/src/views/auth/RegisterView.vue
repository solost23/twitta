<template>
  <div class="auth-page">
    <el-card class="auth-card">
      <h2>注册 Twitta</h2>
      <el-form :model="form" @submit.prevent="submit" label-position="top">
        <el-form-item label="用户名 *">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码 *">
          <el-input v-model="form.password" type="password" placeholder="至少6位" show-password />
        </el-form-item>
        <el-form-item label="邮箱 *">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="可选" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.mobile" placeholder="可选" />
        </el-form-item>
        <el-form-item label="个人简介">
          <el-input v-model="form.introduce" type="textarea" :rows="2" placeholder="可选" />
        </el-form-item>
        <el-form-item label="头像">
          <el-upload action="#" :before-upload="uploadAvatar" :show-file-list="false" accept="image/*">
            <el-button>上传头像</el-button>
          </el-upload>
          <el-avatar v-if="form.avatar" :src="form.avatar" style="margin-left:12px" />
        </el-form-item>
        <el-button type="primary" native-type="submit" :loading="loading" style="width:100%">注册</el-button>
      </el-form>
      <div class="links">
        <router-link to="/login">已有账号？去登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const form = reactive({
  username: '', password: '', email: '',
  nickname: '', mobile: '', introduce: '', avatar: ''
})

async function uploadAvatar(file: File) {
  const url = await authApi.uploadAvatar(file)
  form.avatar = url
  return false
}

async function submit() {
  if (!form.username || !form.password || !form.email) {
    ElMessage.warning('请填写必填项')
    return
  }
  loading.value = true
  try {
    await authApi.register(form)
    // 记录注册历史
    const list = JSON.parse(localStorage.getItem('register_records') || '[]')
    list.unshift({ username: form.username, time: new Date().toLocaleString('zh-CN') })
    localStorage.setItem('register_records', JSON.stringify(list.slice(0, 5)))
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f7fa; }
.auth-card { width: 420px; }
h2 { text-align: center; margin-bottom: 24px; color: #1da1f2; }
.links { margin-top: 16px; text-align: center; font-size: 13px; }
</style>
