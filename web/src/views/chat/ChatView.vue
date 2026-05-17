<template>
  <div class="chat-page">
    <el-card class="chat-card">
      <template #header>
        <div class="chat-header">
          <el-button :icon="ArrowLeft" link @click="router.back()" />
          <el-avatar :src="ossUrl(targetUser?.avatar)" size="small" />
          <span>{{ targetUser?.nickname || targetUser?.username || targetId }}</span>
          <el-tag :type="wsStatus === 'open' ? 'success' : 'danger'" size="small">
            {{ wsStatus === 'open' ? '已连接' : '未连接' }}
          </el-tag>
        </div>
      </template>

      <div class="load-more" v-if="hasMore">
        <el-button size="small" link @click="loadMore">加载更多历史消息</el-button>
      </div>

      <div class="messages" ref="msgContainer">
        <div
          v-for="(msg, i) in messages"
          :key="i"
          :class="['msg-row', isMine(msg) ? 'mine' : 'theirs']"
        >
          <!-- 对方消息：头像在左 -->
          <el-avatar v-if="!isMine(msg)" :src="ossUrl(targetUser?.avatar)" size="small" class="avatar" />
          <div class="bubble-wrap">
            <div class="sender-name" v-if="!isMine(msg)">
              {{ targetUser?.nickname || targetUser?.username }}
            </div>
            <div class="bubble">
              <div class="text">{{ msg.content || msg.msg }}</div>
            </div>
            <div class="meta">
              <span class="time">{{ formatTime(msg.createdAt) }}</span>
              <span v-if="isMine(msg)" class="read-status">
                {{ msg.read ? '已读' : '未读' }}
              </span>
            </div>
          </div>
          <!-- 自己消息：头像在右 -->
          <el-avatar v-if="isMine(msg)" :src="ossUrl(myAvatar)" size="small" class="avatar" />
        </div>
      </div>

      <div class="input-bar">
        <el-input
          v-model="inputText"
          placeholder="输入消息..."
          @keyup.enter="send"
          :disabled="wsStatus !== 'open'"
        />
        <el-button type="primary" :disabled="wsStatus !== 'open'" @click="send">发送</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { chatApi } from '@/api/chat'
import { authApi, type UserDetail } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { ArrowLeft } from '@element-plus/icons-vue'

interface DisplayMsg {
  fromId?: string; userId?: string
  content?: string; msg?: string
  createdAt: string; read: boolean
  _pending?: boolean
}

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const targetId = route.params.id as string
const myId = auth.user?.id || ''
const myAvatar = auth.user?.avatar || ''

const messages = ref<DisplayMsg[]>([])
const inputText = ref('')
const wsStatus = ref<'connecting' | 'open' | 'closed'>('connecting')
const msgContainer = ref<HTMLElement | null>(null)
const hasMore = ref(false)
const page = ref(1)
const targetUser = ref<UserDetail | null>(null)

let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null

function normalizeId(id: string): string {
  const m = id.match(/^ObjectID\("([a-f0-9]+)"\)$/i)
  return m ? m[1] : id
}

function isMine(msg: DisplayMsg): boolean {
  const uid = normalizeId(msg.fromId || msg.userId || '')
  return uid === normalizeId(myId)
}

function ossUrl(path?: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:9000/${path}`
}

function formatTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  if (isNaN(d.getTime())) return t
  const now = new Date()
  const isToday = d.toDateString() === now.toDateString()
  return isToday
    ? d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    : d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function buildWsUrl() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  const token = localStorage.getItem('token') || ''
  return `${proto}://${location.host}/api/twitta/chats/${targetId}/ws?token=${token}`
}

function connectWs() {
  if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null }
  ws = new WebSocket(buildWsUrl())
  wsStatus.value = 'connecting'

  ws.onopen = () => { wsStatus.value = 'open' }

  ws.onclose = () => {
    wsStatus.value = 'closed'
    // 断线自动重连
    reconnectTimer = setTimeout(connectWs, 3000)
  }

  ws.onerror = () => { wsStatus.value = 'closed' }

  ws.onmessage = e => {
    const msg = JSON.parse(e.data) as DisplayMsg
    msg.read = false
    if (isMine(msg)) {
      // 用服务端确认消息替换乐观消息，避免重复
      const idx = messages.value.findIndex(m => m._pending && (m.content || m.msg) === (msg.content || msg.msg))
      if (idx !== -1) {
        messages.value[idx] = { ...msg, _pending: false }
        return
      }
    }
    messages.value.push(msg)
    scrollBottom()
    // 收到对方消息立即标记已读
    if (!isMine(msg)) chatApi.markRead(targetId).catch(() => {})
  }
}

