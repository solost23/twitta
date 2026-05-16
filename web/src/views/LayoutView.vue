<template>
  <el-container class="layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">🐦 Twitta</div>
      <el-menu :router="true" :default-active="route.path" class="menu">
        <el-menu-item index="/home"><el-icon><House /></el-icon>首页</el-menu-item>
        <el-menu-item index="/search"><el-icon><Search /></el-icon>搜索</el-menu-item>
        <el-menu-item index="/own"><el-icon><Document /></el-icon>我的推文</el-menu-item>
        <el-menu-item index="/favorites"><el-icon><Star /></el-icon>收藏</el-menu-item>
        <el-menu-item index="/friends" @click="notify.clear()">
          <el-icon><ChatDotRound /></el-icon>
          好友
          <el-badge v-if="notify.unread > 0" :value="notify.unread" class="nav-badge" />
        </el-menu-item>
        <el-menu-item index="/fans"><el-icon><User /></el-icon>关注/粉丝</el-menu-item>
        <el-menu-item index="/profile"><el-icon><Setting /></el-icon>个人资料</el-menu-item>
      </el-menu>
      <div class="user-bar">
        <el-avatar :src="ossUrl(auth.user?.avatar)" size="small" />
        <span class="username">{{ auth.user?.nickname || auth.user?.username }}</span>
        <el-button link @click="auth.logout()"><el-icon><SwitchButton /></el-icon></el-button>
      </div>
    </el-aside>
    <el-main class="main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotifyStore } from '@/stores/notify'
import { ElNotification } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const notify = useNotifyStore()

let ws: WebSocket | null = null

function buildNotifyWsUrl() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  const host = location.host
  const token = localStorage.getItem('token') || ''
  return `${proto}://${host}/api/twitta/notifications/ws?token=${token}`
}

function normalizeId(id: string): string {
  const m = id.match(/^ObjectID\("([a-f0-9]+)"\)$/i)
  return m ? m[1] : id
}

function connectNotifyWs() {
  if (!auth.isLoggedIn) return
  ws = new WebSocket(buildNotifyWsUrl())

  ws.onmessage = e => {
    try {
      const n = JSON.parse(e.data)
      // 如果当前已在该聊天页，不弹通知
      const chatPath = `/chat/${normalizeId(n.fromId)}`
      const altPath = `/chat/${n.fromId}`
      if (route.path === chatPath || route.path === altPath) return

      notify.increment()
      ElNotification({
        title: n.fromName || '新消息',
        message: n.content,
        type: 'info',
        duration: 4000,
        onClick: () => {
          notify.clear()
          router.push(`/chat/${n.fromId}`)
        }
      })
    } catch {}
  }

  ws.onclose = () => {
    // 断线后 3 秒重连
    setTimeout(() => { if (auth.isLoggedIn) connectNotifyWs() }, 3000)
  }
}

function disconnectNotifyWs() {
  ws?.close()
  ws = null
}

onMounted(() => { connectNotifyWs() })
onUnmounted(() => { disconnectNotifyWs() })

// 登录/登出时重连或断开
watch(() => auth.isLoggedIn, (loggedIn) => {
  if (loggedIn) connectNotifyWs()
  else disconnectNotifyWs()
})

function ossUrl(path?: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}
</script>

<style scoped>
.layout { height: 100vh; }
.sidebar {
  display: flex; flex-direction: column;
  background: #fff; border-right: 1px solid #e4e7ed;
}
.logo { font-size: 22px; font-weight: 700; padding: 20px 24px; color: #1da1f2; }
.menu { border: none; flex: 1; }
.user-bar {
  display: flex; align-items: center; gap: 8px;
  padding: 16px; border-top: 1px solid #e4e7ed;
}
.username { flex: 1; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.main { background: #f5f7fa; padding: 24px; overflow-y: auto; }
.nav-badge { margin-left: 6px; }
.nav-badge :deep(.el-badge__content) { transform: translateY(-4px); }
</style>
