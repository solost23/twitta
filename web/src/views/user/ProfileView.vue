<template>
  <el-card>
    <h3>个人资料</h3>
    <el-form :model="form" label-width="80px" style="max-width:480px;margin-top:16px">
      <el-form-item label="头像">
        <el-upload action="#" :before-upload="uploadAvatar" :show-file-list="false" accept="image/*">
          <el-avatar :src="ossUrl(form.avatar)" :size="64" style="cursor:pointer" />
        </el-upload>
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item label="简介">
        <el-input v-model="form.introduce" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const auth = useAuthStore()
const saving = ref(false)
const form = reactive({ username: '', nickname: '', introduce: '', avatar: '' })

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

async function uploadAvatar(file: File) {
  const url = await authApi.uploadAvatar(file)
  form.avatar = url
  return false
}

async function save() {
  saving.value = true
  try {
    await authApi.userUpdate(form)
    ElMessage.success('保存成功')
    if (auth.user) {
      auth.user.nickname = form.nickname
      auth.user.avatar = form.avatar
      localStorage.setItem('user', JSON.stringify(auth.user))
    }
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (auth.user) {
    form.username = auth.user.username || ''
    form.nickname = auth.user.nickname || ''
    form.avatar = auth.user.avatar || ''
  }
})
</script>
