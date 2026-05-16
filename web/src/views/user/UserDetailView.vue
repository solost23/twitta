<template>
  <div v-if="user">
    <el-card class="profile-card">
      <div class="profile-header">
        <el-avatar :src="ossUrl(user.avatar)" :size="72" />
        <div class="profile-info">
          <div class="name">{{ user.nickname || user.username }}</div>
          <div class="username">@{{ user.username }}</div>
          <div class="intro">{{ user.introduce }}</div>
          <div class="stats">
            <span>关注 {{ user.whatCount }}</span>
            <span>粉丝 {{ user.fansCount }}</span>
          </div>
        </div>
        <div class="actions" v-if="!isMe">
          <el-button type="primary" size="small" @click="follow">关注</el-button>
          <el-button size="small" @click="showApply = true">发好友申请</el-button>
          <el-button size="small" @click="router.push(`/chat/${route.params.id}`)">发消息</el-button>
        </div>
      </div>
    </el-card>

    <el-dialog v-model="showApply" title="发送好友申请" width="360px">
      <el-input v-model="applyContent" type="textarea" :rows="3" placeholder="附言（必填）" />
      <template #footer>
        <el-button @click="showApply = false">取消</el-button>
        <el-button type="primary" @click="sendApply">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authApi, type UserDetail } from '@/api/auth'
import { socialApi } from '@/api/social'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const user = ref<UserDetail | null>(null)
const showApply = ref(false)
const applyContent = ref('')
const isMe = computed(() => auth.user?.id === route.params.id)

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

async function follow() {
  await socialApi.follow(route.params.id as string)
  ElMessage.success('已关注')
}

async function sendApply() {
  if (!applyContent.value.trim()) { ElMessage.warning('请填写附言'); return }
  await socialApi.sendApplication(route.params.id as string, applyContent.value)
  ElMessage.success('申请已发送')
  showApply.value = false
}

onMounted(async () => {
  user.value = await authApi.userDetail(route.params.id as string)
})
</script>

<style scoped>
.profile-card { margin-bottom: 16px; }
.profile-header { display: flex; gap: 16px; align-items: flex-start; }
.profile-info { flex: 1; }
.name { font-size: 18px; font-weight: 700; }
.username { color: #999; font-size: 13px; }
.intro { margin: 6px 0; font-size: 14px; }
.stats { display: flex; gap: 16px; font-size: 13px; color: #555; }
.actions { display: flex; flex-direction: column; gap: 6px; }
</style>
