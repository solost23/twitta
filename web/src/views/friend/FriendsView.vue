<template>
  <div>
    <el-tabs v-model="tab">
      <!-- 好友列表 -->
      <el-tab-pane label="好友" name="friends">
        <div v-for="f in friends" :key="f.userId" class="user-row">
          <el-avatar :src="ossUrl(f.avatar)" size="small" @click="router.push(`/user/${f.userId}`)" style="cursor:pointer" />
          <div class="info" style="flex:1">
            <div class="name">{{ f.nickname || f.username }}</div>
            <div class="intro">{{ f.introduce }}</div>
          </div>
          <el-button size="small" @click="router.push(`/chat/${f.userId}`)">聊天</el-button>
          <el-button size="small" type="danger" plain @click="deleteFriend(f.userId)">删除</el-button>
        </div>
        <el-empty v-if="!friends.length" description="还没有好友" />
      </el-tab-pane>

      <!-- 好友申请 -->
      <el-tab-pane label="好友申请" name="applications">
        <div v-for="a in applications" :key="a.userId + a.createdAt" class="app-row">
          <el-avatar :src="ossUrl(a.avatar)" size="small" />
          <div class="info" style="flex:1">
            <div class="name">{{ a.username }}</div>
            <div class="intro">{{ a.content }}</div>
            <div class="time">{{ a.createdAt }}</div>
          </div>
          <template v-if="a.type === 0">
            <el-button size="small" type="primary" @click="accept(a.userId)">通过</el-button>
            <el-button size="small" @click="reject(a.userId)">拒绝</el-button>
          </template>
          <el-tag v-else-if="a.type === 1" type="success" size="small">已通过</el-tag>
          <el-tag v-else type="info" size="small">已拒绝</el-tag>
        </div>
        <el-empty v-if="!applications.length" description="暂无申请" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { socialApi, type FriendItem, type FriendApplication } from '@/api/social'
import { ElMessage } from 'element-plus'

const router = useRouter()
const tab = ref('friends')
const friends = ref<FriendItem[]>([])
const applications = ref<FriendApplication[]>([])

function ossUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

async function accept(id: string) {
  await socialApi.acceptApplication(id)
  ElMessage.success('已通过')
  loadAll()
}

async function reject(id: string) {
  await socialApi.rejectApplication(id)
  ElMessage.success('已拒绝')
  loadAll()
}

async function deleteFriend(id: string) {
  await socialApi.deleteFriend(id)
  ElMessage.success('已删除')
  friends.value = friends.value.filter(f => f.userId !== id)
}

async function loadAll() {
  const [f, a] = await Promise.all([socialApi.friendList(), socialApi.applicationList()])
  friends.value = f
  applications.value = a
}

onMounted(loadAll)
</script>

<style scoped>
.user-row, .app-row {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 0; border-bottom: 1px solid #f0f0f0;
}
.info .name { font-weight: 600; font-size: 14px; }
.info .intro { font-size: 12px; color: #999; }
.info .time { font-size: 11px; color: #bbb; }
</style>