function send() {
  if (!inputText.value.trim() || !ws || ws.readyState !== WebSocket.OPEN) return
  const optimistic: DisplayMsg = {
    fromId: myId,
    content: inputText.value,
    createdAt: new Date().toISOString(),
    read: false,
    _pending: true,
  }
  messages.value.push(optimistic)
  scrollBottom()
  ws.send(JSON.stringify({ content: inputText.value }))
  inputText.value = ''
}

async function loadHistory(p = 1) {
  const res = await chatApi.history(targetId, p, 20)
  // 后端降序返回，reverse 后得到时间升序（旧→新），prepend 到消息列表顶部
  const older = [...res.records].reverse().map(r => ({
    userId: r.userId, msg: r.msg, createdAt: r.createdAt, read: r.read
  }))
  messages.value = [...older, ...messages.value]
  hasMore.value = p < res.pages
  page.value = p
}

async function loadMore() {
  await loadHistory(page.value + 1)
}

function scrollBottom() {
  nextTick(() => {
    if (msgContainer.value) msgContainer.value.scrollTop = msgContainer.value.scrollHeight
  })
}

// 处理来自通知 WS 的已读回执，更新消息状态
function handleReadReceipt(roomId: string) {
  messages.value.forEach(m => {
    if (isMine(m)) m.read = true
  })
}

// 监听全局通知 WS 的已读回执（通过 window 事件桥接）
function onGlobalNotify(e: Event) {
  const n = (e as CustomEvent).detail
  if (n?.type === 'read_receipt' && n.roomId) handleReadReceipt(n.roomId)
}

onMounted(async () => {
  const [_, detail] = await Promise.allSettled([
    loadHistory(1),
    authApi.userDetail(targetId)
  ])
  if (detail.status === 'fulfilled') targetUser.value = detail.value
  scrollBottom()
  connectWs()
  // 进入聊天页立即标记对方发给我的消息为已读
  chatApi.markRead(targetId).catch(() => {})
  window.addEventListener('twitta:notify', onGlobalNotify)
})

onUnmounted(() => {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  ws?.close()
  window.removeEventListener('twitta:notify', onGlobalNotify)
})
</script>

<style scoped>
.chat-page { height: calc(100vh - 96px); display: flex; flex-direction: column; }
.chat-card { flex: 1; display: flex; flex-direction: column; height: 100%; }
.chat-card :deep(.el-card__body) { flex: 1; display: flex; flex-direction: column; padding: 0; overflow: hidden; }
.chat-header { display: flex; align-items: center; gap: 8px; }
.load-more { text-align: center; padding: 8px; }
.messages {
  flex: 1; overflow-y: auto; padding: 16px;
  display: flex; flex-direction: column; gap: 12px;
}
.msg-row { display: flex; align-items: flex-end; gap: 8px; }
.msg-row.mine { flex-direction: row-reverse; }
.avatar { flex-shrink: 0; }
.bubble-wrap { display: flex; flex-direction: column; max-width: 60%; }
.mine .bubble-wrap { align-items: flex-end; }
.theirs .bubble-wrap { align-items: flex-start; }
.sender-name { font-size: 11px; color: #999; margin-bottom: 2px; }
.bubble { display: inline-block; }
.mine .bubble .text {
  background: #1da1f2; color: #fff;
  padding: 8px 12px; border-radius: 16px 16px 4px 16px;
  word-break: break-word;
}
.theirs .bubble .text {
  background: #f0f0f0; color: #333;
  padding: 8px 12px; border-radius: 16px 16px 16px 4px;
  word-break: break-word;
}
.meta { display: flex; align-items: center; gap: 4px; margin-top: 2px; }
.time { font-size: 11px; color: #bbb; }
.read-status { font-size: 11px; color: #bbb; }
.input-bar { display: flex; gap: 8px; padding: 12px 16px; border-top: 1px solid #e4e7ed; }
.input-bar .el-input { flex: 1; }
</style>
